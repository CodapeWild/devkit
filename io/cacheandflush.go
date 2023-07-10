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
)

type CacheAndFlush struct {
	cur, maxSize int
	tick         time.Ticker
	msgchan      chan *IOMessage
	batch        []*IOMessage
	handler      func(batch []*IOMessage) *IOResponse
	closer       chan struct{}
}

func (cf *CacheAndFlush) Start() {
	select {
	case <-cf.closer:
		return
	default:
	}
	go func() {
		for {
			select {
			case <-cf.closer:
				return
			case <-cf.tick.C:
				if cf.cur > 0 {
					cf.handler(cf.batch)
					cf.batch = make([]*IOMessage, cf.maxSize)
					cf.cur = 0
				}
			case msg := <-cf.msgchan:
				cf.batch[cf.cur] = msg
				if cf.cur++; cf.cur == cf.maxSize {
					cf.handler(cf.batch)
					cf.batch = make([]*IOMessage, cf.maxSize)
					cf.cur = 0
				}
			}
		}
	}()
}

func (cf *CacheAndFlush) Publish(ctx context.Context, message *IOMessage) (*IOResponse, error) {
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
	case cf.msgchan <- message:
	}

	return InputSuccess, nil
}

func (cf *CacheAndFlush) SubscribeBatch(ctx context.Context, handler func(batch []*IOMessage) *IOResponse) error {
	cf.handler = handler

	return nil
}

func (cf *CacheAndFlush) Close() {
	select {
	case <-cf.closer:
	default:
		close(cf.closer)
	}
}

func NewCacheAndFlush(maxSize int, d time.Duration) *CacheAndFlush {
	cache := maxSize / 4
	if cache == 0 {
		cache = 1
	}

	return &CacheAndFlush{
		maxSize: maxSize,
		tick:    *time.NewTicker(d),
		msgchan: make(chan *IOMessage, cache),
		batch:   make([]*IOMessage, maxSize),
	}
}
