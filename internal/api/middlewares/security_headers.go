package middlewares

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-DNS-Prefetch-Control", "off") // disable DNS prefetching
		w.Header().Set("X-Frame-Options", "off")        // prevent clickjacking
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")             // prevent MIME sniffing
		w.Header().Set("Content-Security-Policy", "default-src 'self'") // enable HSTS
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}
