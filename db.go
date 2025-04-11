package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func getUserBySession(ctx context.Context, token string) (int, string, error) {
	var id int
	var username string
	err := db.QueryRow(ctx,
		`SELECT users.id, users.username
		 FROM sessions JOIN users ON sessions.user_id = users.id
		 WHERE sessions.token = $1`, token).Scan(&id, &username)
	if err != nil {
		return 0, "", err
	}
	return id, username, nil
}

func saveSession(ctx context.Context, userID int, token string) error {
	_, err := db.Exec(ctx,
		`INSERT INTO sessions (user_id, token)
		 VALUES ($1, $2)
		 ON CONFLICT (user_id) DO UPDATE SET token = $2`,
		userID, token)
	return err
}

func deleteSession(ctx context.Context, token string) error {
	_, err := db.Exec(ctx, `DELETE FROM sessions WHERE token = $1`, token)
	return err
}
