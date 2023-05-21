package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type StatisticsCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewStatisticsCrm(config config.Config, cache cache.Cache) StatisticsCrm {
	return StatisticsCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (s StatisticsCrm) TicketOpenStat(ctx context.Context, userModel domain.User) ([]domain.HelpDesk, error) {
	query := s.generateTicketStatsQuery(userModel, "Open")
	return s.executeTicketStatsQuery(ctx, query)
}

func (s StatisticsCrm) TicketInProgressStat(ctx context.Context, userModel domain.User) ([]domain.HelpDesk, error) {
	query := s.generateTicketStatsQuery(userModel, "In Progress")
	return s.executeTicketStatsQuery(ctx, query)
}

func (s StatisticsCrm) TicketWaitForResponseStat(ctx context.Context, userModel domain.User) ([]domain.HelpDesk, error) {
	query := s.generateTicketStatsQuery(userModel, "Wait For Response")
	return s.executeTicketStatsQuery(ctx, query)
}

func (s StatisticsCrm) TicketClosedStat(ctx context.Context, userModel domain.User) ([]domain.HelpDesk, error) {
	query := s.generateTicketStatsQuery(userModel, "Closed")
	return s.executeTicketStatsQuery(ctx, query)
}

func (s StatisticsCrm) CalcTicketTotal(ctx context.Context, userModel domain.User) (int, error) {
	query := s.generateTicketsQuery(userModel, "")
	return s.vtiger.ExecuteCount(ctx, query)
}

func (s StatisticsCrm) CalcTicketOpen(ctx context.Context, userModel domain.User) (int, error) {
	query := s.generateTicketsQuery(userModel, "Open")
	return s.vtiger.ExecuteCount(ctx, query)
}

func (s StatisticsCrm) CalcTicketInProgress(ctx context.Context, userModel domain.User) (int, error) {
	query := s.generateTicketsQuery(userModel, "In Progress")
	return s.vtiger.ExecuteCount(ctx, query)
}

func (s StatisticsCrm) CalcTicketWaitForResponse(ctx context.Context, userModel domain.User) (int, error) {
	query := s.generateTicketsQuery(userModel, "Wait For Response")
	return s.vtiger.ExecuteCount(ctx, query)
}

func (s StatisticsCrm) CalcTicketClosed(ctx context.Context, userModel domain.User) (int, error) {
	query := s.generateTicketsQuery(userModel, "Closed")
	return s.vtiger.ExecuteCount(ctx, query)
}

func (s StatisticsCrm) generateTicketsQuery(userModel domain.User, status string) string {
	query := "SELECT COUNT(*) FROM HelpDesk WHERE parent_id = " + userModel.AccountId
	if status != "all" && status != "" && status != "total" {
		query += " AND ticketstatus = '" + status + "'"
	}
	query += ";"
	return query
}

func (s StatisticsCrm) generateTicketStatsQuery(userModel domain.User, status string) string {
	query := "SELECT hours, days FROM HelpDesk WHERE parent_id = " + userModel.AccountId
	if status != "" {
		query += " AND ticketstatus = '" + status + "'"
	}
	query += ";"
	return query
}

func (s StatisticsCrm) executeTicketStatsQuery(ctx context.Context, query string) ([]domain.HelpDesk, error) {
	result, err := s.vtiger.Query(ctx, query)
	tickets := make([]domain.HelpDesk, 0)
	if err != nil {
		return nil, e.Wrap("can not execute query "+query+", got error", err)
	}
	for _, data := range result.Result {
		ticket, err := domain.ConvertMapToHelpDesk(data)
		if err != nil {
			return nil, e.Wrap("can not convert map to helpdesk", err)
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}
