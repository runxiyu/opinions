package main

import (
	"net/http"
	"time"
)

type PostListItem struct {
	ID        int
	Title     string
	Source    string
	Author    string
	CreatedAt time.Time
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("forum_session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	_, _, err = getUserBySession(r.Context(), cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	ctx := r.Context()
	rows, err := db.Query(ctx, `
		SELECT posts.id, posts.title, posts.source, COALESCE(users.username, '[deleted]'), posts.created_at
		FROM posts
		LEFT JOIN users ON posts.author_id = users.id
		ORDER BY posts.created_at DESC`)
	if err != nil {
		http.Error(w, "query error "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var posts []PostListItem
	for rows.Next() {
		var p PostListItem
		if err := rows.Scan(&p.ID, &p.Title, &p.Source, &p.Author, &p.CreatedAt); err != nil {
			http.Error(w, "scan error", 500)
			return
		}
		posts = append(posts, p)
	}

	templates.ExecuteTemplate(w, "index", struct{ Posts []PostListItem }{posts})
}
