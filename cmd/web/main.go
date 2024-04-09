package main

import (
	"log"
	"net/http"
	"time"

	"github.com/coltonmosier/api-v1/internal/handlers"
	"github.com/coltonmosier/api-v1/internal/helpers"
	"github.com/coltonmosier/api-v1/internal/middleware"
	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	devices := handlers.DeviceHandler{}
	manufactuerers := handlers.ManufactuerHandler{}
    equipment := handlers.EquipmentHandler{}

	r := http.NewServeMux()

	r.HandleFunc("GET /api/v1/health", HealthHandler)

	// NOTE: Device Type routes
	r.HandleFunc("GET /api/v1/device", devices.GetDeviceTypes)
	//r.HandleFunc("POST /api/v1/device/{id}", devices.CreateDeviceType)

	// NOTE: Manufacturer routes
	r.HandleFunc("GET /api/v1/manufacturer", manufactuerers.GetManufacturers)
    r.HandleFunc("GET /api/v1/manufacturer/{id}", manufactuerers.GetManufacturerByID)
    r.HandleFunc("PATCH /api/v1/manufacturer/{id}/name/{name}", manufactuerers.UpdateManufacturerName)
    //r.HandleFunc("PATCH /api/v1/manufacturer/{id}/status/{status}", manufactuerers.UpdateManufacturerStatus)

	// NOTE: Equipment routes



    //NOTE: Handle not found endpoints
    //r.HandleFunc("/", equipment.BadEndpointHandler)

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
    helpers.JsonResponseSuccess(w, http.StatusOK, "API is healthy")
}
