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

type ManufactuerHandler struct{}

// GetManufacturers Getting all manufacturers
//	@Summary		get all manufacturers
//	@Description	get all manufacturers from the database
//	@Tags			manufacturer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.JsonResponse{MSG=models.Manufacturer}
//	@Failure		500	{object}	models.JsonResponse
//	@Router			/manufacturer [get]
func (h *ManufactuerHandler) GetManufacturers(w http.ResponseWriter, r *http.Request) {
	var out []models.Manufacturer
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/manufacturer")
		return
	}
	d, err := q.GetManufacturersActive(r.Context())
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "something went wrong with query"+err.Error(), "GET /api/v1/manufacturer")
		return
	}
	for _, v := range d {
		out = append(out, models.Manufacturer{
			ID:     v.ID,
			Name:   v.Name,
			Status: string(v.Status),
		})
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

// GetManufacturerByID Get a manufacturer by ID
//	@Summary		get a manufacturer by ID
//	@Description	get a manufacturer by ID from the database
//	@Tags			manufacturer
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Manufacturer ID" minimum(1)
//	@Success		200	{object}	models.JsonResponse{MSG=models.Manufacturer}
//	@Failure		400	{object}	models.JsonResponse
//	@Failure		500	{object}	models.JsonResponse
//	@Router			/manufacturer/{id} [get]
func (h *ManufactuerHandler) GetManufacturerByID(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/manufacturer")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing manufacturer id", "GET /api/v1/manufacturer/{id}")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "GET /api/v1/manufacturer")
		return
	}

	d, err := q.GetManufacturerById(r.Context(), int32(i))
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id does not exists in database", "GET /api/v1/manufacturer")
		return
	}
	out := models.Manufacturer{
		ID:     int32(i),
		Name:   d.Name,
		Status: string(d.Status),
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

// UpdateManufacturer Update a manufacturer by ID
//	@Summary		update a manufacturer by ID
//	@Description	update a manufacturer by ID from the database
//	@Tags			manufacturer
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Manufacturer ID" minimum(1)
//	@Param			name	query		string	true	"Manufacturer Name"
//	@Param			status	query		string	true	"Manufacturer Status" Enums(active,inactive)
//	@Success		200		{object}	models.JsonResponse{MSG=models.Manufacturer}
//	@Failure		400		{object}	models.JsonResponse
//	@Failure		500		{object}	models.JsonResponse
//	@Router			/manufacturer/{id} [patch]
func (h *ManufactuerHandler) UpdateManufacturer(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
		return
	}
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing device id", "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer id is not a number", "GET /api/v1/manufacturer")
		return
	}

	resp, err := http.Get("http://localhost:8081/api/v1/manufacturer/" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer/"+id, "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
		return
	}

	if req.Status == "ERROR" {
		helpers.JsonResponseError(w, http.StatusBadRequest, req.Message, "GET /api/v1/manufacturer")
		return
	}

	name := r.FormValue("name")
	status := r.FormValue("status")
	if name == "" && status == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing name and/or status", "PATH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
		return
	} else if name != "" && status != "" {
		err = q.UpdateManufacturer(r.Context(), sqlc.UpdateManufacturerParams{ID: int32(i), Name: name})
		if err != nil {
			helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update manufacturer name", "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
			return
		}
		if status == string(sqlc.ManufacturerStatusActive) {
			err = q.UpdateManufacturerStatus(r.Context(), sqlc.UpdateManufacturerStatusParams{ID: int32(i), Status: sqlc.ManufacturerStatusActive})
			if err != nil {
				helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update manufacturer status"+err.Error(), "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
				return
			}
		} else if status == string(sqlc.ManufacturerStatusInactive) {
			err = q.UpdateManufacturerStatus(r.Context(), sqlc.UpdateManufacturerStatusParams{ID: int32(i), Status: sqlc.ManufacturerStatusInactive})
			if err != nil {
				helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update manufacturer status"+err.Error(), "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
				return
			}
		} else {
			helpers.JsonResponseError(w, http.StatusBadRequest, "status must be either 'active' or 'inactive'", "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
			return
		}
		helpers.JsonResponseSuccess(w, http.StatusOK, "manufacturer updated with name of "+name+" and status of "+status)
		return
	} else if name != ""  && status == ""{
		err = q.UpdateManufacturer(r.Context(), sqlc.UpdateManufacturerParams{ID: int32(i), Name: name})
		if err != nil {
			helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update manufacturer name", "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
			return
		}
		helpers.JsonResponseSuccess(w, http.StatusOK, "manufacturer updated with name of "+name)
		return
	} else if status != "" && name == ""{
		if status == string(sqlc.ManufacturerStatusActive) {
			err = q.UpdateManufacturerStatus(r.Context(), sqlc.UpdateManufacturerStatusParams{ID: int32(i), Status: sqlc.ManufacturerStatusActive})
			if err != nil {
				helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update manufacturer status"+err.Error(), "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
				return
			}
		} else if status == string(sqlc.ManufacturerStatusInactive) {
			err = q.UpdateManufacturerStatus(r.Context(), sqlc.UpdateManufacturerStatusParams{ID: int32(i), Status: sqlc.ManufacturerStatusInactive})
			if err != nil {
				helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update manufacturer status"+err.Error(), "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
				return
			}
		} else {
			helpers.JsonResponseError(w, http.StatusBadRequest, "status must be either 'active' or 'inactive'", "PATCH /api/v1/manufacturer/{id}?name={newName}&status={newStatus}")
			return
		}
		helpers.JsonResponseSuccess(w, http.StatusOK, "manufacturer updated with status of "+status)
		return
	}
}


//  CreateManufacturer Creating a manufacturer
//	@Summary		create manufacturer  
//	@Description	create manufacturer for the database
//	@Tags			manufacturer 
//	@x-order		4
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	true	"Manufacturer Name"
//	@Success		200		{object}	models.JsonResponse
//	@Failure		400		{object}	models.JsonResponse
//	@Failure		500		{object}	models.JsonResponse
//	@Router			/manufacturer [post]
func (h *ManufactuerHandler) CreateManufacturer(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/manufacturer")
		return
	}

	name := r.FormValue("name")
	if name == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing name", "POST /api/v1/manufacturer?name={newManufacturerName}")
		return
	}

    req, err := http.Get("http://localhost:8081/api/v1/manufacturer")
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer", "POST /api/v1/manufacturer?name={newName}")
        return
    }
    defer req.Body.Close()

    var resp models.JsonResponse
    err = json.NewDecoder(req.Body).Decode(&resp)
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "POST /api/v1/manufacturer?name={newName}")
        return
    }

    if resp.Status == "ERROR" {
        helpers.JsonResponseError(w, http.StatusBadRequest, resp.Message, "POST /api/v1/manufacturer?name={newName}")
        return
    }

    var manufacturer []models.Manufacturer 
    for _, v := range resp.Message.([]interface{}) {
        if m, ok := v.(map[string]interface{}); ok {
        // Extract values from the map and create a Manufacturer instance
        var d models.Manufacturer
        if id, ok := m["id"].(float64); ok {
            d.ID = int32(id)
        } else {
            // Handle error when id cannot be converted to float64
            helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to convert id to float64", "POST /api/v1/manufacturer?name={newName}")
            return
        }
        if name, ok := m["name"].(string); ok {
            d.Name = name
        } else {
            // Handle error when name cannot be converted to string
            helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to convert name to string", "POST /api/v1/manufacturer?name={newName}")
            return
        }
        if status, ok := m["status"].(string); ok {
            d.Status = status
        } else {
            // Handle error when status cannot be converted to string
            helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to convert status to string", "POST /api/v1/manufacturer?name={newName}")
            return
        }
        // Add the created Manufacturer instance to the manufacturer slice
        manufacturer = append(manufacturer, d)
    } else {
        // Type assertion failed for an element in the slice
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to assert map[string]interface{} type", "POST /api/v1/manufacturer?name={newName}")
        return
    }
    }

    for _, v := range manufacturer {
        if v.Name == name {
            helpers.JsonResponseError(w, http.StatusBadRequest, "manufacturer type already exists", "POST /api/v1/manufacturer?name={newName}")
            return
        }
    }
	err = q.CreateManufacturer(r.Context(), name)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to write new manufacturer to database"+err.Error(), "POST /api/v1/manufacturer/{newManufacturerName}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusCreated, "manufacturer created with name of "+name)
}
