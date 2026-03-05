package handler

import (
	"log"
	"net/http"
	"school21/internal/domain"
	"school21/internal/infrastructure/datasource"
	"school21/internal/web/dto"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GameHandler struct {
	service domain.GameService
	repo    domain.GameRepository
}

func NewGameHandler(service domain.GameService, db *datasource.DB) *GameHandler {
	return &GameHandler{
		service: service,
		repo:    datasource.NewGameRepository(db),
	}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	VsAi := strings.ToLower(c.Query("vs")) == "computer"
	userId, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	game := domain.NewGame(userId, VsAi)

	if err = h.repo.Save(*game); err != nil {
		log.Printf("Ошибка при сохранении игры %s: %v", game.ID.String(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant create game"})
		return
	}

	response := dto.FromDomain(*game, false, nil)
	c.JSON(http.StatusCreated, response)
}

func (h *GameHandler) JoinGame(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	game, err := h.repo.Get(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	err = h.service.JoinGame(&game, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(game.ID.String())
	err = h.repo.Save(game)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.FromDomain(game, false, nil))
}

func (h *GameHandler) MakeMove(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id игры"})
		return
	}

	//Берем id пользователя, который нажал на ход
	userId, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req dto.GameRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в обработке json"})
		return
	}

	if req.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID в пути и в теле должны совпадать"})
		return
	}

	game := req.ToDomain()
	log.Printf("vs ai = %v", game.VsAI)
	err = h.service.ValidateMove(&game, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.repo.Save(game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant save game"})
		return
	}

	game, err = h.repo.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось загрузить актуальную версию игры"})
		return
	}

	if game.VsAI && game.State == domain.TurnPlayer2 {
		updatedGame, err := h.service.GetNextMove(&game)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при ходе компьютера"})
			return
		}

		if err = h.repo.Save(updatedGame); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant save game"})
			return
		}
		game = updatedGame
	}
	isOver, err, Winner := h.service.IsGameOver(&game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверка после хода игрока"})
		return
	}
	response := dto.FromDomain(game, isOver, Winner)
	c.JSON(http.StatusOK, response)
}

func (h *GameHandler) GetGame(c *gin.Context) {
	accept := c.GetHeader("Accept")

	if strings.Contains(accept, "application/json") {
		idStr := c.Param("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID игры не указан"})
			return
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID игры"})
			return
		}

		game, err := h.repo.Get(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "игра не найдена"})
			return
		}

		isOver, err, winner := h.service.IsGameOver(&game)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка проверки состояния игры"})
			return
		}

		response := dto.FromDomain(game, isOver, winner)
		c.JSON(http.StatusOK, response)
		return
	}

	c.File("frontend/index.html")
}

func (h *GameHandler) AvailableGames(c *gin.Context) {
	games, err := h.service.ShowAllAvailableGames()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить список доступных игр"})
		return
	}
	gameListResponse := dto.FromDomainList(games)
	c.JSON(http.StatusOK, gameListResponse)
}
