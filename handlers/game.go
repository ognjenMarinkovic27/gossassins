package handlers

import (
	"context"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameRepo GameRepo
}

type GameRepo interface {
	GetAll(ctx context.Context) ([]models.Game, apierrors.StatusError)
	GetById(id int) (*models.Game, apierrors.StatusError)
	Create(game *models.Game) apierrors.StatusError
	Patch(id int, patch *models.Game) apierrors.StatusError
	Delete(id int) apierrors.StatusError
}

func NewGameHandler(gameRepo GameRepo) *GameHandler {
	return &GameHandler{gameRepo}
}

func (h *GameHandler) GetAll(context *gin.Context) {
	games, err := h.gameRepo.GetAll(context)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, games)
}

func (h *GameHandler) GetById(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	games, err := h.gameRepo.GetById(id)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, games)
}

func (h *GameHandler) Create(context *gin.Context) {
	var game models.Game
	if err := context.BindJSON(&game); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	game.State = "OPEN"

	err := h.gameRepo.Create(&game)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusCreated, "")
}

func (h *GameHandler) Patch(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var patch models.Game
	if err := context.BindJSON(&patch); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	nullifyUnpatchable(&patch)

	err := h.gameRepo.Patch(id, &patch)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}

func nullifyUnpatchable(patch *models.Game) {
	patch.CreatedBy = ""
	patch.State = ""
}

func (h *GameHandler) Delete(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	err := h.gameRepo.Delete(id)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}
