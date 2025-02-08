package repos

import (
	"context"
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"net/http"
	"strconv"

	"github.com/supabase-community/supabase-go"
)

type GameRepo struct {
	db *supabase.Client
}

func NewGameRepo(db *supabase.Client) *GameRepo {
	return &GameRepo{db}
}

func (r *GameRepo) GetAll(ctx context.Context) ([]models.Game, apierrors.StatusError) {
	var games []models.Game
	query := r.db.
		From("games").
		Select("*", "exact", false)

	_, err := execeuteSelect(query, &games)
	if err != nil {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return games, nil
}

func (r *GameRepo) GetById(id int) (*models.Game, apierrors.StatusError) {
	var games []models.Game
	query := r.db.
		From("games").
		Select("*", "exact", false).
		Eq("id", strconv.Itoa(id))

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

func (r *GameRepo) Patch(id int, patch *models.GamePatch) apierrors.StatusError {
	_, _, err := r.db.
		From("games").
		Update(patch, "", "").
		Eq("id", strconv.Itoa(id)).
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}

func (r *GameRepo) Delete(id int) apierrors.StatusError {
	_, _, err := r.db.
		From("games").
		Delete("", "").
		Eq("id", strconv.Itoa(id)).
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}
