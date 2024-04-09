package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/coltonmosier/api-v1/internal/database"
	"github.com/coltonmosier/api-v1/internal/helpers"
	"github.com/coltonmosier/api-v1/internal/models"
	"github.com/coltonmosier/api-v1/internal/sqlc"
	//"github.com/coltonmosier/api-v1/internal/sqlc"
)

type ManufactuerHandler struct{}

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
			Status: v.Status,
		})
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

// TODO: Handle unexpected input. i.e. can't be a number must be a string.
// if error, action is GET /api/v1/manufacturer
func (h *ManufactuerHandler) GetManufacturerByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "GET /api/v1/manufacturer")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "GET /api/v1/manufacturer")
		return
	}

	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/manufacturer")
		return
	}
	d, err := q.GetManufacturerById(r.Context(), int32(i))
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id does not exists in database", "GET /api/v1/manufacturer")
		return
	}
	out := models.Manufacturer{
		ID:     int32(i),
		Name:   d.Name,
		Status: d.Status,
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

func (h *ManufactuerHandler) UpdateManufacturerName(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/manufacturer")
		return
	}
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "GET /api/v1/manufacturer")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "GET /api/v1/manufacturer")
		return
	}
	// NOTE: Must make a request to /api/v1/manufacturer to see if id exists...
	resp, err := http.Get("http://localhost:8081/api/v1/manufacturer/" + id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer/"+id, "PATH /api/v1/{id}/name/{newName}")
		return
	}
	defer resp.Body.Close()

	var req models.JsonResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "PATH /api/v1/{id}/name/{newName}")
		return
	}

    if req.Status == "ERROR" {
        helpers.JsonResponseError(w, http.StatusBadRequest, "id does not exist", "GET /api/v1/manufacturer")
        return
    }
        
	name := r.PathValue("name")
	if name == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing name", "PATH /api/v1/{id}/name/{newName}")
		return
	}

    err = q.UpdateManufacturer(r.Context(), sqlc.UpdateManufacturerParams{ID: int32(i), Name: name})
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update manufacturer name", "PATCH /api/v1/manufacturer/{id}/name/{name}")
        return
    }

    helpers.JsonResponseSuccess(w, http.StatusOK, "manufacturer updated with name of "+name)
}

func (h *ManufactuerHandler) CreateManufacturer(w http.ResponseWriter, r *http.Request) {
	q, err := database.InitEquipmentDatabase()
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/manufacturer")
		return
	}
	name := r.PathValue("name")
	if name == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing name", "POST /api/v1/manufacturer/{name}")
	}
	// check if name already exists within db
	resp, err := http.Get("http://localhost:8081/api/v1/manufacturer")
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to get response from /api/v1/manufacturer", "PATCH /api/v1/manufacturer/{name}")
		return
	}
	defer resp.Body.Close()
	var req []models.Manufacturer
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to decode response from /api/v1/manufacturer", "PATCH /api/v1/manufacturer/{name}")
		return
	}

	for _, v := range req {
		if v.Name == name {
            // NOTE: I have the ID and name of the manufacturer and can update the name
			err = q.UpdateManufacturer(r.Context(), sqlc.UpdateManufacturerParams{ID: v.ID, Name: name})
            if err != nil {
                helpers.JsonResponseError(w, http.StatusInternalServerError, "failed to update manufacturer", "PATCH /api/v1/manufacturer/{name}")
                return
            }
            helpers.JsonResponseSuccess(w, http.StatusOK, "manufacturer updated with name of "+name)
            return
		}
	}

	err = q.CreateManufacturer(r.Context(), name)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not create", "POST /api/v1/manufacturer/{name}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, "manufacturer created with name of "+name)
}
