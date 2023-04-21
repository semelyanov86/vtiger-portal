package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"time"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Insert(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error)
	GetById(ctx context.Context, id int64) (domain.User, error)
}

type UsersCrm interface {
	FindByEmail(ctx context.Context, email string) ([]domain.User, error)
	RetrieveById(ctx context.Context, id string) (domain.User, error)
	ClearUserCodeField(ctx context.Context, id string) (domain.User, error)
}

type Tokens interface {
	Insert(ctx context.Context, token *domain.Token) error
	DeleteAllForUser(ctx context.Context, scope string, userId int64) error
	New(ctx context.Context, userId int64, ttl time.Duration, scope string) (*domain.Token, error)
}

type Managers interface {
	RetrieveById(ctx context.Context, id string) (domain.Manager, error)
}

type Modules interface {
	GetModuleInfo(ctx context.Context, module string) (vtiger.Module, error)
}

type Company interface {
	GetCompanyInfo(ctx context.Context) (domain.Company, error)
}

var ErrRecordNotFound = errors.New("record not found")
var ErrEditConflict = errors.New("edit conflict")
var ErrWrongCrmId = errors.New("wrong crm id")

type Repositories struct {
	Users    Users
	UsersCrm UsersCrm
	Tokens   *TokensRepo
	Managers Managers
	Modules  ModulesCrm
	Company  Company
}

func NewRepositories(db *sql.DB, config config.Config, cache cache.Cache) *Repositories {
	return &Repositories{
		Users:    NewUsersRepo(db),
		UsersCrm: NewUsersVtiger(config, cache),
		Tokens:   NewTokensRepo(db),
		Managers: NewManagersCrm(config, cache),
		Modules:  NewModulesCrm(config, cache),
		Company:  NewCompanyCrm(config, cache),
	}
}
