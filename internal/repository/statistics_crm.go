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

func (s StatisticsCrm) CalcProjectsTotal(ctx context.Context, userModel domain.User) (int, error) {
	query := s.generateProjectsQuery(userModel, "")
	return s.vtiger.ExecuteCount(ctx, query)
}

func (s StatisticsCrm) CalcProjectsOpen(ctx context.Context, userModel domain.User) (int, error) {
	query := s.generateProjectsQuery(userModel, "open")
	return s.vtiger.ExecuteCount(ctx, query)
}

func (s StatisticsCrm) CalcProjectsClosed(ctx context.Context, userModel domain.User) (int, error) {
	query := s.generateProjectsQuery(userModel, "closed")
	return s.vtiger.ExecuteCount(ctx, query)
}

func (s StatisticsCrm) InvoicesOpenStat(ctx context.Context, userModel domain.User) ([]domain.Invoice, error) {
	return s.executeInvoiceStatsQuery(ctx, s.generateInvoiceStatsQuery(userModel, "Open"))
}

func (s StatisticsCrm) InvoicesClosedStat(ctx context.Context, userModel domain.User) ([]domain.Invoice, error) {
	return s.executeInvoiceStatsQuery(ctx, s.generateInvoiceStatsQuery(userModel, "Closed"))
}

func (s StatisticsCrm) InvoicesTotalStat(ctx context.Context, userModel domain.User) ([]domain.Invoice, error) {
	return s.executeInvoiceStatsQuery(ctx, s.generateInvoiceStatsQuery(userModel, ""))
}

func (s StatisticsCrm) TasksFromInProgressProjects(ctx context.Context, userModel domain.User) ([]domain.ProjectTask, error) {
	projects, err := s.getInProgressProjects(ctx, userModel)
	tasks := make([]domain.ProjectTask, 0)
	if err != nil {
		return tasks, e.Wrap("can not receive in progress projects", err)
	}
	for _, project := range projects {
		query := "SELECT * FROM ProjectTask WHERE projectid = " + project.Id + ";"
		pt, err := executeQuery[domain.ProjectTask](ctx, query, s.vtiger, domain.ConvertMapToProjectTask)

		if err != nil {
			return tasks, e.Wrap("can not retrieve project tasks", err)
		}
		tasks = append(tasks, pt...)
	}

	return tasks, err
}

func (s StatisticsCrm) generateTicketsQuery(userModel domain.User, status string) string {
	query := "SELECT COUNT(*) FROM HelpDesk WHERE parent_id = " + userModel.AccountId
	if status != "all" && status != "" && status != "total" {
		query += " AND ticketstatus = '" + status + "'"
	}
	query += ";"
	return query
}

func (s StatisticsCrm) generateProjectsQuery(userModel domain.User, status string) string {
	query := "SELECT COUNT(*) FROM Project WHERE linktoaccountscontacts = " + userModel.Crmid + " OR linktoaccountscontacts = " + userModel.AccountId + " "
	if status == "" {
		query += ";"
	} else if status == "open" {
		query += " AND projectstatus IN ('prospecting', 'initiated', 'in progress', 'waiting for feedback');"
	} else if status == "closed" {
		query += " AND projectstatus IN ('completed', 'delivered');"
	}

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
	return executeQuery[domain.HelpDesk](ctx, query, s.vtiger, domain.ConvertMapToHelpDesk)
}

func (s StatisticsCrm) executeInvoiceStatsQuery(ctx context.Context, query string) ([]domain.Invoice, error) {
	return executeQuery[domain.Invoice](ctx, query, s.vtiger, domain.ConvertMapToInvoice)
}

func (s StatisticsCrm) getInProgressProjects(ctx context.Context, userModel domain.User) ([]domain.Project, error) {
	query := "SELECT projectname FROM Project WHERE linktoaccountscontacts = " + userModel.Crmid + " OR linktoaccountscontacts = " + userModel.AccountId + " AND projectstatus IN ('in progress', 'Выполняется');"
	return executeQuery[domain.Project](ctx, query, s.vtiger, domain.ConvertMapToProject)
}

func executeQuery[T domain.HelpDesk | domain.Invoice | domain.Project | domain.ProjectTask](ctx context.Context, query string, c vtiger.VtigerConnector, fn func(map[string]any) (T, error)) ([]T, error) {
	result, err := c.Query(ctx, query)
	tickets := make([]T, 0)
	if err != nil {
		return nil, e.Wrap("can not execute query "+query+", got error", err)
	}
	for _, data := range result.Result {
		ticket, err := fn(data)
		if err != nil {
			return nil, e.Wrap("can not convert map to type", err)
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func (s StatisticsCrm) generateInvoiceStatsQuery(userModel domain.User, status string) string {
	query := "SELECT hdnGrandTotal FROM Invoice WHERE account_id = " + userModel.AccountId
	if status == "Open" {
		query += " AND invoicestatus IN ('Created', 'Approved', 'Sent')"
	} else if status == "Closed" {
		query += " AND invoicestatus IN ('Paid')"
	}

	query += ";"
	return query
}
