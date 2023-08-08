/*
 *   Copyright (c) 2023 CodapeWild
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package io

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/CodapeWild/devkit/directory"
	"google.golang.org/protobuf/proto"
)

var _ PubPubBatchAndFetchFetchBatch = (*FileCache)(nil)

type FileCache struct {
	sync.Mutex
	path                      string
	seqDir                    *directory.SequentialDirectory // cache data in sequential read/write directory
	readChan, writeChan       chan *IOMessage
	pageSize                  int          // number of entries count
	readPageName              string       // the current file name where the data of readPageBuf from
	readPageBuf, writePageBuf []*IOMessage // buffer for reading and writing
	readIndex                 int          // indicating the index position for reading start from 0 to pageSize-1
	writeIndex                int          // indicating the index position for writing start from 0 to pageSize-1
	writePause, writeResume   chan struct{}
	closer                    chan struct{}
}

func (fc *FileCache) Publish(ctx context.Context, message *IOMessage) (*IOResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			if _, ok := ctx.Deadline(); ok {
				return InputTimeout, nil
			} else {
				return InputFailed, err
			}
		}
	case <-fc.closer:
		return nil, ErrIOClosed
	case fc.writeChan <- message:
	}

	return InputSuccess, nil
}

func (fc *FileCache) PublishBatch(ctx context.Context, batch *IOMessageBatch) (*IOResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	for _, msg := range batch.IOMessageBatch {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				if _, ok := ctx.Deadline(); ok {
					return InputTimeout, nil
				} else {
					return InputFailed, err
				}
			}
		case <-fc.closer:
			return nil, ErrIOClosed
		case fc.writeChan <- msg:
		}
	}

	return InputSuccess, nil
}

func (fc *FileCache) Fetch(ctx context.Context) (*IOMessage, *IOResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	select {
	case <-fc.closer:
		return nil, OutputFailed, ErrIOClosed
	default:
	}

	fc.Lock()
	defer fc.Unlock()

	// load new page from SequentialDirectory(disk) if empty load data from writePageBuf into readPageBuf
	if (fc.readIndex == fc.pageSize-1) || (fc.readPageBuf[fc.readIndex] == nil) {
		fname, bts, err := fc.seqDir.OpenAndDelete("")
		if err != nil {
			// directory is empty try to load data from writePageBuf
			if errors.Is(err, directory.ErrDirEmpty) {
				if fc.writeIndex != 0 {
					fc.writePause <- struct{}{}
					fc.readPageBuf = fc.writePageBuf
					fc.readIndex = -1
					fc.writePageBuf = make([]*IOMessage, fc.pageSize)
					fc.writeIndex = -1
					fc.writeResume <- struct{}{}
				} else {
					return nil, OutputDataEmpty, nil
				}
			} else {
				return nil, OutputFailed, err
			}
		} else {
			batch := &IOMessageBatch{}
			if err = proto.Unmarshal(bts.Bytes(), batch); err != nil {
				return nil, OutputFailed, err
			}
			if len(batch.IOMessageBatch) != fc.pageSize {
				return nil, OutputFailed, ErrWrongDataSetLength
			}
			fc.readPageName = fname
			fc.readPageBuf = batch.IOMessageBatch
			fc.readIndex = -1
		}
	}

	fc.readIndex++

	return fc.readPageBuf[fc.readIndex], OutputSuccess, nil
}

// FetchBatch returns messages batch the number of message count depends:
// - readPageBuf not empty and SequentialDirectory not empty then return readIndex + pageSize
// - readPageBuf empty and SequentialDirectory not empty then return pageSize
// - readPageBuf not empty and SequentialDirectory empty and writePageBuf not empty then return readIndex + writeIndex
// - readPageBuf empty and SequentialDirectory empty and writePageBuf not empty then return writeIndex
// - readPageBuf empty and SequentialDirectory empty and writePageBuf empty return 0
// the order of returning data is readPageBuf, SequentialDirectory, writePageBuf
func (fc *FileCache) FetchBatch(ctx context.Context) (*IOMessageBatch, *IOResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	select {
	case <-fc.closer:
		return nil, OutputFailed, ErrIOClosed
	default:
	}

	fc.Lock()
	defer fc.Unlock()

	var list []*IOMessage = fc.readPageBuf
	_, bts, err := fc.seqDir.OpenAndDelete("")
	if err != nil {
		if errors.Is(err, directory.ErrDirEmpty) {
			if fc.writeIndex != -1 {
				fc.writePause <- struct{}{}
				list = append(list, fc.writePageBuf[:fc.writeIndex+1]...)
				fc.writePageBuf = make([]*IOMessage, fc.pageSize)
				fc.writeIndex = -1
				fc.writeResume <- struct{}{}

				fc.readPageBuf = make([]*IOMessage, fc.pageSize)
				fc.readIndex = -1

				return &IOMessageBatch{IOMessageBatch: list}, OutputSuccess, nil
			} else {
				return nil, OutputDataEmpty, nil
			}
		} else {
			return nil, OutputFailed, err
		}
	}

	batch := &IOMessageBatch{}
	if err = proto.Unmarshal(bts.Bytes(), batch); err != nil {
		return nil, OutputFailed, err
	}
	list = append(list, batch.IOMessageBatch...)

	fc.readPageBuf = make([]*IOMessage, fc.pageSize)
	fc.readIndex = -1

	return &IOMessageBatch{IOMessageBatch: list}, OutputSuccess, nil
}

func (fc *FileCache) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-fc.closer:
		return ErrIOClosed
	default:
	}

	// start write thread
	go func() {
	BEFORE_EXITS:
		for {
			select {
			case <-fc.closer:
				break BEFORE_EXITS
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Println(err.Error())
				}
				break BEFORE_EXITS
			case <-fc.writePause:
				<-fc.writeResume
			case msg := <-fc.writeChan:
				if err := fc.writeRoutine(msg); err != nil {
					log.Println(err.Error())
				}
			}
		}

		if err := fc.bufferToDisk(); err != nil {
			log.Println(err.Error())
		}
	}()

	return nil
}

func (fc *FileCache) Close() {
	select {
	case <-fc.closer:
	default:
		close(fc.closer)
	}
}

func (fc *FileCache) writeRoutine(message *IOMessage) error {
	fc.writeIndex++
	fc.writePageBuf[fc.writeIndex] = message
	// move data into SequentialDirectory(disk)
	if fc.writeIndex == fc.pageSize-1 {
		if bts, err := proto.Marshal(&IOMessageBatch{IOMessageBatch: fc.writePageBuf}); err != nil {
			return err
		} else {
			if err = fc.seqDir.Save("", bytes.NewBuffer(bts)); err != nil {
				return err
			}
		}
		fc.writePageBuf = make([]*IOMessage, fc.pageSize)
		fc.writeIndex = -1
	}

	return nil
}

// bufferToDisk writes data in readPageBuf and writePageBuf back to directory(disk)
func (fc *FileCache) bufferToDisk() error {
	if fc.readIndex != -1 && fc.readPageName != "" {
		bts, err := proto.Marshal(&IOMessageBatch{IOMessageBatch: fc.readPageBuf})
		if err != nil {
			return err
		}
		if err = os.WriteFile(fc.readPageName, bts, 0644); err != nil {
			return err
		}
	}
	if fc.writeIndex != -1 {
		if bts, err := proto.Marshal(&IOMessageBatch{IOMessageBatch: fc.writePageBuf}); err != nil {
			return err
		} else {
			if err = fc.seqDir.Save("", bytes.NewBuffer(bts)); err != nil {
				return err
			}
		}
	}

	return nil
}

func OpenFileCache(dir string, pageSize int) (*FileCache, error) {
	seqDir, err := directory.OpenSequentialDirectory(dir)
	if err != nil {
		return nil, err
	}

	cache := pageSize / 2
	if cache == 0 {
		pageSize = 20
		cache = 10
	}

	return &FileCache{
		path:         dir,
		seqDir:       seqDir,
		readChan:     make(chan *IOMessage, cache),
		writeChan:    make(chan *IOMessage, cache),
		pageSize:     pageSize,
		readPageBuf:  make([]*IOMessage, pageSize),
		writePageBuf: make([]*IOMessage, pageSize),
		readIndex:    -1,
		writeIndex:   -1,
		closer:       make(chan struct{}),
	}, nil
}
