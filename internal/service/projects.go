package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strconv"
)

const CacheProjectTtl = 5000

type ProjectsService struct {
	repository     repository.Project
	taskRepository repository.ProjectTask
	cache          cache.Cache
	document       DocumentServiceInterface
	comment        CommentServiceInterface
	module         ModulesService
	config         config.Config
}

func NewProjectsService(repository repository.Project, cache cache.Cache, comments CommentServiceInterface, document DocumentServiceInterface, module ModulesService, config config.Config, taskRepository repository.ProjectTask) ProjectsService {
	return ProjectsService{
		repository:     repository,
		cache:          cache,
		document:       document,
		module:         module,
		config:         config,
		comment:        comments,
		taskRepository: taskRepository,
	}
}

func (p ProjectsService) GetProjectById(ctx context.Context, id string, calcStat bool, userModel *domain.User) (domain.Project, error) {
	project := &domain.Project{}
	err := GetFromCache[*domain.Project](id, project, p.cache)
	if err == nil {
		if userModel.AccountId != project.Linktoaccountscontacts && userModel.Crmid != project.Linktoaccountscontacts {
			return *project, ErrOperationNotPermitted
		}
		return *project, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		projectData, err := p.repository.RetrieveById(ctx, id)
		if err != nil {
			return projectData, e.Wrap("can not get a project", err)
		}
		if userModel.AccountId != projectData.Linktoaccountscontacts && userModel.Crmid != projectData.Linktoaccountscontacts {
			return projectData, ErrOperationNotPermitted
		}
		if calcStat {
			stat, err2 := p.calcProjectStatistics(ctx, id)

			if err2 == nil {
				projectData.Statistics = stat
			}
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

func (p ProjectsService) calcProjectStatistics(ctx context.Context, id string) (domain.CurrentProjectStatistics, error) {
	var stat domain.CurrentProjectStatistics
	tasks, err := p.taskRepository.GetFromProject(ctx, vtiger.PaginationQueryFilter{
		Page:     1,
		PageSize: 100,
		Client:   "",
		Contact:  "",
		Parent:   id,
		Sort:     "",
		Filters:  nil,
		Search:   "",
	})
	if err != nil {
		return stat, e.Wrap("can not calculate project statistics", err)
	}

	for _, task := range tasks {
		stat.TotalTasks++
		if s, err := strconv.ParseFloat(task.Projecttaskhours, 64); err == nil {
			stat.TotalHours += s
		}
		if task.Projecttaskstatus == "Open" {
			stat.OpenTasks++
		}
		if task.Projecttaskstatus == "In Progress" {
			stat.InProgressTasks++
		}
		if task.Projecttaskstatus == "Completed" {
			stat.ClosedTasks++
		}
		if task.Projecttaskstatus == "Deferred" {
			stat.DeferredTasks++
		}
		if task.Projecttaskstatus == "Canceled" {
			stat.CancelledTasks++
		}
		if task.Projecttaskpriority == "low" {
			stat.LowTasks++
		}
		if task.Projecttaskpriority == "normal" {
			stat.NormalTasks++
		}
		if task.Projecttaskpriority == "high" {
			stat.HighTasks++
		}
	}
	return stat, nil
}

func (p ProjectsService) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.Project, int, error) {
	projects, err := p.repository.GetAll(ctx, filter)
	if err != nil {
		return projects, 0, err
	}
	count, err := p.repository.Count(ctx, filter.Client, filter.Contact)
	return projects, count, err
}

func (p ProjectsService) GetRelatedComments(ctx context.Context, id string, userModel *domain.User) ([]domain.Comment, error) {
	_, err := p.GetProjectById(ctx, id, false, userModel)
	if err != nil {
		return []domain.Comment{}, err
	}

	return p.comment.GetRelated(ctx, id)
}

func (p ProjectsService) GetRelatedDocuments(ctx context.Context, id string, userModel *domain.User) ([]domain.Document, error) {
	_, err := p.GetProjectById(ctx, id, false, userModel)
	if err != nil {
		return []domain.Document{}, err
	}

	return p.document.GetRelated(ctx, id)
}

func (p ProjectsService) AddComment(ctx context.Context, content string, related string, userModel *domain.User) (domain.Comment, error) {
	_, err := p.GetProjectById(ctx, related, false, userModel)
	if err != nil {
		return domain.Comment{}, err
	}

	return p.comment.Create(ctx, content, related, userModel.Crmid)
}
