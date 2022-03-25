package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type reqTest struct {
	// error message
	msg string

	// request data
	method, uri, body string

	// expected response data
	wantCode int
	wantBody string
}

func TestServer(t *testing.T) {
	testdata := []reqTest{
		{msg: "error hello handler", method: "GET", uri: "/hello", wantCode: 200, wantBody: "Hello World!"},
		{msg: "error postUser", method: "POST", uri: "/users", body: `{"username":"hedwig100","password":"abcde12345"}`, wantCode: 200, wantBody: ""},
		{msg: "error postUser", method: "POST", uri: "/users", body: `{"username":"John"}`, wantCode: 400},
	}

	mux := getMux()
	var w *httptest.ResponseRecorder
	for i, rt := range testdata {
		body := strings.NewReader(rt.body)
		req, err := http.NewRequest(rt.method, rt.uri, body)
		if err != nil {
			t.Fatalf("[test%d] invalid request: %v", i, err)
		}
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != rt.wantCode {
			t.Fatalf("[test%d] %s, code expected: %d,actual: %d", i, rt.msg, rt.wantCode, w.Code)
		}
		actBody := string(w.Body.Bytes())
		if actBody != rt.wantBody {
			t.Fatalf("[test%d] %s, body expected: %s,actual: %s", i, rt.msg, rt.wantBody, actBody)
		}
	}
}
