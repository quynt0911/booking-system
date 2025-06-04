package main

import (
	"log"
	"net/http"

	"services/api-gateway/internal/config"
	"services/api-gateway/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	router := routes.SetupRoutes(cfg)

	log.Println("API Gateway listening on", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
