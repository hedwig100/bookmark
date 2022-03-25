package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/hedwig100/bookmark/backend/db"
	"github.com/hedwig100/bookmark/backend/middleware"
	"github.com/hedwig100/bookmark/backend/slog"
	"golang.org/x/crypto/bcrypt"
)

func readBody(r *http.Request) ([]byte, error) {
	len := r.ContentLength
	body := make([]byte, len)

	if _, err := r.Body.Read(body); err != nil && err != io.EOF {
		return nil, err
	}
	return body, nil
}

// id create uuid
func id() (string, error) {
	ret, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return ret.String(), nil
}

// /users
// POST
func postUser(w http.ResponseWriter, r *http.Request) {

	// read request body
	body, err := readBody(r)
	if err != nil {
		slog.Infof("error while parsing request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		slog.Infof("expect User model: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// generate uuid
	user_id, err := id()
	if err != nil {
		slog.Infof("uuid generation failure: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// hash password
	hashedPassword := []byte(user.Password)
	hashedPassword, err = bcrypt.GenerateFromPassword(hashedPassword, 10)
	if err != nil {
		slog.Infof("internal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.Pool.Exec(context.Background(),
		"INSERT INTO users (user_id,username,password) VALUES ($1,$2,$3);",
		user_id, user.Username, hashedPassword)

	if err != nil {
		slog.Infof("db insert error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate JWT and send it to client
	middleware.GenJWT(w, user_id)
}
