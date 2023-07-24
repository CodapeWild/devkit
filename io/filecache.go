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
	"time"

	"google.golang.org/protobuf/proto"
)

var _ PubPubBatchAndSubSubBatch = (*FileCache)(nil)

type FileCache struct {
	dir           string // pages path
	pageSize      int    // kb
	pageBuf       chan proto.Message
	pior          int // current page index of reading
	flushTick     time.Ticker
	handleMessage SubscribeMessageHandler
	handleBatch   SubscribeMessageBatchHandler
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
	case fc.pageBuf <- message:
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
		case fc.pageBuf <- msg:
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
