package services

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"mognjen/gossassins/repos"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type GameActionsService struct {
	gameRepo *repos.GameRepo
	client   *supabase.Client
}

func (s *GameActionsService) Start(callerUserId string, gameId int) apierrors.StatusError {
	game, err := s.gameRepo.GetById(gameId)
	if err != nil {
		return err
	}

	if game.CreatedBy != callerUserId {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("Can't start a game you didn't create"))
	}

	if game.State == "STARTED" {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("Game already started"))
	}

	err = s.assignTargets(game)
	if err != nil {
		return err
	}

	err = s.setGameStatusToRunning(err, gameId, game)
	if err != nil {
		return err
	}

	return nil
}

func (s *GameActionsService) setGameStatusToRunning(err apierrors.StatusError, gameId int, game *models.Game) apierrors.StatusError {
	err = s.gameRepo.Patch(gameId, &models.Game{
		Id:        game.Id,
		Name:      game.Name,
		CreatedBy: game.CreatedBy,
		State:     models.RUNNING,
	})

	if err != nil {
		return apierrors.NewStatusError(
			http.StatusInternalServerError,
			err,
		)
	}
	return nil
}

func (s *GameActionsService) assignTargets(game *models.Game) apierrors.StatusError {
	/* Stored procedures are the only option for these kinds of transactions */
	errMsg := s.client.Rpc("assign_kill_codes_and_targets", "", gin.H{"p_game_id": game.Id})
	if errMsg != "" {
		return apierrors.NewStatusError(
			http.StatusInternalServerError,
			errors.New(errMsg),
		)
	}
	return nil
}
