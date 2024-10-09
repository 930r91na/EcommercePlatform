package handlers

import (
	"encoding/json"
	_ "encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	_ "net/http"
	"user-services/models"
	_ "user-services/respository"
	repository "user-services/respository"
	"user-services/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// TODO: Review security risks
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Save user in DB
	result := repository.DB.Create(&user)
	if result.Error != nil {
		log.Printf("Error creating user: %v", result.Error)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode("User registered successfully")
	if err != nil {
		return
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	var dbUser models.User
	result := repository.DB.Where("email = ?", user.Email).First(&dbUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("User %v does not exist", user.Email)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check passwords
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		log.Println("Invalid password", err)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := utils.GenerateUserJWT(dbUser.ID)
	if err != nil {
		log.Printf("Error creating token: %v \n", err)
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"token": token})

	if err != nil {
		return
	}
}
