package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/peoxia/auth-api/auth"
	"github.com/peoxia/auth-api/config"
)

func CurrentUser(c config.Config, s auth.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		email, err := ExtractJWTSubject(r, []byte(c.JWTSecret))
		if err != nil {
			http.Error(w, "Error extracting subject from jwt", http.StatusBadRequest)
			log.Println(err)
			return
		}

		user, err := s.FindUserByEmail(ctx, email)
		if err != nil {
			http.Error(w, "Error finding user by email", http.StatusInternalServerError)
			return
		}

		JSONdata, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Error marshalling user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(JSONdata)
	}
}

func UpdateUser(c config.Config, s auth.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var payload auth.Profile
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Error invalid payload", http.StatusBadRequest)
			return
		}

		email, err := ExtractJWTSubject(r, []byte(c.JWTSecret))
		if err != nil {
			http.Error(w, "Error extracting subject from jwt", http.StatusBadRequest)
			return
		}
		if payload.Email != email {
			http.Error(w, "Error invalid payload", http.StatusBadRequest)
			return
		}

		err = s.UpsertUser(ctx, payload)
		if err != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteUser(c config.Config, s auth.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		email, err := ExtractJWTSubject(r, []byte(c.JWTSecret))
		if err != nil {
			http.Error(w, "Error extracting subject from jwt", http.StatusBadRequest)
			log.Println(err)
			return
		}

		err = s.DeleteUser(ctx, email)
		if err != nil {
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func ExtractJWTSubject(r *http.Request, secret []byte) (string, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return "", err
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		// Provide the key used to validate the signature
		return secret, nil
	})

	if err != nil {
		return "", err
	}

	// Ensure the token is valid and contains the "sub" claim
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("Sub claim not found or not a string")
		}
		return sub, nil
	} else {
		return "", errors.New("Invalid token")
	}
}
