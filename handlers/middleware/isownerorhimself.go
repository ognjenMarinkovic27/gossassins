package middleware

import (
	"mognjen/gossassins/repos"

	"github.com/gin-gonic/gin"
)

func IsGameOwnerOrHimselfMiddleware(gameRepo *repos.GameRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := IsHimself(ctx)
		if err == nil {
			ctx.Next()
			return
		}

		json, err := IsGameOwner(ctx, gameRepo)
		if err == nil {
			ctx.Next()
			return
		}

		AbortAppropriately(ctx, json, err)
	}
}
