package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
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

func (h HelpDesk) GetHelpDeskById(ctx context.Context, id string, userModel domain.User) (domain.HelpDesk, error) {
	helpDesk := &domain.HelpDesk{}
	err := GetFromCache[*domain.HelpDesk](id, helpDesk, h.cache)
	if err == nil {
		if userModel.AccountId != helpDesk.ParentID {
			return *helpDesk, ErrOperationNotPermitted
		}
		return *helpDesk, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		helpDeskData, err := h.retrieveHelpDesk(ctx, id)
		if err != nil {
			return helpDeskData, e.Wrap("can not get a helpDesk", err)
		}
		if userModel.AccountId != helpDeskData.ParentID {
			return helpDeskData, ErrOperationNotPermitted
		}
		err = StoreInCache[*domain.HelpDesk](id, &helpDeskData, CacheHelpDeskTtl, h.cache)
		if err != nil {
			return helpDeskData, err
		}
		return helpDeskData, nil
	} else {
		if userModel.AccountId != helpDesk.ParentID {
			return *helpDesk, ErrOperationNotPermitted
		}
		return *helpDesk, e.Wrap("can not convert caches data to helpDesk", err)
	}
}

func (h HelpDesk) retrieveHelpDesk(ctx context.Context, id string) (domain.HelpDesk, error) {
	return h.repository.RetrieveById(ctx, id)
}

func (h HelpDesk) GetRelatedComments(ctx context.Context, id string, userModel domain.User) ([]domain.Comment, error) {
	_, err := h.GetHelpDeskById(ctx, id, userModel)
	if err != nil {
		return []domain.Comment{}, err
	}

	return h.comment.GetRelated(ctx, id)
}

func (h HelpDesk) AddComment(ctx context.Context, content string, related string, userModel domain.User) (domain.Comment, error) {
	_, err := h.GetHelpDeskById(ctx, related, userModel)
	if err != nil {
		return domain.Comment{}, err
	}
	return h.comment.Create(ctx, content, related, userModel.Crmid)
}

func (h HelpDesk) GetRelatedDocuments(ctx context.Context, id string, userModel domain.User) ([]domain.Document, error) {
	_, err := h.GetHelpDeskById(ctx, id, userModel)
	if err != nil {
		return []domain.Document{}, err
	}

	return h.document.GetRelated(ctx, id)
}

func (h HelpDesk) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.HelpDesk, int, error) {
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
	var helpDesk domain.HelpDesk

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

	err := h.validateInputFields(ctx, &helpDesk)
	if err != nil {
		return helpDesk, err
	}

	return h.repository.Create(ctx, helpDesk)
}

func (h HelpDesk) validateInputFields(ctx context.Context, helpDesk *domain.HelpDesk) error {
	module, err := h.module.Describe(ctx, "HelpDesk")
	if err != nil {
		return e.Wrap("can not get module info", err)
	}
	var fields = module.Fields
	for _, field := range fields {
		switch field.Name {
		case "ticketstatus":
			helpDesk.TicketStatus = field.Type.DefaultValue
		case "ticketseverities":
			if !field.Type.IsPicklistExist(helpDesk.TicketSeverities) {
				return e.Wrap("Wrong value for field ticketseverities", ErrValidation)
			}
		case "ticketpriorities":
			if !field.Type.IsPicklistExist(helpDesk.TicketPriorities) {
				return e.Wrap("Wrong value for field ticketpriorities", ErrValidation)
			}
		case "ticketcategories":
			if !field.Type.IsPicklistExist(helpDesk.TicketCategories) {
				return e.Wrap("Wrong value for field ticketcategories", ErrValidation)
			}
		}
	}
	return nil
}

func (h HelpDesk) UpdateTicket(ctx context.Context, input CreateTicketInput, id string, user domain.User) (domain.HelpDesk, error) {
	ticket, err := h.retrieveHelpDesk(ctx, id)
	if err != nil {
		return ticket, e.Wrap("can not retrieve helpdesk during update", err)
	}
	if user.AccountId != ticket.ParentID {
		return domain.HelpDesk{}, ErrOperationNotPermitted
	}
	if input.TicketTitle != "" {
		ticket.TicketTitle = input.TicketTitle
	}
	if input.Description != "" {
		ticket.Description = input.Description
	}
	if input.Ticketpriorities != "" {
		ticket.TicketPriorities = input.Ticketpriorities
	}
	if input.Ticketseverities != "" {
		ticket.TicketSeverities = input.Ticketseverities
	}

	err = h.validateInputFields(ctx, &ticket)
	if err != nil {
		return ticket, err
	}
	ticket, err = h.repository.Update(ctx, ticket)
	if err != nil {
		return ticket, err
	}
	err = StoreInCache[*domain.HelpDesk](id, &ticket, CacheHelpDeskTtl, h.cache)
	return ticket, err
}

func (h HelpDesk) Revise(ctx context.Context, input map[string]any, id string, user domain.User) (domain.HelpDesk, error) {
	ticket, err := h.retrieveHelpDesk(ctx, id)
	if err != nil {
		return ticket, e.Wrap("can not retrieve helpdesk during update", err)
	}
	if user.AccountId != ticket.ParentID {
		return domain.HelpDesk{}, ErrOperationNotPermitted
	}
	input["id"] = id

	ticket, err = h.repository.Revise(ctx, input)
	if err != nil {
		return ticket, err
	}
	err = StoreInCache[*domain.HelpDesk](id, &ticket, CacheHelpDeskTtl, h.cache)
	return ticket, err
}
