package workerpool

import "net/http"

func WorkerPoolHTTPWrapper(wp *WorkerPool, next http.Handler) http.Handler {
	if wp == nil {
		return next
	} else {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		})
	}
}
