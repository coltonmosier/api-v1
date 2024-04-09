package handlers

import (
	"net/http"

	"github.com/coltonmosier/api-v1/internal/database"
	"github.com/coltonmosier/api-v1/internal/helpers"
	"github.com/coltonmosier/api-v1/internal/models"
)

type DeviceHandler struct{}

func (h *DeviceHandler) GetDeviceTypes(w http.ResponseWriter, r *http.Request) {
    var out []models.DeviceType
    q, err := database.InitEquipmentDatabase()
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/manufacturer")
        return
    }
    d, err := q.GetDeviceTypesActive(r.Context())
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "something went wrong", "GET /api/v1/manufacturer")
        return
    }
    for _, v := range d {
        out = append(out, models.DeviceType{
            ID:     v.ID,
            Name:   v.Name,
            Status: "active",
        })
    }
    
    helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

func (h *DeviceHandler) GetDeviceByValue(w http.ResponseWriter, r *http.Request) {
}


func (h *DeviceHandler) UpdateDeviceTypeName(w http.ResponseWriter, r *http.Request) {
}

func (h *DeviceHandler) UpdateDeviceTypeStatus(w http.ResponseWriter, r *http.Request) {
}

func (h *DeviceHandler) CreateDeviceType(w http.ResponseWriter, r *http.Request) {
}

func (h *DeviceHandler) DeleteDeviceType(w http.ResponseWriter, r *http.Request) {
}
