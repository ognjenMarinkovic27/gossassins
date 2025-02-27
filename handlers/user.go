package handlers

import (
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo UserRepo
}

type UserRepo interface {
	GetById(id string) (*models.User, apierrors.StatusError)
}

func NewUserHandler(repo UserRepo) *UserHandler {
	return &UserHandler{repo}
}

func (h *UserHandler) GetById(context *gin.Context) {
	id := context.Param("user_id")
	user, err := h.userRepo.GetById(id)

	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, *user)
}

func (h *UserHandler) GetMe(context *gin.Context) {
	id := context.GetString("userId")
	user, err := h.userRepo.GetById(id)

	if err != nil {
		context.AbortWithError(err.Status(), err)
		return
	}

	context.JSON(http.StatusOK, *user)
}
