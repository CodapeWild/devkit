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

package workerpool

import (
	"context"
)

type JobProcess func(ctx context.Context, out chan interface{}) error

type JobCallback func(out interface{}, err error)

type Job interface {
	Process(ctx context.Context, out chan interface{}) error
	Callback(out interface{}, err error)
}

type JobWrapper struct {
	proc JobProcess
	cb   JobCallback
}

func (jw *JobWrapper) Process(ctx context.Context, out chan interface{}) error {
	return jw.proc(ctx, out)
}

func (jw *JobWrapper) Callback(out interface{}, err error) {
	jw.cb(out, err)
}

func NewJobWrapper(job Job) *JobWrapper {
	return &JobWrapper{
		proc: job.Process,
		cb:   job.Callback,
	}
}

func NewJobWrapperFromFunc(process JobProcess, callback JobCallback) *JobWrapper {
	return &JobWrapper{
		proc: process,
		cb:   callback,
	}
}

func NewJobWrapperWithContext(ctx context.Context, job Job) *JobWrapper {
	return &JobWrapper{
		proc: func(_ context.Context, out chan interface{}) error { return job.Process(ctx, out) },
		cb:   job.Callback,
	}
}
