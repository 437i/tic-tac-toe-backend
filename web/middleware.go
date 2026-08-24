package web

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	if lrw.wroteHeader {
		return
	}
	lrw.statusCode = code
	lrw.wroteHeader = true
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{w, http.StatusOK, false}
		next.ServeHTTP(lrw, r)
		route := mux.CurrentRoute(r)
		name := "unknown route"
		if route != nil {
			name = route.GetName()
		}
		log.Printf("%s %s %s %d %s\n", r.Method, r.RequestURI, name, lrw.statusCode, time.Since(start))
	})
}
