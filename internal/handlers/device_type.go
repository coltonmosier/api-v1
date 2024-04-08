package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coltonmosier/api-v1/internal/models"
)

type Handler struct{}

var devices = []models.DeviceType{
	{
        ID: 1, 
        Name: "apple", 
        Status: "active",
    },
	{
        ID: 2, 
        Name: "samsung", 
        Status: "active",
    },
}

func (h *Handler) GetDeviceType(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
    // I want to return the devices slice as a JSON response
    output, err := json.Marshal(devices)
    if err != nil {
        log.Println("Error marshalling devices")
        return
    }
    w.Write([]byte(output))
}

func (h *Handler) CreateDeviceType(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	log.Println("Creating device type")
	log.Println("Checking if device type exists")
	log.Println("Created device type")
	log.Println("Device type name: ", name)
	w.Write([]byte(`{"id": "3", "name": ` + name + `}"}`))
}
