package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/email"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"sync"
	"time"
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
	Faqs      Faqs
	Invoices  Invoices
}

var ErrOperationNotPermitted = errors.New("you are not permitted to view this record")

func NewServices(repos repository.Repositories, email email.Sender, wg *sync.WaitGroup, config config.Config, cache cache.Cache) *Services {
	emailService := *NewEmailsService(email, config.Email, cache)
	companyService := NewCompanyService(repos.Company, cache)
	commentsService := NewComments(repos.Comments, cache, config)
	documentService := NewDocuments(repos.Documents, cache)
	modulesService := NewModulesService(repos.Modules, cache)
	return &Services{
		Users:     NewUsersService(repos.Users, repos.UsersCrm, wg, emailService, companyService, repos.Tokens, repos.Documents),
		Emails:    *NewEmailsService(email, config.Email, cache),
		Tokens:    NewTokensService(repos.Tokens, repos.Users, emailService, config, companyService),
		Context:   NewContextService(),
		Managers:  NewManagerService(repos.Managers, cache),
		Modules:   modulesService,
		Company:   companyService,
		HelpDesk:  NewHelpDeskService(repos.HelpDesk, cache, commentsService, documentService, modulesService, config),
		Comments:  commentsService,
		Documents: documentService,
		Faqs:      NewFaqsService(repos.Faqs, cache, modulesService, config),
		Invoices:  NewInvoiceService(repos.Invoice, cache, modulesService, config),
	}
}

type ContextServiceInterface interface {
	ContextSetUser(c *gin.Context, user *domain.User) *gin.Context
	ContextGetUser(c *gin.Context) *domain.User
}

type CommentServiceInterface interface {
	GetRelated(ctx context.Context, id string) ([]domain.Comment, error)
	Create(ctx context.Context, content string, related string, userId string) (domain.Comment, error)
}

type DocumentServiceInterface interface {
	GetRelated(ctx context.Context, id string) ([]domain.Document, error)
	GetFile(ctx context.Context, id string, relatedId string) (vtiger.File, error)
}

type SupportedTypes interface {
	*domain.HelpDesk | *domain.Company | *domain.Manager | *vtiger.Module | *[]domain.Document | *domain.Invoice
}

func GetFromCache[T SupportedTypes](key string, dest T, c cache.Cache) error {
	cachedData, err := c.Get(key)
	if err != nil || cachedData == nil {
		return cache.ErrItemNotFound
	}

	err = json.Unmarshal(cachedData, dest)
	if err != nil {
		if jsonErr, ok := err.(*json.SyntaxError); ok {
			problemPart := cachedData[jsonErr.Offset-10 : jsonErr.Offset+10]
			return fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
		}
		return err
	}
	return nil
}

func StoreInCache[T SupportedTypes](key string, value T, ttl time.Duration, c cache.Cache) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(key, data, int64(ttl))
}
