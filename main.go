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
	r.Use(cors.Default())
	// Setup routes
	r.GET("/tools", routes.GetTools)
	r.GET("/tools/:name", routes.GetToolByName)
	r.POST("/vote/:tool", routes.VoteForTool)

	// Start server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
