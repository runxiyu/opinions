package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostViewData struct {
	Post struct {
		ID        int
		Title     string
		Author    string
		CreatedAt time.Time
		Body      string
		Source    string
	}
	Replies []struct {
		Author    string
		CreatedAt time.Time
		Body      string
		Type      string
	}
}

func handlePostIDGet(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/post/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ctx := r.Context()
	var pv PostViewData

	err = db.QueryRow(ctx, `
		SELECT posts.id, posts.title, COALESCE(users.username, '[deleted]'), posts.created_at, posts.body, posts.source
		FROM posts LEFT JOIN users ON posts.author_id = users.id
		WHERE posts.id = $1`, id).Scan(
		&pv.Post.ID, &pv.Post.Title, &pv.Post.Author, &pv.Post.CreatedAt, &pv.Post.Body, &pv.Post.Source)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	rows, err := db.Query(ctx, `
		SELECT COALESCE(users.username, '[deleted]'), replies.created_at, replies.body, replies.type
		FROM replies LEFT JOIN users ON replies.author_id = users.id
		WHERE replies.post_id = $1
		ORDER BY replies.created_at ASC`, id)
	if err != nil {
		http.Error(w, "reply query error", 500)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rep struct {
			Author    string
			CreatedAt time.Time
			Body      string
			Type      string
		}
		if err := rows.Scan(&rep.Author, &rep.CreatedAt, &rep.Body, &rep.Type); err != nil {
			http.Error(w, "scan error", 500)
			return
		}
		pv.Replies = append(pv.Replies, rep)
	}

	templates.ExecuteTemplate(w, "postview", pv)
}
