package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coltonmosier/api-v1/internal/models"
)

func JsonResponseSuccess(w http.ResponseWriter, status int, message string) {
	out := models.JsonResponse{
		Status:  "ok",
		Message: message,
		Action:  "none",
	}
	output, err := json.Marshal(out)
	if err != nil {
		log.Println("Error marshalling JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(output)
}

func JsonResponseError(w http.ResponseWriter, status int, message string, action string) {
	out := models.JsonResponse{
		Status:  "ERROR",
		Message: message,
		Action:  action,
	}
	output, err := json.Marshal(out)
	if err != nil {
		log.Println("Error marshalling JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(output)
}
