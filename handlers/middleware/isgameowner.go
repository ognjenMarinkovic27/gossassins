package middleware

import (
	"errors"
	"mognjen/gossassins/repos"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IsGameOwnerMiddleware(gameRepo *repos.GameRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gameId, _ := strconv.Atoi(ctx.Param("game_id"))
		game, err := gameRepo.GetById(gameId)
		if err != nil {
			ctx.AbortWithError(err.Status(), err)
			return
		}

		callerId := ctx.GetString("userId")
		if callerId == "" {
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("userId not present in context"))
			return
		}

		if game.CreatedBy != callerId {
			ctx.Error(errors.New("not owner tried performing owner action")) // <--- great error message
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"type": "not_owner", "message": "Only the game owner is authorized for this action"})
			return
		}

		ctx.Next()
	}
}
