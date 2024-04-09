package handlers

import (
	"net/http"
	"strconv"

	"github.com/coltonmosier/api-v1/internal/database"
	"github.com/coltonmosier/api-v1/internal/helpers"
	"github.com/coltonmosier/api-v1/internal/models"
	"github.com/coltonmosier/api-v1/internal/sqlc"
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
        helpers.JsonResponseError(w, http.StatusInternalServerError, "something went wrong", "GET /api/v1/manufacturer")
        return
    }
    for _, v := range d {
        out = append(out, models.Manufacturer{
            ID:     v.ID,
            Name:   v.Name,
            Status: "active",
        })
    }
    
    helpers.JsonResponseSuccess(w, http.StatusOK, out)
}

// TODO: Handle unexpected input. i.e. can't be a number must be a string.
// if error, action is GET /api/v1/manufacturer
func (h *ManufactuerHandler) GetManufacturerWithValue(w http.ResponseWriter, r *http.Request) {
}

func (h *ManufactuerHandler) UpdateManufacturer(w http.ResponseWriter, r *http.Request) {
    q, err := database.InitEquipmentDatabase()
    if err != nil {
        helpers.JsonResponseError(w, http.StatusInternalServerError, "could not connect to database", "GET /api/v1/manufacturer")
        return
    }
	id := r.PathValue("id")
	value := r.PathValue("value")

	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "GET /api/v1/manufacturer")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, "id is not a number", "GET /api/v1/manufacturer")
		return
	}

	d, err := q.GetManufacturerById(r.Context(), int32(i))
	if err != nil {
		helpers.JsonResponseError(w, http.StatusBadRequest, value+" does not exists within database", "GET /api/v1/manufacturer")
		return
	}

	if value == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing update value", "PATCH /api/v1/manufacturer/{value}")
		return
	}

	if value == "active" || value == "inactive" {
		err = q.UpdateManufacturerStatus(r.Context(),
			sqlc.UpdateManufacturerStatusParams{
				ID:     int32(i),
				Status: sqlc.ManufacturerStatusActive,
			},
		)
	} else {
		if value == d.Name {
			helpers.JsonResponseError(w, http.StatusBadRequest, "that name already exists", "GET /api/v1/manufacturer")
			return
		}
		err = q.UpdateManufacturer(r.Context(),
			sqlc.UpdateManufacturerParams{
				Name: value,
				ID:   int32(i),
			},
		)
		if err != nil {
			helpers.JsonResponseError(w, http.StatusInternalServerError, "error with sql", "none")
			return
		}
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, "updated data "+id+value)
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
	d, _ := q.GetManufacturerByName(r.Context(), name)
	if d.Name == name {
		helpers.JsonResponseError(w, http.StatusBadRequest, "name already exists", "GET /api/v1/manufacturer")
		return
	}

	err = q.CreateManufacturer(r.Context(), name)
	if err != nil {
		helpers.JsonResponseError(w, http.StatusInternalServerError, "could not create", "POST /api/v1/manufacturer/{name}")
		return
	}

	helpers.JsonResponseSuccess(w, http.StatusOK, "manufacturer created with name of "+name)
}

// TODO: do we need to delete these?
func (h *ManufactuerHandler) DeleteManufacturer(w http.ResponseWriter, r *http.Request) {
}
