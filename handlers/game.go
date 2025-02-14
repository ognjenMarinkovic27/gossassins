package handlers

import (
	"context"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/dto"
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
	Create(game *models.GameCreation) apierrors.StatusError
	Patch(id int, patch *models.GamePatch) apierrors.StatusError
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
	id, _ := strconv.Atoi(context.Param("game_id"))
	games, err := h.gameRepo.GetById(id)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, games)
}

func (h *GameHandler) Create(context *gin.Context) {
	var request dto.CreateGameRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	userId := context.GetString("userId")

	game := models.GameCreation{
		Name:      request.Name,
		CreatedBy: userId,
		State:     models.OPEN,
	}

	err := h.gameRepo.Create(&game)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusCreated, "")
}

func (h *GameHandler) Patch(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("game_id"))
	var request dto.PatchGameRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.Name == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing name"})
		return
	}

	err := h.patchGame(request, id)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}

func (h *GameHandler) patchGame(request dto.PatchGameRequest, id int) apierrors.StatusError {
	patch := models.GamePatch{
		Name: request.Name,
	}

	err := h.gameRepo.Patch(id, &patch)
	if err != nil {
		return err
	}
	return nil
}

func (h *GameHandler) Delete(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("game_id"))
	err := h.gameRepo.Delete(id)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}
