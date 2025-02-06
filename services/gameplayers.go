package services

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"mognjen/gossassins/repos"
	"net/http"
)

type GamePlayerService struct {
	gamePlayerRepo *repos.GamePlayerRepo
}

func NewGamePlayerService(repo *repos.GamePlayerRepo) *GamePlayerService {
	return &GamePlayerService{repo}
}

func (s *GamePlayerService) GetAllByGameId(gameId int) ([]models.GamePlayer, apierrors.StatusError) {
	return s.gamePlayerRepo.GetAllByGameId(gameId)
}

func (s *GamePlayerService) GetByGameIdUserId(gameId int, userId string) (*models.GamePlayer, apierrors.StatusError) {
	return s.gamePlayerRepo.GetByGameIdUserId(gameId, userId)
}

func (s *GamePlayerService) Create(gamePlayer *models.GamePlayer) apierrors.StatusError {
	return s.gamePlayerRepo.Create(gamePlayer)
}

func (s *GamePlayerService) Patch(gameId int, userId string, patch *models.GamePlayerPatch) apierrors.StatusError {
	player, err := s.gamePlayerRepo.GetByGameIdUserId(gameId, userId)
	if err != nil {
		return err
	}

	if player.Status == models.DEAD {
		return apierrors.NewStatusError(
			http.StatusUnprocessableEntity,
			errors.New("Can't patch dead player"),
		)
	}

	return s.Patch(gameId, userId, patch)
}

func (s *GamePlayerService) Delete(gameId int, userId string) apierrors.StatusError {
	return s.Delete(gameId, userId)
}
