package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type HelpDeskMockRepository struct {
}

func (m HelpDeskMockRepository) RetrieveById(ctx context.Context, id string) (domain.HelpDesk, error) {
	return domain.MockedHelpDesk, nil
}

func (m HelpDeskMockRepository) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.HelpDesk, error) {
	return []domain.HelpDesk{}, nil
}

func (m HelpDeskMockRepository) Count(ctx context.Context, client string) (int, error) {
	return 1, nil
}

func (m HelpDeskMockRepository) Create(ctx context.Context, ticket domain.HelpDesk) (domain.HelpDesk, error) {
	return domain.MockedHelpDesk, nil
}

func (m HelpDeskMockRepository) Update(ctx context.Context, ticket domain.HelpDesk) (domain.HelpDesk, error) {
	return domain.MockedHelpDesk, nil
}

func (m HelpDeskMockRepository) Revise(ctx context.Context, ticket map[string]any) (domain.HelpDesk, error) {
	return domain.MockedHelpDesk, nil
}
