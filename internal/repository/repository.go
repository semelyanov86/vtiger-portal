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
	SaveOtp(ctx context.Context, otpSecret string, otpUrl string, userId int64) error
	EnableAndVerifyOtp(ctx context.Context, userId int64) error
	VerifyOrInvalidateOtp(ctx context.Context, userId int64, valid bool) error
	DisableOtp(ctx context.Context, userId int64) error
	GetAllByAccountId(ctx context.Context, account string) ([]domain.User, error)
}

type UsersCrm interface {
	FindByEmail(ctx context.Context, email string) ([]domain.User, error)
	RetrieveById(ctx context.Context, id string) (domain.User, error)
	ClearUserCodeField(ctx context.Context, id string) (domain.User, error)
	FindContactsInAccount(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]string, error)
	Update(ctx context.Context, id string, user domain.User) (domain.User, error)
	RetrieveContactMap(ctx context.Context, id string) (map[string]any, error)
	ChangeSettingField(ctx context.Context, id string, field string, value bool) error
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

type HelpDesk interface {
	RetrieveById(ctx context.Context, id string) (domain.HelpDesk, error)
	GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.HelpDesk, error)
	Count(ctx context.Context, client string) (int, error)
	Create(ctx context.Context, ticket domain.HelpDesk) (domain.HelpDesk, error)
	Update(ctx context.Context, ticket domain.HelpDesk) (domain.HelpDesk, error)
	Revise(ctx context.Context, ticket map[string]any) (domain.HelpDesk, error)
}

type Comment interface {
	RetrieveFromModule(ctx context.Context, id string) ([]domain.Comment, error)
	Create(ctx context.Context, comment domain.Comment) (domain.Comment, error)
}

type Lead interface {
	Create(ctx context.Context, comment domain.Lead) (domain.Lead, error)
}

type Document interface {
	RetrieveFromModule(ctx context.Context, id string) ([]domain.Document, error)
	RetrieveFile(ctx context.Context, id string) (vtiger.File, error)
	AttachFile(ctx context.Context, doc domain.Document, parent string) (domain.Document, error)
	DeleteFile(ctx context.Context, id string) error
}

type Faq interface {
	GetAllFaqs(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.Faq, error)
	Count(ctx context.Context, client string) (int, error)
}

type Invoice interface {
	RetrieveById(ctx context.Context, id string) (domain.Invoice, error)
	GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.Invoice, error)
	Count(ctx context.Context, client string) (int, error)
	GetFromSalesOrder(ctx context.Context, soId string) ([]domain.Invoice, error)
}

type ServiceContract interface {
	RetrieveById(ctx context.Context, id string) (domain.ServiceContract, error)
	Count(ctx context.Context, client string, contact string) (int, error)
	GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.ServiceContract, error)
}

type Currency interface {
	RetrieveById(ctx context.Context, id string) (domain.Currency, error)
}

type Product interface {
	RetrieveById(ctx context.Context, id string) (domain.Product, error)
	GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.Product, error)
	Count(ctx context.Context, filters map[string]any) (int, error)
}

type Service interface {
	RetrieveById(ctx context.Context, id string) (domain.Service, error)
	GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.Service, error)
	Count(ctx context.Context, filters map[string]any) (int, error)
}

type Project interface {
	RetrieveById(ctx context.Context, id string) (domain.Project, error)
	GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.Project, error)
	Count(ctx context.Context, client string, contact string) (int, error)
}

type ProjectTask interface {
	RetrieveById(ctx context.Context, id string) (domain.ProjectTask, error)
	GetFromProject(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.ProjectTask, error)
	Count(ctx context.Context, parent string) (int, error)
	Create(ctx context.Context, task domain.ProjectTask) (domain.ProjectTask, error)
	Revise(ctx context.Context, task map[string]any) (domain.ProjectTask, error)
}

type Account interface {
	RetrieveById(ctx context.Context, id string) (domain.Account, error)
}

var ErrRecordNotFound = errors.New("record not found")
var ErrEditConflict = errors.New("edit conflict")
var ErrWrongCrmId = errors.New("wrong crm id")
var ErrCanNotParseCountObject = errors.New("can not parse count object")

type Repositories struct {
	Users            Users
	UsersCrm         UsersCrm
	Tokens           *TokensRepo
	Managers         Managers
	Modules          ModulesCrm
	Company          Company
	HelpDesk         HelpDesk
	Comments         Comment
	Documents        Document
	Faqs             Faq
	Invoice          Invoice
	SalesOrder       SalesOrderCrm
	ServiceContract  ServiceContract
	Currency         CurrencyCrm
	Product          ProductCrm
	Service          ServicesCrm
	Projects         ProjectCrm
	ProjectTasks     ProjectTaskCrm
	Statistics       StatisticsCrm
	Leads            LeadCrm
	Account          AccountCrm
	Search           SearchCrm
	Payment          *PaymentsRepo
	Notifications    *NotificationsRepo
	NotificationsCrm NotificationsCrm
	CustomModule     CustomModuleCrm
}

func NewRepositories(db *sql.DB, config config.Config, cache cache.Cache) *Repositories {
	return &Repositories{
		Users:            NewUsersRepo(db),
		UsersCrm:         NewUsersVtiger(config, cache),
		Tokens:           NewTokensRepo(db),
		Managers:         NewManagersCrm(config, cache),
		Modules:          NewModulesCrm(config, cache),
		Company:          NewCompanyCrm(config, cache),
		HelpDesk:         NewHelpDeskCrm(config, cache),
		Comments:         NewCommentCrm(config, cache),
		Documents:        NewDocumentCrm(config, cache),
		Faqs:             NewFaqsCrm(config, cache),
		Invoice:          NewInvoiceCrm(config, cache),
		SalesOrder:       NewSalesOrderCrm(config, cache),
		ServiceContract:  NewServiceContractCrm(config, cache),
		Currency:         NewCurrencyCrm(config, cache),
		Product:          NewProductCrm(config, cache),
		Service:          NewServicesCRM(config, cache),
		Projects:         NewProjectCrm(config, cache),
		ProjectTasks:     NewProjectTaskCrm(config, cache),
		Statistics:       NewStatisticsCrm(config, cache),
		Leads:            NewLeadCrm(config, cache),
		Account:          NewAccountCrm(config, cache),
		Search:           NewSearchCrm(config, cache),
		Payment:          NewPaymentsRepo(db),
		Notifications:    NewNotificationsRepo(db),
		NotificationsCrm: NewNotificationsCrm(config, cache),
		CustomModule:     NewCustomModuleCrm(config, cache),
	}
}
