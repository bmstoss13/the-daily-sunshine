package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/config"
	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/handler"
	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/middleware"
)

func main() {
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.DB.Close()

	router := gin.Default()
	api := router.Group("/api/v1")

	// PUBLIC ROUTES (No auth required)
	publicProfiles := api.Group("/profiles")
	{
		// Anyone can view a profile
		publicProfiles.GET("/id/:id", handler.GetProfile)
		publicProfiles.GET("/username/:username", handler.GetProfileFromUsernameHandler)
		publicProfiles.GET("/check-username", handler.GetUsernameAvailability)
	}

	// publicPosts := api.Group("/posts")
	// {

	// }

	// PROTECTED ROUTES (Require valid JWT)
	protected := api.Group("/")
	protected.Use(middleware.RequireAuth())
	{

	}

	log.Println("Starting The Daily Sunshine API on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
