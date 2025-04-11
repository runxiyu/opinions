package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	var err error
	db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := loadTemplates(); err != nil {
		log.Fatal("template load error:", err)
	}

	http.HandleFunc("GET /", handleRoot)
	http.HandleFunc("GET /post", handlePostGet)
	http.HandleFunc("POST /post", handlePostPost)
	http.HandleFunc("GET /post/", handlePostIDGet)
	http.HandleFunc("POST /post/", handlePostIDPost)
	http.HandleFunc("GET /login", handleLogin)
	http.HandleFunc("POST /login", handleLogin)
	http.HandleFunc("GET /logout", handleLogout)
	http.HandleFunc("GET /signup", handleSignup)
	http.HandleFunc("POST /signup", handleSignup)
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
