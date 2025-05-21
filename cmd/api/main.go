package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	mw "github.com/gepestudy/go-rest-api/internal/api/middlewares"
)

type Teacher struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Class     string `json:"class"`
	Subject   string `json:"subject"`
}

var teachers = []Teacher{
	{ID: 1, FirstName: "John", LastName: "Doe", Class: "10", Subject: "Math"},
	{ID: 2, FirstName: "Jane", LastName: "Doe", Class: "10", Subject: "Science"},
	{ID: 3, FirstName: "Jack", LastName: "Doe", Class: "10", Subject: "English"},
	{ID: 4, FirstName: "Jill", LastName: "Doe", Class: "10", Subject: "History"},
	{ID: 5, FirstName: "Joe", LastName: "Doe", Class: "10", Subject: "Geography"},
	{ID: 6, FirstName: "Judy", LastName: "Doe", Class: "10", Subject: "Biology"},
	{ID: 7, FirstName: "Jenny", LastName: "Doe", Class: "10", Subject: "Chemistry"},
	{ID: 8, FirstName: "Jake", LastName: "Doe", Class: "10", Subject: "Physics"},
	{ID: 9, FirstName: "Jill", LastName: "Doe", Class: "10", Subject: "Music"},
	{ID: 10, FirstName: "Jim", LastName: "Doe", Class: "10", Subject: "Art"},
}
var (
	mutex  = &sync.Mutex{}
	nextID = 1
)

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	firstNameQuery := r.URL.Query().Get("first_name")
	lastNameQuery := r.URL.Query().Get("last_name")

	teacherList := make([]Teacher, 0, len(teachers))

	for _, teacher := range teachers {
		if (firstNameQuery == "" || teacher.FirstName == firstNameQuery) && (lastNameQuery == "" || teacher.LastName == lastNameQuery) {
			teacherList = append(teacherList, teacher)
		}
	}

	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teacherList),
		Data:   teacherList,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(response))
}

func main() {
	cert := "cert.pem"
	key := "key.pem"

	rl := mw.NewRatelimiter(10, 10*time.Second)

	port := 8080

	// multiplexer
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Root route!"))
	})
	mux.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTeachersHandler(w, r)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Fuck Off"))
		}
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
		Handler:   ApplyMiddleware(mux, mw.Cors, rl.Middleware, mw.ResponseTime, mw.Compression, mw.SecurityHeaders),
	}

	fmt.Println("Starting server on port", port)
	if err := server.ListenAndServeTLS(cert, key); err != nil {
		log.Fatalln(err)
	}
}

type Middleware func(http.Handler) http.Handler

func ApplyMiddleware(next http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		next = middleware(next)
	}
	return next
}
