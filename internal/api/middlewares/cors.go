package middlewares

import (
	"fmt"
	"net/http"
)

var allowedOrigins = []string{
	"https://localhost:8080",
	"https://localhost:8081",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Println(origin)

		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")

			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
		}
	})
}

func isOriginAllowed(origin string) bool {
	for _, ao := range allowedOrigins {
		if ao == origin {
			return true
		}
	}
	return false
}
