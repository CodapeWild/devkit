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

package threadpool

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

type WorkerPool struct {
	sync.Once
	maxThreads int
	jobchan    chan Job
	closer     chan struct{}
}

func (wp *WorkerPool) Start() {
	wp.Do(func() {
		for i := 0; i < wp.maxThreads; i++ {
			go wp.workLoop()
		}
	})
}

func (wp *WorkerPool) SendJob(ctx context.Context, job Job) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	select {
	case <-wp.closer:
		return ErrWorkerPoolClosed
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			return err
		}
	case wp.jobchan <- NewJobWrapperWithContext(ctx, job):
	}

	return nil
}

func (wp *WorkerPool) workLoop() {
	for {
		select {
		case <-wp.closer:
			return
		case job := <-wp.jobchan:
			go func(job Job) {
				out := make(chan interface{})
				err := job.Process(nil, out)
				if jw, ok := job.(*JobWrapper); ok && jw.cb != nil {
					select {
					case <-wp.closer:
						return
					case o := <-out:
						job.Callback(o, err)
					}
				}
			}(job)
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

func NewWorkerPool(n int) *WorkerPool {
	return &WorkerPool{
		maxThreads: n,
		jobchan:    make(chan Job),
		closer:     make(chan struct{}),
	}
}
