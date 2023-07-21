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

package http

import (
	"net/http"

	"github.com/CodapeWild/devkit/iterator"
)

func CheckHeaders(next, failed http.HandlerFunc, target map[string][]string) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			if ss, ok := target[k]; ok {
				if !iterator.Include(v, ss) {
					failed(resp, req)

					return
				}
			}
		}
		next(resp, req)
	}
}
