package main

import (
	"log"
	"net/http"
	"os"

	"incident-triage-assistant/internal/config"
	"incident-triage-assistant/internal/handler"
	"incident-triage-assistant/internal/repository"
	"incident-triage-assistant/internal/service"
	"incident-triage-assistant/internal/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database configuration
	dbConfig := config.NewDatabaseConfig()
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	incidentRepo := repository.NewMySQLIncidentRepository(db)

	// Initialize services
	aiService := service.NewOpenAIService()

	// Initialize use cases
	incidentUseCase := usecase.NewIncidentUseCase(incidentRepo, aiService)

	// Initialize handlers
	incidentHandler := handler.NewIncidentHandler(incidentUseCase)

	// Initialize Echo server
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Setup routes
	api := e.Group("/api/v1")
	
	// Health check
	api.GET("/health", incidentHandler.HealthCheck)
	
	// Incident routes
	incidents := api.Group("/incidents")
	incidents.POST("", incidentHandler.CreateIncident)
	incidents.GET("", incidentHandler.GetAllIncidents)
	incidents.GET("/:id", incidentHandler.GetIncident)
	incidents.PUT("/:id", incidentHandler.UpdateIncident)
	incidents.DELETE("/:id", incidentHandler.DeleteIncident)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
