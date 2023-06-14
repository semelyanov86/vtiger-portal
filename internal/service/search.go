package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"sync"
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
	var results []domain.Search
	var wg sync.WaitGroup
	errChan := make(chan error, 3)
	resultChan := make(chan []domain.Search, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		faqs, e := s.repository.SearchFaqs(ctx, query)
		if e != nil {
			errChan <- errors.New("can not get faqs: " + e.Error())
			return
		}
		resultChan <- faqs
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tickets, e := s.repository.SearchTickets(ctx, query, user)
		if e != nil {
			errChan <- errors.New("can not get tickets: " + e.Error())
			return
		}
		resultChan <- tickets
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		projects, e := s.repository.SearchProjects(ctx, query, user)
		if e != nil {
			errChan <- errors.New("can not get projects: " + e.Error())
			return
		}
		resultChan <- projects
	}()

	wg.Wait()
	close(errChan)
	close(resultChan)

	for r := range resultChan {
		results = append(results, r...)
	}

	if len(errChan) > 0 {
		return results, <-errChan
	}

	return results, nil
}
