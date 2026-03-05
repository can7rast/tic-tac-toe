package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"school21/internal/application"
	"school21/internal/web/dto"
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
	var req dto.JwtRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Login == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login or password is empty"})
		return
	}

	ctx := c.Request.Context()
	if err := a.service.SignUp(ctx, req); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"login":   req.Login})
}

func (a *AuthHandler) Login(c *gin.Context) {
	var req dto.JwtRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Проблема в обработке json", "details": err.Error()})
		return
	}
	if req.Login == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login or password is empty"})
		return
	}
	resp, err := a.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	log.Println("Зарегистрировался user:", req.Login)
	c.JSON(http.StatusOK, resp)

}
