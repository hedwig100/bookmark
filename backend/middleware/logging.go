package middleware

import (
	"fmt"
	"log"
	"net/http"
)

// statusResponseWriter implements http.ResponseWriter Interface
// this is used for keeping statusCode
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w statusResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// LogWrap wraps handlers to output logs
func LogWrap(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("[INFO] Method: %s,URL: %s,Protocol: %s,RemoteIP: %s", r.Method, r.URL, r.Proto, r.RemoteAddr))
		sw := statusResponseWriter{
			ResponseWriter: w,
		}
		handler(sw, r)
		log.Println(fmt.Sprintf("[INFO] Status: %d,Header: %v)", sw.statusCode, w.Header()))
	}
}
