package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"user-services/handlers"
	"user-services/middleware"
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

	// OAuth Routes
	oauthRouter := r.PathPrefix("/auth/").Subrouter()
	oauthRouter.Path("/google/login").HandlerFunc(handlers.GoogleLogin).Methods("GET")
	oauthRouter.Path("/google/callback").HandlerFunc(handlers.GoogleCallback).Methods("GET")
	oauthRouter.Path("/github/login").HandlerFunc(handlers.GithubLogin).Methods("GET")
	oauthRouter.Path("/github/callback").HandlerFunc(handlers.GithubCallback).Methods("GET")

	// Protected routes
	profileRouter := r.PathPrefix("/users/profile").Subrouter()
	profileRouter.Use(middleware.JWTMiddleware)
	profileRouter.HandleFunc("/{id}", handlers.GetUserProfile).Methods("GET")
	// profileRouter.HandleFunc("/{id}", handlers.UpdateUserProfile).Methods("PUT")

	// Start Server
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
