package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
)

type Comments struct {
	repository repository.Comment
	cache      cache.Cache
}

func NewComments(repository repository.Comment, cache cache.Cache) Comments {
	return Comments{
		repository: repository,
		cache:      cache,
	}
}

func (c Comments) GetRelated(ctx context.Context, id string) ([]domain.Comment, error) {
	return c.repository.RetrieveFromModule(ctx, id)
}
