package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"user-services/models"
	"user-services/respository"
)

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	userIDFromToken := r.Context().Value("userID").(uint)

	// Convert idParam to uint
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Ensure the requested ID matches the ID from the token
	if uint(id) != userIDFromToken {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := repository.DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}
