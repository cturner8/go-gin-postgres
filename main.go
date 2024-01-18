package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"cturner8/go-gin-postgres/middlewares"
	"cturner8/go-gin-postgres/routes"
)

func connectDatabase() (*sql.DB, error) {
	passwordFile, err := os.ReadFile("/run/secrets/db_password")

	username := os.Getenv("DB_USERNAME")
	password := string(passwordFile)
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	if err != nil {
		return nil, err
	}

	conn_str := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", username, password, host, dbName)
	return sql.Open("postgres", conn_str)
}

func setupRouter(db *sql.DB) *gin.Engine {
	gin.ForceConsoleColor()
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./web/dist", false)))
	router.Use(middlewares.DatabaseMiddleware(db))

	api := router.Group("/api")
	{
		routes.RegisterAlbums(api)
	}

	return router
}

func main() {
	db, err := connectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")

	router := setupRouter(db)
	router.Run()
}
