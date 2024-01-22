package utils

import (
	"errors"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ExtractUserIDFromToken(tokenString string) (string, error) {

	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 {
		return "", errors.New("Token no válido")
	}
	tokenPart := parts[1]

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	// Decodificar el token
	token, err := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return "", err
	}

	// Verificar si el token es válido
	if !token.Valid {
		return "", errors.New("Token no válido")
	}

	// Extraer el ID de usuario de las reclamaciones del token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Error al obtener reclamaciones del token")
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return "", errors.New("ID de usuario no válido en el token")
	}

	return userID, nil
}
