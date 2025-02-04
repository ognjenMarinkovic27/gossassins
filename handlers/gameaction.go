package handlers

import (
	"errors"
	"mognjen/gossassins/apierrors"
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

	if request.CallerUserId == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing user_id"})
		return
	}

	err := h.gameActionService.Start(*request.GameId, *request.CallerUserId)
	if err != nil {
		context.AbortWithStatusJSON(err.Status(), err)
		return
	}

	context.AbortWithStatus(http.StatusOK)
}

func (h *GameActionHandler) Approve(context *gin.Context) {
	var request dto.GameActionApproveRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := validateApprovalRequest(&request); err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	err := h.gameActionService.ApprovePlayer(*request.GameId, *request.ApproveeUserId, *request.CallerUserId)
	if err != nil {
		context.AbortWithStatusJSON(err.Status(), err)
		return
	}

	context.AbortWithStatus(http.StatusOK)
}

func (h *GameActionHandler) Unapprove(context *gin.Context) {
	var request dto.GameActionApproveRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := validateApprovalRequest(&request); err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	err := h.gameActionService.UnapprovePlayer(*request.GameId, *request.ApproveeUserId, *request.CallerUserId)
	if err != nil {
		context.AbortWithStatusJSON(err.Status(), err)
		return
	}

	context.AbortWithStatus(http.StatusOK)
}

func validateApprovalRequest(request *dto.GameActionApproveRequest) apierrors.StatusError {
	if request.GameId == nil {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("missing game_id"))
	}

	if request.CallerUserId == nil {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("missing user_id"))
	}

	if request.ApproveeUserId == nil {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("missing approvee_id"))
	}

	return nil
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

	if request.KillerUserId == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing killer_id"})
		return
	}

	if request.KillCode == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing kill_code"})
		return
	}

	err := h.gameActionService.Kill(*request.GameId, *request.KillerUserId, *request.KillCode)
	if err != nil {
		context.AbortWithStatusJSON(err.Status(), err)
		return
	}

	context.AbortWithStatus(http.StatusOK)
}
