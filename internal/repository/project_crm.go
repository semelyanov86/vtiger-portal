package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
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

func (p ProjectCrm) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.Project, error) {
	items, err := p.vtiger.GetAll(ctx, filter, vtiger.QueryFieldsProps{
		DefaultSort:  "-project_no",
		SearchFields: []string{"projectname", "project_no", "projecttype"},
		ClientField:  "linktoaccountscontacts",
		AccountField: "linktoaccountscontacts",
		TableName:    "Project",
	})

	if err != nil {
		return nil, err
	}

	projects := make([]domain.Project, 0, len(items))

	for _, data := range items {
		project, err := domain.ConvertMapToProject(data)
		if err != nil {
			return projects, e.Wrap("can not convert map to project", err)
		}
		if project.Linktoaccountscontacts == filter.Client || project.Linktoaccountscontacts == filter.Contact {
			projects = append(projects, project)
		}
	}
	return projects, nil
}

func (p ProjectCrm) Count(ctx context.Context, client string, contact string) (int, error) {
	body := make(map[string]string)
	body["linktoaccountscontacts"] = client
	body["_linktoaccountscontacts"] = contact

	return p.vtiger.Count(ctx, "Project", body)
}
