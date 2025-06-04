package main

import (
	"log"
	"net/http"
	"github.com/yourname/BOOKING-SYSTEM/internal/config"
	"github.com/yourname/BOOKING-SYSTEM/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	router := routes.SetupRoutes(cfg)

	log.Println("API Gateway listening on", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
