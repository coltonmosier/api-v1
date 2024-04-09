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

func (h *DeviceHandler) GetDeviceTypes(w http.ResponseWriter, r *http.Request) {
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

func (h *DeviceHandler) GetDeviceByName(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    w.WriteHeader(http.StatusOK)
    log.Println("Getting device type by name")
    log.Println("Device type name: ", name)
    w.Write([]byte(`{"id": "3", "name": ` + name + `}"}`))
}

func (h *DeviceHandler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    w.WriteHeader(http.StatusOK)
    log.Println("Getting device type by ID")
    log.Println("Device type ID: ", id)
    w.Write([]byte(`{"id": ` + id + `, "name": "test"}`))
}

func (h *DeviceHandler) UpdateDeviceTypeName(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    name := r.PathValue("name")
    w.WriteHeader(http.StatusOK)
    log.Println("Updating device type name")
    log.Println("Device type ID: ", id)
    w.Write([]byte(`{"id": ` + id + `, "name": ` + name + `}`))
}

func (h *DeviceHandler) UpdateDeviceTypeStatus(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    status := r.PathValue("status")
    w.WriteHeader(http.StatusOK)
    log.Println("Updating device type status")
    log.Println("Device type ID: ", id)
    w.Write([]byte(`{"id": ` + id + `, "status": ` + status + `}`))
}

func (h *DeviceHandler) CreateDeviceType(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    w.WriteHeader(http.StatusOK)
    log.Println("Creating device type")
    log.Println("Device type name: ", name)
    w.Write([]byte(`{"name": ` + name + `}`))
}

func (h *DeviceHandler) DeleteDeviceType(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    w.WriteHeader(http.StatusOK)
    log.Println("Deleting device type")
    log.Println("Device type ID: ", id)
    w.Write([]byte(`{"id": ` + id + `}`))
}
