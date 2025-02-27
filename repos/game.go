package repos

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"net/http"

	"github.com/supabase-community/supabase-go"
)

type GameRepo struct {
	db *supabase.Client
}

func NewGameRepo(db *supabase.Client) *GameRepo {
	return &GameRepo{db}
}

func (r *GameRepo) GetAllCreated(userId string) ([]models.Game, apierrors.StatusError) {
	var games []models.Game
	query := r.db.
		From("games").
		Select("*", "exact", false).
		Eq("created_by", userId)

	_, err := execeuteSelect(query, &games)
	if err != nil {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return games, nil
}

func (r *GameRepo) GetAllJoined(userId string) ([]models.Game, apierrors.StatusError) {
	var games []models.Game
	query := r.db.
		From("games").
		Select("id, name, state, created_by, join_code, game_players!inner(game_id, user_id)", "exact", false).
		Eq("game_players.user_id", userId)

	_, err := execeuteSelect(query, &games)
	if err != nil {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return games, nil
}

func (r *GameRepo) GetById(id string) (*models.Game, apierrors.StatusError) {
	var games []models.Game
	query := r.db.
		From("games").
		Select("*", "exact", false).
		Eq("id", id)

	count, err := execeuteSelect(query, &games)
	if err != nil {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	if count == 0 {
		return nil, apierrors.NewStatusError(http.StatusNotFound, errors.New("Game not found"))
	}

	return &games[0], nil
}

func (r *GameRepo) GetByJoinCode(joinCode string) (*models.Game, apierrors.StatusError) {
	var games []models.Game
	query := r.db.
		From("games").
		Select("*", "exact", false).
		Eq("join_code", joinCode)

	count, err := execeuteSelect(query, &games)
	if err != nil {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	if count == 0 {
		return nil, apierrors.NewStatusError(http.StatusNotFound, errors.New("Game not found"))
	}

	return &games[0], nil
}

func (r *GameRepo) Create(game *models.GameCreation) apierrors.StatusError {
	_, _, err := r.db.
		From("games").
		Insert(game, false, "", "", "").
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}

func (r *GameRepo) Patch(id string, patch *models.GamePatch) apierrors.StatusError {
	_, _, err := r.db.
		From("games").
		Update(patch, "", "").
		Eq("id", id).
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}

func (r *GameRepo) Delete(id string) apierrors.StatusError {
	_, _, err := r.db.
		From("games").
		Delete("", "").
		Eq("id", id).
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}
