package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
)

type Faqs struct {
	repository repository.Faq
	cache      cache.Cache
	module     ModulesService
	config     config.Config
}

func NewFaqsService(repository repository.Faq, cache cache.Cache, module ModulesService, config config.Config) Faqs {
	return Faqs{
		repository: repository,
		cache:      cache,
		module:     module,
		config:     config,
	}
}

func (f Faqs) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.Faq, int, error) {
	faqs, err := f.repository.GetAllFaqs(ctx, filter)
	if err != nil {
		return faqs, 0, err
	}
	count, err := f.repository.Count(ctx, filter.Client)
	return faqs, count, err
}
