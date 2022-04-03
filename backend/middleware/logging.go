package middleware

import (
	"net/http"

	"github.com/hedwig100/bookmark/backend/slog"
)

// statusResponseWriter implements http.ResponseWriter Interface
// this is used for keeping statusCode
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *statusResponseWriter) Write(buf []byte) (int, error) {
	return w.ResponseWriter.Write(buf)
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// LogWrap wraps handlers to output logs and set 'Access-Control-Allow-Origin'
func LogWrap(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Infof("Method: %s,URL: %s,Protocol: %s,RemoteIP: %s", r.Method, r.URL, r.Proto, r.RemoteAddr)
		sw := statusResponseWriter{
			ResponseWriter: w,
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler(&sw, r)
		slog.Infof("Status: %d,Header: %v)", sw.statusCode, w.Header())
	}
}
