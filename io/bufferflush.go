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

var _ PubAndSubBatch = (*BufferFlush)(nil)

type BufferFlush struct {
	cur, maxSize int
	flushTick    time.Ticker
	msgchan      chan proto.Message
	batch        []proto.Message
	handler      SubscribeMessageBatchHandler
	closer       chan struct{}
}

func (bf *BufferFlush) Start() {
	select {
	case <-bf.closer:
		return
	default:
	}
	go func() {
		for {
			select {
			case <-bf.closer:
				return
			case <-bf.flushTick.C:
				if bf.cur > 0 {
					bf.handler(bf.batch)
					bf.batch = make([]proto.Message, bf.maxSize)
					bf.cur = 0
				}
			case msg := <-bf.msgchan:
				bf.batch[bf.cur] = msg
				if bf.cur++; bf.cur == bf.maxSize {
					bf.handler(bf.batch)
					bf.batch = make([]proto.Message, bf.maxSize)
					bf.cur = 0
				}
			}
		}
	}()
}

func (bf *BufferFlush) Publish(ctx context.Context, message proto.Message) (*IOResponse, error) {
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
	case bf.msgchan <- message:
	}

	return InputSuccess, nil
}

func (bf *BufferFlush) SubscribeBatch(handler SubscribeMessageBatchHandler) error {
	bf.handler = handler

	return nil
}

func (bf *BufferFlush) Close() {
	select {
	case <-bf.closer:
	default:
		close(bf.closer)
	}
}

func NewCacheAndFlush(maxSize int, d time.Duration) *BufferFlush {
	cache := maxSize / 4
	if cache == 0 {
		cache = 1
	}

	return &BufferFlush{
		maxSize:   maxSize,
		flushTick: *time.NewTicker(d),
		msgchan:   make(chan proto.Message, cache),
		batch:     make([]proto.Message, maxSize),
	}
}
