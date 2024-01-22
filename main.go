package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Michael-Wilburn/go_movie_api/auth"
	"github.com/Michael-Wilburn/go_movie_api/controllers"
	"github.com/Michael-Wilburn/go_movie_api/db"
	"github.com/Michael-Wilburn/go_movie_api/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/treblle/treblle-go"
)

func main() {

	APIKeyTreblle := os.Getenv("TREBLLE_API_KEY")
	ProjectIDTreblle := os.Getenv("TREBLLE_PROJECT_ID")

	treblle.Configure(treblle.Configuration{
		APIKey:    APIKeyTreblle,
		ProjectID: ProjectIDTreblle,
	})

	// Establecer conexión a la base de datos
	db, err := db.EstablishDbConnection()
	if err != nil {
		log.Fatal("Error establishing database connection:", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Println("Error closing database connection:", cerr)
		}
	}()

	// controladores
	movieController := controllers.NewMovieController(db)
	userController := controllers.NewUserController(db)
	authController := auth.NewAuthController(db)
	commentController := controllers.NewCommentController(db)

	// Crear un nuevo enrutador
	r := mux.NewRouter()

	r.Use(treblle.Middleware)

	// Rutas de la API v1
	r.HandleFunc("/api/v1/movies/{movie_id}", movieController.ViewMovieDetails).Methods("GET")
	r.HandleFunc("/api/v1/movies/popularity/{n}", movieController.MostViewMovies).Methods("GET")
	r.HandleFunc("/api/v1/users", userController.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", authController.Login).Methods("POST")
	r.HandleFunc("/api/v1/users/{id}", userController.GetUserById).Methods("GET")
	r.HandleFunc("/api/v1/users/{id}", middleware.AuthMiddleware(userController.UpdateUserById)).Methods("PUT")
	r.HandleFunc("/api/v1/comments", middleware.AuthMiddleware(commentController.CreateComment)).Methods("POST")
	r.HandleFunc("/api/v1/comments/{id}", middleware.AuthMiddleware(commentController.DeleteCommentById)).Methods("DELETE")
	r.HandleFunc("/api/v1/comments/{id}", middleware.AuthMiddleware(commentController.UpdateCommentById)).Methods("PUT")

	// Configurar el middleware CORS
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Iniciar el servidor
	fmt.Println("Servidor iniciado en http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}
