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
	// meta data
	name string

	// request data
	method, uri, body string

	// expected response data
	wantCode int
	wantBody string
}

func testServer(t *testing.T) {
	testdata := []reqTest{
		{name: "hello", method: "GET", uri: "/hello", wantCode: 200, wantBody: "Hello World!"},
		{name: "postUserExpectSucess", method: "POST", uri: "/users", body: `{"username":"hedwig100","password":"abcde12345"}`, wantCode: 200, wantBody: ""},
		{name: "postUserExpectFailure", method: "POST", uri: "/users", body: `{"username":"John"}`, wantCode: 400},
	}

	// server initialization
	mux := server.GetMux()

	var w *httptest.ResponseRecorder
	for i, rt := range testdata {
		t.Run(rt.name, func(t *testing.T) {
			body := strings.NewReader(rt.body)
			req, err := http.NewRequest(rt.method, rt.uri, body)
			if err != nil {
				t.Fatalf("[test%d] invalid request: %v", i, err)
			}
			w = httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != rt.wantCode {
				t.Fatalf("[test%d] code expected: %d,actual: %d", i, rt.wantCode, w.Code)
			}
			actBody := string(w.Body.Bytes())
			if actBody != rt.wantBody {
				t.Fatalf("[test%d] body expected: %s,actual: %s", i, rt.wantBody, actBody)
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
