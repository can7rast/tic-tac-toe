package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"school21/internal/application"
	"school21/internal/web/dto"
	"school21/pkg"
)

func UserAuthenticator(service application.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		login, password, err := pkg.ParseBasicAuth(authHeader)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": err.Error(),
			})
			return
		}
		req := dto.SignUpRequest{
			Login:    login,
			Password: password,
		}

		user, err := service.Login(c.Request.Context(), req)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("user_id", user.ID.String())
		c.Set("user_login", user.Login)

		c.Next()
	}
}
