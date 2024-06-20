package main

import (
	"encoding/json"
	"net/http"
)

// Define additional handler functions as needed


func userHandler(w http.ResponseWriter, r *http.Request) {
    // Handle GET and POST requests for /api/v1/items
    switch r.Method {
    case http.MethodGet:
        getUser(w, r)
    case http.MethodPost:
    case http.MethodPatch:
    case http.MethodPut:
    case http.MethodDelete:
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func getUser(w http.ResponseWriter, r *http.Request) {
    // Implement logic to get user details
    user := map[string]interface{}{
        "name": "John Doe",
        "age":  30,
    }

    // Convert user map to JSON
    jsonResponse, err := json.Marshal(user)
    if err != nil {
        http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
        return
    }

    // Set response headers and write JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}
