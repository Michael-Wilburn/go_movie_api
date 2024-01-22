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
	"github.com/gorilla/mux"
)

type MovieController struct {
	DB *sql.DB
}

func NewMovieController(db *sql.DB) *MovieController {
	return &MovieController{DB: db}
}

func (mc *MovieController) ViewMovieDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["movie_id"]

	// Obtener detalles de la película utilizando la API externa
	movieDetails, statusCode := helpers.GetMovieDetails(movieID)
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

	visualizations, err := services.AddViews(mc.DB, movieID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error registering visualization: %v", err), http.StatusInternalServerError)
		return
	}

	comments, err := services.GetMovieComments(mc.DB, visualizations.MovieID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting movie comments: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.ResponseData{
		Success:  true,
		MovieID:  visualizations.MovieID,
		Views:    visualizations.Views,
		Details:  helpers.ParseMovieDetails(movieDetails),
		Comments: comments,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (mc *MovieController) MostViewMovies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n := vars["n"]
	nInt, err := strconv.Atoi(n)

	if err != nil {
		http.Error(w, "Invalid parameter N", http.StatusBadRequest)
		return
	}

	mostViewedMovies, err := services.GetMostViewedMovies(mc.DB, nInt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting most viewed movies: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"movies":  mostViewedMovies,
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
