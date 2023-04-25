package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

const CacheHelpDeskTtl = 500

var ErrValidation = errors.New("validation error")

type HelpDesk struct {
	repository repository.HelpDesk
	cache      cache.Cache
	comment    CommentServiceInterface
	document   DocumentServiceInterface
	module     ModulesService
	config     config.Config
}

func NewHelpDeskService(repository repository.HelpDesk, cache cache.Cache, comments CommentServiceInterface, document DocumentServiceInterface, module ModulesService, config config.Config) HelpDesk {
	return HelpDesk{
		repository: repository,
		cache:      cache,
		comment:    comments,
		document:   document,
		module:     module,
		config:     config,
	}
}

func (h HelpDesk) GetHelpDeskById(ctx context.Context, id string) (domain.HelpDesk, error) {
	helpDesk := &domain.HelpDesk{}
	err := GetFromCache[*domain.HelpDesk](id, helpDesk, h.cache)
	if err == nil {
		return *helpDesk, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		helpDeskData, err := h.retrieveHelpDesk(ctx, id)
		if err != nil {
			return helpDeskData, e.Wrap("can not get a helpDesk", err)
		}
		err = StoreInCache[*domain.HelpDesk](id, &helpDeskData, CacheHelpDeskTtl, h.cache)
		if err != nil {
			return helpDeskData, err
		}
		return helpDeskData, nil
	} else {
		return *helpDesk, e.Wrap("can not convert caches data to helpDesk", err)
	}
}

func (h HelpDesk) retrieveHelpDesk(ctx context.Context, id string) (domain.HelpDesk, error) {
	return h.repository.RetrieveById(ctx, id)
}

func (h HelpDesk) GetRelatedComments(ctx context.Context, id string, companyId string) ([]domain.Comment, error) {
	helpDesk, err := h.GetHelpDeskById(ctx, id)
	if err != nil {
		return []domain.Comment{}, err
	}
	if helpDesk.ParentID != companyId {
		return []domain.Comment{}, ErrOperationNotPermitted
	}
	return h.comment.GetRelated(ctx, id)
}

func (h HelpDesk) GetRelatedDocuments(ctx context.Context, id string, companyId string) ([]domain.Document, error) {
	helpDesk, err := h.GetHelpDeskById(ctx, id)
	if err != nil {
		return []domain.Document{}, err
	}
	if helpDesk.ParentID != companyId {
		return []domain.Document{}, ErrOperationNotPermitted
	}
	return h.document.GetRelated(ctx, id)
}

func (h HelpDesk) GetAll(ctx context.Context, filter repository.TicketsQueryFilter) ([]domain.HelpDesk, int, error) {
	tickets, err := h.repository.GetAll(ctx, filter)
	if err != nil {
		return tickets, 0, err
	}
	count, err := h.repository.Count(ctx, filter.Client)
	return tickets, count, err
}

type CreateTicketInput struct {
	TicketTitle      string `json:"ticket_title" binding:"required"`
	Ticketpriorities string `json:"ticketpriorities" binding:"required"`
	Ticketseverities string `json:"ticketseverities"`
	Ticketcategories string `json:"ticketcategories"`
	Description      string `json:"description" binding:"required"`
}

func (h HelpDesk) CreateTicket(ctx context.Context, input CreateTicketInput, user domain.User) (domain.HelpDesk, error) {
	module, err := h.module.Describe(ctx, "HelpDesk")
	var helpDesk domain.HelpDesk
	if err != nil {
		return domain.HelpDesk{}, e.Wrap("can not get module info", err)
	}
	helpDesk.TicketTitle = input.TicketTitle
	helpDesk.AssignedUserID = h.config.Vtiger.Business.DefaultUser
	helpDesk.TicketPriorities = input.Ticketpriorities
	helpDesk.TicketSeverities = input.Ticketseverities
	helpDesk.TicketCategories = input.Ticketcategories
	helpDesk.Description = input.Description
	helpDesk.ParentID = user.AccountId
	helpDesk.ContactID = user.Crmid
	helpDesk.FromPortal = true
	helpDesk.Source = "PORTAL"
	var fields = module.Fields
	for _, field := range fields {
		switch field.Name {
		case "ticketstatus":
			helpDesk.TicketStatus = field.Type.DefaultValue
		case "ticketseverities":
			if !field.Type.IsPicklistExist(helpDesk.TicketSeverities) {
				return helpDesk, e.Wrap("Wrong value for field ticketseverities", ErrValidation)
			}
		case "ticketpriorities":
			if !field.Type.IsPicklistExist(helpDesk.TicketPriorities) {
				return helpDesk, e.Wrap("Wrong value for field ticketpriorities", ErrValidation)
			}
		case "ticketcategories":
			if !field.Type.IsPicklistExist(helpDesk.TicketCategories) {
				return helpDesk, e.Wrap("Wrong value for field ticketcategories", ErrValidation)
			}
		}
	}
	return h.repository.Create(ctx, helpDesk)
}
