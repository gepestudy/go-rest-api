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
	mux.HandleFunc("/teachers/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTeachersHandler(w, r, db)
		case http.MethodPost:
			handlers.AddTeacherHandler(w, r, db)
		}
	})
	mux.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from students route!"))
	})
	mux.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from execs route!"))
	})

	return mux
}
