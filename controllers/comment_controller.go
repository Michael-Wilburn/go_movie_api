package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Michael-Wilburn/go_movie_api/helpers"
	"github.com/Michael-Wilburn/go_movie_api/models"
	"github.com/Michael-Wilburn/go_movie_api/services"
	"github.com/Michael-Wilburn/go_movie_api/utils"
	"github.com/gorilla/mux"
)

type CommentController struct {
	DB *sql.DB
}

func NewCommentController(db *sql.DB) *CommentController {
	return &CommentController{DB: db}
}

func (cc *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.CommentRequest
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Error decoding request body1", http.StatusBadRequest)
		return
	}

	movie_id := strconv.Itoa(int(comment.MovieID))

	// Obtener detalles de la película utilizando la API externa
	_, statusCode := helpers.GetMovieDetails(movie_id)
	if statusCode == http.StatusNotFound {
		response := map[string]interface{}{
			"success":        false,
			"status_code":    http.StatusNotFound,
			"status_message": "The resource you requested could not be found.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	} else if statusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Error from external API: %d", statusCode), http.StatusInternalServerError)
		return
	}

	if movie_id == "" || comment.CommentText == "" {
		errorResponse := models.ErrorLogin{
			Success: false,
			Message: "All fields movie_id and comment_text must be provided for adding a comment",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

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

	createdComment, err := services.AddComment(cc.DB, comment, userIDFromToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Comment: %v", err), http.StatusInternalServerError)
		return
	}

	CommentResponse := models.CommentResponse{
		ID:          createdComment.ID,
		UserID:      createdComment.UserID,
		MovieID:     createdComment.MovieID,
		CommentText: createdComment.CommentText,
		CreatedDate: createdComment.CreatedDate,
		UpdatedDate: createdComment.UpdatedDate,
	}

	response := map[string]interface{}{
		"success": true,
		"comment": CommentResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
func (cc *CommentController) DeleteCommentById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	commentID := vars["id"]
	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment id", http.StatusUnauthorized)
		return
	}

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

	err = services.DeleteCommentByID(cc.DB, commentIDInt, userIDFromToken)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Comment not found or does not belong to the user", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error deleting comment: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Comment deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (cc *CommentController) UpdateCommentById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["id"]
	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment id", http.StatusUnauthorized)
		return
	}

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

	var updatedComment models.CommentRequest
	err = json.NewDecoder(r.Body).Decode(&updatedComment)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	err = services.UpdateCommentByID(cc.DB, commentIDInt, userIDFromToken, updatedComment)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Comment not found or does not belong to the user", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error updating comment: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Comment updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
