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

type JoinRequestHandler struct {
	joinRequestService JoinRequestService
}

func NewJoinRequestHandler(service JoinRequestService) *JoinRequestHandler {
	return &JoinRequestHandler{service}
}

type JoinRequestService interface {
	GetAllByGameId(gameId int) ([]models.JoinRequest, apierrors.StatusError)
	Create(joinRequest *models.JoinRequest) apierrors.StatusError
	Patch(gameId int, userId string, patch *models.JoinRequest) apierrors.StatusError
	Approve(gameId int, userId string) apierrors.StatusError
	Unapprove(gameId int, userId string) apierrors.StatusError
}

func (h *JoinRequestHandler) GetAllByGameId(context *gin.Context) {
	gameId, _ := strconv.Atoi(context.Param("game_id"))
	joinRequests, err := h.joinRequestService.GetAllByGameId(gameId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, joinRequests)
}

func (h *JoinRequestHandler) Create(context *gin.Context) {
	var request dto.CreateJoinRequestRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.UserId == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing user_id"})
		return
	}

	gameId, _ := strconv.Atoi(context.Param("game_id"))

	joinRequest := models.JoinRequest{
		GameId: gameId,
		UserId: *request.UserId,
		Status: models.NotApproved,
	}

	err := h.joinRequestService.Create(&joinRequest)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusCreated, "")
}

func (h *JoinRequestHandler) Patch(context *gin.Context) {
	gameId, _ := strconv.Atoi(context.Param("game_id"))
	userId := context.Param("user_id")
	var request dto.PatchJoinRequestRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err := validateRequest(request)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	patch := models.JoinRequest{
		GameId: gameId,
		UserId: userId,
		Status: models.JoinRequestStatus(*request.Status),
	}

	err = h.joinRequestService.Patch(gameId, userId, &patch)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}

func validateRequest(patch dto.PatchJoinRequestRequest) apierrors.StatusError {
	if patch.Status == nil {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("missing status"))
	} else if !models.IsValidJoinRequestStatus(*patch.Status) {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("invalid status value, status can be APPROVED, NOT_APPROVED"))
	}

	return nil
}

func (h *JoinRequestHandler) Approve(context *gin.Context) {
	gameId, _ := strconv.Atoi(context.Param("game_id"))
	userId := context.Param("user_id")

	err := h.joinRequestService.Approve(gameId, userId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.AbortWithStatus(http.StatusOK)
}

func (h *JoinRequestHandler) Unapprove(context *gin.Context) {
	gameId, _ := strconv.Atoi(context.Param("game_id"))
	userId := context.Param("user_id")

	err := h.joinRequestService.Unapprove(gameId, userId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.AbortWithStatus(http.StatusOK)
}
