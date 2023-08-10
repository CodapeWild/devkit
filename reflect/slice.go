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

package reflect

import (
	"reflect"

	"github.com/CodapeWild/devkit/comerr"
)

func Swap(slice any, i, j int) error {
	if i < 0 || j < 0 {
		return comerr.ErrInvalidParameters
	}
	refslc := reflect.ValueOf(slice)
	if k := refslc.Kind(); k != reflect.Slice || k != reflect.Array {
		return comerr.ErrAssertFailed
	}
	if l := refslc.Len(); l < i || l < j {
		return comerr.ErrIndexOverflow
	}

	reflect.Swapper(slice)(i, j)

	return nil
}
