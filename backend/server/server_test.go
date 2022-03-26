package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hedwig100/bookmark/backend/data"
	"github.com/hedwig100/bookmark/backend/server"
)

func testServer(t *testing.T) {
	testdata := []struct {
		// meta data
		name     string
		username string
		needJWT  bool

		// request data
		method, uri, body string

		// expected response data
		wantCode int
		wantBody string
		wantJWT  bool // if JWT should be set or not
	}{
		{name: "hello", method: "GET", uri: "/hello", wantCode: 200, wantBody: "Hello World!"},
		{name: "postUserExpectSucess", username: "hedwig100", method: "POST", uri: "/users", body: `{"username":"hedwig100","password":"abcde12345"}`, wantCode: 201, wantJWT: true},
		{name: "postUserExpectFailure", method: "POST", uri: "/users", body: `{"username":"John"}`, wantCode: 400},
		{name: "read", username: "hedwig100", needJWT: true, method: "POST", uri: "/users/hedwig100/books",
			body:     `{"bookName":"Harry Potter","authorName":"J.K.Rowling","genres":["fantasy","for children"],"thoughts":"Voldemort scared me a lot.","readAt":"2021-10-30T21:07"}`,
			wantCode: 201},
		{name: "readErrorWhenUsernameIsDifferent", username: "hedwig100", needJWT: true, method: "POST", uri: "/users/hedwig/books",
			body:     `{"bookName":"Harry Potter","authorName":"J.K.Rowling","genres":["fantasy","for children"],"thoughts":"Voldemort scared me a lot.","readAt":"2021-10-30T21:07"}`,
			wantCode: 401},
		{name: "login", username: "hedwig100", method: "POST", uri: "/login", body: `{"username":"hedwig100","password":"abcde12345"}`, wantCode: 200, wantJWT: true},
		{name: "loginFailureWithWrongPassword", method: "POST", uri: "/login", body: `{"username":"hedwig100","password":"abc45"}`, wantCode: 500},
		{name: "loginFailureWithUnregisteredUser", method: "POST", uri: "/login", body: `{"username":"he100","password":"abc45"}`, wantCode: 500},
	}
	// maps username to jwt
	mp := make(map[string]string, 1)

	// server initialization
	mux := server.GetMux()

	var w *httptest.ResponseRecorder
	for i, rt := range testdata {
		t.Run(rt.name, func(t *testing.T) {
			/*
				create request body
			*/
			body := strings.NewReader(rt.body)
			req, err := http.NewRequest(rt.method, rt.uri, body)
			if err != nil {
				t.Fatalf("[test%d] invalid request: %v", i, err)
			}
			if rt.needJWT {
				token, ok := mp[rt.username]
				if !ok {
					t.Fatalf("[test%d] token not found", i)
				}
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			/*
				send request
			*/
			w = httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			/*
				check response
			*/

			// statuscode
			if w.Code != rt.wantCode {
				t.Fatalf("[test%d] code expected: %d,actual: %d", i, rt.wantCode, w.Code)
			}

			// response body
			actBody := string(w.Body.Bytes())
			if actBody != rt.wantBody {
				t.Fatalf("[test%d] body expected: %s,actual: %s", i, rt.wantBody, actBody)
			}

			// authorization
			if rt.wantJWT {
				auth, ok := w.Header()["Authorization"]
				if !ok || !strings.HasPrefix(auth[0], "Bearer ") {
					t.Fatalf("[test%d] Authorization expected but not found", i)
				}
				mp[rt.username] = strings.TrimPrefix(auth[0], "Bearer ")
			}
		})
	}
}

func TestIntegrateServer(t *testing.T) {
	server.Db = data.NewDbReal()
	testServer(t)
}

func TestUnitServer(t *testing.T) {
	server.Db = data.NewDbMock()
	testServer(t)
}
