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
	router.HandleFunc("/auth/register", handler.HandleAuth).Methods("POST")
	router.HandleFunc("/auth/login", handler.HandleAuth).Methods("POST")
	router.HandleFunc("/auth/refresh", handler.HandleAuth).Methods("POST")

	// Secured routes group
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.RateLimitMiddleware)
	secured.Use(middleware.AuthMiddleware)

	// User service routes
	secured.PathPrefix("/users").Handler(middleware.NewReverseProxy(cfg.UserURL))

	// Other service routes
	secured.PathPrefix("/bookings").Handler(middleware.NewReverseProxy(cfg.BookingURL))
	secured.PathPrefix("/experts").Handler(middleware.NewReverseProxy(cfg.ExpertURL))
	secured.PathPrefix("/notifications").Handler(middleware.NewReverseProxy(cfg.NotifyURL))

	return router
}
