package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// Step 1: Get the Authorization header
	// Convention: "Authorization: Bearer <token>"
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	// Step 2: The header looks like "Bearer eyJhbG..."
	// We split on space and take the second part
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
		return
	}
	tokenString := parts[1]

	// Step 3: Parse and validate the token
	// The callback here is how we tell the JWT library which secret to use for verification
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Enforce that the signing method is HMAC (HS256)
		// This prevents the "alg:none" attack — where an attacker sends a token
		// with no signature and claims it's valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	// Step 4: Extract claims and store user ID on the context
	// MapClaims is a map[string]interface{} — we read "sub" which we set to user ID in generateJWT
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token subject"})
		return
	}

	// Store user ID in context so handlers can access it
	c.Set("user_id", userID)

	// Pass control to the next handler
	c.Next()
}
