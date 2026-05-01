package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTClaims represents the structure provided by the user
type JWTClaims struct {
	Role   string `json:"role"`
	UserID string `json:"userId"`
	Sub    string `json:"sub"`
	Iat    int64  `json:"iat"`
	Exp    int64  `json:"exp"`
}

// AuthMiddleware extracts the user role from the JWT in the Authorization header.
// It assumes the JWT is already verified by an API Gateway.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := parseJWTUnverified(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Store claims in context
		c.Set("role", claims.Role)
		c.Set("userId", claims.UserID)
		
		c.Next()
	}
}

// parseJWTUnverified performs a base64 decode of the JWT payload without signature verification.
// Use only when the token has been verified upstream (e.g. by an API Gateway).
func parseJWTUnverified(tokenString string) (*JWTClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
