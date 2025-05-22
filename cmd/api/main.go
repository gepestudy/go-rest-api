package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	mw "github.com/gepestudy/go-rest-api/internal/api/middlewares"
	"github.com/gepestudy/go-rest-api/internal/api/router"
	"github.com/gepestudy/go-rest-api/internal/repository/sqlconnect"
	"github.com/gepestudy/go-rest-api/pkg/config"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	cert := "cert.pem"
	key := "key.pem"

	port := 8080

	// multiplexer
	mux := http.NewServeMux()
	mux = router.InitRouter(mux, db)

	tslConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	secureMux := mw.SecurityHeaders(mux)
	// create server with TLS
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: tslConfig,
		Handler:   secureMux,
	}

	fmt.Println("Starting server on port", port)
	if err := server.ListenAndServeTLS(cert, key); err != nil {
		log.Fatalln(err)
	}
}
