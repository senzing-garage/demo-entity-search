package httpserver

import (
	"net/http"
	"time"
)

func addIncomingRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		then := time.Now()

		defer func() {
			if recovered := recover(); recovered != nil {
				createRequestLog(request).Info("request errored out")
			}
		}()

		next.ServeHTTP(writer, request)

		duration := time.Since(then)
		createRequestLog(
			request,
		).Infof("request completed in %vms", float64(duration.Nanoseconds())/NanosecondsPerMillisecond)
	})
}
