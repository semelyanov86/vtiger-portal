package repository

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Insert(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

var ErrRecordNotFound = errors.New("record not found")
var ErrEditConflict = errors.New("edit conflict")
