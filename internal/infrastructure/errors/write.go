package apierrors

import (
	"encoding/json"
	"log"
	"net/http"
)


func WriteHTTPError(w http.ResponseWriter, err any) {
    w.Header().Set("Content-Type", "application/json")

    switch err := err.(type) {
    case *APIError:
        if err.Code == 0 {
            err.Code = http.StatusInternalServerError
        }
        w.WriteHeader(err.Code)
        json.NewEncoder(w).Encode(map[string]any{
            "error": err.Message,
        })

    case error:
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": err.Error(),
        })

    default:
        log.Printf("Internal Server Error: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]any{
            "error":   "internal server error",
            "details": "An unexpected error occurred",
        })
    }
}