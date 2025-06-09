package handler

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"services/api-gateway/internal/config"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	// Read the request body into a buffer
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Determine the target path based on the request path
	targetPath := ""
	switch r.URL.Path {
	case "/auth/register":
		targetPath = "/user/register"
	case "/auth/login":
		targetPath = "/user/login"
	case "/auth/refresh":
		targetPath = "/user/refresh"
	default:
		http.Error(w, "Unknown authentication path", http.StatusNotFound)
		return
	}

	// Forward request to user service
	cfg := config.NewConfig()
	client := &http.Client{}

	// Create new request with body from buffer
	userServiceReq, err := http.NewRequest(r.Method, cfg.UserURL+targetPath, bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		// Exclude headers that should not be forwarded, e.g., Host, Content-Length
		if strings.ToLower(key) == "host" || strings.ToLower(key) == "content-length" {
			continue
		}
		for _, value := range values {
			userServiceReq.Header.Add(key, value)
		}
	}

	// Set Content-Length for the new request
	userServiceReq.ContentLength = int64(len(bodyBytes))

	// Send request to user service
	resp, err := client.Do(userServiceReq)
	if err != nil {
		http.Error(w, "Error connecting to user service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Error copying response", http.StatusInternalServerError)
		return
	}
}
