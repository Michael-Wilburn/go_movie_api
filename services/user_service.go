package services

import (
	"database/sql"
	"errors"

	"github.com/Michael-Wilburn/go_movie_api/models"
	"github.com/Michael-Wilburn/go_movie_api/utils"
)

func AddUserDb(db *sql.DB, user models.User) (*models.User, error) {
	// Hash de la contraseña antes de almacenarla en la base de datos
	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)", user.Username, user.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	// Consultar los detalles del usuario recién insertado
	var createdUser models.User
	row := db.QueryRow("SELECT id, username, email, password_hash, created_date, updated_date FROM users WHERE email = $1", user.Email)
	if err := row.Scan(&createdUser.ID, &createdUser.Username, &createdUser.Email, &createdUser.PasswordHash, &createdUser.CreatedDate, &createdUser.UpdatedDate); err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func GetUserByID(db *sql.DB, userID string) (*models.User, error) {
	var user models.User

	row := db.QueryRow("SELECT id, username, email, created_date, updated_date FROM users WHERE id = $1", userID)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedDate, &user.UpdatedDate)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(db *sql.DB, userID string, updateUserReq models.User) error {

	if updateUserReq.Username == "" || updateUserReq.Email == "" || updateUserReq.PasswordHash == "" {
		return errors.New("All fields username, email and password must be provided for update")
	}

	query := "UPDATE users SET username = $1, email = $2, password_hash = $3 WHERE id = $4"

	passwordHash, err := utils.HashPassword(updateUserReq.PasswordHash)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, updateUserReq.Username, updateUserReq.Email, passwordHash, userID)
	if err != nil {
		return err
	}

	return nil
}
