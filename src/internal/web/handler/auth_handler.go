package handler

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"school21/internal/application"
	"school21/internal/web/dto"
	"strings"
)

type AuthHandler struct {
	service application.AuthService
}

func NewAuthHandler(service application.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (a *AuthHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Login == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login or password is empty"})
	}

	ctx := c.Request.Context()
	if err := a.service.SignUp(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"login":   req.Login})
}

func (a *AuthHandler) Login(c *gin.Context) {
	log.Printf("LOGIN ХЕНДЛЕР ВЫЗВАН! Path: %s, Method: %s, Auth: %s",
		c.Request.URL.Path, c.Request.Method, c.GetHeader("Authorization"))

	authHeader := c.GetHeader("Authorization")
	lowerHeader := strings.ToLower(authHeader)
	if authHeader == "" || !strings.HasPrefix(lowerHeader, "basic") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	payload := authHeader[6:]

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	arrLogAndPassword := strings.SplitN(string(decoded), ":", 2)
	if len(arrLogAndPassword) != 2 {
		c.JSON(http.StatusOK, gin.H{"error": "login or password is empty"})
	}
	login := arrLogAndPassword[0]
	password := arrLogAndPassword[1]

	ctx := c.Request.Context()

	req := dto.SignUpRequest{
		Login:    login,
		Password: password,
	}
	user, err := a.service.Login(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": user.ID,
		"login":   login,
	})

}
