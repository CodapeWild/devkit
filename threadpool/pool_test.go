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
	"log"
	"testing"
	"time"
)

type mockJob struct{}

func (mj *mockJob) Process(ctx context.Context, out chan interface{}) error {
	time.Sleep(30 * time.Millisecond)
	if out != nil {
		out <- 123
	}

	return nil
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
	time.Sleep(3 * time.Second)

	wp.Close()
}
