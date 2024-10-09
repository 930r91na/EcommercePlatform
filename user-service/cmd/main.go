package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"user-services/handlers"
	repository "user-services/respository"
)

func main() {
	// Initialize the database
	err := repository.InitDB()

	if err != nil {
		return
	}

	// Initialize router
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("User Service is running"))
		if err != nil {
			return
		}
	}).Methods("GET")

	// User Routes
	r.HandleFunc("/users/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/users/login", handlers.LoginUser).Methods("POST")

	// Start Server
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
