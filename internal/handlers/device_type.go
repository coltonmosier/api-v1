package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coltonmosier/api-v1/internal/models"
)

type DeviceHandler struct{}

var devices = []models.DeviceType{
	{
		ID:     1,
		Name:   "apple",
		Status: "active",
	},
	{
		ID:     2,
		Name:   "samsung",
		Status: "active",
	},
}

func (h *DeviceHandler) GetDeviceType(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)

	out := models.JsonResponse{
		Status: "ok",
        Message: devices,
		Action: "none",
	}

	output, err := json.Marshal(out)
	if err != nil {
		log.Println("Error marshalling JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(output)
}

func (h *DeviceHandler) CreateDeviceType(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	log.Println("Creating device type")
	log.Println("Checking if device type exists")
	log.Println("Created device type")
	log.Println("Device type name: ", name)
	w.Write([]byte(`{"id": "3", "name": ` + name + `}"}`))
}
