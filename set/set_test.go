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
	"log"
	"testing"
)

func TestIntDataSet(t *testing.T) {
	var temp = []int{2, 4, 6, 8, 3, 4, 3, 6, 7}
	ds := IntDataSet(temp)
	ok := ds.Remove(4, 0)
	log.Println(ok)
	log.Println(ds)
}