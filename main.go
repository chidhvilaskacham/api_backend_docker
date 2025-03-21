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

	// Trust proxies (e.g., AKS load balancer)
	r.SetTrustedProxies([]string{"0.0.0.0/0"})
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get CORS origins from environment variable
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000" // Default for local development
	}

	// Split CORS origins by comma
	origins := strings.Split(corsOrigins, ",")

	// Apply custom CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add centralized error handling middleware
	r.Use(errorHandler())

	// Setup routes
	r.GET("/tools", routes.GetTools)
	r.GET("/tools/:name", routes.GetToolByName)
	r.POST("/vote/:tool", routes.VoteForTool)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
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

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
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
