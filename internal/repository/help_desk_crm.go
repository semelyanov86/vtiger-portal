package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strconv"
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

func (m HelpDeskCrm) GetAll(ctx context.Context, filter PaginationQueryFilter) ([]domain.HelpDesk, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM HelpDesk WHERE parent_id = " + filter.Client + " "
	sort := filter.Sort
	if sort == "" {
		sort = "-ticket_no"
	}
	if filter.Search != "" {
		query += " AND ticket_no LIKE '%" + filter.Search + "%' OR ticket_title LIKE '%" + filter.Search + "%' "
	}
	query += GenerateOrderByClause(sort)
	query += " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"
	tickets := make([]domain.HelpDesk, 0)
	result, err := m.vtiger.Query(ctx, query)
	if err != nil {
		return tickets, e.Wrap("can not execute query "+query+", got error", err)
	}
	for _, data := range result.Result {
		ticket, err := domain.ConvertMapToHelpDesk(data)
		if err != nil {
			return tickets, e.Wrap("can not convert map to helpdesk", err)
		}
		tickets = append(tickets, ticket)
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
