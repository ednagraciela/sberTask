package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"t/db"
	"t/models"

	"github.com/gorilla/mux"
)

func GetTodos(w http.ResponseWriter, r *http.Request) {

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	status := r.URL.Query().Get("status")
	date := r.URL.Query().Get("date")

	query := "SELECT id, title, description, date, completed, created_at, updated_at FROM todos"
	var args []interface{}
	var conditions []string

	if status != "" {
		completed := status == "completed"
		conditions = append(conditions, "completed = $"+strconv.Itoa(len(args)+1))
		args = append(args, completed)
	}

	if date != "" {
		_, err := time.Parse("2006-01-02", date)
		if err == nil {
			conditions = append(conditions, "date = $"+strconv.Itoa(len(args)+1))
			args = append(args, date)
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " LIMIT $" + strconv.Itoa(len(args)+1)
	args = append(args, limit)
	query += " OFFSET $" + strconv.Itoa(len(args)+1)
	args = append(args, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Date,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todoReq models.TodoRequest
	err := json.NewDecoder(r.Body).Decode(&todoReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if todoReq.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	err = db.DB.QueryRow(
		"INSERT INTO todos (title, description, date, completed) VALUES ($1, $2, $3, $4) RETURNING id, title, description, date, completed, created_at, updated_at",
		todoReq.Title,
		todoReq.Description,
		todoReq.Date,
		todoReq.Completed,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Date,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	err = db.DB.QueryRow(
		"SELECT id, title, description, date, completed, created_at, updated_at FROM todos WHERE id = $1",
		id,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Date,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var todoReq models.TodoRequest
	err = json.NewDecoder(r.Body).Decode(&todoReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if todoReq.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	err = db.DB.QueryRow(
		"UPDATE todos SET title = $1, description = $2, date = $3, completed = $4, updated_at = NOW() WHERE id = $5 RETURNING id, title, description, date, completed, created_at, updated_at",
		todoReq.Title,
		todoReq.Description,
		todoReq.Date,
		todoReq.Completed,
		id,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Date,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
