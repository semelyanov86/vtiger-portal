package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Insert(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type UsersCrm interface {
	FindByEmail(ctx context.Context, email string) ([]domain.User, error)
}

var ErrRecordNotFound = errors.New("record not found")
var ErrEditConflict = errors.New("edit conflict")
var ErrWrongCrmId = errors.New("wrong crm id")

type Repositories struct {
	Users    Users
	UsersCrm UsersCrm
}

func NewRepositories(db *sql.DB, config config.Config, cache cache.Cache) *Repositories {
	return &Repositories{
		Users:    NewUsersRepo(db),
		UsersCrm: NewUsersVtiger(config, cache),
	}
}
