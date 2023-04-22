package repository

import (
	"context"
	"errors"
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

type TicketsQueryFilter struct {
	Page     int
	PageSize int
	Client   string
}

var ErrCanNotParseCountObject = errors.New("can not parse count object")

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

func (m HelpDeskCrm) GetAll(ctx context.Context, filter TicketsQueryFilter) ([]domain.HelpDesk, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM HelpDesk WHERE parent_id = " + filter.Client + " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"
	tickets := make([]domain.HelpDesk, 0)
	result, err := m.vtiger.Query(ctx, query)
	if err != nil {
		return tickets, e.Wrap("can not execute query "+query+", got error: "+result.Error.Message, err)
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
	query := "SELECT COUNT(*) FROM HelpDesk WHERE parent_id = " + client + ";"
	result, err := m.vtiger.Query(ctx, query)
	if err != nil {
		return 0, e.Wrap("can not execute query "+query+", got error: "+result.Error.Message, err)
	}
	countObject := result.Result[0]
	if countObject == nil {
		return 0, ErrCanNotParseCountObject
	}
	count, err := strconv.Atoi(countObject["count"].(string))
	if err != nil {
		return 0, err
	}
	return count, nil
}
