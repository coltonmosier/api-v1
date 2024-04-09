package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/coltonmosier/api-v1/internal/helpers"
	"github.com/coltonmosier/api-v1/internal/models"
)

type ManufactuerHandler struct{}

var manufacturers = []models.Manufacturer{
	{
		ID:     1,
		Name:   "watch",
		Status: "active",
	},
	{
		ID:     2,
		Name:   "smart watch",
		Status: "active",
	},
}

func (h *ManufactuerHandler) GetManufacturers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	out := models.JsonResponse{
		Status:  "ok",
		Message: manufacturers,
		Action:  "none",
	}

	output, err := json.Marshal(out)
	if err != nil {
		log.Println("Error marshalling JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(output)
}

// TODO: Handle unexpected input. i.e. can't be a number must be a string.
// if error, action is GET /api/v1/manufacturer
func (h *ManufactuerHandler) GetManufacturerByName(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	out := models.JsonResponse{
		Status:  "ok",
		Message: "manufacturer info based on name",
		Action:  "none",
	}

	output, err := json.Marshal(out)
	if err != nil {
		log.Println("Error marshalling JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(output)
}

func (h *ManufactuerHandler) GetManufacturerByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	out := models.JsonResponse{
		Status:  "ok",
		Message: "manufacturer info based on ID",
		Action:  "none",
	}

	output, err := json.Marshal(out)
	if err != nil {
		log.Println("Error marshalling JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(output)
}

func (h *ManufactuerHandler) UpdateManufacturer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing id", "GET /api/v1/manufacturer")
		return
	}

	name := r.PathValue("name")
	// name does appear in request
	if name != "" {
		// do the updating to the name
	}

	status := r.PathValue("status")
	//name does appear in request
	if status != "" {
		// do the updated to the status
	}

	//if here then we check if this manufacturer exists in DB
	// if it does, we update based on information provied

	helpers.JsonResponseSuccess(w, http.StatusOK, "updated data "+id+name)
}

func (h *ManufactuerHandler) CreateManufacturer(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		helpers.JsonResponseError(w, http.StatusBadRequest, "missing name", "retry")
	}
	// check if name already exists within db

	// NOTE: ALL OKAY!
	helpers.JsonResponseSuccess(w, http.StatusOK, "manufacturer created with name of "+name)
}

// TODO: do we need to delete these?
func (h *ManufactuerHandler) DeleteManufacturer(w http.ResponseWriter, r *http.Request) {
}
