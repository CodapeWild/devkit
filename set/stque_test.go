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
	"crypto/rand"
	"testing"
)

func TestSTQuePushAndPop(t *testing.T) {
	stq := NewSingleThreadQueue(10)
	for i := 0; i < 100; i++ {
		t.Run("stq:push", func(t *testing.T) {
			buf := make([]byte, 1000)
			rand.Read(buf)
			if err := stq.Push(buf); err != nil {
				t.Fatal(err.Error())
			}
		})
	}
	for i := 0; i < 100; i++ {
		t.Run("stq:pop", func(t *testing.T) {
			if _, err := stq.Pop(); err != nil {
				t.Fatal(err.Error())
			}
		})
	}
	if stq.Peek() != nil {
		t.Fatal("single thread queue not work as expeccted")
	}
}
