package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hedwig100/bookmark/backend/data"
	"github.com/hedwig100/bookmark/backend/server"
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

func testServer(t *testing.T) {
	testdata := []reqTest{
		{msg: "error hello handler", method: "GET", uri: "/hello", wantCode: 200, wantBody: "Hello World!"},
		{msg: "error postUser", method: "POST", uri: "/users", body: `{"username":"hedwig100","password":"abcde12345"}`, wantCode: 200, wantBody: ""},
		{msg: "error postUser", method: "POST", uri: "/users", body: `{"username":"John"}`, wantCode: 400},
	}

	// server initialization
	mux := server.GetMux()

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

func TestIntegrateServer(t *testing.T) {
	server.Db = data.NewDbReal()
	testServer(t)
}

func TestUnitServer(t *testing.T) {
	server.Db = data.NewDbMock()
	testServer(t)
}
