package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
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

func respErr(w http.ResponseWriter, status int, message string) {
	body, _ := json.Marshal(Error{Message: message})
	w.WriteHeader(status)
	w.Write(body)
}

// /hello
// GET
// this is a test handler
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
	w.WriteHeader(http.StatusOK)
}

// cors
// OPTIONS
func cors(w http.ResponseWriter, r *http.Request) {
	// TODO: security
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Max-Age", "86400")
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
		slog.Errf("user create error: %v", err)
		respErr(w, http.StatusInternalServerError, "The username is already registered.")
		return
	}

	// generate JWT and send it to client
	middleware.GenJWT(w, user_id, user.Username)
	w.WriteHeader(http.StatusCreated)
}

// /login
// POST
func login(w http.ResponseWriter, r *http.Request) {
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

	userId, err := Db.Login(user)
	if err != nil {
		slog.Errf("internal server error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate JWT and send it to client
	middleware.GenJWT(w, userId, user.Username)
	w.WriteHeader(http.StatusOK)
}

// /users/:username/books
// POST
func read(w http.ResponseWriter, r *http.Request) {

	// read request body
	body, err := readBody(r)
	if err != nil {
		slog.Infof("error while parsing request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var read data.Read
	if err = json.Unmarshal(body, &read); err != nil {
		slog.Infof("expect Read model: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if read.BookName == "" || read.AuthorName == "" {
		slog.Infof("expect valid Read model")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := httptreemux.ContextParams(r.Context())
	username := params["username"]
	err = Db.ReadCreate(username, read)
	if err != nil {
		slog.Infof("internal server error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// /users/:username/books
// GET
func readGet(w http.ResponseWriter, r *http.Request) {
	params := httptreemux.ContextParams(r.Context())
	username := params["username"]
	reads, err := Db.ReadGet(username)
	if err != nil {
		slog.Errf("internal server error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(data.Reads{Reads: reads})
	if err != nil {
		slog.Errf("internal server error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	return
}
