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

type GameActionService struct {
	gameRepo *repos.GameRepo
	db       *supabase.Client
}

func NewGameActionService(gameRepo *repos.GameRepo, db *supabase.Client) *GameActionService {
	return &GameActionService{gameRepo, db}
}

func (s *GameActionService) Start(gameId int) apierrors.StatusError {
	game, err := helpers.GetValidatedGame(s.gameRepo, gameId, models.OPEN)
	if err != nil {
		return err
	}

	errMsg := s.db.Rpc("start_game", "", gin.H{"p_game_id": game.Id})
	if errMsg != "" {
		return apierrors.NewStatusError(
			http.StatusInternalServerError,
			errors.New(errMsg),
		)
	}

	return nil
}

func (s *GameActionService) Kill(gameId int, killerUserId string, killCode string) apierrors.StatusError {
	game, err := s.gameRepo.GetById(gameId)
	if err != nil {
		return err
	}

	err = helpers.ValidateGameState(game, models.RUNNING)
	if err != nil {
		return err
	}

	/* Stored procedures are the only option for these kinds of transactions */
	errMsg := s.db.Rpc("kill_player", "", gin.H{
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
