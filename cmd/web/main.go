package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coltonmosier/api-v1/internal/middleware"
    "github.com/coltonmosier/api-v1/internal/handlers"
)

var DB *sql.DB

func main() {
	fmt.Println("This is the api server")

    handlers := handlers.Handler{}

	r := http.NewServeMux()

	r.HandleFunc("GET /api/v1/health", HealthHandler)
    r.HandleFunc("GET /api/v1/device_type", handlers.GetDeviceType)
    r.HandleFunc("POST /api/v1/device_type/{name}", handlers.CreateDeviceType)

	http.Handle("/", r)

	s := &http.Server{
		Addr:         ":8081",
		Handler:      middleware.LoggingMiddleware(r),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}


func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
