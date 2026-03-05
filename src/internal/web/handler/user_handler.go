package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"school21/internal/infrastructure/datasource"
)

type UserHandler struct {
	repo datasource.UserRepository
}

func (u *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := u.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID.String(),
		"login": user.Login,
	})
}

func NewUserHandler(repo datasource.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}
