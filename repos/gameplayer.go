package repos

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"net/http"
	"strings"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type GamePlayerRepo struct {
	db *supabase.Client
}

func NewGamePlayerRepo(db *supabase.Client) *GamePlayerRepo {
	return &GamePlayerRepo{db}
}

func (r *GamePlayerRepo) GetAllByGameId(gameId string) ([]models.GamePlayer, apierrors.StatusError) {
	var players []models.GamePlayer
	query := r.db.
		From("game_players").
		Select("*", "exact", false).
		Eq("game_id", gameId).
		Order("user_id", &postgrest.DefaultOrderOpts)

	_, err := execeuteSelect(query, &players)
	if err != nil {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return players, nil
}

func (r *GamePlayerRepo) GetByGameIdUserId(gameId string, userId string) (*models.GamePlayer, apierrors.StatusError) {
	var players []models.GamePlayer
	query := r.db.
		From("game_players").
		Select("*", "exact", false).
		Eq("game_id", gameId).
		Eq("user_id", userId)

	count, err := execeuteSelect(query, &players)
	if err != nil {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	if count == 0 {
		return nil, apierrors.NewStatusError(http.StatusNotFound, errors.New("GamePlayer not found"))
	}

	return &players[0], nil
}

func (r *GamePlayerRepo) Create(gamePlayer *models.GamePlayer) apierrors.StatusError {
	_, _, err := r.db.
		From("game_players").
		Insert(gamePlayer, false, "", "", "").
		Execute()

	if err != nil {
		// MMMM, such nice code
		duplicateKeyErr := "(23505)"
		if strings.Contains(err.Error(), duplicateKeyErr) {
			return apierrors.NewStatusError(http.StatusOK, err)
		}

		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}

func (r *GamePlayerRepo) Patch(gameId string, userId string, patch *models.GamePlayerPatch) apierrors.StatusError {
	_, _, err := r.db.
		From("game_players").
		Update(patch, "", "").
		Eq("game_id", gameId).
		Eq("user_id", userId).
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}

func (r *GamePlayerRepo) Delete(gameId string, userId string) apierrors.StatusError {
	_, _, err := r.db.
		From("game_players").
		Delete("", "").
		Eq("game_id", gameId).
		Eq("user_id", userId).
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}
