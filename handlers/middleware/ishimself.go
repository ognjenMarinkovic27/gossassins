package middleware

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsHimself(ctx *gin.Context) (interface{}, apierrors.StatusError) {
	userId := ctx.Param("user_id")
	callerId := ctx.GetString("userId")

	if callerId == "" {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, errors.New("userId not present in context"))
	}

	if userId != callerId {
		// stupid error, but I'm leaving it for comedic purposes
		return gin.H{"type": "not_himself", "message": "Only HIM is authorized for actions upon himself"},
			apierrors.NewStatusError(http.StatusUnauthorized, errors.New("himself tried performing action over not himself"))
	}

	return nil, nil
}

func IsHimselfMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		json, err := IsHimself(ctx)
		if err != nil {
			AbortAppropriately(ctx, json, err)
			return
		}
		ctx.Next()
	}
}
