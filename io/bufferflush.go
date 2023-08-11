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
	"log"
	"time"

	"github.com/CodapeWild/devkit/message"
)

var _ PubAndSub = (*BufferFlush)(nil)

type BufferFlush struct {
	cur, maxSize int
	flushTick    time.Ticker
	msgChan      chan SubscribeMessageHandler
	buffer       []SubscribeMessageHandler
	handler      SubscribeMessageHandler
	closer       chan struct{}
}

func (bf *BufferFlush) Publish(ctx context.Context, msg message.Message) *IOResponse {
	if err := ctx.Err(); err != nil {
		return InputFailed.With(IORespWithMessage(err.Error()))
	}
	// if bf.handler == nil {
	// 	return InputFailed.With(IORespWithMessage(ErrSubscribeHandlerUnset.Error()))
	// }

	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			if _, ok := ctx.Deadline(); ok {
				return InputTimeout
			} else {
				return InputFailed.With(IORespWithMessage(err.Error()))
			}
		}
	case <-bf.closer:
		return IOClosed
	case bf.msgChan <- bf.handler.BindContext(ctx, msg):
	}

	return InputSuccess
}

func (bf *BufferFlush) Subscribe(handler SubscribeMessageHandler) error {
	bf.handler = handler

	return nil
}

func (bf *BufferFlush) Start(ctx context.Context) error {
	if bf.handler == nil {
		return ErrIOUncompleted
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	select {
	case <-bf.closer:
		return ErrIOClosed
	default:
	}

	go func() {
		for {
			select {
			case <-bf.closer:
				return
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Println(err.Error())
				}

				return
			case <-bf.flushTick.C:
				if bf.cur > 0 {
					bf.doFlush()
				}
			case msg := <-bf.msgChan:
				bf.buffer[bf.cur] = msg
				bf.cur++
				if bf.cur == bf.maxSize {
					bf.doFlush()
				}
			}
		}
	}()

	return nil
}

func (bf *BufferFlush) Close() {
	select {
	case <-bf.closer:
	default:
		close(bf.closer)
	}
}

func (bf *BufferFlush) doFlush() {
	for _, handler := range bf.buffer {
		resp := handler(nil, nil)
		if resp != nil && resp != OutputSuccess {
			log.Printf("do flush failed: %#v", resp)
			continue
		}
	}
	bf.cur = 0
	bf.buffer = make([]SubscribeMessageHandler, bf.maxSize)
}

func NewBufferFlush(maxSize int, d time.Duration) *BufferFlush {
	cache := maxSize / 2
	if cache == 0 {
		maxSize = 20
		cache = 10
	}

	return &BufferFlush{
		maxSize:   maxSize,
		flushTick: *time.NewTicker(d),
		msgChan:   make(chan SubscribeMessageHandler, cache),
		buffer:    make([]SubscribeMessageHandler, maxSize),
		closer:    make(chan struct{}),
	}
}
