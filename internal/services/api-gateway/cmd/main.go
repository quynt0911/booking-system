package main

import (
	"log"
	"net/http"

	"booking-system/internal/config"
	"booking-system/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	router := routes.SetupRoutes(cfg)

	log.Println("API Gateway listening on", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
