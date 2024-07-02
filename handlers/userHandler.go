package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/KOLAAAAA1/go-personal-project/models"
	"github.com/KOLAAAAA1/go-personal-project/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Collection *mongo.Collection
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

    var existingUser models.User
    err := h.Collection.FindOne(ctx, bson.M{"username": newUser.Username}).Decode(&existingUser)
    if err == nil {
        http.Error(w, "Username already exists", http.StatusBadRequest)
        return
    } else if err != mongo.ErrNoDocuments {
        http.Error(w, "Failed to check existing user", http.StatusInternalServerError)
        return
    }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)
	newUser.ID = primitive.NewObjectID()
    newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	success := false

	_, err = h.Collection.InsertOne(ctx, newUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	success = true

	response := struct {
		Success bool `json:"success"`
	}{
		Success: success,
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.Collection.FindOne(ctx, bson.M{"username": loginRequest.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	access_token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refresh_token, err := utils.GenerateJWT(user.Username, 60)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := struct {
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken: access_token,
		RefreshToken: refresh_token,
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

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := h.Collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			http.Error(w, "Failed to decode user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		users = []models.User{}
	}

	response := struct {
		Users []models.User `json:"users"`
	}{
		Users: users,
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
