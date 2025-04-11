package main

import (
	"net/http"
	"strconv"
	"strings"
)

func handlePostIDPost(w http.ResponseWriter, r *http.Request) {
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

	idStr := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	replyType := r.FormValue("type")
	body := strings.TrimSpace(r.FormValue("body"))
	if body == "" {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	validTypes := map[string]bool{
		"opinion": true, "concur": true,
		"concurj": true, "dissent": true,
	}
	if !validTypes[replyType] {
		http.Error(w, "invalid type", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(r.Context(), `
		INSERT INTO replies (post_id, author_id, body, type)
		VALUES ($1, $2, $3, $4)`, postID, userID, body, replyType)
	if err != nil {
		http.Error(w, "insert error", 500)
		return
	}

	http.Redirect(w, r, "/post/"+idStr, http.StatusSeeOther)
}
