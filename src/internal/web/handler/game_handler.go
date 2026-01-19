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
	db      *datasource.DB
	repo    datasource.GameRepository
}

func NewGameHandler(service domain.GameService, db *datasource.DB) *GameHandler {
	return &GameHandler{
		service: service,
		db:      db,
		repo:    datasource.NewGameRepository(db),
	}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	game := domain.NewGame()

	if err := h.repo.Save(*game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant create game"})
		return
	}

	response := dto.FromDomain(*game, false, nil)

	c.JSON(http.StatusCreated, response)
}

func (h *GameHandler) MakeMove(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id игры"})
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
	if err = h.service.ValidateMove(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.repo.Save(game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant save game"})
		return
	}

	gameOver, err, PlayerWinner := h.service.IsGameOver(&game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверка после хода игрока"})
		return
	}
	if gameOver {
		response := dto.FromDomain(game, true, PlayerWinner)
		c.JSON(http.StatusOK, response)
		return
	}

	updatedGame, err := h.service.GetNextMove(&game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при ходе компьютера"})
		return
	}

	if err = h.repo.Save(updatedGame); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant save game"})
		return
	}

	gameOver, err, FinalWinner := h.service.IsGameOver(&updatedGame)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке хода компьютера"})
		return
	}
	response := dto.FromDomain(updatedGame, gameOver, FinalWinner)
	c.JSON(http.StatusOK, response)
}

func (h *GameHandler) GetGame(c *gin.Context) {
	accept := c.GetHeader("Accept")
	log.Printf("GetGame: Accept='%s', Path='%s'", accept, c.Request.URL.Path)

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
