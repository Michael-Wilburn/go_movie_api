package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Michael-Wilburn/go_movie_api/models"
)

// GetMovieDetails realiza una consulta a la API de TMDb para obtener detalles de una película por su ID.
func GetMovieDetails(movieID string) ([]byte, int) {
	// Obtener el TOKEN de TMDb desde las variables de entorno
	token := os.Getenv("TMDB_TOKEN")
	if token == "" {
		return nil, http.StatusInternalServerError
	}

	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s", movieID)

	// Crear una solicitud HTTP a la url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	// Agregar el token al encabezado
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", token)

	// Realizar la solicitud a la API de TMDb
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return body, res.StatusCode
}

func ParseMovieDetails(details []byte) models.MovieDetails {
	var movieDetails models.MovieDetails
	json.Unmarshal(details, &movieDetails)
	movieDetails.FullPosterUrl = "https://image.tmdb.org/t/p/w600_and_h900_bestv2" + movieDetails.PosterPath
	return movieDetails
}
