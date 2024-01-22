package models

// Estruta de la tabla visualizations de la base de datos
type Visualizations struct {
	MovieID int `json:"movie_id"`
	Views   int `json:"views"`
}

// Estructura para deserializar detalles de la película
type MovieDetails struct {
	Title  string `json:"title"`
	ImdbID string `json:"imdb_id"`

	Overview      string `json:"overview"`
	PosterPath    string `json:"poster_path"`
	FullPosterUrl string `json:"poster_url"`
	ReleaseDate   string `json:"release_date"`
	Genres        []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Comments CommentResponse
}

// Respuesta estructurada
type ResponseData struct {
	Success  bool             `json:"success"`
	MovieID  int              `json:"movie_id"`
	Views    int              `json:"views"`
	Details  MovieDetails     `json:"details"`
	Comments []CommentDetails `json:"comments"`
}
