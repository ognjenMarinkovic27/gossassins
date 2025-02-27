package handlers

import (
	"math/rand"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/dto"
	"mognjen/gossassins/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService GameService
}

type GameService interface {
	GetAllCreated(userId string) ([]models.Game, apierrors.StatusError)
	GetAllJoined(userId string) ([]models.Game, apierrors.StatusError)
	GetIdByJoinCode(joinCode string) (*string, apierrors.StatusError)
	GetById(callerId string, id string) (*models.GameWithJoinStatus, apierrors.StatusError)
	Create(game *models.GameCreation) apierrors.StatusError
	Patch(id string, patch *models.GamePatch) apierrors.StatusError
	Delete(id string) apierrors.StatusError
}

func NewGameHandler(gameService GameService) *GameHandler {
	return &GameHandler{gameService}
}

func (h *GameHandler) GetAllCreated(context *gin.Context) {
	userId := context.GetString("userId")
	games, err := h.gameService.GetAllCreated(userId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, games)
}

func (h *GameHandler) GetAllJoined(context *gin.Context) {
	userId := context.GetString("userId")
	games, err := h.gameService.GetAllJoined(userId)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, games)
}

func (h *GameHandler) GetById(context *gin.Context) {
	id := context.Param("game_id")
	callerId := context.GetString("userId")
	games, err := h.gameService.GetById(callerId, id)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, games)
}

func (h *GameHandler) GetIdByJoinCode(context *gin.Context) {
	joinCode := strings.ToUpper(context.Param("join_code"))
	gameId, err := h.gameService.GetIdByJoinCode(joinCode)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, gameId)
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
		JoinCode:  randSeq(7),
	}

	err := h.gameService.Create(&game)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusCreated, "")
}

// This function lives here who gives a shit
var chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func (h *GameHandler) Patch(context *gin.Context) {
	id := context.Param("game_id")
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

func (h *GameHandler) patchGame(request dto.PatchGameRequest, id string) apierrors.StatusError {
	patch := models.GamePatch{
		Name: request.Name,
	}

	err := h.gameService.Patch(id, &patch)
	if err != nil {
		return err
	}
	return nil
}

func (h *GameHandler) Delete(context *gin.Context) {
	id := context.Param("game_id")
	err := h.gameService.Delete(id)
	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, "")
}
