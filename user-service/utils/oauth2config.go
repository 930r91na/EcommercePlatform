package utils

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// GoogleOauthConfig Google OAuth2 configuration
var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// GithubOauthConfig GitHub OAuth2 configuration
var GithubOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/github/callback",
	ClientID:     os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"user:email"},
	Endpoint:     github.Endpoint,
}
