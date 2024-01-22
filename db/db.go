package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Driver de conexión con Postgres
)

type DbData struct {
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
}

func LoadEnv() (DbData, error) {
	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		return DbData{}, fmt.Errorf("error loading .env file: %w", err)
	}
	host, exists := os.LookupEnv("DB_HOST")
	if !exists {
		return DbData{}, errors.New("DB_HOST not found in environment variables")
	}

	port, exists := os.LookupEnv("DB_PORT")
	if !exists {
		return DbData{}, errors.New("DB_PORT not found in environment variables")
	}

	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		return DbData{}, errors.New("DB_NAME not found in environment variables")
	}

	username, exists := os.LookupEnv("DB_USER")
	if !exists {
		return DbData{}, errors.New("DB_USER not found in environment variables")
	}

	password, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		return DbData{}, errors.New("DB_PASSWORD not found in environment variables")
	}

	return DbData{
		Host:     host,
		Port:     port,
		DbName:   dbName,
		Username: username,
		Password: password,
	}, nil
}

var Db *sql.DB

func EstablishDbConnection() (*sql.DB, error) {
	dbData, err := LoadEnv()
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbData.Host, dbData.Port, dbData.Username, dbData.Password, dbData.DbName)

	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if err = Db.Ping(); err != nil {
		Db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	fmt.Println("Conexión exitosa a la base de datos")

	return Db, nil
}
