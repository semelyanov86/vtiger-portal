package service

import (
	"errors"
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
	Modules  ModulesService
	Company  Company
	HelpDesk HelpDesk
}

var ErrOperationNotPermitted = errors.New("you are not permitted to view this record")

func NewServices(repos repository.Repositories, email email.Sender, wg *sync.WaitGroup, config config.Config, cache cache.Cache) *Services {
	emailService := *NewEmailsService(email, config.Email, cache)
	companyService := NewCompanyService(repos.Company, cache)
	return &Services{
		Users:    NewUsersService(repos.Users, repos.UsersCrm, wg, emailService, companyService, repos.Tokens),
		Emails:   *NewEmailsService(email, config.Email, cache),
		Tokens:   NewTokensService(repos.Tokens, repos.Users, emailService, config, companyService),
		Context:  NewContextService(),
		Managers: NewManagerService(repos.Managers, cache),
		Modules:  NewModulesService(repos.Modules, cache),
		Company:  companyService,
		HelpDesk: NewHelpDeskService(repos.HelpDesk, cache),
	}
}

type ContextServiceInterface interface {
	ContextSetUser(c *gin.Context, user *domain.User) *gin.Context
	ContextGetUser(c *gin.Context) *domain.User
}
