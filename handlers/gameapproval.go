package handlers

import (
	"mognjen/gossassins/apierrors"
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
	var approval models.GameApproval
	if err := context.BindJSON(&approval); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	gameId, _ := strconv.Atoi(context.Param("game_id"))
	approval.GameId = &gameId
	approval.Status = "NOT_APPROVED"

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
	var patch models.GameApproval
	if err := context.BindJSON(&patch); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	err := h.gameApprovalRepo.Patch(gameId, userId, &patch)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}
