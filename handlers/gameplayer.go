package handlers

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/dto"
	"mognjen/gossassins/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type GamePlayerHandler struct {
	db             *supabase.Client
	gamePlayerRepo GamePlayerRepo
}

func NewGamePlayerHandler(db *supabase.Client, repo GamePlayerRepo) *GamePlayerHandler {
	return &GamePlayerHandler{db, repo}
}

type GamePlayerRepo interface {
	GetAllByGameId(gameId string) ([]models.GamePlayer, apierrors.StatusError)
	GetByGameIdUserId(gameId string, userId string) (*models.GamePlayer, apierrors.StatusError)
	Create(player *models.GamePlayer) apierrors.StatusError
	Patch(gameId string, suerId string, player *models.GamePlayerPatch) apierrors.StatusError
	Delete(gameId string, userId string) apierrors.StatusError
}

func (h *GamePlayerHandler) GetAllByGameId(context *gin.Context) {
	gameId := context.Param("game_id")
	players, err := h.gamePlayerRepo.GetAllByGameId(gameId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, players)
}

func (h *GamePlayerHandler) GetByGameIdUserId(context *gin.Context) {
	gameId := context.Param("game_id")
	userId := context.Param("user_id")
	players, err := h.gamePlayerRepo.GetByGameIdUserId(gameId, userId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, players)
}

func (h *GamePlayerHandler) GetMe(context *gin.Context) {
	gameId := context.Param("game_id")
	userId := context.GetString("userId")
	players, err := h.gamePlayerRepo.GetByGameIdUserId(gameId, userId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, players)
}

func (h *GamePlayerHandler) Create(context *gin.Context) {
	gameId := context.Param("game_id")
	userId := context.GetString("userId")

	player := models.GamePlayer{
		GameId:   gameId,
		UserId:   userId,
		KillCode: nil,
		TargetId: nil,
		Status:   models.NOT_APPROVED,
	}

	err := h.gamePlayerRepo.Create(&player)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusCreated, "")
}

func (h *GamePlayerHandler) Patch(context *gin.Context) {
	gameId := context.Param("game_id")
	userId := context.Param("user_id")
	var request dto.PatchGamePlayerRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err := validateGamePlayerPatchRequest(request)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	newStatus := models.PlayerStatus(*request.Status)
	patch := models.GamePlayerPatch{
		Status: &newStatus,
	}

	err = h.gamePlayerRepo.Patch(gameId, userId, &patch)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}

func validateGamePlayerPatchRequest(patch dto.PatchGamePlayerRequest) apierrors.StatusError {
	if patch.Status == nil {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("missing status"))
	} else if !isValidGamePlayerPatchStatus(*patch.Status) {
		return apierrors.NewStatusError(http.StatusUnprocessableEntity, errors.New("invalid status value, status can only be patched to ALIVE, NOT_APPROVED"))
	}

	return nil
}

func isValidGamePlayerPatchStatus(value string) bool {
	switch models.PlayerStatus(value) {
	case models.ALIVE, models.NOT_APPROVED:
		return true
	default:
		return false
	}
}

func (h *GamePlayerHandler) Delete(context *gin.Context) {
	gameId := context.Param("game_id")
	userId := context.Param("user_id")
	err := h.gamePlayerRepo.Delete(gameId, userId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}
