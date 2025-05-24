package router

import (
	"database/sql"
	"net/http"

	"github.com/gepestudy/go-rest-api/internal/api/handlers"
)

func InitRouter(mux *http.ServeMux, db *sql.DB) *http.ServeMux {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Root route!"))
	})

	//* teacher handler
	mux.HandleFunc("GET /teachers", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTeachersHandler(w, r, db)
	})
	mux.HandleFunc("GET /teachers/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTeacherHandler(w, r, db)
	})

	mux.HandleFunc("POST /teachers", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddTeacherHandler(w, r, db)
	})

	mux.HandleFunc("PUT /teachers/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateTeacherHandler(w, r, db)
	})

	mux.HandleFunc("PATCH /teachers/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.PatchTeacherHandler(w, r, db)
	})

	mux.HandleFunc("DELETE /teachers/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteTeacherHandler(w, r, db)
	})

	//* student handler
	mux.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from students route!"))
	})

	//* exec handler
	mux.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from execs route!"))
	})

	return mux
}
