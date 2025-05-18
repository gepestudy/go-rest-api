package middlewares

import (
	"log"
	"net/http"
	"time"
)

func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		now := time.Now()
		next.ServeHTTP(rw, r)
		responseTime := time.Since(now)
		// rw.Header().Set("X-Response-Time", responseTime.String()) // masih ga bisa buat set header ketika sudah selesai. mungkin butuh middleware lagi di akhir
		log.Printf("path: %s, method: %s, status: %d, response time: %v\n",
			r.URL.Path, r.Method, rw.statusCode, responseTime)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
