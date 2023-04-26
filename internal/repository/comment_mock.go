package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
)

type CommentMock struct {
}

func NewCommentMock() CommentMock {
	return CommentMock{}
}

func (c CommentMock) RetrieveFromModule(ctx context.Context, id string) ([]domain.Comment, error) {
	return []domain.Comment{domain.MockedComment}, nil
}

func (c CommentMock) Create(ctx context.Context, comment domain.Comment) (domain.Comment, error) {
	return domain.MockedComment, nil
}
