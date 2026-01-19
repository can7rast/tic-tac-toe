package web

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"school21/internal/application"
	"school21/internal/web/middleware"
	"strings"

	"school21/internal/domain"
	"school21/internal/infrastructure/datasource"
	"school21/internal/web/handler"
)

func SetupRouter(service domain.GameService, db *datasource.DB, authService application.AuthService) *gin.Engine {
	gameH := handler.NewGameHandler(service, db)
	authH := handler.NewAuthHandler(authService)

	router := gin.Default()

	// ← Добавь это для дебага
	log.Printf("Регистрируем роуты: /login, /signup, /game и т.д.")

	router.POST("/login", authH.Login)
	router.POST("/signup", authH.SignUp)

	authorized := router.Group("/")
	authorized.Use(middleware.UserAuthenticator(authService))

	authorized.POST("/game", gameH.CreateGame)
	authorized.GET("/game/:id", gameH.GetGame)
	authorized.POST("game/:id", gameH.MakeMove)

	// ← Добавь лог после регистрации
	log.Printf("Все API-роуты зарегистрированы. NoRoute идёт последним.")

	router.NoRoute(func(c *gin.Context) {
		log.Printf("NoRoute triggered! Path: %s", c.Request.URL.Path)
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
