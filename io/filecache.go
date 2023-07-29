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

	"github.com/CodapeWild/devkit/directory"
	"google.golang.org/protobuf/proto"
)

var _ PubPubBatchAndFetchFetchBatch = (*FileCache)(nil)

type FileCache struct {
	seqDir                    *directory.SequentialDirectory // cache data in sequential read/write directory
	readChan, writeChan       chan *IOMessage
	pageSize                  int          // number of entries count
	readPageBuf, writePageBuf []*IOMessage // buffer for reading and writing
	readIndex                 int          // indicating the index position for reading start from 0 to pageSize
	writeIndex                int          // indicating the index position for writing start from 0 to pageSize
	closer                    chan struct{}
}

func (fc *FileCache) Publish(ctx context.Context, message *IOMessage) (*IOResponse, error) {
	if err := ctx.Err(); err != nil {
		return InputFailed, err
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
		return IOClosed, nil
	case fc.writeChan <- message:
	}

	return InputSuccess, nil
}

func (fc *FileCache) PublishBatch(ctx context.Context, batch *IOMessageBatch) (*IOResponse, error) {
	if err := ctx.Err(); err != nil {
		return InputFailed, err
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
			return IOClosed, nil
		case fc.writeChan <- msg:
		}
	}

	return InputSuccess, nil
}

func (fc *FileCache) Fetch(ctx context.Context) (*IOMessage, *IOResponse, error) {

}

func (fc *FileCache) FetchBatch(ctx context.Context) (*IOMessageBatch, *IOResponse, error) {

}

func (fc *FileCache) Start() {
	select {
	case <-fc.closer:
		return
	default:
	}

	// start write thread
	go func() {
		for {
			select {
			case <-fc.closer:
				return
			case msg := <-fc.writeChan:
				fc.writeRoutine(msg)
			}
		}
	}()
}

func (fc *FileCache) Close() {
	select {
	case <-fc.closer:
	default:
		close(fc.closer)
	}
}

func (fc *FileCache) writeRoutine(message *IOMessage) error {
	fc.writePageBuf[fc.writeIndex] = message
	fc.writeIndex++
	if fc.writeIndex == fc.pageSize {
		batch := &IOMessageBatch{IOMessageBatch: fc.writePageBuf}
		if bts, err := proto.Marshal(batch); err != nil {
			return err
		} else {
			if err = fc.seqDir.Save("", bytes.NewBuffer(bts)); err != nil {
				return err
			}
		}
		fc.writePageBuf = make([]*IOMessage, fc.pageSize)
		fc.writeIndex = 0
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
		seqDir:       seqDir,
		readChan:     make(chan *IOMessage, cache),
		writeChan:    make(chan *IOMessage, cache),
		pageSize:     pageSize,
		readPageBuf:  make([]*IOMessage, pageSize),
		writePageBuf: make([]*IOMessage, pageSize),
		closer:       make(chan struct{}),
	}, nil
}
