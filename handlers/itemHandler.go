package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/KOLAAAAA1/go-personal-project/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ItemHandler struct {
	Collection *mongo.Collection
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var newItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newItem.ID = primitive.NewObjectID()
	newItem.CreatedAt = time.Now()
	newItem.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := h.Collection.InsertOne(ctx, newItem)
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	// Retrieve the newly created item from the database
	createdItem := models.Item{}
	err = h.Collection.FindOne(ctx, bson.M{"_id": newItem.ID}).Decode(&createdItem)
	if err != nil {
		http.Error(w, "Failed to fetch created item", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(createdItem)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
    var items []models.Item

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := h.Collection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var item models.Item
        if err := cursor.Decode(&item); err != nil {
            http.Error(w, "Failed to decode item", http.StatusInternalServerError)
            return
        }
        items = append(items, item)
    }

    if len(items) == 0 {
        items = []models.Item{}
    }

    response := struct {
        Items []models.Item `json:"items"`
    }{
        Items: items,
    }

    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := h.Collection.FindOne(ctx, bson.M{"_id": objID})
	if result.Err() != nil {
		http.Error(w, "Failed to find item", http.StatusInternalServerError)
		return
	}

	updatedItem := models.Item{}
	err = result.Decode(&updatedItem)
	if err != nil {
		http.Error(w, "Failed to decode in item", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(updatedItem)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure updated_at is always set
	updates["updated_at"] = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	result := h.Collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, options.FindOneAndUpdate().SetReturnDocument(1))
	if result.Err() != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	updatedItem := models.Item{}
	err = result.Decode(&updatedItem)
	if err != nil {
		http.Error(w, "Failed to decode updated item", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(updatedItem)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.Collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}