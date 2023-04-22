package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type HelpDeskCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewHelpDeskCrm(config config.Config, cache cache.Cache) HelpDeskCrm {
	return HelpDeskCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (m HelpDeskCrm) RetrieveById(ctx context.Context, id string) (domain.HelpDesk, error) {
	result, err := m.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.HelpDesk{}, e.Wrap("can not retrieve help desk with id "+id, err)
	}
	return domain.ConvertMapToHelpDesk(result.Result)
}
