package workerpool

import (
	"context"
	"log"
	"testing"
	"time"
)

type mockJob struct{}

func (mj *mockJob) Process(ctx context.Context, out chan interface{}) {
	time.Sleep(30 * time.Millisecond)
	if out != nil {
		out <- 123
	}
}

func (mj *mockJob) Callback(out interface{}, err error) {
	if err != nil {
		log.Println(err.Error())

		return
	}

	log.Printf("job done with output: %v", out)
}

func TestRunWorkerPool(t *testing.T) {
	wp := NewWorkerPool(10)
	wp.Start()

	for i := 0; i < 100; i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		wp.SendJob(ctx, &mockJob{})
	}

	wp.Close()
	time.Sleep(30 * time.Millisecond)
}
