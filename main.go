package main

import (
	"api/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Set trusted proxies dynamically
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies == "" {
		trustedProxies = "0.0.0.0/0" // Default (not recommended for production)
	}
	proxies := strings.Split(trustedProxies, ",")
	if err := r.SetTrustedProxies(proxies); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Get CORS origins from environment variable
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://react-frontend-service:80" // Default for local development
	}
	origins := strings.Split(corsOrigins, ",")

	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allow Authorization header
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add centralized error handling middleware
	r.Use(errorHandler())

	// Setup routes
	r.GET("/tools", routes.GetTools)
	r.GET("/tools/:name", routes.GetToolByName)
	r.POST("/vote/:tool", routes.VoteForTool)

	// Health check endpoint (for Kubernetes probes)
	r.GET("/health", func(c *gin.Context) {
		if isHealthy() {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy"})
		}
	})

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Start server with graceful shutdown
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Starting server on port %s", port)
		serverErrors <- server.ListenAndServe() // Use ListenAndServeTLS if needed
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)
	case <-quit:
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
		log.Println("Server exited gracefully")
	}
}

// Centralized error handling middleware
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("Error: %v", err)
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
		}
	}
}

// Health check function for Kubernetes probes
func isHealthy() bool {
	// Here, you can check DB connection, cache status, etc.
	return true // Update this based on actual health conditions
}
