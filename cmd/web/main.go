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
    r.HandleFunc("GET /api/v1/device/{id}", devices.GetDeviceByID)
    r.HandleFunc("PATCH /api/v1/device/{id}/name/{name}", devices.UpdateDeviceTypeName)
    r.HandleFunc("PATCH /api/v1/device/{id}/status/{status}", devices.UpdateDeviceTypeStatus)
	r.HandleFunc("POST /api/v1/device/{name}", devices.CreateDeviceType)

	// NOTE: Manufacturer routes
	r.HandleFunc("GET /api/v1/manufacturer", manufactuerers.GetManufacturers)
    r.HandleFunc("GET /api/v1/manufacturer/{id}", manufactuerers.GetManufacturerByID)
    r.HandleFunc("PATCH /api/v1/manufacturer/{id}/name/{name}", manufactuerers.UpdateManufacturerName)
    r.HandleFunc("PATCH /api/v1/manufacturer/{id}/status/{status}", manufactuerers.UpdateManufacturerStatus)
    r.HandleFunc("POST /api/v1/manufacturer/{name}", manufactuerers.CreateManufacturer)

	// NOTE: Equipment routes
    r.HandleFunc("GET /api/v1/equipment/limit/{limit}/offset/{offset}", equipment.GetEquipments)
    r.HandleFunc("GET /api/v1/equipment/{sn}", equipment.GetEquipmentBySN)
    r.HandleFunc("GET /api/v1/equipment/sn-like/{sn}/limit/{limit}/offset/{offset}", equipment.GetEquipmentLikeSN)
    r.HandleFunc("GET /api/v1/equipment/manufacturer/{id}/limit/{limit}/offset/{offset}", equipment.GetEquipmentByManufacturerID)
    r.HandleFunc("GET /api/v1/equipment/device/{id}/limit/{limit}/offset/{offset}", equipment.GetEquipmentByDeviceID)
    r.HandleFunc("GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}/limit/{limit}/offset/{offset}", equipment.GetEquipmentByDeviceIDAndManufacturerID)
    r.HandleFunc("GET /api/v1/equipment/sn/{sn}/device/{device_id}", equipment.GetEquipmentByDeviceIDAndSN)
    // get by manufacturer id & sn
    // r.HandleFunc("GET /api/v1/equipment/manufacturer/{manufacturer_id}/sn/{sn}", equipment.GetEquipmentByManufacturerIDAndSN)
    // get by manufacturer id & devive id & sn
    // r.HandleFunc("GET /api/v1/equipment/manufacturer/{manufacturer_id}/device/{device_id}/sn/{sn}", equipment.GetEquipmentByManufacturerIDAndDeviceIDAndSN)

    // create equipment
    // r.HandleFunc("POST /api/v1/equipment", equipment.CreateEquipment)


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
