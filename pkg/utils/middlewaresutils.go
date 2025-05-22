package utils

import "net/http"

type Middleware func(http.Handler) http.Handler

func ApplyMiddleware(next http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		next = middleware(next)
	}
	return next
}
