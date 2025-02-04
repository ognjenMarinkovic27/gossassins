package repos

import (
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"net/http"
	"strconv"

	"github.com/supabase-community/supabase-go"
)

type JoinRequestRepo struct {
	db *supabase.Client
}

func NewJoinRequestRepo(db *supabase.Client) *JoinRequestRepo {
	return &JoinRequestRepo{db}
}

func (r *JoinRequestRepo) GetAllByGameId(gameId int) ([]models.JoinRequest, apierrors.StatusError) {
	var joinRequests []models.JoinRequest
	query := r.db.
		From("join_requests").
		Select("*", "exact", false).
		Eq("game_id", strconv.Itoa(gameId))

	_, err := execeuteSelect(query, &joinRequests)
	if err != nil {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return joinRequests, nil
}

func (r *JoinRequestRepo) Create(game *models.JoinRequest) apierrors.StatusError {
	_, _, err := r.db.
		From("join_requests").
		Insert(game, false, "", "", "").
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}

func (r *JoinRequestRepo) Patch(gameId int, userId string, patch *models.JoinRequest) apierrors.StatusError {
	_, _, err := r.db.
		From("join_requests").
		Update(patch, "", "").
		Eq("game_id", strconv.Itoa(gameId)).
		Eq("user_id", userId).
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}
