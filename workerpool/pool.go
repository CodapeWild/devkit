/*
 *   Copyright (c) 2023 CodapeWild
 *   All rights reserved.

 *   Permission is hereby granted, free of charge, to any person obtaining a copy
 *   of this software and associated documentation files (the "Software"), to deal
 *   in the Software without restriction, including without limitation the rights
 *   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *   copies of the Software, and to permit persons to whom the Software is
 *   furnished to do so, subject to the following conditions:

 *   The above copyright notice and this permission notice shall be included in all
 *   copies or substantial portions of the Software.

 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *   SOFTWARE.
 */

package workerpool

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrWorkerPoolClosed = errors.New("worker pool has closed")
	ErrSendJobTimeout   = errors.New("sending job to worker pool timeout")
	ErrTaskTimeout      = errors.New("task timeout")
)

func NewWorkerPool(n int) *WorkerPool {
	return &WorkerPool{
		n:      n,
		tasks:  make(chan jobContext),
		closer: make(chan struct{}),
	}
}

type WorkerPool struct {
	sync.Once
	n      int
	tasks  chan jobContext
	closer chan struct{}
}

func (wp *WorkerPool) Start() {
	wp.Do(func() {
		for i := 0; i < wp.n; i++ {
			go wp.workLoop()
		}
	})
}

func (wp *WorkerPool) SendJob(ctx context.Context, job Job) error {
	if ctx == nil {
		ctx = context.Background()
	} else {
		if err := ctx.Err(); err != nil {
			return err
		}
	}

	select {
	case <-wp.closer:
		return ErrWorkerPoolClosed
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			return err
		}
	case wp.tasks <- jobWrapper(ctx, job):
	}

	return nil
}

func (wp *WorkerPool) workLoop() {
	for {
		select {
		case <-wp.closer:
			return
		case jobctx := <-wp.tasks:
			ctx, job := jobctx()
			ctxc, canceler := context.WithCancel(ctx)
			out := make(chan interface{})
			go job.Process(ctxc, out)

			select {
			case <-wp.closer:
				canceler()

				return
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					canceler()
					job.Callback(nil, err)
				}
			case rslt := <-out:
				job.Callback(rslt, nil)
			}
		}
	}
}

func (wp *WorkerPool) Close() {
	select {
	case <-wp.closer:
	default:
		close(wp.closer)
	}
}
