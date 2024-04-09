package handlers

import (
	"net/http"

	"github.com/coltonmosier/api-v1/internal/helpers"
)



type EquipmentHandler struct{}

func (h *EquipmentHandler) BadEndpointHandler(w http.ResponseWriter, r *http.Request) {
    helpers.JsonResponseError(w, http.StatusNotFound, "endpoint not found", "none")
}
