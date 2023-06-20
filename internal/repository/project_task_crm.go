package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type ProjectTaskCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewProjectTaskCrm(config config.Config, cache cache.Cache) ProjectTaskCrm {
	return ProjectTaskCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (p ProjectTaskCrm) RetrieveById(ctx context.Context, id string) (domain.ProjectTask, error) {
	result, err := p.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.ProjectTask{}, e.Wrap("can not retrieve project task with id "+id+" got error", err)
	}
	return domain.ConvertMapToProjectTask(result.Result)
}

func (p ProjectTaskCrm) GetFromProject(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.ProjectTask, error) {
	items, err := p.vtiger.GetByWhereClause(ctx, filter, "projectid", filter.Parent, "ProjectTask")
	if err != nil {
		return nil, err
	}
	projectTasks := make([]domain.ProjectTask, 0, len(items))
	for _, data := range items {
		projectTask, err := domain.ConvertMapToProjectTask(data)
		if err != nil {
			return projectTasks, e.Wrap("can not convert map to projectTask", err)
		}
		projectTasks = append(projectTasks, projectTask)
	}
	return projectTasks, nil
}

func (p ProjectTaskCrm) Count(ctx context.Context, parent string) (int, error) {
	body := make(map[string]string)
	body["projectid"] = parent

	return p.vtiger.Count(ctx, "ProjectTask", body)
}

func (p ProjectTaskCrm) Create(ctx context.Context, task domain.ProjectTask) (domain.ProjectTask, error) {
	ticketMap, err := task.ConvertToMap()
	if err != nil {
		return task, e.Wrap("can not convert to map", err)
	}
	result, err := p.vtiger.Create(ctx, "ProjectTask", ticketMap)
	if err != nil {
		return task, e.Wrap("can not create task", err)
	}
	return domain.ConvertMapToProjectTask(result.Result)
}

func (p ProjectTaskCrm) Revise(ctx context.Context, task map[string]any) (domain.ProjectTask, error) {
	result, err := p.vtiger.Revise(ctx, task)
	if err != nil {
		return domain.ProjectTask{}, e.Wrap("can send update map to vtiger", err)
	}
	return domain.ConvertMapToProjectTask(result.Result)
}
