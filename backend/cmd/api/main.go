package main

import (
	"log"

	"github.com/bmstoss13/the-daily-sunshine/internal/config"
	"github.com/bmstoss13/the-daily-sunshine/internal/handler"
	"github.com/bmstoss13/the-daily-sunshine/internal/middleware"
	"github.com/bmstoss13/the-daily-sunshine/internal/repository"
	"github.com/bmstoss13/the-daily-sunshine/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.DB.Close()

	//Create DB layer
	profileRepo := repository.NewPostgresProfileRepository(config.DB)
	profileSvc := service.NewProfileService(profileRepo)
	profileHandler := handler.NewProfileHandler(profileSvc)

	router := gin.Default()
	api := router.Group("/api/v1")

	// PUBLIC ROUTES (No auth required)
	publicProfiles := api.Group("/profiles")
	{
		// Anyone can view a profile
		publicProfiles.GET("/id/:id", profileHandler.GetProfile)
		publicProfiles.GET("/username/:username", profileHandler.GetProfileFromUsernameHandler)
		publicProfiles.GET("/check-username", profileHandler.GetUsernameAvailability)
		publicProfiles.GET("/list", profileHandler.GetProfileList) // example: /api/v1/profiles/list?limit=20&offset=20
	}

	// publicPosts := api.Group("/posts")
	// {

	// }

	// PROTECTED ROUTES (Require valid JWT)
	protected := api.Group("/")
	protected.Use(middleware.RequireAuth())
	{
		protected.POST("/profiles", profileHandler.CreateProfileHandler)
	}

	log.Println("Starting The Daily Sunshine API on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
