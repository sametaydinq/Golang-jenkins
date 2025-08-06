package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// JSON response struct
type response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Ping handler
func pingHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, response{Message: "pong"})
}

// Middleware: Logger
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

// Middleware: Recover from panic
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				respondJSON(w, http.StatusInternalServerError, response{
					Error: "Internal Server Error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Helper: Write JSON response
func respondJSON(w http.ResponseWriter, status int, payload response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// Entry point
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)

	// Wrap with middleware
	handler := loggerMiddleware(recoverMiddleware(mux))

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
