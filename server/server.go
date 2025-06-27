package server

import (
	"context"
	"fmt"
	"go-microservice-template/config"
	"go-microservice-template/handlers"
	"go-microservice-template/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	router     *chi.Mux
	handlers   *handlers.Handlers
}

// NewServer creates a new HTTP server instance
func NewServer() *Server {
	router := chi.NewRouter()
	handlers := handlers.NewHandlers()

	httpConfig := config.GetHttpConfig()

	server := &Server{
		router:   router,
		handlers: handlers,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", httpConfig.Host, httpConfig.Port),
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

// setupMiddleware configures the middleware stack
func (s *Server) setupMiddleware() {
	// Built-in middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware for development
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})
}

// setupRoutes configures all the routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.Get("/health", s.handlers.HealthCheck)

	// Documentation endpoint
	s.router.Get("/", s.serveDocumentation)
}

// serveDocumentation serves basic API documentation
func (s *Server) serveDocumentation(w http.ResponseWriter, r *http.Request) {
	documentation := `
<!DOCTYPE html>
<html>
<head>
    <title>Game Result Microservice API</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { margin: 20px 0; padding: 15px; border-left: 4px solid #007acc; background-color: #f5f5f5; }
        .method { font-weight: bold; color: #007acc; }
        .path { font-family: monospace; background-color: #e8e8e8; padding: 2px 4px; }
        pre { background-color: #f0f0f0; padding: 10px; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>Game Result Microservice API</h1>
    <p>This microservice handles game results, player actions, and chat messages via Redis pub/sub.</p>
    
    <h2>Endpoints</h2>
    
    <div class="endpoint">
        <div><span class="method">GET</span> <span class="path">/health</span></div>
        <p>Health check endpoint - returns service status and timestamp</p>
        <h4>Response:</h4>
        <pre>{
  "status": "healthy",
  "timestamp": "2025-06-27T10:00:00Z",
  "service": "game-result-microservice"
}</pre>
    </div>
    
    <h2>Response Format</h2>
    <p>All endpoints return JSON responses. The health endpoint confirms the service is running.</p>
    
    <h2>Future Endpoints</h2>
    <p>Additional endpoints for game results, player actions, and chat messages will be added as needed.</p>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(documentation))
}

// Start starts the HTTP server
func (s *Server) Start() error {
	httpConfig := config.GetHttpConfig()
	logger.Info(fmt.Sprintf("Starting HTTP server on %s:%d", httpConfig.Host, httpConfig.Port))

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
