package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type ManagersCrm struct {
	vtiger vtiger.Connector
	config config.Config
}

func NewManagersCrm(config config.Config, cache cache.Cache) ManagersCrm {
	return ManagersCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func NewManagersConcrete(config config.Config, vtiger vtiger.Connector) ManagersCrm {
	return ManagersCrm{
		vtiger: vtiger,
		config: config,
	}
}

func (m ManagersCrm) RetrieveById(ctx context.Context, id string) (domain.Manager, error) {
	result, err := m.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.Manager{}, e.Wrap("can not retrieve manager with id "+id, err)
	}
	manager := domain.ConvertMapToManager(result.Result)
	return manager, nil
}
