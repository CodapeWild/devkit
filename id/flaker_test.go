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

package id

import (
	"testing"
)

func TestGenNextID(t *testing.T) {
	gen := NewIDFlaker()
	n := 100000
	saved := make(map[string]bool)
	for i := 0; i < n; i++ {
		t.Run("TestMultiThreadsGenNextID", func(t *testing.T) {
			id := gen.NextID().String('-')
			if saved[id] {
				t.Fatal("duplicated id")
			}
			saved[id] = true
		})
	}
}

var id *ID

func BenchmarkGenNextID(b *testing.B) {
	gen := NewIDFlaker()
	for i := 0; i < b.N; i++ {
		id = gen.NextID()
	}
}
