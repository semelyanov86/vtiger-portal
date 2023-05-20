package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

const CacheProjectTaskTtl = 5000

type ProjectTasksService struct {
	repository repository.ProjectTask
	cache      cache.Cache
	document   DocumentServiceInterface
	comment    CommentServiceInterface
	module     ModulesService
	config     config.Config
	project    ProjectsService
}

func NewProjectTasksService(repository repository.ProjectTask, cache cache.Cache, comments CommentServiceInterface, document DocumentServiceInterface, module ModulesService, config config.Config, project ProjectsService) ProjectTasksService {
	return ProjectTasksService{
		repository: repository,
		cache:      cache,
		document:   document,
		module:     module,
		config:     config,
		comment:    comments,
		project:    project,
	}
}

func (p ProjectTasksService) GetProjectTaskById(ctx context.Context, id string) (domain.ProjectTask, error) {
	projectTask := &domain.ProjectTask{}
	err := GetFromCache[*domain.ProjectTask](id, projectTask, p.cache)
	if err == nil {
		return *projectTask, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		projectTaskData, err := p.repository.RetrieveById(ctx, id)
		if err != nil {
			return projectTaskData, e.Wrap("can not get a projectTask", err)
		}
		err = StoreInCache[*domain.ProjectTask](id, &projectTaskData, CacheProjectTaskTtl, p.cache)
		if err != nil {
			return projectTaskData, err
		}
		return projectTaskData, nil
	} else {
		return *projectTask, e.Wrap("can not convert caches data to projectTask", err)
	}
}

func (p ProjectTasksService) GetAllFromProject(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.ProjectTask, int, error) {
	err := p.validateProjectPermissions(ctx, filter.Parent, filter.Client, filter.Contact)
	if err != nil {
		return nil, 0, err
	}
	projects, err := p.repository.GetFromProject(ctx, filter)
	if err != nil {
		return projects, 0, err
	}
	count, err := p.repository.Count(ctx, filter.Parent)
	return projects, count, err
}

func (p ProjectTasksService) GetRelatedComments(ctx context.Context, id string, companyId string, contactId string) ([]domain.Comment, error) {
	projectTask, err := p.GetProjectTaskById(ctx, id)
	if err != nil {
		return []domain.Comment{}, err
	}
	err = p.validateProjectPermissions(ctx, projectTask.Projectid, companyId, contactId)
	if err != nil {
		return []domain.Comment{}, err
	}

	return p.comment.GetRelated(ctx, id)
}

func (p ProjectTasksService) GetRelatedDocuments(ctx context.Context, id string, companyId string, contactId string) ([]domain.Document, error) {
	projectTask, err := p.GetProjectTaskById(ctx, id)
	if err != nil {
		return []domain.Document{}, err
	}
	err = p.validateProjectPermissions(ctx, projectTask.Projectid, companyId, contactId)
	if err != nil {
		return []domain.Document{}, err
	}

	return p.document.GetRelated(ctx, id)
}

func (p ProjectTasksService) AddComment(ctx context.Context, content string, related string, userModel domain.User) (domain.Comment, error) {
	projectTask, err := p.GetProjectTaskById(ctx, related)
	if err != nil {
		return domain.Comment{}, err
	}
	err = p.validateProjectPermissions(ctx, projectTask.Projectid, userModel.AccountId, userModel.Crmid)
	if err != nil {
		return domain.Comment{}, err
	}

	return p.comment.Create(ctx, content, related, userModel.Crmid)
}

func (p ProjectTasksService) validateProjectPermissions(ctx context.Context, id string, client string, contact string) error {
	project, err := p.project.GetProjectById(ctx, id)
	if err != nil {
		return err
	}
	if project.Linktoaccountscontacts != client && project.Linktoaccountscontacts != contact {
		return ErrOperationNotPermitted
	}
	return nil
}
