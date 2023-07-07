package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type CommentCrm struct {
	vtiger vtiger.Connector
	config config.Config
}

func NewCommentCrm(config config.Config, cache cache.Cache) CommentCrm {
	return CommentCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func NewCommentConcrete(config config.Config, vtiger vtiger.Connector) CommentCrm {
	return CommentCrm{
		vtiger: vtiger,
		config: config,
	}
}

func (c CommentCrm) RetrieveFromModule(ctx context.Context, id string) ([]domain.Comment, error) {
	result, err := c.vtiger.RetrieveRelated(ctx, id, "ModComments")
	if err != nil {
		return []domain.Comment{}, e.Wrap("can not retrieve related comments from id "+id, err)
	}
	comments := make([]domain.Comment, 0)
	for _, m := range result.Result {
		comment := domain.ConvertMapToComment(m)
		if !comment.IsPrivate {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (c CommentCrm) Create(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	commentMap, err := comment.ConvertToMap()
	if err != nil {
		return comment, e.Wrap("can not convert to map", err)
	}
	result, err := c.vtiger.Create(ctx, "ModComments", commentMap)
	if err != nil {
		return comment, e.Wrap("can not create comment", err)
	}
	newComment := domain.ConvertMapToComment(result.Result)
	return newComment, nil
}
