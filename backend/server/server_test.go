package server_test

import (
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
		{name: "postUserExpectSucess1", username: "hedwig100", method: "POST", uri: "/users", body: `{"username":"hedwig100","password":"abcde12345"}`, wantCode: 201, wantJWT: true},
		{name: "postUserExpectSucess2", username: "Kate", method: "POST", uri: "/users", body: `{"username":"Kate","password":"01234pol"}`, wantCode: 201, wantJWT: true},
		{name: "postUserFailWithInvalidJson", method: "POST", uri: "/users", body: `{"username":"John`, wantCode: 400,
			wantBody: `{"message":"Invalid json format.","code":1}`},
		{name: "postUserFailWithJson", method: "POST", uri: "/users", body: `{"username":"John"}`, wantCode: 400,
			wantBody: `{"message":"Empty username or password.","code":0}`},
		{name: "postUserFailWithEmpty", method: "POST", uri: "/users", body: `{"username":"","password":""}`, wantCode: 400,
			wantBody: `{"message":"Empty username or password.","code":0}`},
		{name: "postUserFailuresWithSameUsername", method: "POST", uri: "/users", body: `{"username":"hedwig100","password":"aaa"}`, wantCode: 500,
			wantBody: `{"message":"The username is already registered.","code":0}`},
		{name: "read1", username: "hedwig100", needJWT: true, method: "POST", uri: "/users/hedwig100/books",
			body:     `{"bookName":"Harry Potter","authorName":"J.K.Rowling","genres":["fantasy","for children"],"thoughts":"Voldemort scared me a lot.","readAt":"2021-10-30T21:07"}`,
			wantCode: 201},
		{name: "read2", username: "Kate", needJWT: true, method: "POST", uri: "/users/Kate/books",
			body:     `{"bookName":"Who Moved My Cheese?","authorName":"Spencer Johnson","genres":["life"],"readAt":"2022-03-29T21:07"}`,
			wantCode: 201},
		{name: "read3", username: "Kate", needJWT: true, method: "POST", uri: "/users/Kate/books",
			body:     `{"bookName":"Harry Potter","authorName":"J.K.Rowling","genres":["fantasy","for children"],"thoughts":"Very Exciting!","readAt":"2021-10-30T21:07"}`,
			wantCode: 201},
		{name: "readErrorWhenUsernameIsDifferent", username: "hedwig100", needJWT: true, method: "POST", uri: "/users/hedwig/books",
			body:     `{"bookName":"Harry Potter","authorName":"J.K.Rowling","genres":["fantasy","for children"],"thoughts":"Voldemort scared me a lot.","readAt":"2021-10-30T21:07"}`,
			wantCode: 401},
		{name: "login", username: "hedwig100", method: "POST", uri: "/login", body: `{"username":"hedwig100","password":"abcde12345"}`, wantCode: 200, wantJWT: true},
		{name: "loginFailureWithWrongPassword", method: "POST", uri: "/login", body: `{"username":"hedwig100","password":"abc45"}`, wantCode: 500,
			wantBody: `{"message":"Invalid user or password.","code":0}`},
		{name: "loginFailureWithUnregisteredUser", method: "POST", uri: "/login", body: `{"username":"he100","password":"abc45"}`, wantCode: 500,
			wantBody: `{"message":"Invalid user or password.","code":0}`},
		{name: "readGet", username: "hedwig100", needJWT: true, method: "GET", uri: "/users/hedwig100/books", wantCode: 200,
			wantBody: `{"reads":[{"bookName":"Harry Potter","authorName":"J.K.Rowling","genres":["fantasy","for children"],"thoughts":"Voldemort scared me a lot.","readAt":"2021-10-30T21:07"}]}`},
		{name: "readGet", username: "Kate", needJWT: true, method: "GET", uri: "/users/Kate/books", wantCode: 200,
			wantBody: `{"reads":[{"bookName":"Who Moved My Cheese?","authorName":"Spencer Johnson","genres":["life"],"thoughts":"","readAt":"2022-03-29T21:07"},{"bookName":"Harry Potter","authorName":"J.K.Rowling","genres":["fantasy","for children"],"thoughts":"Very Exciting!","readAt":"2021-10-30T21:07"}]}`},
	}
	// maps username to jwt
	mp := make(map[string]*http.Cookie, 1)

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
				c, ok := mp[rt.username]
				if !ok {
					t.Fatalf("[test%d] token not found", i)
				}
				req.AddCookie(c)
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
				cStr, ok := w.Header()["Set-Cookie"]
				if !ok {
					t.Fatalf("[test%d] want jwt but cookie isn't set.", i)
				}
				parser := &http.Request{Header: http.Header{"Cookie": cStr}}
				c, err := parser.Cookie("bookmark_auth")
				if err != nil {
					t.Fatalf("[test%d] cookie parse error", i)
				}
				mp[rt.username] = c
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
