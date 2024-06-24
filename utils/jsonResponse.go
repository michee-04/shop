package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse structure to hold response data
type JSONResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondWithJSON writes JSON response to the client
func RespondWithJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := JSONResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	respJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}
	w.Write(respJSON)
}
