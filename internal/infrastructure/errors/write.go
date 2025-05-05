package apierrors

import (
	"encoding/json"
	"net/http"
)


func WriteHTTPError(w http.ResponseWriter, err any) {
	w.Header().Set("Content-Type", "application/json")

	switch err := err.(type) {
	case *APIError:
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   err.Message,
		})
	case error:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   err.Error(),
		})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "internal server error",
			"details": "An unexpected error occurred",
		})
	}
}