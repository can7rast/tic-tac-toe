package web

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"school21/internal/application"
	"school21/internal/web/jwt"
	"school21/internal/web/middleware"
	"strings"

	"school21/internal/domain"
	"school21/internal/infrastructure/datasource"
	"school21/internal/web/handler"
)

func SetupRouter(service domain.GameService, db *datasource.DB, authService application.AuthService, repo datasource.UserRepository, provider *jwt.Provider) *gin.Engine {
	gameH := handler.NewGameHandler(service, db)
	authH := handler.NewAuthHandler(authService)
	userH := handler.NewUserHandler(repo)
	router := gin.Default()

	log.Printf("Регистрируем роуты: /login, /signup, /game и т.д.")

	router.POST("/login", authH.Login)
	router.POST("/signup", authH.SignUp)

	authorized := router.Group("/")
	authorized.Use(middleware.JwtUserAuthenticator(provider))

	authorized.POST("/game", gameH.CreateGame)
	authorized.GET("/game/:id", gameH.GetGame)
	authorized.POST("/game/:id", gameH.MakeMove)
	authorized.POST("/game/:id/join", gameH.JoinGame)
	authorized.GET("/games", gameH.AvailableGames)
	authorized.GET("/user/:id", userH.GetUser)

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/game") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		path := "frontend/index.html"
		if _, err := os.Stat(path); err == nil {
			log.Printf("Отдаём index.html для пути: %s", c.Request.URL.Path)
			c.File(path)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "HTML файл не найден"})
	})

	return router
}
