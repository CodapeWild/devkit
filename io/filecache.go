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
	"context"

	"github.com/CodapeWild/devkit/directory"
	"google.golang.org/protobuf/proto"
)

var _ PubPubBatchAndSubSubBatch = (*FileCache)(nil)

type FileCache struct {
	seqDir        *directory.SequentialDirectory // cache data in sequential read/write directory
	pageSize      int                            // kb
	msgChan       chan proto.Message
	buffer        []proto.Message
	cur           int
	handleMessage SubscribeMessageHandler
	handleBatch   SubscribeMessageBatchHandler
	closer        chan struct{}
}

func (fc *FileCache) Publish(ctx context.Context, message proto.Message) (*IOResponse, error) {
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
	case fc.msgChan <- message:
	}

	return InputSuccess, nil
}

func (fc *FileCache) PublishBatch(ctx context.Context, batch []proto.Message) (*IOResponse, error) {
	if err := ctx.Err(); err != nil {
		return InputFailed, err
	}

	for _, msg := range batch {
		select {
		case <-ctx.Done():
		case fc.msgChan <- msg:
		}
	}

	return InputSuccess, nil
}

func (fc *FileCache) Subscribe(handler SubscribeMessageHandler) error {
	fc.handleMessage = handler

	return nil
}

func (fc *FileCache) SubscribeBatch(handler SubscribeMessageBatchHandler) error {
	fc.handleBatch = handler

	return nil
}

func (fc *FileCache) Start() {

}

func (fc *FileCache) Close() {
	select {
	case <-fc.closer:
	default:
		close(fc.closer)
	}
}

func OpenFileCache(dir string, pageSize int) (*FileCache, error) {
	seqDir, err := directory.OpenSequentialDirectory(dir)
	if err != nil {
		return nil, err
	}

	cache := pageSize / 2
	if cache == 0 {
		cache = 10
	}
	return &FileCache{
		seqDir:   seqDir,
		pageSize: pageSize,
		msgChan:  make(chan proto.Message, cache),
		buffer:   make([]proto.Message, pageSize),
		closer:   make(chan struct{}),
	}, nil
}
