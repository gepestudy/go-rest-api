package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gepestudy/go-rest-api/internal/models"
)

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	_, ok := validFields[field]
	return ok
}

func GetTeachersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1`
	var args []any

	query, args = addFilters(r, query, args)

	// /teachers?sortby=first_name:asc,class:desc menjadi [first_name:asc, class:desc] kalo pakai .Get("sortby") bakal nge return string, kalo pakai Query()["xxx"] bakal nge return array of string
	query = addSorting(r, query)

	rows, err := db.Query(query, args...)
	var teacherList []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to scan teacher data"))
			return
		}
		teacherList = append(teacherList, teacher)
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
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong 1"))
		return
	}
}

func GetTeacherHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong 2"))
		return
	}

	query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?`
	stmt, err := db.Prepare(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong 3"))
		return
	}
	defer stmt.Close()

	var teacher models.Teacher
	if err := stmt.QueryRow(id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			response := struct {
				Status string          `json:"status"`
				Count  int             `json:"count"`
				Data   *models.Teacher `json:"data"`
			}{
				Status: "failed",
				Count:  0,
				Data:   nil,
			}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong 4"))
		return
	}
	response := struct {
		Status string         `json:"status"`
		Count  int            `json:"count"`
		Data   models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  1,
		Data:   teacher,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong 5"))
		return
	}
}

func addSorting(r *http.Request, query string) string {
	sortParam := r.URL.Query()["sortby"]
	if len(sortParam) > 0 {
		query += " ORDER BY "

		for i, param := range sortParam {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}

			field, order := parts[0], parts[1]
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}

			if i > 0 {
				query += ", "
			}

			query += fmt.Sprintf("%s %s", field, order)
		}
	}
	return query
}

func addFilters(r *http.Request, query string, args []any) (string, []any) {
	params := map[string]string{
		"first_name": r.URL.Query().Get("first_name"),
		"last_name":  r.URL.Query().Get("last_name"),
		"email":      r.URL.Query().Get("email"),
		"class":      r.URL.Query().Get("class"),
		"subject":    r.URL.Query().Get("subject"),
	}
	for key, value := range params {
		if value != "" {
			query += fmt.Sprintf(" AND %s = ?", key)
			args = append(args, value)
		}
	}
	return query, args
}

// POST /teachers
func AddTeacherHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var teachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&teachers)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body is not valid"))
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}
	defer stmt.Close()

	for i, teacher := range teachers {
		res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong"))
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong"))
			return
		}
		teacher.ID = int(id)
		teachers[i].ID = int(id)
	}
	json.NewEncoder(w).Encode(teachers)
}

// PUT /teachers/{id}
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	var teacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&teacher)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body is not valid"))
		return
	}

	stmt, err := db.Prepare("UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject, id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}
	json.NewEncoder(w).Encode(teacher)
}

func PatchTeacherHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	var existingTeacher models.Teacher
	stmt, err := db.Prepare("SELECT * FROM teachers WHERE id=?")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error preparing query"))
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Teacher not found"))
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error retrieving teacher"))
		return
	}

	var patchData models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	existingTeacherElem := reflect.ValueOf(&existingTeacher).Elem()
	patch := reflect.ValueOf(patchData)

	for i := 0; i < patch.NumField(); i++ {
		field := patch.Field(i)
		if !field.IsZero() {
			existingTeacherElem.Field(i).Set(field)
		}
	}

	stmt, err = db.Prepare("UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error preparing query"))
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error updating teacher"))
		return
	}
	json.NewEncoder(w).Encode(existingTeacher)
}

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 || idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	stmt, err := db.Prepare("DELETE FROM teachers WHERE id=?")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
