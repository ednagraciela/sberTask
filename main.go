package main

import (
	"log"
	"net/http"
	"os"
	"t/db"
	"t/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/todos", handlers.GetTodos).Methods("GET")
	r.HandleFunc("/todos", handlers.CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", handlers.GetTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", handlers.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", handlers.DeleteTodo).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Todo API is running. Use /todos endpoints"))
	})
}
