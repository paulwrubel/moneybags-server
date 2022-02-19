package routing

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type logrusResponseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func (lrw *logrusResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

func (lrw *logrusResponseWriter) WriteHeader(statusCode int) {
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.status = statusCode
}

func logrusMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lrw := &logrusResponseWriter{ResponseWriter: rw}
			next.ServeHTTP(lrw, r)

			log.WithFields(log.Fields{
				"path":     r.URL.RequestURI(),
				"method":   r.Method,
				"status":   lrw.status,
				"size":     lrw.size,
				"duration": time.Since(start),
			}).Debug("request completed")
		})
	}
}
