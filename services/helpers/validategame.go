package helpers

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"mognjen/gossassins/repos"
	"net/http"
)

/* TODO: This validate game bs is kinda ugly */
func GetValidatedGame(gameRepo *repos.GameRepo, gameId int, state models.GameState) (*models.Game, apierrors.StatusError) {
	game, err := gameRepo.GetById(gameId)
	if err != nil {
		return nil, err
	}

	err = ValidateGameState(game, state)
	if err != nil {
		return nil, err
	}

	return game, nil
}
func ValidateGameState(game *models.Game, state models.GameState) apierrors.StatusError {
	if game.State != state {
		return apierrors.NewStatusError(http.StatusBadRequest, errors.New("Game already started"))
	}

	return nil
}

func ValidateGame(gameRepo *repos.GameRepo, gameId int, state models.GameState) apierrors.StatusError {
	_, err := GetValidatedGame(gameRepo, gameId, state)
	return err
}
