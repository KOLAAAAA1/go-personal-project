package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
    // Create a new ServeMux instance (router)
    mux := http.NewServeMux()

    // Define your API routes
    mux.HandleFunc("/api/v1/hello", helloHandler)
    mux.HandleFunc("/api/v1/items", itemsHandler)
		mux.HandleFunc("/api/v1/user", userHandler)

    // Start the server
    log.Println("Server starting on port 3001")
    log.Fatal(http.ListenAndServe(":3001", mux))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    // Respond with a simple message
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "Hello, World!"}`))
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
    // Handle GET and POST requests for /api/v1/items
    switch r.Method {
    case http.MethodGet:
        getAllItems(w, r)
    case http.MethodPost:
        createItem(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func getAllItems(w http.ResponseWriter, r *http.Request) {
    // Simulate fetching items from a database or other source
    items := []string{"item1", "item2", "item3"}

    // Create a struct to hold the items with a "json" tag to format as "items"
    response := struct {
        Items []string `json:"items"`
    }{
        Items: items,
    }

    // Convert response struct to JSON
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
        return
    }

    // Set response headers and write JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func createItem(w http.ResponseWriter, r *http.Request) {
    // Implement your logic to create an item (dummy response)
    newItem := "newItem"

    // Convert newItem to JSON
    jsonResponse, err := json.Marshal(newItem)
    if err != nil {
        http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
        return
    }

    // Set response headers and write JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(jsonResponse)
}
