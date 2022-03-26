package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hedwig100/bookmark/backend/data"
	"github.com/hedwig100/bookmark/backend/middleware"
	"github.com/hedwig100/bookmark/backend/slog"
)

func readBody(r *http.Request) ([]byte, error) {
	len := r.ContentLength
	body := make([]byte, len)

	if _, err := r.Body.Read(body); err != nil && err != io.EOF {
		return nil, err
	}
	return body, nil
}

// /hello
// this is a test handler
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
	w.WriteHeader(http.StatusOK)
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

	var user data.User
	if err = json.Unmarshal(body, &user); err != nil {
		slog.Infof("expect User model: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		slog.Infof("expect valid User model")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user_id, err := Db.UserCreate(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate JWT and send it to client
	middleware.GenJWT(w, user_id)
}
