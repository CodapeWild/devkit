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

import "net/http"

func WorkerPoolHTTPWrapper(wp *WorkerPool, next http.Handler) http.Handler {
	if wp == nil {
		return next
	} else {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			// todo: integrate with net.http message to complete this job
			// copy body from request and send it to next handler
		})
	}
}
