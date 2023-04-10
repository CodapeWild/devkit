package workerpool

import (
	"context"
)

type Job interface {
	Process(ctx context.Context, out chan interface{})
	Callback(out interface{}, err error)
}

type jobContext func() (ctx context.Context, job Job)

func jobWrapper(ctx context.Context, job Job) jobContext {
	return func() (context.Context, Job) {
		return ctx, job
	}
}
