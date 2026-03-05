package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"school21/internal/web/jwt"
	"strings"
)

func JwtUserAuthenticator(JwtProvider *jwt.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if authHeader == tokenStr || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "error in authorization"})
			return
		}

		claims, err := JwtProvider.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "error in jwt token"})
			return
		}

		c.Set("user_id", claims.UserID.String())
		c.Set("Login", claims.Login)
		c.Next()
	}
}
