package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

const CacheCustomModuleTtl = 500

var ErrModuleNotSupported = errors.New("module not supported")

type CustomModule struct {
	repository repository.CustomModuleCrm
	cache      cache.Cache
	comment    CommentServiceInterface
	document   DocumentServiceInterface
	module     ModulesService
	config     config.Config
}

func NewCustomModuleService(repository repository.CustomModuleCrm, cache cache.Cache, comments CommentServiceInterface, document DocumentServiceInterface, module ModulesService, config config.Config) CustomModule {
	return CustomModule{
		repository: repository,
		cache:      cache,
		comment:    comments,
		document:   document,
		module:     module,
		config:     config,
	}
}

func (c CustomModule) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter, moduleName string) ([]map[string]any, int, error) {
	module, err := c.module.Describe(ctx, moduleName)
	if err != nil {
		return nil, 0, e.Wrap("can not describe module "+moduleName, err)
	}
	cfgModule, ok := c.config.Vtiger.Business.CustomModules[moduleName]
	if !ok {
		return nil, 0, ErrModuleNotSupported
	}
	if filter.Sort == "" {
		filter.Sort = cfgModule[0]
	}
	entities, err := c.repository.GetAll(ctx, filter, module, cfgModule)
	if err != nil {
		return entities, 0, err
	}
	count, err := c.repository.Count(ctx, filter.Client, filter.Contact, module)
	return entities, count, err
}

func (c CustomModule) GetById(ctx context.Context, moduleName string, id string, user domain.User) (map[string]any, error) {
	module, err := c.module.Describe(ctx, moduleName)
	if err != nil {
		return nil, e.Wrap("can not describe module "+moduleName, err)
	}
	_, ok := c.config.Vtiger.Business.CustomModules[moduleName]
	if !ok {
		return nil, ErrModuleNotSupported
	}
	return c.repository.GetById(ctx, id, module, user)
}
