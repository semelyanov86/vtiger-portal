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
	repository      repository.Comment
	cache           cache.Cache
	config          config.Config
	usersService    UsersService
	managersService ManagerService
}

func NewComments(repository repository.Comment, cache cache.Cache, config config.Config, usersService UsersService, managerService ManagerService) Comments {
	return Comments{
		repository:      repository,
		cache:           cache,
		config:          config,
		usersService:    usersService,
		managersService: managerService,
	}
}

func (c Comments) GetRelated(ctx context.Context, id string) ([]domain.Comment, error) {
	comments, err := c.repository.RetrieveFromModule(ctx, id)
	if err != nil {
		return comments, err
	}
	for i, comment := range comments {
		if comment.Customer != "" {
			user, err := c.usersService.FindByCrmid(ctx, comment.Customer)
			if err != nil {
				continue
			}
			comment.Author = domain.CommentAuthor{
				FirstName:    user.FirstName,
				LastName:     user.LastName,
				Email:        user.Email,
				Id:           user.Crmid,
				Imagecontent: user.Imagecontent,
			}
		} else if comment.AssignedUserId != "" {
			manager, err := c.managersService.GetManagerById(ctx, comment.AssignedUserId)
			if err != nil {
				continue
			}
			comment.Author = domain.CommentAuthor{
				FirstName: manager.FirstName,
				LastName:  manager.LastName,
				Email:     manager.Email,
				Id:        manager.Id,
			}

		}
		comments[i] = comment
	}
	return comments, err
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
