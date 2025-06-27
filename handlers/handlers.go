package handlers

import (
	"encoding/json"
	"go-microservice-template/logger"
	"net/http"
	"time"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// Handlers struct holds dependencies
type Handlers struct {
}

// NewHandlers creates a new handlers instance
func NewHandlers() *Handlers {
	return &Handlers{}
}

// HealthCheck handles health check requests
func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "game-result-microservice",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	logger.Info("Health check requested")
}
