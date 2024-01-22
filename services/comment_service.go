package services

import (
	"database/sql"

	"github.com/Michael-Wilburn/go_movie_api/models"
)

func AddComment(db *sql.DB, comment models.CommentRequest, userID string) (*models.CommentResponse, error) {

	var commentCreated models.CommentResponse

	query := `INSERT INTO comments (user_id, movie_id, comment_text) VALUES ($1, $2, $3)
	RETURNING id, user_id, movie_id, comment_text, created_date, updated_date`

	err := db.QueryRow(query, userID, comment.MovieID, comment.CommentText).Scan(
		&commentCreated.ID,
		&commentCreated.UserID,
		&commentCreated.MovieID,
		&commentCreated.CommentText,
		&commentCreated.CreatedDate,
		&commentCreated.UpdatedDate,
	)

	if err != nil {
		return nil, err
	}

	return &commentCreated, nil
}

func DeleteCommentByID(db *sql.DB, commentID int, userID string) error {
	// Verificar si el comentario pertenece al usuario
	var count int
	err := db.QueryRow(`SELECT COUNT (*) FROM comments WHERE id = $1 AND user_id = $2`, commentID, userID).Scan(&count)

	if err != nil {
		return err
	}

	// El comentario no pertenece al usuario o no existe
	if count == 0 {
		return sql.ErrNoRows
	}

	// Borrar el comentario
	_, err = db.Exec(`DELETE FROM comments WHERE id = $1`, commentID)

	if err != nil {
		return err
	}

	return nil
}

func UpdateCommentByID(db *sql.DB, commentID int, userID string, updatedComment models.CommentRequest) error {
	// Verificar si el comentario pertenece al usuario
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM comments
		WHERE id = $1 AND user_id = $2
	`, commentID, userID).Scan(&count)

	if err != nil {
		return err
	}

	// El comentario no pertenece al usuario o no existe
	if count == 0 {
		return sql.ErrNoRows
	}

	// Actualizar el comentario
	_, err = db.Exec(`
		UPDATE comments
		SET comment_text = $1
		WHERE id = $2
	`, updatedComment.CommentText, commentID)

	if err != nil {
		return err
	}

	return nil
}
