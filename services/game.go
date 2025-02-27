package services

import (
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"mognjen/gossassins/repos"
)

type GameService struct {
	gameRepo       *repos.GameRepo
	gamePlayerRepo *repos.GamePlayerRepo
}

func NewGameService(gameRepo *repos.GameRepo, gamePlayerRepo *repos.GamePlayerRepo) *GameService {
	return &GameService{gameRepo, gamePlayerRepo}
}

func (s *GameService) GetAllCreated(userId string) ([]models.Game, apierrors.StatusError) {
	return s.gameRepo.GetAllCreated(userId)
}

func (s *GameService) GetAllJoined(userId string) ([]models.Game, apierrors.StatusError) {
	return s.gameRepo.GetAllJoined(userId)
}

func (s *GameService) GetById(callerId string, id string) (*models.GameWithJoinStatus, apierrors.StatusError) {
	game, err := s.gameRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	players, err := s.gamePlayerRepo.GetAllByGameId(game.Id)
	if err != nil {
		return nil, err
	}

	joinedStatus := false
	for _, p := range players {
		if p.UserId == callerId {
			joinedStatus = true
			break
		}
	}

	return &models.GameWithJoinStatus{
		Game:   *game,
		Joined: joinedStatus,
	}, nil
}

func (s *GameService) GetIdByJoinCode(joinCode string) (*string, apierrors.StatusError) {
	game, err := s.gameRepo.GetByJoinCode(joinCode)
	if err != nil {
		return nil, err
	}

	return &game.Id, nil
}

func (s *GameService) Create(game *models.GameCreation) apierrors.StatusError {
	return s.gameRepo.Create(game)
}

func (s *GameService) Patch(id string, patch *models.GamePatch) apierrors.StatusError {
	return s.gameRepo.Patch(id, patch)
}

func (s *GameService) Delete(id string) apierrors.StatusError {
	return s.gameRepo.Delete(id)
}
