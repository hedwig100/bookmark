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
		{msg: "err hello handler", method: "GET", uri: "/hello", wantCode: 200, wantBody: "Hello World!"},
	}

	mux := getMux()
	w := httptest.NewRecorder()
	for i, rt := range testdata {
		body := strings.NewReader(rt.body)
		req, err := http.NewRequest(rt.method, rt.uri, body)
		if err != nil {
			t.Errorf("[test%d] invalid request: %v", i, err)
		}
		w.Flush()
		mux.ServeHTTP(w, req)

		if w.Code != rt.wantCode {
			t.Errorf("[test%d] %s, code expected: %d,actual: %d", i, rt.msg, rt.wantCode, w.Code)
		}
		actBody := string(w.Body.Bytes())
		if actBody != rt.wantBody {
			t.Errorf("[test%d] %s, body expected: %s,actual: %s", i, rt.msg, rt.wantBody, actBody)
		}
	}
}
