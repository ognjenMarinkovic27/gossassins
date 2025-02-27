package middleware

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/repos"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsGameOwner(ctx *gin.Context, gameRepo *repos.GameRepo) (interface{}, apierrors.StatusError) {
	gameId := ctx.Param("game_id")
	game, err := gameRepo.GetById(gameId)
	if err != nil {
		return nil, err
	}

	callerId := ctx.GetString("userId")
	if callerId == "" {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, errors.New("userId not present in context"))
	}

	if game.CreatedBy != callerId {
		return gin.H{"type": "not_owner", "message": "Only the game owner is authorized for this action"},
			apierrors.NewStatusError(http.StatusUnauthorized, errors.New("not owner tried performing owner action")) // <-- great error message
	}

	return nil, nil
}

func IsGameOwnerMiddleware(gameRepo *repos.GameRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		json, err := IsGameOwner(ctx, gameRepo)
		if err != nil {
			AbortAppropriately(ctx, json, err)
			return
		}
		ctx.Next()
	}
}
