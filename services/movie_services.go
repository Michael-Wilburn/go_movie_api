package services

import (
	"database/sql"

	"github.com/Michael-Wilburn/go_movie_api/models"
)

// AddViews incrementa el contador de visualizaciones para una película y devuelve los datos actualizados.
func AddViews(db *sql.DB, movieID string) (*models.Visualizations, error) {
	// Verificar si la fila ya existe
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM visualizations WHERE movie_id = $1", movieID).Scan(&count)
	if err != nil {
		return nil, err
	}

	// Si la fila no existe, insertarla con views inicializado a 1
	if count == 0 {
		_, err := db.Exec("INSERT INTO visualizations (movie_id, views) VALUES ($1, 1)", movieID)
		if err != nil {
			return nil, err
		}
	}

	// Actualizar el contador de visualizaciones en 1
	_, err = db.Exec("UPDATE visualizations SET views = views + 1 WHERE movie_id = $1", movieID)
	if err != nil {
		return nil, err
	}

	// Obtener los datos actualizados después de la actualización
	visualization := &models.Visualizations{}
	err = db.QueryRow("SELECT movie_id, views FROM visualizations WHERE movie_id = $1", movieID).
		Scan(&visualization.MovieID, &visualization.Views)
	if err != nil {
		return nil, err
	}

	return visualization, nil
}

func GetMovieComments(db *sql.DB, movieID int) ([]models.CommentDetails, error) {
	// consulta a la base de datos utilizando JOIN
	rows, err := db.Query(`
		SELECT u.username, c.comment_text, c.created_date, c.updated_date
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.movie_id = $1
	`, movieID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.CommentDetails

	// Iterar sobre las filas y escanear los comentarios
	for rows.Next() {
		var comment models.CommentDetails
		err := rows.Scan(
			&comment.Username,
			&comment.CommentText,
			&comment.CreatedDate,
			&comment.UpdatedDate,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	// Manejar cualquier error que pueda haber ocurrido durante la iteración de las filas
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func GetMostViewedMovies(db *sql.DB, n int) ([]models.Visualizations, error) {

	query := `SELECT movie_id, views FROM visualizations ORDER BY views DESC LIMIT $1`

	// Ejecutar la consulta y obtener los n resultados
	rows, err := db.Query(query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterar sobre los resultados y mapearlos a objetos models.Visualizations
	var mostViewedMovies []models.Visualizations
	for rows.Next() {
		var visualization models.Visualizations
		if err := rows.Scan(&visualization.MovieID, &visualization.Views); err != nil {
			return nil, err
		}
		mostViewedMovies = append(mostViewedMovies, visualization)
	}

	// Manejar cualquier error que ocurra durante la iteración
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mostViewedMovies, nil
}
