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

type ProjectCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewProjectCrm(config config.Config, cache cache.Cache) ProjectCrm {
	return ProjectCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (p ProjectCrm) RetrieveById(ctx context.Context, id string) (domain.Project, error) {
	result, err := p.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.Project{}, e.Wrap("can not retrieve project with id "+id+" got error", err)
	}
	return domain.ConvertMapToProject(result.Result)
}

func (p ProjectCrm) GetAll(ctx context.Context, filter PaginationQueryFilter) ([]domain.Project, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM Project WHERE linktoaccountscontacts = " + filter.Client + " "
	if filter.Search == "" {
		query += " OR linktoaccountscontacts = " + filter.Contact + " "
	}

	sort := filter.Sort
	if sort == "" {
		sort = "-ticket_no"
	}
	if filter.Search != "" {
		query += " AND project_no LIKE '%" + filter.Search + "%' OR projectname LIKE '%" + filter.Search + "%' OR projecttype LIKE '%" + filter.Search + "%' "
	}
	query += GenerateOrderByClause(sort)
	query += " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"

	projects := make([]domain.Project, 0)
	result, err := p.vtiger.Query(ctx, query)
	if err != nil {
		return projects, e.Wrap("can not execute query "+query+", got error", err)
	}
	for _, data := range result.Result {
		project, err := domain.ConvertMapToProject(data)
		if err != nil {
			return projects, e.Wrap("can not convert map to project", err)
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (p ProjectCrm) Count(ctx context.Context, client string, contact string) (int, error) {
	body := make(map[string]string)
	body["linktoaccountscontacts"] = client
	body["_linktoaccountscontacts"] = contact

	return p.vtiger.Count(ctx, "Project", body)
}
