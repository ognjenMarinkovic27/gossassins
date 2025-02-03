package handlers

import (
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/dto"
	"mognjen/gossassins/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type GamePlayerHandler struct {
	db             *supabase.Client
	gamePlayerRepo GamePlayerRepo
}

type GamePlayerRepo interface {
	GetAllByGameId(gameId int) ([]models.GamePlayer, apierrors.StatusError)
	Create(player *models.GamePlayer) apierrors.StatusError
	Delete(gameId int, userId string) apierrors.StatusError
}

func (h *GamePlayerHandler) GetAllByGameId(context *gin.Context) {
	gameId, _ := strconv.Atoi(context.Param("game_id"))
	players, err := h.gamePlayerRepo.GetAllByGameId(gameId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, players)
}

func (h *GamePlayerHandler) Create(context *gin.Context) {
	var request dto.CreateGamePlayerRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	gameId, _ := strconv.Atoi(context.Param("game_id"))

	if request.UserId == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing user_id"})
		return
	}

	player := models.GamePlayer{
		GameId:   gameId,
		UserId:   *request.UserId,
		KillCode: nil,
		TargetId: nil,
		Status:   models.ALIVE,
	}

	err := h.gamePlayerRepo.Create(&player)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusCreated, "")
}

func (h *GamePlayerHandler) Delete(context *gin.Context) {
	gameId, _ := strconv.Atoi(context.Param("game_id"))
	userId := context.Param("user_id")
	err := h.gamePlayerRepo.Delete(gameId, userId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}
