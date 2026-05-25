package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/bmstoss13/the-daily-sunshine/internal/config"
	"github.com/bmstoss13/the-daily-sunshine/internal/handler"
	"github.com/bmstoss13/the-daily-sunshine/internal/middleware"
	"github.com/bmstoss13/the-daily-sunshine/internal/redis"
	"github.com/bmstoss13/the-daily-sunshine/internal/repository"
	"github.com/bmstoss13/the-daily-sunshine/internal/service"
	"github.com/bmstoss13/the-daily-sunshine/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// We use a context to initialize our services so they can timeout if the network is down
	ctx := context.Background()

	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.DB.Close()

	//Initialize Cloudflare R2
	r2Client, err := config.InitR2(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudflare R2: %v", err)
	}

	// Initialize Upstash Redis
	redisClient, err := config.InitRedis(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	imageStorage := storage.NewCloudflareR2Storage(r2Client)
	serverCache := redis.NewRedisPostCache(redisClient)

	//Create DB layer
	profileRepo := repository.NewPostgresProfileRepository(config.DB)
	profileSvc := service.NewProfileService(profileRepo, imageStorage)
	profileHandler := handler.NewProfileHandler(profileSvc)

	postRepo := repository.NewPostgresPostRepository(config.DB)
	postImageRepo := repository.NewPostgresPostImageRepository(config.DB)
	postVideoRepo := repository.NewPostgresPostVideoRepository(config.DB)
	postService := service.NewPostService(postRepo, postImageRepo, postVideoRepo, imageStorage, serverCache)
	postHandler := handler.NewPostHandler(postService)

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

	publicPosts := api.Group("/posts")
	{
		publicPosts.GET("/id/:id", postHandler.GetPostByID)
		publicPosts.GET("/slug/:slug", postHandler.GetPostBySlug)
		publicPosts.GET("/check-slug", postHandler.GetSlugAvailability)
		publicPosts.GET("/list", postHandler.GetListOfPosts)
		publicPosts.GET("/top", postHandler.GetTopPostsOfDay)
	}

	// PROTECTED ROUTES (Require valid JWT)
	protected := api.Group("/")
	protected.Use(middleware.RequireAuth())
	{
		protected.POST("/profiles", profileHandler.CreateProfileHandler)
		protected.PUT("/profiles", profileHandler.UpdateProfileHandler)
		protected.DELETE("/profiles", profileHandler.PermanentlyDeleteProfile)
		protected.PATCH("/profiles/deactivate", profileHandler.SoftDeleteProfileHandler)

		protected.POST("/posts", postHandler.CreatePost)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
		protected.PATCH("/posts/:id", postHandler.SoftDeletePost)
	}

	protectedAdmin := api.Group("/admin")
	protectedAdmin.Use(middleware.RequireAuth())
	protectedAdmin.Use(middleware.RequireAdmin(profileSvc))
	{
		// admin-only routes
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Starting The Daily Sunshine API on port 8080...")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
	// log.Println("Starting The Daily Sunshine API on port 8080...")
	// if err := router.Run(":8080"); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
}
