package repos

import (
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"net/http"
	"strconv"

	"github.com/supabase-community/supabase-go"
)

type GameApprovalRepo struct {
	db *supabase.Client
}

func NewGameApprovalRepo(db *supabase.Client) *GameApprovalRepo {
	return &GameApprovalRepo{db}
}

func (r *GameApprovalRepo) GetAllByGameId(gameId int) ([]models.GameApproval, apierrors.StatusError) {
	var approvals []models.GameApproval
	query := r.db.
		From("game_approvals").
		Select("*", "exact", false).
		Eq("game_id", strconv.Itoa(gameId))

	_, err := execeuteSelect(query, &approvals)
	if err != nil {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return approvals, nil
}

func (r *GameApprovalRepo) Create(game *models.GameApproval) apierrors.StatusError {
	_, _, err := r.db.
		From("game_approvals").
		Insert(game, false, "", "", "").
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}

func (r *GameApprovalRepo) Patch(gameId int, userId string, patch *models.GameApproval) apierrors.StatusError {
	_, _, err := r.db.
		From("game_approvals").
		Update(patch, "", "").
		Eq("game_id", strconv.Itoa(gameId)).
		Eq("user_id", userId).
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}
