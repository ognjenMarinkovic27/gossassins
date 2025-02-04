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

func (s *GameActionsService) Start(gameId int, callerUserId string) apierrors.StatusError {
	game, err := s.getValidatedGame(gameId, callerUserId, models.OPEN)

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

func (s *GameActionsService) ApprovePlayer(gameId int, userId string, callerUserId string) apierrors.StatusError {
	err := s.validateGame(gameId, callerUserId, models.OPEN)
	if err != nil {
		return err
	}

	/* Stored procedures are the only option for these kinds of transactions */
	errMsg := s.client.Rpc("approve_player", "", gin.H{"p_game_id": gameId, "p_user_id": userId})
	if errMsg != "" {
		return apierrors.NewStatusError(
			http.StatusInternalServerError,
			errors.New(errMsg),
		)
	}
	return nil
}

func (s *GameActionsService) Kill(gameId int, killerUserId string, killCode string) apierrors.StatusError {
	game, err := s.gameRepo.GetById(gameId)
	if err != nil {
		return err
	}

	err = s.validateGameState(game, models.RUNNING)
	if err != nil {
		return err
	}

	/* Stored procedures are the only option for these kinds of transactions */
	errMsg := s.client.Rpc("kill_player", "", gin.H{
		"p_game_id":   gameId,
		"p_killer_id": killerUserId,
		"p_kill_code": killCode,
	})

	if errMsg != "" {
		return apierrors.NewStatusError(
			http.StatusInternalServerError,
			errors.New(errMsg),
		)
	}
	return nil
}

/* TODO: This validate game bs is kinda ugly */
func (s *GameActionsService) getValidatedGame(gameId int, callerUserId string, state models.GameState) (*models.Game, apierrors.StatusError) {
	game, err := s.gameRepo.GetById(gameId)
	if err != nil {
		return nil, err
	}

	err = s.validateCreator(game, callerUserId)
	if err != nil {
		return nil, err
	}

	err = s.validateGameState(game, state)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (s *GameActionsService) validateGame(gameId int, callerUserId string, state models.GameState) apierrors.StatusError {
	_, err := s.getValidatedGame(gameId, callerUserId, state)
	return err
}

func (s *GameActionsService) validateCreator(game *models.Game, callerUserId string) apierrors.StatusError {
	if game.CreatedBy != callerUserId {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("Can't modify a game you didn't create"))
	}

	return nil
}

func (s *GameActionsService) validateGameState(game *models.Game, state models.GameState) apierrors.StatusError {
	if game.State != state {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("Game already started"))
	}

	return nil
}
