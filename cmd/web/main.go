package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/coltonmosier/api-v1/internal/database"
	"github.com/coltonmosier/api-v1/internal/handlers"
	"github.com/coltonmosier/api-v1/internal/middleware"
	"github.com/joho/godotenv"
)

var (
	EqDB  *sql.DB
	LogDB *sql.DB
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	EqDB = database.InitEquipmentDatabase()
	LogDB = database.InitLoggingDatabase()

	devices := handlers.DeviceHandler{}
	manufactuerers := handlers.ManufactuerHandler{}

	r := http.NewServeMux()

	r.HandleFunc("GET /api/v1/health", HealthHandler)
	// NOTE: Device Type routes
	r.HandleFunc("GET /api/v1/device_type", devices.GetDeviceTypes)
	r.HandleFunc("POST /api/v1/device_type/{name}", devices.CreateDeviceType)
	// NOTE: Manufacturer routes
	r.HandleFunc("GET /api/v1/manufacturer", manufactuerers.GetManufacturers)

	// NOTE: Equipment routes

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
