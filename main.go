package main

import (
	"log"
	"net/http"
	"time"

	"github.com/coltonmosier/api-v1/internal/handlers"
	"github.com/coltonmosier/api-v1/internal/helpers"
	"github.com/coltonmosier/api-v1/internal/middleware"
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger/v2"
    _ "github.com/coltonmosier/api-v1/docs"
)

//	@title			Equipment API
//	@version		1.0
//	@description	This is the API to interact with Equipment database
//	@termsOfService	http://swagger.io/terms/
//	@basePath		/api/v1
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	devices := handlers.DeviceHandler{}
	manufactuerers := handlers.ManufactuerHandler{}
    equipment := handlers.EquipmentHandler{}

	r := http.NewServeMux()

    r.HandleFunc("GET /swagger/*", httpSwagger.Handler(
            httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
        ))

	r.HandleFunc("GET /api/v1/health", HealthHandler)

	// NOTE: Device Type routes
	r.HandleFunc("GET /api/v1/device", devices.GetDeviceTypes)
    r.HandleFunc("GET /api/v1/device/{id}", devices.GetDeviceByID)
    r.HandleFunc("PATCH /api/v1/device/{id}", devices.UpdateDeviceType)
	r.HandleFunc("POST /api/v1/device", devices.CreateDeviceType)

	// NOTE: Manufacturer routes
	r.HandleFunc("GET /api/v1/manufacturer", manufactuerers.GetManufacturers)
    r.HandleFunc("GET /api/v1/manufacturer/{id}", manufactuerers.GetManufacturerByID)
    r.HandleFunc("PATCH /api/v1/manufacturer/{id}", manufactuerers.UpdateManufacturer)
    r.HandleFunc("POST /api/v1/manufacturer", manufactuerers.CreateManufacturer)

	// NOTE: Equipment routes
    r.HandleFunc("GET /api/v1/equipment", equipment.GetEquipments)
    r.HandleFunc("GET /api/v1/equipment/id", equipment.GetEquipmentByID)
    r.HandleFunc("GET /api/v1/equipment/sn", equipment.GetEquipmentBySN)
    r.HandleFunc("GET /api/v1/equipment/sn-like/{sn}", equipment.GetEquipmentLikeSN)
    r.HandleFunc("GET /api/v1/equipment/manufacturer/{id}", equipment.GetEquipmentByManufacturerID)
    r.HandleFunc("GET /api/v1/equipment/device/{id}", equipment.GetEquipmentByDeviceID)
    r.HandleFunc("GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}", equipment.GetEquipmentByDeviceIDAndManufacturerID)
    r.HandleFunc("GET /api/v1/equipment/sn/{sn}/device/{device_id}", equipment.GetEquipmentByDeviceIDAndSN)
    r.HandleFunc("GET /api/v1/equipment/sn/{sn}/manufacturer/{manufacturer_id}", equipment.GetEquipmentByManufacturerIDAndSN)
    r.HandleFunc("GET /api/v1/equipment/sn/{sn}/manufacturer/{manufacturer_id}/device/{device_id}", equipment.GetEquipmentByManufacturerIDAndDeviceIDAndSN)
    r.HandleFunc("PATCH /api/v1/equipment/sn", equipment.UpdateSerialNumber)
    r.HandleFunc("PATCH /api/v1/equipment", equipment.UpdateEquipment)
    r.HandleFunc("POST /api/v1/equipment", equipment.CreateEquipment)

    // NOTE: Serial number routes


    //NOTE: Handle not found endpoints
    r.HandleFunc("/", equipment.BadEndpointHandler)

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
