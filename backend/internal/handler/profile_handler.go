package handler

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/bmstoss13/the-daily-sunshine/internal/helpers"
	"github.com/bmstoss13/the-daily-sunshine/internal/service"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	// hold a pointer to the service so all methods can use it
	profileService *service.ProfileService
}

func NewProfileHandler(svc *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: svc,
	}
}

// `binding:"required"` tags tell Gin to automatically reject the request
// if frontend forgets to send those fields
type CreateProfileRequest struct {
	FirstName string  `form:"first_name" binding:"required"`
	LastName  string  `form:"last_name" binding:"required"`
	Username  string  `form:"username" binding:"required"`
	Role      string  `form:"role" binding:"required"`
	Bio       *string `form:"bio"` // Pointer makes it optional in JSON
}

type UpdateProfileRequest struct {
	FirstName string  `form:"first_name" binding:"required"`
	LastName  string  `form:"last_name" binding:"required"`
	Username  string  `form:"username" binding:"required"`
	Role      string  `form:"role" binding:"required"`
	Bio       *string `form:"bio"` // Pointer makes it optional in JSON
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	targetProfileID := c.Param("id")
	requestingUserID := c.GetString("userID")

	profile, err := h.profileService.FetchProfileByID(c.Request.Context(), requestingUserID, targetProfileID)

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

func (h *ProfileHandler) GetProfileFromUsernameHandler(c *gin.Context) {
	targetUsername := c.Param("username")
	requestingUserID := c.GetString("userID")

	profile, err := h.profileService.FetchProfileByUsername(c.Request.Context(), requestingUserID, targetUsername)

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

func (h *ProfileHandler) GetUsernameAvailability(c *gin.Context) {
	usernameToCheck := c.Query("username") // Extracts from ?username=xxx
	requestingUserID := c.GetString("userID")

	isTaken, err := h.profileService.IsUsernameTaken(c.Request.Context(), requestingUserID, usernameToCheck)

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

func (h *ProfileHandler) GetProfileList(c *gin.Context) {

	var limit int32 = 20
	var offset int32 = 0

	limitStr := c.Query("limit") //limit?=5 for example
	if limitStr != "" {
		parsedLimit, err := helpers.StringToBase10Int32(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter. Must be a positive integer."})
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
	profileList, err := h.profileService.FetchListOfProfiles(c.Request.Context(), requestingUserID, limit, offset)

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

func (h *ProfileHandler) CreateProfileHandler(c *gin.Context) {
	requestingUserID := c.GetString("userID")

	var req CreateProfileRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON payload",
			"details": err.Error(),
		})
		return
	}

	newProfile := domain.Profile{
		ID:        requestingUserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Role:      domain.MembershipRole(req.Role),
		Bio:       req.Bio,
	}

	var imageBytes []byte
	var fileExtension string

	fileHeader, err := c.FormFile("profile_image")
	if err == nil {
		// Open the file
		file, openErr := fileHeader.Open()
		if openErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer file.Close()

		// Read it into bytes for the Service layer
		imageBytes, _ = io.ReadAll(file)
		// Grab the extension (e.g., ".jpg")
		fileExtension = filepath.Ext(fileHeader.Filename)
	}

	createdProfile, err := h.profileService.CreateUserProfile(
		c.Request.Context(),
		requestingUserID,
		newProfile,
		imageBytes,
		fileExtension,
	)

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

func (h *ProfileHandler) UpdateProfileHandler(c *gin.Context) {
	requestingUserID := c.GetString("userID")

	var req UpdateProfileRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON payload",
			"details": err.Error(),
		})
		return
	}

	newProfile := domain.Profile{
		ID:        requestingUserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Role:      domain.MembershipRole(req.Role),
		Bio:       req.Bio,
	}

	var imageBytes []byte
	var fileExtension string

	fileHeader, err := c.FormFile("profile_image")
	if err == nil {
		// Open the file
		file, openErr := fileHeader.Open()
		if openErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer file.Close()

		// Read it into bytes for the Service layer
		imageBytes, _ = io.ReadAll(file)
		// Grab the extension (e.g., ".jpg")
		fileExtension = filepath.Ext(fileHeader.Filename)
	}

	updatedProfile, err := h.profileService.UpdateUserProfile(
		c.Request.Context(),
		requestingUserID,
		newProfile,
		imageBytes,
		fileExtension,
	)

	if err != nil {
		if strings.Contains(err.Error(), "already taken") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": updatedProfile,
	})
}

func (h *ProfileHandler) SoftDeleteProfileHandler(c *gin.Context) {
	requestingUserID := c.GetString("userID")
	softDeletedProfile, err := h.profileService.SoftDeleteUserProfile(c.Request.Context(), requestingUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to soft delete profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": softDeletedProfile,
	})
}

func (h *ProfileHandler) PermanentlyDeleteProfile(c *gin.Context) {
	requestingUserID := c.GetString("userID")

	deletedProfile, err := h.profileService.PermanentlyDeleteUserProfile(c.Request.Context(), requestingUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to permanently delete profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": deletedProfile,
	})
}
