package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("This is the api server")

	s := &http.Server{
		Addr:         ":8081",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/api/v1/health", HealthHandler).Methods("GET")


	http.Handle("/", r)

	log.Fatal(s.ListenAndServe())
}

// This function should be logging to a database and also to a file
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-IP")
		msg := fmt.Sprintf("%s %s %s %s\n", ip, r.Method, r.RequestURI, r.UserAgent())
		log.Print(msg)

		next.ServeHTTP(w, r)
	})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

