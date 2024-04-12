package handlers

import (
	"database/sql"
	"encoding/json"
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

func (h *EquipmentHandler) GetEquipments(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment?limit={limit}&offset={offset}")
		return
	}

    limit := r.FormValue("limit")
    offset := r.FormValue("offset")

	if limit == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing limit", "GET /api/v1/equipment?limit={limit}&offset={offset}")
		return
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "limit is not a number", "GET /api/v1/equipment?limit={limit}&offset={offset}")
		return
	}

	if offset == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing offset", "GET /api/v1/equipment?limit={limit}&offset={offset}")
		return
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "offset is not a number", "GET /api/v1/equipment?limit={limit}&offset={offset}")
		return
	}

	d, err := q.GetAllEquipment(r.Context(), sqlc.GetAllEquipmentParams{Limit: int32(l), Offset: int32(o)})
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment", "GET /api/v1/equipment?limit={limit}&offset={offset}")
		return
	}

	var e []models.Equipment
	for _, v := range d {
		e = append(e, models.Equipment{
			AutoID:         v.AutoID,
			DeviceTypeID:   v.DeviceTypeID,
			ManufacturerID: v.ManufacturerID,
			SerialNumber:   v.SerialNumber,
		})
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

func (h *EquipmentHandler) GetEquipmentBySN(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/sn/{sn}")
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
		helpers.JsonResponseError(w, http.StatusNotFound, "equipment does not exist in database", "GET /api/v1/equipment?limit={limit}&offset={offset}")
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
        helpers.JsonResponseError(w, http.StatusNotFound, "equipment id does not exist", "GET /api/v1/equipment/id?id={id}")
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
    }

    helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

func (h *EquipmentHandler) GetEquipmentLikeSN(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/sn-like/{sn}?limit={limit}&offset={offset}")
		return
	}

	sn := r.PathValue("sn")
	if sn == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "GET /api/v1/equipment/sn-like/{sn}?limit={limit}&offset={offset}")
		return
	}

	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	if limit == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing limit", "GET /api/v1/equipment/sn-like/{sn}?limit={limit}&offset={offset}")
		return
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "limit is not a number", "GET /api/v1/equipment/sn-like/{sn}?limit={limit}&offset={offset}")
		return
	}

	if offset == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing offset", "GET /api/v1/equipment/sn-like/{sn}?limit={limit}&offset={offset}")
		return
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "offset is not a number", "GET /api/v1/equipment/sn-like/{sn}?limit={limit}&offset={offset}")
		return
	}
	sn = "%" + sn + "%"

	d, err := q.GetEquipmentLikeSerialNumber(r.Context(), sqlc.GetEquipmentLikeSerialNumberParams{SerialNumber: sn, Limit: int32(l), Offset: int32(o)})
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment "+err.Error(), "GET /api/v1/equipment/sn-like/{sn}?limit={limit}&offset={offset}")
		return
	}

	var e []models.Equipment
	for _, v := range d {
		e = append(e, models.Equipment{
			AutoID:         v.AutoID,
			DeviceTypeID:   v.DeviceTypeID,
			ManufacturerID: v.ManufacturerID,
			SerialNumber:   v.SerialNumber,
		})
	}
	if e == nil {
		helpers.JsonResponseError(w, http.StatusNotFound, "no equipment found", "GET /api/v1/equipment/sn-like/{sn}?limit={limit}&offset={offset}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}
func (h *EquipmentHandler) GetEquipmentByManufacturerID(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/manufacturer/{id}?limit={limit}&offset={offset}")
		return
	}

	manufacturerID := r.PathValue("id")
	if manufacturerID == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing manufacturer id", "GET /api/v1/equipment/manufacturer/{id}?limit={limit}&offset={offset}")
		return
	}
	id, err := strconv.Atoi(manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "GET /api/v1/equipment/manufacturer/{id}?limit={limit}&offset={offset}")
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

	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	if limit == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "limit missing", "GET /api/v1/equipment/manufacturer/{id}?limit={limit}&offset={offset}")
		return
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "limit is not a number", "GET /api/v1/equipment/manufacturer/{id}?limit={limit}&offset={offset}")
		return
	}

	if offset == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing offset", "GET /api/v1/equipment/manufacturer/{id}?limit={limit}&offset={offset}")
		return
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "offset is not a number", "GET /api/v1/equipment/manufacturer/{id}?limit={limit}&offset={offset}")
		return
	}

	d, err := q.GetEquipmentByManufacturer(r.Context(), sqlc.GetEquipmentByManufacturerParams{ManufacturerID: int32(id), Limit: int32(l), Offset: int32(o)})
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment", "GET /api/v1/equipment/manufacturer/{id}?limit={limit}&offset={offset}")
		return
	}

	var e []models.Equipment
	for _, v := range d {
		e = append(e, models.Equipment{
			AutoID:         v.AutoID,
			DeviceTypeID:   v.DeviceTypeID,
			ManufacturerID: v.ManufacturerID,
			SerialNumber:   v.SerialNumber,
		})
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

func (h *EquipmentHandler) GetEquipmentByDeviceID(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/device/{id}?limit={limit}&offset={offset}")
		return
	}

	deviceID := r.PathValue("id")
	if deviceID == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "GET /api/v1/equipment/device/{id}?limit={limit}&offset={offset}")
		return
	}
	id, err := strconv.Atoi(deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id is not a number", "GET /api/v1/equipment/device/{id}?limit={limit}&offset={offset}")
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

	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	if limit == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing limit", "GET /api/v1/equipment/device/{id}?limit={limit}&offset={offset}")
		return
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "limit is not a number", "GET /api/v1/equipment/device/{id}?limit={limit}&offset={offset}")
		return
	}

	if offset == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing offset", "GET /api/v1/equipment/device/{id}?limit={limit}&offset={offset}")
		return
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "offset is not a number", "GET /api/v1/equipment/device/{id}?limit={limit}&offset={offset}")
		return
	}

	d, err := q.GetEquipmentByDeviceType(r.Context(), sqlc.GetEquipmentByDeviceTypeParams{DeviceTypeID: int32(id), Limit: int32(l), Offset: int32(o)})
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment"+err.Error(), "GET /api/v1/equipment/device/{id}?limit={limit}&offset={offset}")
		return
	}

	var e []models.Equipment
	for _, v := range d {
		e = append(e, models.Equipment{
			AutoID:         v.AutoID,
			DeviceTypeID:   v.DeviceTypeID,
			ManufacturerID: v.ManufacturerID,
			SerialNumber:   v.SerialNumber,
		})
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

func (h *EquipmentHandler) GetEquipmentByDeviceIDAndManufacturerID(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}?limit={limit}&offset={offset}")
		return
	}

	deviceID := r.PathValue("device_id")
	if deviceID == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}?limit={limit}&offset={offset}")
		return
	}
	did, err := strconv.Atoi(deviceID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "device id is not a number", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}?limit={limit}&offset={offset}")
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
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing manufacturer id", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}?limit={limit}&offset={offset}")
		return
	}
	mid, err := strconv.Atoi(manufacturerID)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}?limit={limit}&offset={offset}")
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
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	if limit == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing limit", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}?limit={limit}&offset={offset}")
		return
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "limit is not a number", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}?limit={limit}&offset={offset}")
		return
	}

	if offset == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing offset", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}?limit={limit}&offset={offset}")
		return
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "offset is not a number", "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}?limit={limit}&offset={offset}")
		return
	}

	d, err := q.GetEquipmentByDeviceTypeAndManufacturer(r.Context(), sqlc.GetEquipmentByDeviceTypeAndManufacturerParams{DeviceTypeID: int32(did), ManufacturerID: int32(mid), Limit: int32(l), Offset: int32(o)})
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "failed to query database for equipment"+err.Error(), "GET /api/v1/equipment/device/{device_id}/manufacturer/{manufacturer_id}?limit={limit}&offset={offset}")
		return
	}

	var e []models.Equipment
	for _, v := range d {
		e = append(e, models.Equipment{
			AutoID:         v.AutoID,
			DeviceTypeID:   v.DeviceTypeID,
			ManufacturerID: v.ManufacturerID,
			SerialNumber:   v.SerialNumber,
		})
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, e)
}

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
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/equipment/sn-like/{sn}/limit/{limit}/offset/{offset}")
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
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/equipment/sn-like/{sn}/limit/{limit}/offset/{offset}")
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
    if err == sql.ErrNoRows{
        helpers.JsonResponseError(w, http.StatusBadRequest, "equipment does not exists", "none")
        return
    } else if err != nil{
        helpers.JsonResponseError(w, http.StatusBadRequest, "something went wrong with query" + err.Error(), "none")
        return
    }

    out := models.Equipment{
        AutoID: d.AutoID,
        DeviceTypeID: d.DeviceTypeID,
        ManufacturerID: d.ManufacturerID,
        SerialNumber: d.SerialNumber,
    }

    helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

func (h *EquipmentHandler) GetEquipmentByManufacturerIDAndDeviceIDAndSN(w http.ResponseWriter, r *http.Request){
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
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/equipment/sn-like/{sn}/limit/{limit}/offset/{offset}")
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
        SerialNumber: sn,
        DeviceTypeID: int32(did),
        ManufacturerID: int32(mid),
    })

    if err == sql.ErrNoRows{
        helpers.JsonResponseError(w, http.StatusBadRequest, "equipment does not exists", "GET /api/v1/equipment/limit/{limit}/offset/{offset}")
        return
    } else if err != nil {
        helpers.JsonResponseError(w, http.StatusBadRequest, "something went wrong with query "+err.Error(), "GET /api/v1/equipment/limit/{limit}/offset/{offset}")
        return
    }

    helpers.JsonResponseSuccess(w, http.StatusOK, d)
}
func (h *EquipmentHandler) UpdateSerialNumber (w http.ResponseWriter, r *http.Request){
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "PATCH /api/v1/serial-number?id={id}&sn={sn}")
		return
	}

    id := r.FormValue("id") 
    if id == "" {
        helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "PATCH /api/v1/serial-numberid={id}&sn={sn}")
        return
    }
    i, err := strconv.Atoi(id)
    if err != nil {
        helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "PATCH /api/v1/serial-number?id={id}&sn={sn}")
        return
    }

    resp, err := http.Get("http://localhost:8081/api/v1/equipment/id?id=" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/serial-number?id="+id, "PATCH /api/v1/serial-number?id={id}&sn={sn}")
		return
	}
	defer resp.Body.Close()

    var req models.JsonResponse
    err = json.NewDecoder(resp.Body).Decode(&req)
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "PATCH /api/v1/serial-number?id={id}&sn={sn}")
        return
    }
    if req.Status == "ERROR" {
        helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "PATCH /api/v1/serial-number?id={id}&sn={sn}")
        return
    }

    sn := r.FormValue("sn")
    if sn == "" {
        helpers.JsonResponseError(w, http.StatusBadRequest, "missing serial number", "PATCH /api/v1/serial-number?id={id}&sn={sn}")
        return
    }

    resp, err = http.Get("http://localhost:8081/api/v1/equipment/sn/" + sn)
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/equipment/sn/"+sn, "PATCH /api/v1/serial-number?id={id}&sn={sn}")
        return
    }
    defer resp.Body.Close()

    err = json.NewDecoder(resp.Body).Decode(&req)
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/equipment", "PATCH /api/v1/serial-number?id={id}&sn={sn}")
        return
    }

    if req.Status == "ERROR" {
        helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "PATCH /api/v1/serial-number?id={id}&sn={sn}")
        return
    }

    err = q.UpdateSerialNumber(r.Context(), sqlc.UpdateSerialNumberParams{AutoID: int32(i), SerialNumber: sn})
    if err != nil {
        helpers.JsonResponseError(w, http.StatusBadRequest, "failed to update serial number in database", "PATCH /api/v1/serial-number?id={id}&sn={sn}")
        return
    }

    helpers.JsonResponseSuccess(w, http.StatusOK, "serial number updated")
}

func (h *EquipmentHandler) UpdateEquipment(w http.ResponseWriter, r *http.Request){
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
    i, err := strconv.Atoi(id)
    if err != nil {
        helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
        return
    }

    resp, err := http.Get("http://localhost:8081/api/v1/equipment/id?id=" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/serial-number?id="+id, "PATCH /api/v1/equipment?id={id}&sn={sn}&device_id={device_id}&manufacturer_id={manufacturer_id}")
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
    
    if req.Status == "ERROR" && req.Message != "equipment does not exist in database"{
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
