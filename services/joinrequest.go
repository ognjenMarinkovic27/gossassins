package services

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"mognjen/gossassins/repos"
	"mognjen/gossassins/services/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type JoinRequestService struct {
	gameRepo        *repos.GameRepo
	joinRequestRepo *repos.JoinRequestRepo
	db              *supabase.Client
}

func NewJoinRequestService(gameRepo *repos.GameRepo, joinRequestRepo *repos.JoinRequestRepo, db *supabase.Client) *JoinRequestService {
	return &JoinRequestService{gameRepo, joinRequestRepo, db}
}

func (s *JoinRequestService) GetAllByGameId(gameId int) ([]models.JoinRequest, apierrors.StatusError) {
	return s.joinRequestRepo.GetAllByGameId(gameId)
}

func (s *JoinRequestService) Create(game *models.JoinRequest) apierrors.StatusError {
	return s.joinRequestRepo.Create(game)
}

func (s *JoinRequestService) Patch(gameId int, userId string, patch *models.JoinRequest) apierrors.StatusError {
	return s.joinRequestRepo.Patch(gameId, userId, patch)
}

func (s *JoinRequestService) Approve(gameId int, userId string) apierrors.StatusError {
	err := helpers.ValidateGame(s.gameRepo, gameId, models.OPEN)
	if err != nil {
		return err
	}

	/* Stored procedures are the only option for these kinds of transactions */
	errMsg := s.db.Rpc("approve_player", "", gin.H{"p_game_id": gameId, "p_user_id": userId})
	if errMsg != "" {
		return apierrors.NewStatusError(
			http.StatusInternalServerError,
			errors.New(errMsg),
		)
	}
	return nil
}

func (s *JoinRequestService) Unapprove(gameId int, userId string) apierrors.StatusError {
	err := helpers.ValidateGame(s.gameRepo, gameId, models.OPEN)
	if err != nil {
		return err
	}

	/* Stored procedures are the only option for these kinds of transactions */
	errMsg := s.db.Rpc("unapprove_player", "", gin.H{"p_game_id": gameId, "p_user_id": userId})
	if errMsg != "" {
		return apierrors.NewStatusError(
			http.StatusInternalServerError,
			errors.New(errMsg),
		)
	}
	return nil
}
