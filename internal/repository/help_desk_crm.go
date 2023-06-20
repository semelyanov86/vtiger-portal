package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type HelpDeskCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewHelpDeskCrm(config config.Config, cache cache.Cache) HelpDeskCrm {
	return HelpDeskCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (m HelpDeskCrm) RetrieveById(ctx context.Context, id string) (domain.HelpDesk, error) {
	result, err := m.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.HelpDesk{}, e.Wrap("can not retrieve help desk with id "+id+" got error:"+result.Error.Message, err)
	}
	return domain.ConvertMapToHelpDesk(result.Result)
}

func (m HelpDeskCrm) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.HelpDesk, error) {
	items, err := m.vtiger.GetAll(ctx, filter, vtiger.QueryFieldsProps{
		DefaultSort:  "-ticket_no",
		SearchFields: []string{"ticket_title", "ticket_no"},
		ClientField:  "",
		AccountField: "parent_id",
		TableName:    "HelpDesk",
	})
	if err != nil {
		return nil, err
	}

	tickets := make([]domain.HelpDesk, 0, len(items))

	for _, data := range items {
		ticket, err := domain.ConvertMapToHelpDesk(data)
		if err != nil {
			return tickets, e.Wrap("can not convert map to helpdesk", err)
		}
		if ticket.ParentID == filter.Client {
			tickets = append(tickets, ticket)
		}
	}
	return tickets, nil
}

func (m HelpDeskCrm) Count(ctx context.Context, client string) (int, error) {
	body := make(map[string]string)
	body["parent_id"] = client
	return m.vtiger.Count(ctx, "HelpDesk", body)
}

func (m HelpDeskCrm) Create(ctx context.Context, ticket domain.HelpDesk) (domain.HelpDesk, error) {
	ticketMap, err := ticket.ConvertToMap()
	if err != nil {
		return ticket, e.Wrap("can not convert to map", err)
	}
	result, err := m.vtiger.Create(ctx, "HelpDesk", ticketMap)
	if err != nil {
		return ticket, e.Wrap("can not create user", err)
	}
	return domain.ConvertMapToHelpDesk(result.Result)
}

func (m HelpDeskCrm) Update(ctx context.Context, ticket domain.HelpDesk) (domain.HelpDesk, error) {
	ticketMap, err := ticket.ConvertToMap()
	if err != nil {
		return ticket, e.Wrap("can not convert to map", err)
	}

	result, err := m.vtiger.Update(ctx, ticketMap)
	if err != nil {
		return ticket, e.Wrap("can send update map to vtiger", err)
	}
	return domain.ConvertMapToHelpDesk(result.Result)
}

func (m HelpDeskCrm) Revise(ctx context.Context, ticket map[string]any) (domain.HelpDesk, error) {
	result, err := m.vtiger.Revise(ctx, ticket)
	if err != nil {
		return domain.HelpDesk{}, e.Wrap("can send update map to vtiger", err)
	}
	return domain.ConvertMapToHelpDesk(result.Result)
}
