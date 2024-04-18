package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/coltonmosier/api-v1/internal/database"
	"github.com/coltonmosier/api-v1/internal/helpers"
	"github.com/coltonmosier/api-v1/internal/models"
	"github.com/coltonmosier/api-v1/internal/sqlc"
)

type EquipmentHandler struct{}

func (h *EquipmentHandler) BadEndpointHandler(w http.ResponseWriter, r *http.Request) {
	helpers.JsonResponseError(w, http.StatusNotFound, "endpoint not found", "none")
}

// GetEquipments get all equipment
//
//	@Summary		get all equipments with
//	@Description	get all equipment with manufacturers from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			all	query		bool	true	"active and inactive"
//	@Success		200	{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		500	{object}	models.JsonResponse
//	@Router			/equipment [get]
func (h *EquipmentHandler) GetEquipments(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment")
		return
	}
	d, err := q.GetAllEquipment(r.Context())
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment", "GET /api/v1/equipment")
		return
	}

    all := r.FormValue("all")
    if all == "true" {
        var e []models.Equipment
        for _, v := range d {
            e = append(e, models.Equipment{
                AutoID:         v.AutoID,
                DeviceTypeID:   v.DeviceTypeID,
                ManufacturerID: v.ManufacturerID,
                SerialNumber:   v.SerialNumber,
                Status:         string(v.Status),
            })
        }

        helpers.JsonResponseSuccess(w, http.StatusOK, e)
        return
    } else {
        var e []models.Equipment
        for _, v := range d {
            if v.Status == sqlc.SerialNumbersStatusActive {
                e = append(e, models.Equipment{
                    AutoID:         v.AutoID,
                    DeviceTypeID:   v.DeviceTypeID,
                    ManufacturerID: v.ManufacturerID,
                    SerialNumber:   v.SerialNumber,
                    Status:         string(v.Status),
                })
            }
        }
        helpers.JsonResponseSuccess(w, http.StatusOK, e)
        return
    }
}

// GetEquipmentBySN get equipment by serial number
//
//	@Summary		get equipment by serial number
//	@Description	get equipment by serial number from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			sn	query		string	true	"serial number"
//	@Success		200	{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400	{object}	models.JsonResponse
//	@Failure		500	{object}	models.JsonResponse
//	@Router			/equipment/sn [get]
func (h *EquipmentHandler) GetEquipmentBySN(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment?sn={sn}")
		return
	}

	sn := r.FormValue("sn")
	if sn == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "GET /api/v1/equipment/sn/{sn}")
		return
	}

	s := strings.Split(sn, "-")[0]
	if s != "SN" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "serial number must start with SN-", "GET /api/v1/equipment/sn/{sn}")
		return
	}

	if len(sn) > 68 {
		helpers.JsonResponseError(w, http.StatusBadRequest, "serial number cannot be longer than 64 characters", "GET /api/v1/equipment/sn/{sn}")
		return
	}

	d, err := q.GetEquipmentBySerialNumber(r.Context(), sn)
	if err == sql.ErrNoRows {
		helpers.JsonResponseError(w, http.StatusBadRequest, "equipment does not exist in database", "GET /api/v1/equipment")
		return
	} else if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment", "GET /api/v1/equipment/sn/{sn}")
		return
	}

	e := models.Equipment{
		AutoID:         d.AutoID,
		DeviceTypeID:   d.DeviceTypeID,
		ManufacturerID: d.ManufacturerID,
		SerialNumber:   d.SerialNumber,
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

// GetEquipmentByID get equipment by auto ID
//
//	@Summary		get equipment by auto ID
//	@Description	get equipment by auto_id from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			id	query		int	true	"auto_id"	minimum(1)
//	@Success		200	{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400	{object}	models.JsonResponse
//	@Failure		500	{object}	models.JsonResponse
//	@Router			/equipment/id [get]
func (h *EquipmentHandler) GetEquipmentByID(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/id?id={id}")
		return
	}
	id := r.FormValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "GET /api/v1/equipment/id?id={id}")
		return
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "GET /api/v1/equipment/id?id={id}")
		return
	}

	d, err := q.GetEquipmentByAutoID(r.Context(), int32(i))
	if err == sql.ErrNoRows {
		helpers.JsonResponseError(w, http.StatusBadRequest, "equipment id does not exist", "GET /api/v1/equipment/id?id={id}")
		return
	}
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment"+err.Error(), "GET /api/v1/equipment/id?id={id}")
		return
	}

	e := models.Equipment{
		AutoID:         d.AutoID,
		DeviceTypeID:   d.DeviceTypeID,
		ManufacturerID: d.ManufacturerID,
		SerialNumber:   d.SerialNumber,
		Status:         string(d.Status),
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

// GetEquipmentLikeSn get equipment like serial number
//
//	@Summary		get equipment like serial number
//	@Description	get equipment like serial number from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			sn		path		string	true	"serial number"
//	@Success		200		{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400		{object}	models.JsonResponse
//	@Failure		500		{object}	models.JsonResponse
//	@Router			/equipment/sn-like/{sn} [get]
func (h *EquipmentHandler) GetEquipmentLikeSN(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/sn-like/{sn}")
		return
	}

	sn := r.PathValue("sn")
	if sn == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "GET /api/v1/equipment/sn-like/{sn}")
		return
	}

	sn = "%" + sn + "%"

	d, err := q.GetEquipmentLikeSerialNumber(r.Context(), sn)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment "+err.Error(), "GET /api/v1/equipment/sn-like/{sn}")
		return
	}

	var e []models.Equipment
	for _, v := range d {
		e = append(e, models.Equipment{
			AutoID:         v.AutoID,
			DeviceTypeID:   v.DeviceTypeID,
			ManufacturerID: v.ManufacturerID,
			SerialNumber:   v.SerialNumber,
			Status:         string(v.Status),
		})
	}
	if e == nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "no equipment found", "GET /api/v1/equipment/sn-like/{sn}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

// GetEquipmentByManufacturerID get equipment by manufacturer id
//
//	@Summary		get equipment by manufacturer id
//	@Description	get equipment by manufacturer id from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"manufacturer id"	minimum(1)
//	@Success		200		{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400		{object}	models.JsonResponse
//	@Failure		500		{object}	models.JsonResponse
//	@Router			/equipment/manufacturer/{id} [get]
func (h *EquipmentHandler) GetEquipmentByManufacturerID(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/manufacturer/{id}")
		return
	}

	manufacturerID := r.PathValue("id")
	if manufacturerID == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing manufacturer id", "GET /api/v1/equipment/manufacturer/{id}")
		return
	}
	id, err := strconv.Atoi(manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "GET /api/v1/equipment/manufacturer/{id}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/manufacturer/" + manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer/"+manufacturerID, "GET /api/v1/manufacturer/{id}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "GET /api/v1/manufacturer/{id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/manufacturer")
		return
	}


	d, err := q.GetEquipmentByManufacturer(r.Context(), int32(id))
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment", "GET /api/v1/equipment/manufacturer/{id}")
		return
	}

	var e []models.Equipment
	for _, v := range d {
		e = append(e, models.Equipment{
			AutoID:         v.AutoID,
			DeviceTypeID:   v.DeviceTypeID,
			ManufacturerID: v.ManufacturerID,
			SerialNumber:   v.SerialNumber,
			Status:         string(v.Status),
		})
	}

	if e == nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "no equipment found", "GET /api/v1/equipment/manufacturer/{id}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

// GetEquipmentByDeviceID get equipment by device id
//
//	@Summary		get equipment by device id
//	@Description	get equipment by device id from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"device id"	minimum(1)
//	@Success		200		{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400		{object}	models.JsonResponse
//	@Failure		500		{object}	models.JsonResponse
//	@Router			/equipment/device/{id} [get]
func (h *EquipmentHandler) GetEquipmentByDeviceID(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/device/{id}")
		return
	}

	deviceID := r.PathValue("id")
	if deviceID == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "GET /api/v1/equipment/device/{id}")
		return
	}
	id, err := strconv.Atoi(deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id is not a number", "GET /api/v1/equipment/device/{id}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/device/" + deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+deviceID, "GET /api/v1/device/{id}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device", "none")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/device")
		return
	}

	d, err := q.GetEquipmentByDeviceType(r.Context(), int32(id))
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment"+err.Error(), "GET /api/v1/equipment/device/{id}")
		return
	}

	var e []models.Equipment
	for _, v := range d {
		e = append(e, models.Equipment{
			AutoID:         v.AutoID,
			DeviceTypeID:   v.DeviceTypeID,
			ManufacturerID: v.ManufacturerID,
			SerialNumber:   v.SerialNumber,
			Status:         string(v.Status),
		})
	}

	if e == nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "no equipment found", "GET /api/v1/equipment/device/{id}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

// GetEquipmentByDeviceIDAndManufacturerID get equipment by device id and manufacturer id
//
//	@Summary		get equipment by device id and manufacturer id
//	@Description	get equipment by device id and manufacturer id from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			device_id		path		int	true	"device id"			minimum(1)
//	@Param			manufacturer_id	path		int	true	"manufacturer id"	minimum(1)
//	@Success		200				{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400				{object}	models.JsonResponse
//	@Failure		500				{object}	models.JsonResponse
//	@Router			/equipment/device/{device_id}/manufacturer/{manufacturer_id} [get]
func (h *EquipmentHandler) GetEquipmentByDeviceIDAndManufacturerID(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}")
		return
	}

	deviceID := r.PathValue("device_id")
	if deviceID == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}")
		return
	}
	did, err := strconv.Atoi(deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id is not a number", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/device/" + deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+deviceID, "GET /api/v1/device/{id}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device", "GET /api/v1/device/{id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/device")
		return
	}

	manufacturerID := r.PathValue("manufacturer_id")
	if manufacturerID == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing manufacturer id", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}")
		return
	}
	mid, err := strconv.Atoi(manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}")
		return
	}
	resp, err = http.Get("http://localhost:8081/api/v1/manufacturer/" + manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer/"+manufacturerID, "GET /api/v1/manufacturer/{id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "GET /api/v1/manufacturer/{id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/manufacturer")
		return
	}
	d, err := q.GetEquipmentByDeviceTypeAndManufacturer(r.Context(), sqlc.GetEquipmentByDeviceTypeAndManufacturerParams{DeviceTypeID: int32(did), ManufacturerID: int32(mid)})
	if err != nil && err != sql.ErrNoRows {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment"+err.Error(), "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}")
		return
	} else if err == sql.ErrNoRows {
		helpers.JsonResponseError(w, http.StatusBadRequest, "no equipment found", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}")
		return
	}

	var e []models.Equipment
	for _, v := range d {
		e = append(e, models.Equipment{
			AutoID:         v.AutoID,
			DeviceTypeID:   v.DeviceTypeID,
			ManufacturerID: v.ManufacturerID,
			SerialNumber:   v.SerialNumber,
			Status:         string(v.Status),
		})
	}

	if e == nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "no equipment found", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

// GetEquipmentByDeviceIDAndSN get equipment by device id and serial number
//
//	@Summary		get equipment by device id and serial number
//	@Description	get equipment by device id and serial number from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			device_id	path		int		true	"device id"	minimum(1)
//	@Param			sn			path		string	true	"serial number"
//	@Success		200			{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400			{object}	models.JsonResponse
//	@Failure		500			{object}	models.JsonResponse
//	@Router			/equipment/sn/{sn}/device/{device_id} [get]
func (h *EquipmentHandler) GetEquipmentByDeviceIDAndSN(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/{sn}")
		return
	}
	// get equipment by sn to see if it exists
	sn := r.PathValue("sn")
	if sn == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "GET /api/v1/equipment/device/{device_id}/sn/{sn}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/equipment/" + sn)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/equipment/"+sn, "GET /api/v1/equipment/{sn}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "GET /api/v1/equipment/{sn}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/equipment/sn-like/{sn}")
		return
	}

	// get equipment by device id to see if it exists
	deviceID := r.PathValue("device_id")
	if deviceID == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "GET /api/v1/equipment/device/{id}/sn/{sn}")
		return
	}

	id, err := strconv.Atoi(deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "device id is not a number", "GET /api/v1/device")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/device/" + deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+deviceID, "GET /api/v1/device/{id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device", "GET /api/v1/device/{id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/device")
		return
	}

	d, err := q.GetEquipmentByDeviceTypeAndSerialNumber(r.Context(), sqlc.GetEquipmentByDeviceTypeAndSerialNumberParams{SerialNumber: sn, DeviceTypeID: int32(id)})
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment", "GET /api/v1/equipment/device/{device_id}/sn/{sn}")
		return
	}

	e := models.Equipment{
		AutoID:         d.AutoID,
		DeviceTypeID:   d.DeviceTypeID,
		ManufacturerID: d.ManufacturerID,
		SerialNumber:   d.SerialNumber,
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

// GetEquipmentByManufacturerIDAndSN get equipment by manufacturer id and serial number
//
//	@Summary		get equipment by manufacturer id and serial number
//	@Description	get equipment by manufacturer id and serial number from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			manufacturer_id	path		int		true	"manufacturer id"	minimum(1)
//	@Param			sn				path		string	true	"serial number"
//	@Success		200				{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400				{object}	models.JsonResponse
//	@Failure		500				{object}	models.JsonResponse
//	@Router			/equipment/sn/{sn}/manufacturer/{manufacturer_id} [get]
func (h *EquipmentHandler) GetEquipmentByManufacturerIDAndSN(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/{sn}")
		return
	}
	// get equipment by sn to see if it exists
	sn := r.PathValue("sn")
	if sn == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "GET /api/v1/equipment/device/{device_id}/sn/{sn}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/equipment/" + sn)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/equipment/"+sn, "GET /api/v1/equipment/{sn}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "GET /api/v1/equipment/{sn}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/equipment/sn-like/{sn}")
		return
	}

	manufacturerID := r.PathValue("manufacturer_id")

	id, err := strconv.Atoi(manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "GET /api/v1/manufacturer/{id}")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/manufacturer/" + manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer/"+manufacturerID, "GET /api/v1/manufacturer/{id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "GET /api/v1/manufacturer/{id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/manufacturer")
		return
	}

	d, err := q.GetEquipmentByManufacturerAndSerialNumber(r.Context(), sqlc.GetEquipmentByManufacturerAndSerialNumberParams{SerialNumber: sn, ManufacturerID: int32(id)})
	if err == sql.ErrNoRows {
		helpers.JsonResponseError(w, http.StatusBadRequest, "equipment does not exists", "none")
		return
	} else if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "something went wrong with query"+err.Error(), "none")
		return
	}

	out := models.Equipment{
		AutoID:         d.AutoID,
		DeviceTypeID:   d.DeviceTypeID,
		ManufacturerID: d.ManufacturerID,
		SerialNumber:   d.SerialNumber,
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

// GetEquipmentByManufacturerIDAndDeviceIDAndSN get equipment by manufacturer id and serial number and device id
//
//	@Summary		get equipment by manufacturer id and serial number and device id
//	@Description	get equipment by manufacturer id and serial number and device id from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			manufacturer_id	path		int		true	"manufacturer id"	minimum(1)
//	@Param			device_id		path		int		true	"device id"			minimum(1)
//	@Param			sn				path		string	true	"serial number"
//	@Success		200				{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400				{object}	models.JsonResponse
//	@Failure		500				{object}	models.JsonResponse
//	@Router			/equipment/sn/{sn}/manufacturer/{manufacturer_id}/device/{device_id} [get]
func (h *EquipmentHandler) GetEquipmentByManufacturerIDAndDeviceIDAndSN(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/{sn}")
		return
	}
	// get equipment by sn to see if it exists
	sn := r.PathValue("sn")
	if sn == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "GET /api/v1/equipment/device/{device_id}/sn/{sn}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/equipment/" + sn)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/equipment/"+sn, "GET /api/v1/equipment/{sn}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "GET /api/v1/equipment/{sn}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/equipment/sn-like/{sn}")
		return
	}

	manufacturerID := r.PathValue("manufacturer_id")

	mid, err := strconv.Atoi(manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "GET /api/v1/manufacturer/{id}")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/manufacturer/" + manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer/"+manufacturerID, "GET /api/v1/manufacturer/{id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "GET /api/v1/manufacturer/{id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/manufacturer")
		return
	}

	deviceID := r.PathValue("device_id")
	if deviceID == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "GET /api/v1/equipment/device/{id}/sn/{sn}")
		return
	}

	did, err := strconv.Atoi(deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "device id is not a number", "GET /api/v1/device")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/device/" + deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+deviceID, "GET /api/v1/device/{id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device", "GET /api/v1/device/{id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/device")
		return
	}

	d, err := q.GetEquipmentByDeviceTypeManufacturerAndSerialNumber(r.Context(), sqlc.GetEquipmentByDeviceTypeManufacturerAndSerialNumberParams{
		SerialNumber:   sn,
		DeviceTypeID:   int32(did),
		ManufacturerID: int32(mid),
	})

	if err == sql.ErrNoRows {
		helpers.JsonResponseError(w, http.StatusBadRequest, "equipment does not exists", "GET /api/v1/equipment")
		return
	} else if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "something went wrong with query "+err.Error(), "GET /api/v1/equipment")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, d)
}

// GetEquipmentByManufacturerIDAndDeviceIDLikeSN get equipment by manufacturer id like serial number and device id
//
//	@Summary		get equipment by manufacturer id like serial number and device id
//	@Description	get equipment by manufacturer id like serial number and device id from the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			manufacturer_id	path		int		true	"manufacturer id"	minimum(1)
//	@Param			device_id		path		int		true	"device id"			minimum(1)
//	@Param			sn				path		string	true	"serial number"
//	@Success		200				{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400				{object}	models.JsonResponse
//	@Failure		500				{object}	models.JsonResponse
//	@Router			/equipment/sn-like/{sn}/manufacturer/{manufacturer_id}/device/{device_id} [get]
func (h *EquipmentHandler) GetEquipmentByManufacturerIDAndDeviceIDLikeSN(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/{sn}")
		return
	}
	// get equipment by sn to see if it exists
	sn := r.PathValue("sn")
	if sn == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "GET /api/v1/equipment/device/{device_id}/sn/{sn}")
		return
	}

	manufacturerID := r.PathValue("manufacturer_id")

	mid, err := strconv.Atoi(manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "GET /api/v1/manufacturer/{id}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/manufacturer/" + manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer/"+manufacturerID, "GET /api/v1/manufacturer/{id}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "GET /api/v1/manufacturer/{id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/manufacturer")
		return
	}

	deviceID := r.PathValue("device_id")
	if deviceID == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "GET /api/v1/equipment/device/{id}/sn/{sn}")
		return
	}

	did, err := strconv.Atoi(deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "device id is not a number", "GET /api/v1/device")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/device/" + deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+deviceID, "GET /api/v1/device/{id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device", "GET /api/v1/device/{id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/device")
		return
	}

	sn = "%" + sn + "%"

	d, err := q.GetEquipmentByDeviceTypeManufacturerLikeSerialNumber(r.Context(), sqlc.GetEquipmentByDeviceTypeManufacturerLikeSerialNumberParams{
		SerialNumber:   sn,
		DeviceTypeID:   int32(did),
		ManufacturerID: int32(mid),
	})

	if err == sql.ErrNoRows {
		helpers.JsonResponseError(w, http.StatusBadRequest, "equipment does not exists", "GET /api/v1/equipment")
		return
	} else if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "something went wrong with query "+err.Error(), "GET /api/v1/equipment")
		return
	}

	var e []models.Equipment
	for _, v := range d {
		e = append(e, models.Equipment{
			AutoID:         v.AutoID,
			DeviceTypeID:   v.DeviceTypeID,
			ManufacturerID: v.ManufacturerID,
			SerialNumber:   v.SerialNumber,
			Status:         string(v.Status),
		})
	}

	if e == nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "no equipment found", "GET /api/v1/equipment")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

// UpdateSerialNumber update equipment serial number
//
//	@Summary		update equipment serial number
//	@Description	update equipment serial number in the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			id	query		int		true	"equipment id"	minimum(1)
//	@Param			sn	query		string	true	"serial number"
//	@Success		200	{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400	{object}	models.JsonResponse
//	@Failure		500	{object}	models.JsonResponse
//	@Router			/equipment/sn [patch]
func (h *EquipmentHandler) UpdateSerialNumber(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}

	id := r.FormValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/equipment/id?id=" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/equipment?id="+id, "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}
	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}

	sn := r.FormValue("sn")
	if sn == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/equipment/sn/" + sn)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/equipment/sn/"+sn, "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}

	err = q.UpdateSerialNumber(r.Context(), sqlc.UpdateSerialNumberParams{AutoID: int32(i), SerialNumber: sn})
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to update serial number in database", "PATCH /api/v1/equipment?id={id}&sn={sn}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, "serial number updated")
}

// UpdateEquipment update equipment
//
//	@Summary		update equipment
//	@Description	update equipment in the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			id				query		int		true	"equipment id"	minimum(1)
//	@Param			sn				query		string	true	"serial number"
//	@Param			manufacturer_id	query		int		true	"manufacturer id"	minimum(1)
//	@Param			device_id		query		int		true	"device id"			minimum(1)
//	@Success		200				{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400				{object}	models.JsonResponse
//	@Failure		500				{object}	models.JsonResponse
//	@Router			/equipment [patch]
func (h *EquipmentHandler) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	id := r.FormValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	_, err = strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/equipment/id?id=" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/equipment?id="+id, "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	sn := r.FormValue("sn")
	if sn == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/equipment/sn/" + sn)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/equipment/sn/"+sn, "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	if req.Status == "ERROR" && req.Message != "equipment does not exist in database" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	tmp := req.Message.(models.Equipment)
	if tmp.SerialNumber != sn {
		helpers.JsonResponseError(w, http.StatusBadRequest, "equipment does not exist in database", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	did := r.FormValue("device_id")
	if did == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	d, err := strconv.Atoi(did)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id is not a number", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/device/" + did)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/device/"+did, "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/device", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/device")
		return
	}

	mid := r.FormValue("manufacturer_id")
	if mid == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing manufacturer id", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	m, err := strconv.Atoi(mid)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/manufacturer/" + mid)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer/"+mid, "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/manufacturer")
		return
	}

	err = q.UpdateEquipment(r.Context(), sqlc.UpdateEquipmentParams{SerialNumber: sn, DeviceTypeID: int32(d), ManufacturerID: int32(m)})
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to update equipment in database", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, "equipment updated")
}

// UpdateEquipmentStatus update equipment status
//
//	@Summary		update equipment status
//	@Description	update equipment status in the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"equipment id"		minimum(1)
//	@Param			status	query		string	true	"equipment status"	Enums("active", "inactive")
//	@Success		200		{object}	models.JsonResponse
//	@Failure		400		{object}	models.JsonResponse
//	@Failure		500		{object}	models.JsonResponse
//	@Router			/equipment/{id}/status [patch]
func (h *EquipmentHandler) UpdateEquipmentStatus(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "PATCH /api/v1/equipment/{id}/status?status={status}")
		return
	}
	id := r.FormValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "PATCH /api/v1/equipment/{id}/status?status={status}")
		return
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "PATCH /api/v1/equipment/{id}/status?status={status}")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/equipment/id?id=" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/equipment?id="+id, "PATCH /api/v1/equipment/{id}/status?status={status}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "PATCH /api/v1/equipment/{id}/status?status={status}")
		return
	}
	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "PATCH /api/v1/equipment/{id}/status?status={status}")
		return
	}

	status := r.FormValue("status")
	if status == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "status cannot be empty", "none")
		return
	}
	if status != "active" && status != "inactive" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "status must be either active or inactive", "PATCH /api/v1/equipment/{id}/status?status={status}")
		return
	} else if status == "active" {
		err = q.UpdateEquipmentStatus(r.Context(), sqlc.UpdateEquipmentStatusParams{AutoID: int32(i), Status: sqlc.SerialNumbersStatusActive})
		if err != nil {
			helpers.JsonResponseError(w, http.StatusBadRequest, "something went wrong with sql statement"+err.Error(), "PATCH /api/v1/equipment/{id}/status?status={status}")
			return
		}
	} else { // status must be inactive here
		err = q.UpdateEquipmentStatus(r.Context(), sqlc.UpdateEquipmentStatusParams{AutoID: int32(i), Status: sqlc.SerialNumbersStatusInactive})
		if err != nil {
			helpers.JsonResponseError(w, http.StatusBadRequest, "something went wrong with sql statement"+err.Error(), "PATCH /api/v1/equipment/{id}/status?status={status}")
			return
		}
	}

	msg := fmt.Sprintf("equipment with id: %v updated status to %v", i, status)

	helpers.JsonResponseSuccess(w, http.StatusOK, msg)
}

// CreateEquipment create equipment
//
//	@Summary		create equipment
//	@Description	create equipment in the database
//	@Tags			equipment
//	@Accept			json
//	@Produce		json
//	@Param			sn				query		string	true	"serial number"
//	@Param			manufacturer	query		int		true	"manufacturer id"	minimum(1)
//	@Param			device			query		int		true	"device id"			minimum(1)
//	@Success		200				{object}	models.JsonResponse{MSG=models.Equipment}
//	@Failure		400				{object}	models.JsonResponse
//	@Failure		500				{object}	models.JsonResponse
//	@Router			/equipment [post]
func (h *EquipmentHandler) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	sn := r.FormValue("sn")
	if sn == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	var req models.JsonResponse
	resp, err := http.Get("http://localhost:8081/api/v1/equipment/sn?sn=" + sn)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/equipment/sn/"+sn, "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	if req.Status == "ERROR" && req.Message != "equipment does not exist in database" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	} else if req.Message != "equipment does not exist in database" {
		log.Println("req.Message: ", req.Message)
		helpers.JsonResponseError(w, http.StatusBadRequest, "equipment already exists in database", "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	did := r.FormValue("device")
	if did == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "POST /api/v1/equipment?sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	d, err := strconv.Atoi(did)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id is not a number", "POST /api/v1/equipment?sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/device/" + did)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1//"+did, "POST /api/v1/equipment?sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/", "POST /api/v1/equipment?sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/device/{id}")
		return
	}

	mid := r.FormValue("manufacturer")
	if mid == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing manufacturer id", "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	m, err := strconv.Atoi(mid)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	resp, err = http.Get("http://localhost:8081/api/v1/manufacturer/" + mid)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer/"+mid, "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/manufacturer")
		return
	}

	err = q.CreateEquipment(r.Context(), sqlc.CreateEquipmentParams{SerialNumber: sn, DeviceTypeID: int32(d), ManufacturerID: int32(m)})
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to create equipment in database", "POST /api/v1/equipment?sn={sn}&_id={device_id}&manufacturer_id={manufacturer_id}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, "equipment created")
}
