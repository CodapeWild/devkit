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

var _ Queue = (*SingleThreadQueue)(nil)

type queopt byte

const (
	que_peek queopt = 101
	que_push queopt = 102
	que_pop  queopt = 103
)

type stqOptWrapper struct {
	opt   queopt
	value any
	out   chan any
}

type SingleThreadQueue struct {
	que    []any
	opts   chan *stqOptWrapper
	closer chan struct{}
}

func (stq *SingleThreadQueue) Peek() any {
	out := make(chan any)
	stq.opts <- &stqOptWrapper{opt: que_peek, out: out}

	return <-out
}

func (stq *SingleThreadQueue) Push(value any) error {
	stq.opts <- &stqOptWrapper{opt: que_push, value: value}

	return nil
}

func (stq *SingleThreadQueue) Pop() (any, error) {
	out := make(chan any)
	stq.opts <- &stqOptWrapper{opt: que_pop, out: out}

	return <-out, nil
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
					stq.thread(wrapper)
				}

				return
			case wrapper := <-stq.opts:
				stq.thread(wrapper)
			}
		}
	}()
}

func (stq *SingleThreadQueue) thread(wrapper *stqOptWrapper) {
	switch wrapper.opt {
	case que_peek:
		if len(stq.que) != 0 {
			wrapper.out <- stq.que[0]
		} else {
			wrapper.out <- nil
		}
	case que_push:
		stq.que = append(stq.que, wrapper.value)
	case que_pop:
		if len(stq.que) != 0 {
			wrapper.out <- stq.que[0]
			stq.que = stq.que[1:]
		} else {
			wrapper.out <- nil
		}
	}
}

func NewSingleThreadQueue(bufSize int) *SingleThreadQueue {
	stq := &SingleThreadQueue{
		opts:   make(chan *stqOptWrapper, bufSize),
		closer: make(chan struct{}),
	}
	stq.startThread()

	return stq
}
