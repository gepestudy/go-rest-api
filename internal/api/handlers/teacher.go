package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gepestudy/go-rest-api/internal/models"
)

var teachers = []models.Teacher{
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

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")

	if idStr == "" {
		firstNameQuery := r.URL.Query().Get("first_name")
		lastNameQuery := r.URL.Query().Get("last_name")

		teacherList := make([]models.Teacher, 0, len(teachers))

		for _, teacher := range teachers {
			if (firstNameQuery == "" || teacher.FirstName == firstNameQuery) && (lastNameQuery == "" || teacher.LastName == lastNameQuery) {
				teacherList = append(teacherList, teacher)
			}
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
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
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}

	for _, teacher := range teachers {
		if teacher.ID == id {
			response := struct {
				Status string         `json:"status"`
				Data   models.Teacher `json:"data"`
			}{
				Status: "success",
				Data:   teacher,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("something went wrong"))
				return
			}
			return
		}
	}

	// Jika teacher tidak ditemukan
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Teacher not found"))
}
