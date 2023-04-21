package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
)

type CompanyMock struct {
}

func NewCompanyMock() CompanyMock {
	return CompanyMock{}
}

func (m CompanyMock) GetCompanyInfo(ctx context.Context) (domain.Company, error) {
	return domain.MockedCompany, nil
}
