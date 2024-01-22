package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// Middleware para verificar el token JWT en las rutas protegidas
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Cargar las variables de entorno desde el archivo .env
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error al cargar el archivo .env")
		}
		tokenString := r.Header.Get("Authorization")

		// Verificar si se proporcionó un token en la cabecera "Authorization"
		if tokenString == "" {
			http.Error(w, "Token no proporcionado", http.StatusUnauthorized)
			return
		}

		// Dividir el token en sus partes
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Obtener la parte del token que contiene los datos
		tokenPart := parts[1]

		jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
		// Verificar la validez del token
		token, err := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Si el token es válido, permite el acceso a la ruta protegida
		next.ServeHTTP(w, r)
	}
}
