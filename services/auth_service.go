package services

import (
	"database/sql"

	"github.com/Michael-Wilburn/go_movie_api/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthUserCredentials verifica las credenciales en la base de datos y devuelve el usuario correspondiente
func AuthUserCredentials(db *sql.DB, credentials models.Credentials) (*models.UserResponse, error) {
	var user models.User

	// Consulta el usuario en la base de datos por nombre de usuario o correo electrónico
	row := db.QueryRow("SELECT id, username, email, password_hash, created_date, updated_date FROM users WHERE username = $1 OR email = $2", credentials.Username, credentials.Email)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedDate, &user.UpdatedDate)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		return nil, err
	}

	responseUser := &models.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		CreatedDate: user.CreatedDate,
		UpdatedDate: user.UpdatedDate,
	}

	return responseUser, nil
}
