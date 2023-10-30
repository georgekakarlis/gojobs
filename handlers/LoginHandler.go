package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

	
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func generateOauthStateString() string {
	bytes := make([]byte, 32)
	for i := range bytes {
		bytes[i] = byte(65 + rand.Intn(25)) // A=65 and Z=90
	}
	return string(bytes)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	oauthState := generateOauthStateString()
	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Path: "/",
		Value:    oauthState,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
	})

	data := map[string]interface{}{
		"GoogleLoginURL": googleOauthConfig.AuthCodeURL(oauthState),
	}

	tmpl, err := template.ParseFiles("views/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, err := r.Cookie("oauthstate")
	if err != nil || oauthState == nil {
    http.Error(w, "Missing or invalid OAuth state cookie", http.StatusBadRequest)
    return
	}

	if r.FormValue("state") != oauthState.Value {
    http.Error(w, "Invalid OAuth2 state", http.StatusUnauthorized)
    return
	}

	content, err := getUserInfo(r.FormValue("code"))
	if err != nil {
		http.Error(w, "Failed to get user information", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Content: %s\n", content)
}

func getUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}
