package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"golang.org/x/crypto/argon2"
)

func handleSignup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templates.ExecuteTemplate(w, "signup", nil)
	case http.MethodPost:
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		if username == "" || password == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}

		hash, err := hashPassword(password)
		if err != nil {
			http.Error(w, "hash error", 500)
			return
		}

		_, err = db.Exec(r.Context(),
			`INSERT INTO users (username, password_hash) VALUES ($1, $2)`,
			username, hash)
		if err != nil {
			http.Error(w, "signup failed", 500)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func hashPassword(password string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	encoded := "$argon2id$v=19$m=65536,t=1,p=4$" +
		base64.RawStdEncoding.EncodeToString(salt) + "$" +
		base64.RawStdEncoding.EncodeToString(hash)
	return []byte(encoded), nil
}
