package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Michael-Wilburn/go_movie_api/models"
	"github.com/Michael-Wilburn/go_movie_api/services"
	"github.com/Michael-Wilburn/go_movie_api/utils"
	"github.com/gorilla/mux"
)

type UserController struct {
	DB *sql.DB
}

func NewUserController(db *sql.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Email == "" || user.PasswordHash == "" {
		errorResponse := models.ErrorLogin{
			Success: false,
			Message: "All fields username, email and password must be provided for creating user",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	createdUser, err := services.AddUserDb(uc.DB, user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	userResponse := models.UserResponse{
		ID:          createdUser.ID,
		Username:    createdUser.Username,
		Email:       createdUser.Email,
		CreatedDate: createdUser.CreatedDate,
		UpdatedDate: createdUser.UpdatedDate,
	}

	response := map[string]interface{}{
		"success": true,
		"user":    userResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		http.Error(w, "'id' parameter is required", http.StatusBadRequest)
		return
	}

	user, err := services.GetUserByID(uc.DB, userID)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"user":    user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (uc *UserController) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	// Obtener el token del encabezado Authorization
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Bearer token not provided", http.StatusUnauthorized)
		return
	}

	// Extraer el ID de usuario del token
	userIDFromToken, err := utils.ExtractUserIDFromToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}

	// Obtener el ID del usuario de la URL
	vars := mux.Vars(r)
	userID := vars["id"]

	// Verificar si el usuario autenticado coincide con el usuario que intenta actualizar
	if userIDFromToken != userID {
		http.Error(w, "Without authorization to update this user", http.StatusForbidden)
		return
	}

	var updateUser models.User
	err = json.NewDecoder(r.Body).Decode(&updateUser)

	// Lógica para actualizar el usuario en la base de datos
	err = services.UpdateUser(uc.DB, userID, updateUser)

	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	// Responder con la actualización exitosa
	response := map[string]interface{}{
		"success": true,
		"message": "User updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
