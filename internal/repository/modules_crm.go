package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type ModulesCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewModulesCrm(config config.Config, cache cache.Cache) ModulesCrm {
	return ModulesCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (m ModulesCrm) GetModuleInfo(ctx context.Context, module string) (vtiger.Module, error) {
	result, err := m.vtiger.Describe(ctx, module)
	if err != nil {
		return vtiger.Module{}, err
	}
	return result.Result, nil
}
