package main

import (
	"log"
	"net/http"

	"services/api-gateway/internal/config"
	"services/api-gateway/internal/routes"
)

func main() {
	cfg := config.NewConfig()
	router := routes.SetupRoutes(cfg)

	log.Println("API Gateway listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
