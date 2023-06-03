package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

type Leads struct {
	repository repository.Lead
	config     config.Config
}

func NewLeads(repository repository.Lead, config config.Config) Leads {
	return Leads{
		repository: repository,
		config:     config,
	}
}

func (l Leads) Create(ctx context.Context, lead domain.Lead) (domain.Lead, error) {
	lead.AssignedUserId = l.config.Vtiger.Business.DefaultUser
	lead.Leadsource = "Existing Customer"
	lead.Leadstatus = "Attempted to Contact"

	createdLead, err := l.repository.Create(ctx, lead)
	if err != nil {
		return createdLead, e.Wrap("can not create lead in repository", err)
	}
	return createdLead, nil
}
