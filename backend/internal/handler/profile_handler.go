package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/helpers"
	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/service"
)

// `binding:"required"` tags tell Gin to automatically reject the request
// if frontend forgets to send those fields
type CreateProfileRequest struct {
	FirstName    string  `json:"first_name" binding:"required"`
	LastName     string  `json:"last_name" binding:"required"`
	Username     string  `json:"username" binding:"required"`
	Role         string  `json:"role" binding:"required"`
	ProfileImage string  `json:"profile_image_url"`
	Bio          *string `json:"bio"` // Pointer makes it optional in JSON
}

func GetProfile(c *gin.Context) {
	targetProfileID := c.Param("id")
	requestingUserID := c.GetString("userID")

	profile, err := service.FetchProfileByID(c.Request.Context(), requestingUserID, targetProfileID)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}

func GetProfileFromUsernameHandler(c *gin.Context) {
	targetUsername := c.Param("username")
	requestingUserID := c.GetString("userID")

	profile, err := service.FetchProfileByUsername(c.Request.Context(), requestingUserID, targetUsername)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch profile from username",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}

func GetUsernameAvailability(c *gin.Context) {
	usernameToCheck := c.Query("username") // Extracts from ?username=xxx
	requestingUserID := c.GetString("userID")

	isTaken, err := service.IsUsernameTaken(c.Request.Context(), requestingUserID, usernameToCheck)

	if err != nil {
		if usernameToCheck == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username query is required"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check username availability",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"taken": isTaken,
	})
}

func GetProfileList(c *gin.Context) {

	var limit int32 = 20
	var offset int32 = 0

	limitStr := c.Query("limit") //limit?=5 for example
	if limitStr != "" {
		parsedLimit, err := helpers.StringToBase10Int32(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid limit parameter. Must be a positive integer."})
			return
		}
		limit = parsedLimit
	}

	offsetStr := c.Query("offset") //offset?=5 for example
	if offsetStr != "" {
		parsedOffset, err := helpers.StringToBase10Int32(offsetStr)
		if err != nil || parsedOffset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter. Cannot be negative."})
			return
		}
		offset = parsedOffset
	}

	requestingUserID := c.GetString("userID")
	profileList, err := service.FetchListOfProfiles(c.Request.Context(), requestingUserID, limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch profile list",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": profileList,
	})
}

func CreateProfileHandler(c *gin.Context) {
	requestingUserID := c.GetString("userID")

	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON payload",
			"details": err.Error(),
		})
		return
	}

	newProfile := domain.Profile{
		ID:              requestingUserID,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Username:        req.Username,
		Role:            domain.MembershipRole(req.Role),
		ProfileImageUrl: req.ProfileImage,
		Bio:             req.Bio,
	}

	createdProfile, err := service.CreateUserProfile(c.Request.Context(), requestingUserID, newProfile)

	if err != nil {
		if strings.Contains(err.Error(), "already taken") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": createdProfile,
	})
}
