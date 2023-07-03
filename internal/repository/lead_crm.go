package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type LeadCrm struct {
	vtiger vtiger.Creator
	config config.Config
}

func NewLeadCrm(config config.Config, cache cache.Cache) LeadCrm {
	return LeadCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func NewLeadCrmConcrete(config config.Config, vtiger vtiger.Creator) LeadCrm {
	return LeadCrm{
		vtiger: vtiger,
		config: config,
	}
}

func (l LeadCrm) Create(ctx context.Context, lead domain.Lead) (domain.Lead, error) {
	leadMap, err := lead.ConvertToMap()
	if err != nil {
		return lead, e.Wrap("can not convert to map", err)
	}
	result, err := l.vtiger.Create(ctx, "Leads", leadMap)
	if err != nil {
		return lead, e.Wrap("can not create lead", err)
	}
	newLead := domain.ConvertMapToLead(result.Result)
	return newLead, nil
}
