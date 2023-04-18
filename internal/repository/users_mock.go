package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
)

type UsersMock struct {
}

func NewUsersMock() *UsersMock {
	return &UsersMock{}
}

func (r *UsersMock) Insert(ctx context.Context, user *domain.User) error {
	user.Id = 1
	user.Version = 1
	return nil
}

func (r *UsersMock) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return domain.User{}, ErrRecordNotFound
}

func (r *UsersMock) Update(ctx context.Context, user *domain.User) error {
	return nil
}

func (r *UsersMock) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	var user domain.User

	return &user, nil
}
