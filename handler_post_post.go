package main

import (
	"net/http"
	"strings"
)

func handlePostPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, err := r.Cookie("forum_session")
	if err != nil || token.Value == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	userID, _, err := getUserBySession(r.Context(), token.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	body := strings.TrimSpace(r.FormValue("body"))
	if title == "" || body == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	source := strings.TrimSpace(r.FormValue("source"))

	_, err = db.Exec(r.Context(), `
		INSERT INTO posts (author_id, title, body, source)
		VALUES ($1, $2, $3, $4)`, userID, title, body, source)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
