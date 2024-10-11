package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	_ "encoding/json"
	"errors"
	"fmt"
	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	_ "net/http"
	"time"
	"user-services/models"
	_ "user-services/respository"
	repository "user-services/respository"
	"user-services/utils"
)

func processUserInfo(w http.ResponseWriter, oauthProvider string, userInfo map[string]interface{}) (string, error) {
	oauthProviderUserID := fmt.Sprintf("%v", userInfo["id"]) // Ensure it's a string
	email := fmt.Sprintf("%v", userInfo["email"])
	name := fmt.Sprintf("%v", userInfo["name"])

	var user models.User
	err := repository.DB.Where("oauth_provider = ? AND oauth_provider_user_id = ?", oauthProvider, oauthProviderUserID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found, create new user
			user = models.User{
				OAuthProvider:       oauthProvider,
				OAuthProviderUserID: oauthProviderUserID,
				Email:               email,
				Name:                name,
			}
			if err := repository.DB.Create(&user).Error; err != nil {
				log.Printf("Error creating user: %v", err)
				http.Error(w, "Error creating user", http.StatusInternalServerError)
				return "", err
			}
		} else {
			log.Printf("Database error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return "", err
		}
	}

	// Generate JWT token
	tokenString, err := utils.GenerateJWT(user.ID)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return "", err
	}

	return tokenString, nil
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(24 * time.Hour)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)
	return state
}

//#region GOOGLE OAUTH

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("GoogleLogin invoked")

	if utils.GoogleOauthConfig.ClientID == "" || utils.GoogleOauthConfig.ClientSecret == "" {
		http.Error(w, "GitHub OAuth client ID or secret not set.", http.StatusInternalServerError)
		return
	}

	state := generateStateOauthCookie(w)
	url := utils.GoogleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Retrieve state from cookie
	log.Println("GoogleCallback invoked")
	stateCookie, err := r.Cookie("oauthstate")
	if err != nil {
		http.Error(w, "Invalid state cookie", http.StatusBadRequest)
		return
	}
	state := stateCookie.Value

	// Compare the state parameter
	if r.FormValue("state") != state {
		http.Error(w, "Invalid OAuth state", http.StatusUnauthorized)
		return
	}

	// Exchange code for token
	token, err := utils.GoogleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	// Retrieve user info
	client := utils.GoogleOauthConfig.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(userInfoResp.Body)

	userInfoData, _ := io.ReadAll(userInfoResp.Body)
	var userInfo map[string]interface{}
	err = json.Unmarshal(userInfoData, &userInfo)
	if err != nil {
		return
	}

	// Process user info and generate JWT token
	tokenString, err := processUserInfo(w, "google", userInfo)
	if err != nil {
		// Errors are already handled inside processUserInfo
		return
	}

	// Send the token to the client
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	if err != nil {
		return
	}
}

//#endregion

//#region GITHUB OAUTH

func GithubLogin(w http.ResponseWriter, r *http.Request) {
	if utils.GithubOauthConfig.ClientID == "" || utils.GithubOauthConfig.ClientSecret == "" {
		http.Error(w, "GitHub OAuth client ID or secret not set.", http.StatusInternalServerError)
		return
	}

	state := generateStateOauthCookie(w)
	url := utils.GithubOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	// Retrieve state from cookie
	stateCookie, err := r.Cookie("oauthstate")
	if err != nil {
		http.Error(w, "Invalid state cookie", http.StatusBadRequest)
		return
	}
	state := stateCookie.Value

	// Compare the state parameter
	if r.FormValue("state") != state {
		http.Error(w, "Invalid OAuth state", http.StatusUnauthorized)
		return
	}

	// Exchange code for token
	token, err := utils.GithubOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	// Retrieve user info
	client := utils.GithubOauthConfig.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(userInfoResp.Body)

	userInfoData, _ := io.ReadAll(userInfoResp.Body)
	var userInfo map[string]interface{}

	err = json.Unmarshal(userInfoData, &userInfo)
	if err != nil {
		return
	}

	// Process user info and generate JWT token
	tokenString, err := processUserInfo(w, "github", userInfo)
	if err != nil {
		// Errors are already handled inside processUserInfo
		return
	}

	// Send the token to the client
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	if err != nil {
		return
	}
}

//#endregion
