package main

import (
	"log"
	"net/http"
	"time"

	"github.com/KOLAAAAA1/go-personal-project/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Connect to MongoDB
	ConnectToMongoDB("mongodb://localhost:27017")

	// Initialize handlers
	itemHandler := &handlers.ItemHandler{
		Collection: GetCollection("items"),
	}
	userHandler := &handlers.UserHandler{
		Collection: GetCollection("users"),
	}

	// Create a new router
	router := mux.NewRouter()

	// Define your API routes
	router.HandleFunc("/api/v1/items", itemHandler.GetItems).Methods("GET")
	router.HandleFunc("/api/v1/item/{id}", itemHandler.GetItem).Methods("GET")
	router.HandleFunc("/api/v1/item/{id}", itemHandler.UpdateItem).Methods("PATCH")
	router.HandleFunc("/api/v1/item/{id}", itemHandler.DeleteItem).Methods("DELETE")
	router.HandleFunc("/api/v1/signup", userHandler.SignUp).Methods("POST")
	router.HandleFunc("/api/v1/login", userHandler.Login).Methods("POST")

	// Start the server
	server := &http.Server{
		Addr:         ":3001",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on port 3001...")
	log.Fatal(server.ListenAndServe())
}
