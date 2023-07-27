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
	"math"
	"sync"
	"time"
)

var max = int64(math.MaxInt64)

type IDFlaker struct {
	sync.Mutex
	ts, seq int64
}

func (flk *IDFlaker) NextID() *ID {
	flk.Lock()
	defer flk.Unlock()

	flk.seq++
	flk.seq &= max
	if flk.seq == 0 {
		now := time.Now().UnixMilli()
		for flk.ts == now {
			now = time.Now().UnixMilli()
		}
		flk.ts = now
	}
	// if now == flk.ts {
	// 	if flk.seq < math.MaxInt64 {
	// 		flk.seq++
	// 	} else {
	// 		for {
	// 			if now = time.Now().UnixMilli(); now != flk.ts {
	// 				flk.ts = now
	// 				flk.seq = 0
	// 				break
	// 			}
	// 		}
	// 	}
	// } else {
	// 	flk.ts = now
	// 	flk.seq = 0
	// }

	return &ID{high: flk.ts, low: flk.seq}
}

func NewIDFlaker() *IDFlaker {
	return &IDFlaker{ts: time.Now().UnixMilli()}
}
