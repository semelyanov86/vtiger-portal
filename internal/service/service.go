package service

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/email/smtp"
	"sync"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Services struct {
	Users   UsersService
	Email   smtp.Mailer
	Tokens  TokensService
	Context ContextServiceInterface
}

func NewServices(repos repository.Repositories, email smtp.Mailer, wg *sync.WaitGroup) *Services {
	return &Services{
		Users:   NewUsersService(repos.Users, repos.UsersCrm, wg),
		Email:   email,
		Tokens:  NewTokensService(repos.Tokens, repos.Users),
		Context: NewContextService(),
	}
}

type ContextServiceInterface interface {
	ContextSetUser(c *gin.Context, user *domain.User) *gin.Context
	ContextGetUser(c *gin.Context) *domain.User
}
