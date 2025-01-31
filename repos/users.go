package repos

import (
	"errors"
	"mognjen/gossassins/apierrors"
	"mognjen/gossassins/models"
	"net/http"

	"github.com/supabase-community/supabase-go"
)

type UserRepo struct {
	db *supabase.Client
}

func NewUserRepo(db *supabase.Client) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) GetById(id string) (*models.User, apierrors.StatusError) {
	var users []models.User
	query := r.db.
		From("users").
		Select("*", "exact", false).
		Eq("id", id)

	count, err := execeuteSelect(query, &users)

	if err != nil {
		return nil, apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	if count == 0 {
		return nil, apierrors.NewStatusError(http.StatusNotFound, errors.New("User not found"))
	}

	return &users[0], nil
}

func (r *UserRepo) Create(user *models.User) apierrors.StatusError {
	_, _, err := r.db.
		From("users").
		Insert(user, false, "", "", "").
		Execute()

	if err != nil {
		return apierrors.NewStatusError(http.StatusInternalServerError, err)
	}

	return nil
}
