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
	"mime/multipart"
	"sync"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Services struct {
	Users            UsersService
	Emails           EmailService
	Tokens           TokensService
	Context          ContextServiceInterface
	Managers         ManagerService
	Modules          ModulesService
	Company          Company
	HelpDesk         HelpDesk
	Comments         Comments
	Documents        DocumentServiceInterface
	Faqs             Faqs
	Invoices         Invoices
	ServiceContracts ServiceContracts
	Currencies       CurrencyService
	Products         ProductService
	Services         ServicesService
	Projects         ProjectsService
	ProjectTasks     ProjectTasksService
	Statistics       StatisticsService
	Leads            Leads
	Auth             AuthService
	Accounts         AccountService
}

var ErrOperationNotPermitted = errors.New("you are not permitted to view this record")

func NewServices(repos repository.Repositories, email email.Sender, wg *sync.WaitGroup, config config.Config, cache cache.Cache) *Services {
	emailService := *NewEmailsService(email, config.Email, cache)
	companyService := NewCompanyService(repos.Company, cache)
	managersService := NewManagerService(repos.Managers, cache)
	accountService := NewAccountService(repos.Account, cache)
	usersService := NewUsersService(repos.Users, repos.UsersCrm, wg, emailService, companyService, repos.Tokens, repos.Documents, cache, accountService)
	commentsService := NewComments(repos.Comments, cache, config, usersService, managersService)
	documentService := NewDocuments(repos.Documents, cache, config)
	modulesService := NewModulesService(repos.Modules, cache)
	currencyService := NewCurrencyService(repos.Currency, cache)
	projectService := NewProjectsService(repos.Projects, cache, commentsService, documentService, modulesService, config, repos.ProjectTasks)
	return &Services{
		Users:            usersService,
		Auth:             NewAuthService(repos.Users, wg, cache, config),
		Emails:           *NewEmailsService(email, config.Email, cache),
		Tokens:           NewTokensService(repos.Tokens, repos.Users, emailService, config, companyService),
		Context:          NewContextService(),
		Managers:         managersService,
		Modules:          modulesService,
		Company:          companyService,
		HelpDesk:         NewHelpDeskService(repos.HelpDesk, cache, commentsService, documentService, modulesService, config),
		Comments:         commentsService,
		Documents:        documentService,
		Faqs:             NewFaqsService(repos.Faqs, cache, modulesService, config),
		Invoices:         NewInvoiceService(repos.Invoice, cache, modulesService, config, currencyService),
		ServiceContracts: NewServiceContractsService(repos.ServiceContract, cache, documentService, modulesService, config),
		Currencies:       currencyService,
		Products:         NewProductService(repos.Product, cache, currencyService, repos.Documents, modulesService, config),
		Services:         NewServicesService(repos.Service, cache, currencyService, modulesService, config),
		Projects:         projectService,
		ProjectTasks:     NewProjectTasksService(repos.ProjectTasks, cache, commentsService, documentService, modulesService, config, projectService),
		Statistics:       NewStatisticsService(repos.Statistics, cache),
		Leads:            NewLeads(repos.Leads, config),
		Accounts:         accountService,
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
	AttachFile(ctx context.Context, file multipart.File, id string, userModel domain.User, header *multipart.FileHeader) (domain.Document, error)
	DeleteFile(ctx context.Context, id string, related string) error
}

type SupportedTypes interface {
	*domain.HelpDesk | *domain.Company | *domain.Manager | *vtiger.Module | *[]domain.Document | *domain.Invoice | *domain.ServiceContract | *domain.Currency | *domain.Product | *domain.Service | *domain.Project | *domain.ProjectTask | *[]domain.User | *domain.Statistics | *domain.User | *domain.Account
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
