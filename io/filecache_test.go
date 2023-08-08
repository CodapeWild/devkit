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

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"testing"
	"time"
)

func mockIOMessage(d time.Duration, l, c int, out chan *IOMessage) {
	for i := 0; i < c; i++ {
		buf := make([]byte, l)
		rand.Read(buf)
		out <- NewIOMessage(IOMessageWithCoding("bytes"), IOMessageWithPayload(buf))
		time.Sleep(d)
	}
	close(out)
}

func TestFileCachePublish(t *testing.T) {
	fc, err := OpenFileCache("./test", 10)
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = fc.Start(context.TODO()); err != nil {
		t.Fatal(err.Error())
	}

	var (
		threads  = 10
		finished = make(chan struct{})
	)
	for i := 0; i < threads; i++ {
		t.Run(fmt.Sprintf("publish_%d", i), func(t *testing.T) {
			t.Parallel()

			out := make(chan *IOMessage)
			go mockIOMessage(10*time.Millisecond, 1000, 10, out)
			for msg := range out {
				if resp, err := fc.Publish(context.TODO(), msg); err != nil {
					t.Fatal(err.Error())
				} else {
					log.Printf("%#v", *resp)
				}
			}
			finished <- struct{}{}
		})
	}

	t.Run("close", func(t *testing.T) {
		t.Parallel()

		var c int
		for range finished {
			if c++; c == threads {
				fc.Close()
				log.Println("FileCache closed")
				break
			}
		}
	})
}
