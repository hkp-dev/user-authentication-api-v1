package middleware

import (
	"app/database"
	"app/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

type LogEntry struct {
	Method   string `json:"method"`
	URL      string `json:"url"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Status   int    `json:"status"`
	Duration string `json:"duration"`
	Time     string `json:"time"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

func (rw *CustomResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func LogRequestMiddleware(next http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		customWriter := &CustomResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		start := time.Now()

		next.ServeHTTP(customWriter, r)

		duration := time.Since(start)
		user, _ := utils.ValidateJWT(r)

		logEntry := LogEntry{
			Method:   r.Method,
			URL:      r.URL.Path,
			Status:   customWriter.StatusCode,
			Duration: duration.String(),
			Time:     start.Format(time.RFC3339),
			Request:  r.RequestURI,
			Response: string(customWriter.Body),
		}
		logEntry.UserID = user.ID.Hex()
		logEntry.Username = user.Username
		logEntry.Request = r.RequestURI
		logEntry.Response = string(customWriter.Body)
		logJSON, err := json.MarshalIndent(logEntry, "", "  ")
		if err != nil {
			log.Printf("Error encoding log to JSON: %v", err)
			return
		}
		log.Println(string(logJSON))
	}
}

func RequireAdminAuth(next http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	Func := func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := utils.ValidateJWT(r)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}
		_, err = database.UserCollection.Find(context.Background(), user)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}
		if user.Role == "admin" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}
	}
	return LogRequestMiddleware(Func)
}

func RequireUserAuth(next http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	Func := func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := utils.ValidateJWT(r)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}
		_, err = database.UserCollection.Find(context.Background(), user)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}
		if user.Role == "user" || user.Role == "admin" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}
	}
	return LogRequestMiddleware(Func)
}
