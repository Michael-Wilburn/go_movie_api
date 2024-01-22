package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Michael-Wilburn/go_movie_api/models"
	"github.com/Michael-Wilburn/go_movie_api/services"
	"github.com/Michael-Wilburn/go_movie_api/utils"
)

type AuthController struct {
	DB *sql.DB
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{DB: db}
}

func (au *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Verificar credenciales y obtener el usuario
	user, err := services.AuthUserCredentials(au.DB, credentials)
	if err != nil {
		// En caso de error, responder con un mensaje de error JSON
		errorResponse := models.ErrorLogin{
			Success: false,
			Message: "Invalid credentials",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Generar token JWT
	tokenString, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
		return
	}

	// Devolver el token JWT en la respuesta
	response := map[string]interface{}{
		"success": true,
		"token":   tokenString,
		"user":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
