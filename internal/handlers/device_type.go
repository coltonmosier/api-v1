package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/coltonmosier/api-v1/internal/database"
	"github.com/coltonmosier/api-v1/internal/helpers"
	"github.com/coltonmosier/api-v1/internal/models"
	"github.com/coltonmosier/api-v1/internal/sqlc"
)

type DeviceHandler struct{}

// GetDeviceTypes Getting all device types
//	@Summary		get all device types
//	@Description	get all device types from the database
//	@Tags			device
//	@x-order		1
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.JsonResponse{MSG=[]models.DeviceType}
//	@Failure		500	{object}	models.JsonResponse 
//	@Router			/device [get]
func (h *DeviceHandler) GetDeviceTypes(w http.ResponseWriter, r *http.Request) {
	var out []models.DeviceType
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/device")
		return
	}
	d, err := q.GetDeviceTypesActive(r.Context())
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "something went wrong with query "+err.Error(), "GET /api/v1/device")
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

// GetDeviceByID Getting a device type by id
//	@Summary		get device type by ID
//	@Description	get device type by ID from the database
//	@Tags			device
//	@x-order		2
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Device ID"	minimum(1)
//	@Success		200	{object}	models.JsonResponse{MSG=models.DeviceType}
//	@Failure		400	{object}	models.JsonResponse
//	@Failure		500	{object}	models.JsonResponse
//	@Router			/device/{id} [get]
func (h *DeviceHandler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/device/{id}")
		return
	}

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

//  UpdateDeviceTypeName Updating a device type name by id
//	@Summary		update device type by name ID
//	@Description	update device type by name ID from the database
//	@Tags			device
//	@x-order		3
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Device ID"	minimum(1)
//	@Param			name	query		string	true	"Device Name"
//	@Success		200		{object}	models.JsonResponse
//	@Failure		400		{object}	models.JsonResponse
//	@Failure		500		{object}	models.JsonResponse
//	@Router			/device/{id}/name [patch]
func (h *DeviceHandler) UpdateDeviceTypeName(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "PATCH /api/v1/device/{id}/name?name={newName}")
		return
	}
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "PATCH /api/v1/device/{id}/name?name={newName}")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id is not a number", "GET /api/v1/device/{id}/name?name={newName}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/device/" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+id, "PATCH /api/v1/device/{id}/name?name={newName}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device/"+id, "PATCH /api/v1/device/{id}/name?name={newName}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/device")
		return
	}

	name := r.FormValue("name")
    if name == ""{
		helpers.JsonResponseError(w, http.StatusBadRequest, "name missing", "PATCH /api/v1/device/{id}/name?name={newName}")
		return
    }

    err = q.UpdateDeviceType(r.Context(), sqlc.UpdateDeviceTypeParams{ID: int32(i), Name: name})
    if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "something went wrong "+err.Error(), "PATCH /api/v1/device/{id}/name?name={newName}")
		return
    }

    msg := fmt.Sprintf("device type %v updated with a name of %v", i, name)

    helpers.JsonResponseSuccess(w, http.StatusOK, msg)
}

//  UpdateDeviceTypeStatus Updating a device type status by id
//	@Summary		update device type by status ID
//	@Description	update device type by status ID from the database
//	@Tags			device
//	@x-order		3
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Device ID"		minimum(1)
//	@Param			status	query		string	true	"Device Status"	Enums(active,inactive)
//	@Success		200		{object}	models.JsonResponse
//	@Failure		400		{object}	models.JsonResponse
//	@Failure		500		{object}	models.JsonResponse
//	@Router			/device/{id}/status [patch]
func (h *DeviceHandler) UpdateDeviceType(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "PATCH /api/v1/device/{id}/status?status={newStatus}")
		return
	}
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "PATCH /api/v1/device/{id}status={newStatus}")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id is not a number", "GET /api/v1/device")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/device/" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+id, "PATCH /api/v1/device/{id}/status?status={newStatus}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device/"+id, "PATCH /api/v1/device/{id}/status?status={newStatus}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "PATCH /api/v1/device/{id}/status?status={newStatus}")
		return
	}

    status := r.FormValue("status")

    if status == ""{
		helpers.JsonResponseError(w, http.StatusBadRequest, "status cannot be empty", "PATCH /api/v1/device/{id}/status?status={newStatus}")
		return
    }else if status != "active" && status != "inactive"{
		helpers.JsonResponseError(w, http.StatusBadRequest, "status must be either active or inactive", "PATCH /api/v1/device/{id}/status?status={newStatus}")
		return
    }else if status == "active"{
        err = q.UpdateDeviceTypeStatus(r.Context(), sqlc.UpdateDeviceTypeStatusParams{ID: int32(i), Status: sqlc.DeviceTypeStatusActive})
        if err != nil{
            helpers.JsonResponseError(w, http.StatusBadRequest, "something went wrong with query "+ err.Error(), "PATCH /api/v1/device/{id}/status?status={newStatus}")
            return
        }
    } else { // status must be inactive here
        err = q.UpdateDeviceTypeStatus(r.Context(), sqlc.UpdateDeviceTypeStatusParams{ID: int32(i), Status: sqlc.DeviceTypeStatusInactive})
        if err != nil{
            helpers.JsonResponseError(w, http.StatusBadRequest, "something went wrong with query "+ err.Error(), "PATCH /api/v1/device/{id}/status?status={newStatus}")
            return
        }
    }
    msg := fmt.Sprintf("device type %v with status of %v", i, status)

    helpers.JsonResponseSuccess(w, http.StatusOK, msg)
}
//  CreateDeviceType Creating a device type
//	@Summary		create device type 
//	@Description	create device type for the database
//	@Tags			device 
//	@x-order		4
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	true	"Device Name"
//	@Success		200		{object}	models.JsonResponse
//	@Failure		400		{object}	models.JsonResponse
//	@Failure		500		{object}	models.JsonResponse
//	@Router			/device [post]
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
        if strings.ToLower(v.Name) == strings.ToLower(name) {
            helpers.JsonResponseError(w, http.StatusBadRequest, "device type already exists", "POST /api/v1/device?name={newName}")
            return
        }
    }

    err = q.CreateDeviceType(r.Context(), name)
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to create device", "POST /api/v1/device?name={newName}")
        return
    }

    helpers.JsonResponseSuccess(w, http.StatusOK, "device created with name of "+name)
}
