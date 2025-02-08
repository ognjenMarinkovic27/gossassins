package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Define custom claims structure
type Claims struct {
	UserId string `json:"sub"`
	jwt.RegisteredClaims
}

// AuthMiddleware creates a gin handler for JWT authentication
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractBearerTokenFromHeader(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		claims, err := validateAndParseToken(token, jwtSecret)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		// Set user ID in context for handlers to use
		c.Set("userId", claims.UserId)
		c.Next()
	}
}

func extractBearerTokenFromHeader(c *gin.Context) (string, error) {
	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", ErrNoHeader
	}

	// Check if the header has the Bearer prefix
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", ErrInvalidFormat
	}

	return tokenParts[1], nil
}

func validateAndParseToken(tokenString string, jwtSecret string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims.UserId == "" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
