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

const CacheProjectTtl = 5000

type ProjectsService struct {
	repository repository.Project
	cache      cache.Cache
	document   DocumentServiceInterface
	comment    CommentServiceInterface
	module     ModulesService
	config     config.Config
}

func NewProjectsService(repository repository.Project, cache cache.Cache, comments CommentServiceInterface, document DocumentServiceInterface, module ModulesService, config config.Config) ProjectsService {
	return ProjectsService{
		repository: repository,
		cache:      cache,
		document:   document,
		module:     module,
		config:     config,
		comment:    comments,
	}
}

func (p ProjectsService) GetProjectById(ctx context.Context, id string) (domain.Project, error) {
	project := &domain.Project{}
	err := GetFromCache[*domain.Project](id, project, p.cache)
	if err == nil {
		return *project, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		projectData, err := p.repository.RetrieveById(ctx, id)
		if err != nil {
			return projectData, e.Wrap("can not get a project", err)
		}
		err = StoreInCache[*domain.Project](id, &projectData, CacheProjectTtl, p.cache)
		if err != nil {
			return projectData, err
		}
		return projectData, nil
	} else {
		return *project, e.Wrap("can not convert caches data to project", err)
	}
}

func (p ProjectsService) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.Project, int, error) {
	projects, err := p.repository.GetAll(ctx, filter)
	if err != nil {
		return projects, 0, err
	}
	count, err := p.repository.Count(ctx, filter.Client, filter.Contact)
	return projects, count, err
}

func (p ProjectsService) GetRelatedComments(ctx context.Context, id string, companyId string, contactId string) ([]domain.Comment, error) {
	project, err := p.GetProjectById(ctx, id)
	if err != nil {
		return []domain.Comment{}, err
	}
	if project.Linktoaccountscontacts != companyId && project.Linktoaccountscontacts != contactId {
		return []domain.Comment{}, ErrOperationNotPermitted
	}
	return p.comment.GetRelated(ctx, id)
}

func (p ProjectsService) GetRelatedDocuments(ctx context.Context, id string, companyId string, contactId string) ([]domain.Document, error) {
	project, err := p.GetProjectById(ctx, id)
	if err != nil {
		return []domain.Document{}, err
	}
	if project.Linktoaccountscontacts != companyId && project.Linktoaccountscontacts != contactId {
		return []domain.Document{}, ErrOperationNotPermitted
	}
	return p.document.GetRelated(ctx, id)
}

func (p ProjectsService) AddComment(ctx context.Context, content string, related string, userModel domain.User) (domain.Comment, error) {
	project, err := p.GetProjectById(ctx, related)
	if err != nil {
		return domain.Comment{}, err
	}
	if project.Linktoaccountscontacts != userModel.AccountId && project.Linktoaccountscontacts != userModel.Crmid {
		return domain.Comment{}, ErrOperationNotPermitted
	}
	return p.comment.Create(ctx, content, related, userModel.Crmid)
}
