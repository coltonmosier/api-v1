package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/coltonmosier/api-v1/internal/database"
	"github.com/coltonmosier/api-v1/internal/helpers"
	"github.com/coltonmosier/api-v1/internal/models"
	"github.com/coltonmosier/api-v1/internal/sqlc"
)

type DeviceHandler struct{}

func (h *DeviceHandler) GetDeviceTypes(w http.ResponseWriter, r *http.Request) {
	var out []models.DeviceType
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/device")
		return
	}
	d, err := q.GetDeviceTypesActive(r.Context())
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "something went wrong", "GET /api/v1/device")
		return
	}
	for _, v := range d {
		out = append(out, models.DeviceType{
			ID:     v.ID,
			Name:   v.Name,
			Status: v.Status,
		})
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

func (h *DeviceHandler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "GET /api/v1/device")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "GET /api/v1/device")
		return
	}

	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/device")
		return
	}
	d, err := q.GetDeviceTypeById(r.Context(), int32(i))
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id does not exists in database", "GET /api/v1/device")
		return
	}
	out := models.DeviceType{
		ID:     int32(i),
		Name:   d.Name,
		Status: string(d.Status),
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

func (h *DeviceHandler) UpdateDeviceTypeName(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/device")
		return
	}
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "GET /api/v1/device")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "GET /api/v1/device")
		return
	}
	// NOTE: Must make a request to /api/v1/device to see if id exists...
	resp, err := http.Get("http://localhost:8081/api/v1/device/" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+id, "PATH /api/v1/device/{id}/name/{newName}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device", "PATH /api/v1/device/{id}/name/{newName}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id does not exist", "GET /api/v1/device")
		return
	}

	name := r.PathValue("name")
	if name == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing name", "PATH /api/v1/device/{id}/name/{newName}")
		return
	}

	err = q.UpdateDeviceType(r.Context(), sqlc.UpdateDeviceTypeParams{ID: int32(i), Name: name})
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update device name", "PATCH /api/v1/device/{id}/name/{newName}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, "device updated with name of "+name)
}

func (h *DeviceHandler) UpdateDeviceTypeStatus(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/device")
		return
	}
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "GET /api/v1/device")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "GET /api/v1/device")
		return
	}
	// NOTE: Must make a request to /api/v1/device to see if id exists...
	resp, err := http.Get("http://localhost:8081/api/v1/device/" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+id, "PATH /api/v1/device/{id}/status/{newStatus}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device", "PATH /api/v1/device/{id}/status/{newStatus}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id does not exist", "GET /api/v1/device")
		return
	}

	status := r.PathValue("status")
	if status == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing name", "PATH /api/v1/device/{id}/status/{newStatus}")
		return
	}

	if status == string(sqlc.DeviceTypeStatusActive) {
		err = q.UpdateDeviceTypeStatus(r.Context(), sqlc.UpdateDeviceTypeStatusParams{ID: int32(i), Status: sqlc.DeviceTypeStatusActive})
	} else {
		err = q.UpdateDeviceTypeStatus(r.Context(), sqlc.UpdateDeviceTypeStatusParams{ID: int32(i), Status: sqlc.DeviceTypeStatusInactive})
	}
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update device status", "PATCH /api/v1/device/{id}/status/{newStatus}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, "device updated with status of "+status)
}

func (h *DeviceHandler) CreateDeviceType(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/device")
		return
	}

	name := r.PathValue("name")
	if name == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing name", "POST /api/v1/device/{newDeviceName}")
		return
	}

	err = q.CreateDeviceType(r.Context(), name)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to write new device to database", "POST /api/v1/device/{newDeviceName}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusCreated, "device created with name of "+name)
}
