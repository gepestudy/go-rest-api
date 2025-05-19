package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	mw "github.com/gepestudy/go-rest-api/internal/api/middlewares"
)

func main() {
	cert := "cert.pem"
	key := "key.pem"

	rl := mw.NewRatelimiter(2, 5*time.Second)

	port := 8080

	// multiplexer
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Root route!"))
	})
	mux.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Teachers route!"))
	})
	mux.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from students route!"))
	})
	mux.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from execs route!"))
	})

	tslConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// create server with TLS
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: tslConfig,
		Handler:   rl.Middleware(mw.ResponseTime(mw.SecurityHeaders(mw.Cors(mw.Compression(mux))))),
	}

	fmt.Println("Starting server on port", port)
	if err := server.ListenAndServeTLS(cert, key); err != nil {
		log.Fatalln(err)
	}
}
