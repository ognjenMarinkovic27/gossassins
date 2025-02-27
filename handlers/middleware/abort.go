package middleware

import (
	"mognjen/gossassins/apierrors"

	"github.com/gin-gonic/gin"
)

func AbortAppropriately(ctx *gin.Context, json interface{}, err apierrors.StatusError) {
	if json != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(err.Status(), json)
	} else {
		ctx.AbortWithError(err.Status(), err)
	}
}
