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
		go func() {
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			wp.SendJob(ctx, &mockJob{})
		}()
	}
	time.Sleep(100 * time.Millisecond)

	wp.Close()
	time.Sleep(30 * time.Millisecond)
}
