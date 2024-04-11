package handlers

import (
	"encoding/json"
	"log"
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
			Status: string(v.Status),
		})
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

func (h *DeviceHandler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing devie id", "GET /api/v1/device/{id}")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id is not a number", "GET /api/v1/device/{id}")
		return
	}

	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/device/{id}")
		return
	}
	d, err := q.GetDeviceTypeById(r.Context(), int32(i))
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id does not exists in database", "GET /api/v1/device")
		return
	}
	out := models.DeviceType{
		ID:     int32(i),
		Name:   d.Name,
		Status: string(d.Status),
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

func (h *DeviceHandler) UpdateDeviceType(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
		return
	}
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id is not a number", "GET /api/v1/device")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/device/" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+id, "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device", "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/device")
		return
	}

	name := r.FormValue("name")
	status := r.FormValue("status")
	if name == "" && status == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing name and/or status", "PATH /api/v1/device/{id}?name={newName}&status={newStatus}")
		return
	} else if name != "" && status != "" {
		err = q.UpdateDeviceType(r.Context(), sqlc.UpdateDeviceTypeParams{ID: int32(i), Name: name})
		if err != nil {
			helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update device name", "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
			return
		}
		if status == string(sqlc.DeviceTypeStatusActive) {
			err = q.UpdateDeviceTypeStatus(r.Context(), sqlc.UpdateDeviceTypeStatusParams{ID: int32(i), Status: sqlc.DeviceTypeStatusActive})
			if err != nil {
				helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update device status"+err.Error(), "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
				return
			}
		} else if status == string(sqlc.DeviceTypeStatusInactive) {
			err = q.UpdateDeviceTypeStatus(r.Context(), sqlc.UpdateDeviceTypeStatusParams{ID: int32(i), Status: sqlc.DeviceTypeStatusInactive})
			if err != nil {
				helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update device status"+err.Error(), "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
				return
			}
		} else {
			helpers.JsonResponseError(w, http.StatusBadRequest, "status must be either 'active' or 'inactive'", "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
			return
		}
		helpers.JsonResponseSuccess(w, http.StatusOK, "device updated with name of "+name+" and status of "+status)
		return
	} else if name != ""  && status == ""{
		err = q.UpdateDeviceType(r.Context(), sqlc.UpdateDeviceTypeParams{ID: int32(i), Name: name})
		if err != nil {
			helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update device name", "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
			return
		}
		helpers.JsonResponseSuccess(w, http.StatusOK, "device updated with name of "+name)
		return
	} else if status != "" && name == ""{
		if status == string(sqlc.DeviceTypeStatusActive) {
			err = q.UpdateDeviceTypeStatus(r.Context(), sqlc.UpdateDeviceTypeStatusParams{ID: int32(i), Status: sqlc.DeviceTypeStatusActive})
			if err != nil {
				helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update device status"+err.Error(), "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
				return
			}
		} else if status == string(sqlc.DeviceTypeStatusInactive) {
			err = q.UpdateDeviceTypeStatus(r.Context(), sqlc.UpdateDeviceTypeStatusParams{ID: int32(i), Status: sqlc.DeviceTypeStatusInactive})
			if err != nil {
				helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update device status"+err.Error(), "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
				return
			}
		} else {
			helpers.JsonResponseError(w, http.StatusBadRequest, "status must be either 'active' or 'inactive'", "PATCH /api/v1/device/{id}?name={newName}&status={newStatus}")
			return
		}
		helpers.JsonResponseSuccess(w, http.StatusOK, "device updated with status of "+status)
		return
	}
}

func (h *DeviceHandler) CreateDeviceType(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "POST /api/v1/device?name={newName}")
		return
	}

    name := r.FormValue("name")
    if name == "" {
        helpers.JsonResponseError(w, http.StatusBadRequest, "missing name", "POST /api/v1/device?name={newName}")
        return
    }

    req, err := http.Get("http://localhost:8081/api/v1/device")
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device", "POST /api/v1/device?name={newName}")
        return
    }
    defer req.Body.Close()

    var resp models.JsonResponse
    err = json.NewDecoder(req.Body).Decode(&resp)
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device", "POST /api/v1/device?name={newName}")
        return
    }

    if resp.Status == "ERROR" {
        helpers.JsonResponseError(w, http.StatusBadRequest, resp.Message, "POST /api/v1/device?name={newName}")
        return
    }

    var device []models.DeviceType 
    for _, v := range resp.Message.([]interface{}) {
        if m, ok := v.(map[string]interface{}); ok {
        // Extract values from the map and create a DeviceType instance
        var d models.DeviceType
        if id, ok := m["id"].(float64); ok {
            d.ID = int32(id)
        } else {
            // Handle error when id cannot be converted to float64
            helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to convert id to float64", "POST /api/v1/device?name={newName}")
            return
        }
        if name, ok := m["name"].(string); ok {
            d.Name = name
        } else {
            // Handle error when name cannot be converted to string
            helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to convert name to string", "POST /api/v1/device?name={newName}")
            return
        }
        if status, ok := m["status"].(string); ok {
            d.Status = status
        } else {
            // Handle error when status cannot be converted to string
            helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to convert status to string", "POST /api/v1/device?name={newName}")
            return
        }
        // Add the created DeviceType instance to the device slice
        device = append(device, d)
    } else {
        // Type assertion failed for an element in the slice
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to assert map[string]interface{} type", "POST /api/v1/device?name={newName}")
        return
    }
    }

    for _, v := range device {
        if v.Name == name {
            helpers.JsonResponseError(w, http.StatusBadRequest, "device type already exists", "POST /api/v1/device?name={newName}")
            return
        }
    }

    log.Println("Creating device type with name of "+name)

    err = q.CreateDeviceType(r.Context(), name)
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to create device", "POST /api/v1/device?name={newName}")
        return
    }

    helpers.JsonResponseSuccess(w, http.StatusOK, "device created with name of "+name)
}
