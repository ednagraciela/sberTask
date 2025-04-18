package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"t/db"
	"t/models"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Настройка тестовой базы данных
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "todo_test")

	db.InitDB()
	defer db.CloseDB()

	// Очистка базы данных перед тестами
	db.DB.Exec("DROP TABLE IF EXISTS todos")
	db.DB.Exec(`
		CREATE TABLE todos (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			date DATE NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)

	code := m.Run()
	os.Exit(code)
}

func TestCreateTodo(t *testing.T) {
	router := setupRouter()

	todo := models.TodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		Date:        time.Now(),
		Completed:   false,
	}

	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.Todo
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, todo.Title, response.Title)
	assert.Equal(t, todo.Description, response.Description)
	assert.Equal(t, todo.Completed, response.Completed)
}

func TestGetTodos(t *testing.T) {
	router := setupRouter()

	// Создаем тестовые данные
	todo := models.TodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		Date:        time.Now(),
		Completed:   false,
	}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)

	// Тестируем получение списка
	req, _ = http.NewRequest("GET", "/todos", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var todos []models.Todo
	json.Unmarshal(rr.Body.Bytes(), &todos)

	assert.GreaterOrEqual(t, len(todos), 1)
}

func TestGetTodo(t *testing.T) {
	router := setupRouter()

	// Создаем тестовые данные
	todo := models.TodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		Date:        time.Now(),
		Completed:   false,
	}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var createdTodo models.Todo
	json.Unmarshal(rr.Body.Bytes(), &createdTodo)

	// Тестируем получение одной задачи
	req, _ = http.NewRequest("GET", "/todos/"+strconv.Itoa(createdTodo.ID), nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Todo
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, createdTodo.ID, response.ID)
	assert.Equal(t, createdTodo.Title, response.Title)
}

func TestUpdateTodo(t *testing.T) {
	router := setupRouter()

	// Создаем тестовые данные
	todo := models.TodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		Date:        time.Now(),
		Completed:   false,
	}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var createdTodo models.Todo
	json.Unmarshal(rr.Body.Bytes(), &createdTodo)

	// Обновляем задачу
	updatedTodo := models.TodoRequest{
		Title:       "Updated Todo",
		Description: "Updated Description",
		Date:        time.Now().Add(24 * time.Hour),
		Completed:   true,
	}
	body, _ = json.Marshal(updatedTodo)
	req, _ = http.NewRequest("PUT", "/todos/"+strconv.Itoa(createdTodo.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Todo
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, createdTodo.ID, response.ID)
	assert.Equal(t, updatedTodo.Title, response.Title)
	assert.Equal(t, updatedTodo.Completed, response.Completed)
}

func TestDeleteTodo(t *testing.T) {
	router := setupRouter()

	// Создаем тестовые данные
	todo := models.TodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		Date:        time.Now(),
		Completed:   false,
	}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var createdTodo models.Todo
	json.Unmarshal(rr.Body.Bytes(), &createdTodo)

	// Удаляем задачу
	req, _ = http.NewRequest("DELETE", "/todos/"+strconv.Itoa(createdTodo.ID), nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Проверяем, что задача удалена
	req, _ = http.NewRequest("GET", "/todos/"+strconv.Itoa(createdTodo.ID), nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/todos", GetTodos).Methods("GET")
	router.HandleFunc("/todos", CreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", GetTodo).Methods("GET")
	router.HandleFunc("/todos/{id}", UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", DeleteTodo).Methods("DELETE")
	return router
}

func TestMain(m *testing.M) {
	// Запуск контейнеров
	compose := testcontainers.NewLocalDockerCompose([]string{"docker-compose.yml"}, "todo-test")
	execError := compose.WithCommand([]string{"up", "-d"}).Invoke()
	if execError.Error != nil {
		log.Fatalf("Failed to start containers: %v", execError.Error)
	}

	// Ожидание готовности БД
	waitForDB()

	// Запуск тестов
	code := m.Run()

	// Остановка контейнеров
	_ = compose.Down()

	os.Exit(code)
}

func waitForDB() {
	var err error
	for i := 0; i < 10; i++ {
		db.DB, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=todo_test sslmode=disable")
		if err == nil {
			err = db.DB.Ping()
			if err == nil {
				return
			}
		}
		time.Sleep(2 * time.Second)
	}
	log.Fatal("Failed to connect to test database")
}
