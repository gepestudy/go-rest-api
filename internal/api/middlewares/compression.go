package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// chec if the client accepts gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// client accepts gzip encoding
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("GEPEGANTENG", "TRUE")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// warp response writter
		w = gzipResponseWriter{ResponseWriter: w, writer: gz}

		next.ServeHTTP(w, r)
	})
}

// gzipResponseWriter is a wrapper around http.ResponseWriter that provides gzip compression.
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer // override the ResponseWriter's Write method with gzip writter for enabling response compression
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}
