package utils

import (
	"os"
	"time"

	"github.com/Michael-Wilburn/go_movie_api/models"
	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(user *models.UserResponse) (string, error) {

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 12).Unix(), // Caducidad del token en 12 horas
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return "Bearer " + signedToken, nil
}
