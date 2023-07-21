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

type Remapper map[string]string

func (rp Remapper) Remap(source interface{}) error {
	refmap := reflect.ValueOf(source)
	if refmap.Kind() != reflect.Map {
		return comerr.ErrAssertFailed
	}

	for old, new := range rp {
		refoldk := reflect.ValueOf(old)
		refoldv := refmap.MapIndex(refoldk)
		if refoldv.IsValid() {
			refmap.SetMapIndex(refoldk, reflect.Value{})
		}
		if len(new) != 0 {
			refmap.SetMapIndex(reflect.ValueOf(new), refoldv)
		}
	}

	return nil
}
