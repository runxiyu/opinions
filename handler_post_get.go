package main

import (
	"net/http"
)

func handlePostGet(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("forum_session")
	if err != nil || token.Value == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	_, _, err = getUserBySession(r.Context(), token.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	templates.ExecuteTemplate(w, "newpost", nil)
}
