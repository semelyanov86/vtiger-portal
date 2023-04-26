package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

type Comments struct {
	repository repository.Comment
	cache      cache.Cache
	config     config.Config
}

func NewComments(repository repository.Comment, cache cache.Cache, config config.Config) Comments {
	return Comments{
		repository: repository,
		cache:      cache,
		config:     config,
	}
}

func (c Comments) GetRelated(ctx context.Context, id string) ([]domain.Comment, error) {
	return c.repository.RetrieveFromModule(ctx, id)
}

func (c Comments) Create(ctx context.Context, content string, related string, userId string) (domain.Comment, error) {
	var comment domain.Comment
	comment.Commentcontent = content
	comment.RelatedTo = related
	comment.AssignedUserId = c.config.Vtiger.Business.DefaultUser
	comment.Source = "PORTAL"
	comment.Customer = userId
	createdComment, err := c.repository.Create(ctx, comment)
	if err != nil {
		return createdComment, e.Wrap("can not create comment in repository", err)
	}
	return createdComment, nil
}
