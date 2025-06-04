package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"booking-system/internal/config"
	"booking-system/internal/handler"
	"booking-system/internal/middleware"
)


func SetupRoutes(cfg *config.Config) http.Handler {
	r := mux.NewRouter()
	r.Use(middleware.RateLimitMiddleware)
	r.Use(middleware.AuthMiddleware)

	r.HandleFunc("/auth/login", handler.AuthHandler).Methods("POST")
	r.PathPrefix("/bookings").Handler(middleware.NewReverseProxy(cfg.BookingURL))
	r.PathPrefix("/experts").Handler(middleware.NewReverseProxy(cfg.ExpertURL))
	r.PathPrefix("/notifications").Handler(middleware.NewReverseProxy(cfg.NotifyURL))

	return r
}
