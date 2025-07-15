// main.go

//	@title			RAG File Service API
//	@version		1.0
//	@description	API for storing and searching embedded files
//	@host			localhost:8080
//	@BasePath		/
//	@schemes		http

package main

import (
	"log"

	"github.com/joho/godotenv"

	api "github.com/fain17/rag-backend/api/routes"
	"github.com/fain17/rag-backend/db"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	queries := db.ConnectDB()
	r := api.NewRouter(queries)

	r.Run(":8080")

}
