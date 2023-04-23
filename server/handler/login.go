package handler

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/peoxia/auth-api/auth"
	"github.com/peoxia/auth-api/config"
)

// Login is called when login is initiated by a user.
// 		GET /api/v1/login
// 		Responds: 307, 500
func Login(p auth.LoginProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authCodeURL, err := p.GetAuthCodeURL()
		if err != nil {
			http.Error(w, "Failed to get auth code URL", http.StatusInternalServerError)
			return
		}

		// Store the state parameter in a session cookie
		cookie := &http.Cookie{
			Name:  "oauth-state",
			Value: authCodeURL.State,
			// Disabled for local testing
			// HttpOnly: true,
			// Secure:   true,
			// SameSite: http.SameSiteStrictMode,
			MaxAge: 60 * 10, // 10 minutes
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, authCodeURL.URL, http.StatusTemporaryRedirect)
	}
}

// LoginCallback is called after user completed login process with third party provider.
// 		GET /api/v1/login/callback
// 		Responds: 200, 400, 500
func LoginCallback(c config.Config, p auth.LoginProvider, s auth.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Retrieve the state parameter from the session cookie
		stateCookie, err := r.Cookie("oauth-state")
		if err != nil {
			http.Error(w, "Failed to retrieve state parameter", http.StatusBadRequest)
			return
		}
		state := stateCookie.Value
		stateCookie.MaxAge = -1
		http.SetCookie(w, stateCookie)

		// Verify that the state parameter returned by the auth provider matches the stored state parameter
		if r.URL.Query().Get("state") != state {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		accessToken, err := p.GetAccessToken(r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}

		user, err := p.GetUserInfo(*accessToken)
		if err != nil {
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}

		// Create or update a user
		err = s.UpsertUser(ctx, *user)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		// Generate a JWT
		claims := jwt.MapClaims{
			"sub": user.Email,
			"iss": "v-p.dev",
			"aud": "auth-api",
			"exp": time.Now().Add(time.Hour * 24).Unix(),
			"iat": time.Now().Unix(),
		}
		tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(c.JWTSecret))
		if err != nil {
			http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
			return
		}
		jwtCookie := &http.Cookie{
			Name:   "jwt",
			Path:   "/",
			Value:  tokenString,
			MaxAge: 60 * 60 * 24, // 24 hrs,
			// Disable for local testing
			// HttpOnly: true,
			// Secure:   true,
			// SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, jwtCookie)

		JSONdata, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Failed to marshal json response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(JSONdata)
	}
}
