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

import "net/http"

func MergeHeaders(h1, h2 http.Header) http.Header {
	dst := make(http.Header)
	for k, v := range h1 {
		dst[k] = make([]string, len(v))
		copy(dst[k], v)
	}
	for k, v := range h2 {
		if _, ok := dst[k]; ok {
			for _, u := range v {
				dst[k] = append(dst[k], u)
			}
		}
	}

	return dst
}
