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

package set

import (
	"github.com/CodapeWild/devkit/comerr"
)

var _ QueueSet = (*SingleThreadQueue)(nil)

type queopt byte

const (
	que_push      queopt = 101
	que_pop       queopt = 102
	que_async_pop queopt = 103
)

type stqOptWrapper struct {
	opt   queopt
	value any
	cb    func(any)
}

type SingleThreadQueue struct {
	que           []any
	opts          chan *stqOptWrapper
	pause, resume chan struct{}
	closer        chan struct{}
}

func (stq *SingleThreadQueue) Push(value any) error {
	stq.opts <- &stqOptWrapper{opt: que_push, value: value}

	return nil
}

func (stq *SingleThreadQueue) Pop() (any, error) {
	if len(stq.que) == 0 {
		return nil, comerr.ErrEmptyValue
	}

	stq.pause <- struct{}{}
	v := stq.que[0]
	stq.que = stq.que[1:]
	stq.resume <- struct{}{}

	return v, nil
}

func (stq *SingleThreadQueue) AsyncPop(callback func(value any)) error {
	if len(stq.que) == 0 {
		return comerr.ErrEmptyValue
	}

	stq.opts <- &stqOptWrapper{opt: que_async_pop, cb: callback}

	return nil
}

func (stq *SingleThreadQueue) Peek() any {
	if len(stq.que) != 0 {
		return stq.que[0]
	} else {
		return nil
	}
}

func (stq *SingleThreadQueue) Close() {
	select {
	case <-stq.closer:
	default:
		close(stq.closer)
	}
}

func (stq *SingleThreadQueue) startThread() {
	select {
	case <-stq.closer:
		return
	default:
	}

	go func() {
		for {
			select {
			case <-stq.closer:
				for wrapper := range stq.opts {
					stq.routine(wrapper)
				}

				return
			case <-stq.pause:
				<-stq.resume
			case wrapper := <-stq.opts:
				stq.routine(wrapper)
			}
		}
	}()
}

func (stq *SingleThreadQueue) routine(wrapper *stqOptWrapper) {
	switch wrapper.opt {
	case que_push:
		stq.que = append(stq.que, wrapper.value)
	case que_async_pop:
		v := stq.que[0]
		stq.que = stq.que[1:]
		go func(value any, cb func(value any)) {
			cb(value)
		}(v, wrapper.cb)
	}
}

func NewSingleThreadQueue(bufSize int) *SingleThreadQueue {
	stq := &SingleThreadQueue{
		opts:   make(chan *stqOptWrapper, bufSize),
		pause:  make(chan struct{}),
		resume: make(chan struct{}),
		closer: make(chan struct{}),
	}
	stq.startThread()

	return stq
}
