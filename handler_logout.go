// handler_logout.go
package main

import (
	"net/http"
)

func handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("forum_session")
	if err == nil && cookie.Value != "" {
		_ = deleteSession(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "forum_session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
