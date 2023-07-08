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

import "time"

type CacheAndFlush struct {
	cur, maxSize int
	tick         time.Ticker
	datachan     chan any
	dataset      []any
	flusher      func(data any)
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
					cf.flusher(cf.dataset)
					cf.dataset = make([]any, cf.maxSize)
					cf.cur = 0
				}
			case data := <-cf.datachan:
				cf.dataset[cf.cur] = data
				if cf.cur++; cf.cur == cf.maxSize {
					cf.flusher(cf.dataset)
					cf.dataset = make([]any, cf.maxSize)
					cf.cur = 0
				}
			}
		}
	}()
}

func (cf *CacheAndFlush) Publish(data any) {
	cf.datachan <- data
}

func (cf *CacheAndFlush) Subscribe(flusher func(data any)) error {
	cf.flusher = flusher

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
		maxSize:  maxSize,
		tick:     *time.NewTicker(d),
		datachan: make(chan any, cache),
		dataset:  make([]any, maxSize),
	}
}
