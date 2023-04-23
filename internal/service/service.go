package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/email"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"sync"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Services struct {
	Users     UsersService
	Emails    EmailService
	Tokens    TokensService
	Context   ContextServiceInterface
	Managers  ManagerService
	Modules   ModulesService
	Company   Company
	HelpDesk  HelpDesk
	Comments  Comments
	Documents DocumentServiceInterface
}

var ErrOperationNotPermitted = errors.New("you are not permitted to view this record")

func NewServices(repos repository.Repositories, email email.Sender, wg *sync.WaitGroup, config config.Config, cache cache.Cache) *Services {
	emailService := *NewEmailsService(email, config.Email, cache)
	companyService := NewCompanyService(repos.Company, cache)
	commentsService := NewComments(repos.Comments, cache)
	documentService := NewDocuments(repos.Documents, cache)
	return &Services{
		Users:     NewUsersService(repos.Users, repos.UsersCrm, wg, emailService, companyService, repos.Tokens, repos.Documents),
		Emails:    *NewEmailsService(email, config.Email, cache),
		Tokens:    NewTokensService(repos.Tokens, repos.Users, emailService, config, companyService),
		Context:   NewContextService(),
		Managers:  NewManagerService(repos.Managers, cache),
		Modules:   NewModulesService(repos.Modules, cache),
		Company:   companyService,
		HelpDesk:  NewHelpDeskService(repos.HelpDesk, cache, commentsService, documentService),
		Comments:  commentsService,
		Documents: documentService,
	}
}

type ContextServiceInterface interface {
	ContextSetUser(c *gin.Context, user *domain.User) *gin.Context
	ContextGetUser(c *gin.Context) *domain.User
}

type CommentServiceInterface interface {
	GetRelated(ctx context.Context, id string) ([]domain.Comment, error)
}

type DocumentServiceInterface interface {
	GetRelated(ctx context.Context, id string) ([]domain.Document, error)
	GetFile(ctx context.Context, id string, relatedId string) (vtiger.File, error)
}
