package main

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	"golang.org/x/crypto/argon2"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templates.ExecuteTemplate(w, "login", nil)
	case http.MethodPost:
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")

		var userID int
		var passwordHash []byte
		err := db.QueryRow(r.Context(),
			`SELECT id, password_hash FROM users WHERE username = $1`, username).
			Scan(&userID, &passwordHash)
		if err != nil {
			http.Error(w, "invalid login", http.StatusUnauthorized)
			return
		}

		ok := verifyPassword(passwordHash, password)
		if !ok {
			http.Error(w, "invalid login", http.StatusUnauthorized)
			return
		}

		token, err := randomToken(32)
		if err != nil {
			http.Error(w, "token error", 500)
			return
		}
		if err := saveSession(r.Context(), userID, token); err != nil {
			http.Error(w, "session error", 500)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "forum_session",
			Value: token,
			Path:  "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func verifyPassword(hash []byte, password string) bool {
	parts := strings.Split(string(hash), "$")
	if len(parts) != 6 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}
	computed := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, uint32(len(expectedHash)))
	return subtle.ConstantTimeCompare(expectedHash, computed) == 1
}
