package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/service"
)

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username parameter is required"})
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
