package service

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/email"
	"sync"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Services struct {
	Users    UsersService
	Emails   EmailService
	Tokens   TokensService
	Context  ContextServiceInterface
	Managers ManagerService
}

func NewServices(repos repository.Repositories, email email.Sender, wg *sync.WaitGroup, config config.Config, cache cache.Cache) *Services {
	emailService := *NewEmailsService(email, config.Email, cache)
	return &Services{
		Users:    NewUsersService(repos.Users, repos.UsersCrm, wg, emailService),
		Emails:   *NewEmailsService(email, config.Email, cache),
		Tokens:   NewTokensService(repos.Tokens, repos.Users, emailService, config),
		Context:  NewContextService(),
		Managers: NewManagerService(repos.Managers, cache),
	}
}

type ContextServiceInterface interface {
	ContextSetUser(c *gin.Context, user *domain.User) *gin.Context
	ContextGetUser(c *gin.Context) *domain.User
}
