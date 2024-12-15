package net

import (
	"log"
	"net/http"
	"time"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *CustomResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customWriter := &CustomResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		start := time.Now()
		log.Printf("\n ==> Request Method: %s\n ==> URL: %s\n ==> Time: %s\n", r.Method, r.URL.Path, start.Format(time.RFC3339))

		next.ServeHTTP(customWriter, r)

		duration := time.Since(start)
		log.Printf("\n ==> Response Status: %d \n ==> Method: %s\n ==> URL: %s\n ==> Duration: %v\n", customWriter.StatusCode, r.Method, r.URL.Path, duration)
	})
}
