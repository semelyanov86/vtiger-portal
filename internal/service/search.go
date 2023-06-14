package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
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

func (s Search) GlobalSearch(ctx context.Context, query string, user domain.User) ([]domain.Search, error) {
	results, err := s.repository.SearchFaqs(ctx, query)
	if err != nil {
		return results, e.Wrap("can not get faqs", err)
	}
	return s.repository.SearchTickets(ctx, query, user)
}
