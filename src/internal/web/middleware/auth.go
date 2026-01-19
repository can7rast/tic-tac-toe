package middleware

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"school21/internal/application"
	"school21/internal/web/dto"
	"strings"
)

func UserAuthenticator(service application.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Требуется авторизация"})
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authHeader[6:])
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат base64"})
			return
		}

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Некорректные даные логин:пароль"})
			return
		}

		req := dto.SignUpRequest{
			Login:    creds[0],
			Password: creds[1],
		}

		user, err := service.Login(c.Request.Context(), req)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_login", user.Login)

		c.Next()
	}
}
