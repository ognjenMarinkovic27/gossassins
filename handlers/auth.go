package handlers

import (
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

type AuthHandler struct {
	db       *supabase.Client
	userRepo UserRepoForAuth
}

type UserRepoForAuth interface {
	Create(user *models.User) apierrors.StatusError
}

func NewAuthHandler(db *supabase.Client, userRepo UserRepoForAuth) *AuthHandler {
	return &AuthHandler{db, userRepo}
}

type AuthRequestBody struct {
	Email string `json:"email"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(context *gin.Context) {
	var body SignInRequest
	if err := context.BindJSON(&body); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res, err := h.db.Auth.SignInWithEmailPassword(body.Email, body.Password)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"token": res.AccessToken,
	})
}

type SignupRequestWithMetadata struct {
	Name     string `json:"name"`
	PhotoUrl string `json:"photo_url"`
	types.SignupRequest
}

func (h *AuthHandler) Signup(context *gin.Context) {
	var body SignupRequestWithMetadata
	if err := context.BindJSON(&body); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res, err := h.db.Auth.Signup(body.SignupRequest)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = h.userRepo.Create(&models.User{
		Uid:      res.User.ID.String(),
		Name:     body.Name,
		PhotoUrl: body.PhotoUrl,
	})

	if err != nil {
		// TODO: Very bad, user exists without metadata
		// For now who gives a shit <- can check on login?
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusOK, "")
}
