package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const userIDContextKey = "userID"

type AdminAuthorizer interface {
	IsAdmin(ctx context.Context, userID string) (bool, error)
}

// RequireAuth intercepts the request, validates the Supabase JWT,
// and injects the user's UUID into the Gin context.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := authenticateRequest(c); !ok {
			return
		}
		c.Next()
	}
}

func RequireAdmin(authorizer AdminAuthorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authorizer == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Admin authorization is not configured"})
			return
		}

		userID := c.GetString(userIDContextKey)
		if userID == "" {
			var ok bool
			userID, ok = authenticateRequest(c)
			if !ok {
				return
			}
		}

		isAdmin, err := authorizer.IsAdmin(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify admin access"})
			return
		}
		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access is required"})
			return
		}

		c.Next()
	}
}

func authenticateRequest(c *gin.Context) (string, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return "", false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
		return "", false
	}

	userID, err := parseSupabaseUserID(parts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return "", false
	}

	c.Set(userIDContextKey, userID)
	return userID, true
}

func parseSupabaseUserID(tokenString string) (string, error) {
	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("server configuration error")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token payload")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in token")
	}

	return userID, nil
}
