package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"services/api-gateway/internal/config"
	"services/api-gateway/internal/handler"
	"services/api-gateway/internal/middleware"
)

func SetupRoutes(cfg *config.Config) http.Handler {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/auth/login", handler.AuthHandler).Methods("POST")

	// Secured routes group
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.RateLimitMiddleware)
	secured.Use(middleware.AuthMiddleware)

	secured.PathPrefix("/bookings").Handler(middleware.NewReverseProxy(cfg.BookingURL))
	secured.PathPrefix("/experts").Handler(middleware.NewReverseProxy(cfg.ExpertURL))
	secured.PathPrefix("/notifications").Handler(middleware.NewReverseProxy(cfg.NotifyURL))

	return router
}
