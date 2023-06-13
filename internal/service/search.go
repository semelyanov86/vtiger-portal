package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
)

type Search struct {
	repository repository.SearchCrm
	cache      cache.Cache
	config     config.Config
}

func NewSearchService(repository repository.SearchCrm, cache cache.Cache, config config.Config) Search {
	return Search{
		repository: repository,
		cache:      cache,
		config:     config,
	}
}

func (s Search) GlobalSearch(ctx context.Context, query string) ([]domain.Search, error) {
	return s.repository.SearchFaqs(ctx, query)
}
