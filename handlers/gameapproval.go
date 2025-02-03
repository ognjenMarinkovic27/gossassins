package handlers

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/dto"
	"mognjen/gossassins/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GameApprovalHandler struct {
	gameApprovalRepo GameApprovalRepo
}

func NewGameApprovalHandler(repo GameApprovalRepo) *GameApprovalHandler {
	return &GameApprovalHandler{repo}
}

type GameApprovalRepo interface {
	GetAllByGameId(gameId int) ([]models.GameApproval, apierrors.StatusError)
	Create(approval *models.GameApproval) apierrors.StatusError
	Patch(gameId int, userId string, patch *models.GameApproval) apierrors.StatusError
}

func (h *GameApprovalHandler) GetAllByGameId(context *gin.Context) {
	gameId, _ := strconv.Atoi(context.Param("game_id"))
	approvals, err := h.gameApprovalRepo.GetAllByGameId(gameId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, approvals)
}

func (h *GameApprovalHandler) Create(context *gin.Context) {
	var request dto.CreateGameApprovalRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	if request.UserId == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing user_id"})
		return
	}

	gameId, _ := strconv.Atoi(context.Param("game_id"))

	approval := models.GameApproval{
		GameId: gameId,
		UserId: *request.UserId,
		Status: models.NotApproved,
	}

	err := h.gameApprovalRepo.Create(&approval)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusCreated, "")
}

func (h *GameApprovalHandler) Patch(context *gin.Context) {
	gameId, _ := strconv.Atoi(context.Param("game_id"))
	userId := context.Param("user_id")
	var request dto.PatchGameApprovalRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	err := validateRequest(request)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	patch := models.GameApproval{
		GameId: gameId,
		UserId: userId,
		Status: models.ApprovalStatus(*request.Status),
	}

	err = h.gameApprovalRepo.Patch(gameId, userId, &patch)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}

func validateRequest(patch dto.PatchGameApprovalRequest) apierrors.StatusError {
	if patch.Status == nil {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("missing status"))
	} else if !models.IsValidApprovalStatus(*patch.Status) {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("invalid status value, status can be APPROVED, NOT_APPROVED"))
	}

	return nil
}
