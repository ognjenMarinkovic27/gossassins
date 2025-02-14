package handlers

import (
	"mognjen/gossassins/dto"
	"mognjen/gossassins/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GameActionHandler struct {
	gameActionService *services.GameActionService
}

func NewGameActionHandler(gameActionService *services.GameActionService) *GameActionHandler {
	return &GameActionHandler{gameActionService}
}

func (h *GameActionHandler) Start(context *gin.Context) {
	var request dto.GameActionStartRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.GameId == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing game_id"})
		return
	}

	err := h.gameActionService.Start(*request.GameId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.AbortWithStatus(http.StatusOK)
}

func (h *GameActionHandler) Kill(context *gin.Context) {
	var request dto.GameActionKillRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.GameId == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing game_id"})
		return
	}

	if request.KillCode == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing kill_code"})
		return
	}

	killerId := context.GetString("userId")

	err := h.gameActionService.Kill(*request.GameId, killerId, *request.KillCode)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.AbortWithStatus(http.StatusOK)
}
