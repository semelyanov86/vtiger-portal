package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

const CacheHelpDeskTtl = 500

type HelpDesk struct {
	repository repository.HelpDesk
	cache      cache.Cache
	comment    CommentServiceInterface
	document   DocumentServiceInterface
}

func NewHelpDeskService(repository repository.HelpDesk, cache cache.Cache, comments CommentServiceInterface, document DocumentServiceInterface) HelpDesk {
	return HelpDesk{
		repository: repository,
		cache:      cache,
		comment:    comments,
		document:   document,
	}
}

func (h HelpDesk) GetHelpDeskById(ctx context.Context, id string) (domain.HelpDesk, error) {
	cachedTicketData, err := h.cache.Get(id)
	if errors.Is(cache.ErrItemNotFound, err) || cachedTicketData == nil {
		ticketData, err := h.retrieveHelpDesk(ctx, id)
		if err != nil {
			return domain.HelpDesk{}, e.Wrap("can not get a helpdesk", err)
		}
		cachedValue, err := json.Marshal(ticketData)
		if err != nil {
			return domain.HelpDesk{}, err
		}
		err = h.cache.Set(id, cachedValue, CacheHelpDeskTtl)
		if err != nil {
			return domain.HelpDesk{}, err
		}
		return ticketData, nil
	} else {
		decodedTicket := &domain.HelpDesk{}
		err = json.Unmarshal(cachedTicketData, decodedTicket)
		if err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := cachedTicketData[jsonErr.Offset-10 : jsonErr.Offset+10]

				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
			}
			return domain.HelpDesk{}, e.Wrap("can not convert caches data to helpdesk", err)
		}
		return *decodedTicket, nil
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
