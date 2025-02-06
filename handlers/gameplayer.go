package handlers

import (
	"errors"
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

func NewGamePlayerHandler(db *supabase.Client, repo GamePlayerRepo) *GamePlayerHandler {
	return &GamePlayerHandler{db, repo}
}

type GamePlayerRepo interface {
	GetAllByGameId(gameId int) ([]models.GamePlayer, apierrors.StatusError)
	GetByGameIdUserId(gameId int, userId string) (*models.GamePlayer, apierrors.StatusError)
	Create(player *models.GamePlayer) apierrors.StatusError
	Patch(gameId int, suerId string, player *models.GamePlayer) apierrors.StatusError
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

func (h *GamePlayerHandler) GetByGameIdUserId(context *gin.Context) {
	gameId, _ := strconv.Atoi(context.Param("game_id"))
	userId := context.Param("user_id")
	players, err := h.gamePlayerRepo.GetByGameIdUserId(gameId, userId)
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
		return
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
	gameId, _ := strconv.Atoi(context.Param("game_id"))
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

	patch := models.GamePlayer{
		GameId: gameId,
		UserId: userId,
		Status: models.PlayerStatus(*request.Status),
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
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("invalid status value, status can only be patched to ALIVE, NOT_APPROVED"))
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
	gameId, _ := strconv.Atoi(context.Param("game_id"))
	userId := context.Param("user_id")
	err := h.gamePlayerRepo.Delete(gameId, userId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}
