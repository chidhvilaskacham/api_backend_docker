package main

import (
	"api/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Apply custom CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},  // Allow requests from frontend
		AllowMethods:     []string{"GET", "POST"},            // Allow only GET and POST requests
		AllowHeaders:     []string{"Origin", "Content-Type"}, // Allow specific headers
		ExposeHeaders:    []string{"Content-Length"},         // Expose specific headers
		AllowCredentials: true,                               // Allow credentials (e.g., cookies)
	}))

	// Setup routes
	r.GET("/tools", routes.GetTools)
	r.GET("/tools/:name", routes.GetToolByName)
	r.POST("/vote/:tool", routes.VoteForTool)

	// Start server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
