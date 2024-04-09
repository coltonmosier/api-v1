package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coltonmosier/api-v1/internal/models"
)

// JsonResponseSuccess returns nothing
// JsonResponseSuccess takes in a http.ResponseWriter, status int, message string
// uses the paramaters to construct a success response for handlers
func JsonResponseSuccess(w http.ResponseWriter, status int, message models.Message) {
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
	return
}

// JsonResponseError returns nothing
// JsonResponseError takes in a http.ResponseWriter, status int, message string, action string
// uses the paramaters to construct an error response for handlers
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
