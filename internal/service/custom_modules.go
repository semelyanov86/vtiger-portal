package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"regexp"
	"strconv"
	"time"
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

type ValidationError struct {
	Field   string
	Message string
}

// Implement the error interface
func (v ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s': %s", v.Field, v.Message)
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

func (c CustomModule) CreateEntity(ctx context.Context, input map[string]any, user domain.User, module string) (map[string]any, error) {
	custom, err := c.module.Describe(ctx, module)
	if err != nil {
		return nil, e.Wrap("can not describe module "+module, err)
	}
	_, ok := c.config.Vtiger.Business.CustomModules[module]
	if !ok {
		return nil, ErrModuleNotSupported
	}

	input["assigned_user_id"] = c.config.Vtiger.Business.DefaultUser
	input["from_portal"] = "1"
	input["source"] = "PORTAL"

	err = c.validateInputFields(input, custom)
	if err != nil {
		return input, err
	}

	return c.repository.Create(ctx, input, custom, user)
}

func (c CustomModule) UpdateEntity(ctx context.Context, input map[string]any, id string, user domain.User, module string) (map[string]any, error) {
	custom, err := c.module.Describe(ctx, module)
	if err != nil {
		return nil, e.Wrap("can not describe module "+module, err)
	}
	_, ok := c.config.Vtiger.Business.CustomModules[module]
	if !ok {
		return nil, ErrModuleNotSupported
	}

	entity, err := c.GetById(ctx, module, id, user)
	if err != nil {
		return entity, e.Wrap("can not retrieve entity during update", err)
	}

	for field, value := range input {
		if value != "" {
			entity[field] = value
		}
	}

	err = c.validateInputFields(entity, custom)
	if err != nil {
		return entity, err
	}
	entity, err = c.repository.Update(ctx, entity)
	if err != nil {
		return entity, err
	}

	return entity, err
}

func (c CustomModule) Revise(ctx context.Context, input map[string]any, id string, user domain.User, module string) (map[string]any, error) {
	ticket, err := c.GetById(ctx, module, id, user)
	if err != nil {
		return ticket, e.Wrap("can not retrieve helpdesk during update", err)
	}
	_, ok := c.config.Vtiger.Business.CustomModules[module]
	if !ok {
		return nil, ErrModuleNotSupported
	}
	input["id"] = id

	ticket, err = c.repository.Revise(ctx, input)
	if err != nil {
		return ticket, err
	}
	return ticket, err
}

func (c CustomModule) GetRelatedDocuments(ctx context.Context, id string, module string, user domain.User) ([]domain.Document, error) {
	_, err := c.GetById(ctx, module, id, user)
	if err != nil {
		return []domain.Document{}, err
	}

	return c.document.GetRelated(ctx, id)
}

func (c CustomModule) GetRelatedComments(ctx context.Context, id string, module string, user domain.User) ([]domain.Comment, error) {
	_, err := c.GetById(ctx, module, id, user)
	if err != nil {
		return []domain.Comment{}, err
	}

	return c.comment.GetRelated(ctx, id)
}

func (c CustomModule) AddComment(ctx context.Context, content string, related string, module string, userModel domain.User) (domain.Comment, error) {
	_, err := c.GetById(ctx, module, related, userModel)
	if err != nil {
		return domain.Comment{}, err
	}

	return c.comment.Create(ctx, content, related, userModel.Crmid)
}

func (c CustomModule) validateInputFields(input map[string]any, module vtiger.Module) error {
	var fields = module.Fields
	for _, field := range fields {
		if field.Mandatory && (input[field.Name] == "" || input[field.Name] == nil) {
			return ValidationError{
				Field:   field.Label,
				Message: "input cannot be empty",
			}
		}
		if field.Type.Name == "date" && input[field.Name] != "" && input[field.Name] != nil {
			dateFormat := "2006-01-02"
			_, err := time.Parse(dateFormat, input[field.Name].(string))
			if err != nil {
				return ValidationError{
					Field:   field.Label,
					Message: "field has wrong date format",
				}
			}
		}
		if field.Type.Name == "picklist" && input[field.Name] != "" && input[field.Name] != nil {
			if !field.Type.IsPicklistExist(input[field.Name].(string)) {
				return ValidationError{
					Field:   field.Label,
					Message: "field has wrong value",
				}
			}
		}
		if field.Type.Name == "boolean" && input[field.Name] != nil {
			if input[field.Name] != true && input[field.Name] != false && input[field.Name] != "1" && input[field.Name] != "0" && input[field.Name] != 0 && input[field.Name] != 1 {
				return ValidationError{
					Field:   field.Label,
					Message: "field has wrong boolean value",
				}
			}
		}
		if field.Type.Name == "reference" && input[field.Name] != "" && input[field.Name] != nil {
			pattern := `^\d{1,2}x\d+$`
			re, err := regexp.Compile(pattern)
			if err != nil {
				return err
			}

			if input[field.Name].(repository.ReferenceField).Id != "" && !re.MatchString(input[field.Name].(repository.ReferenceField).Id) {
				return ValidationError{
					Field:   field.Label,
					Message: "wrong reference value",
				}
			}
		}
		if field.Type.Name == "integer" && input[field.Name] != "" && input[field.Name] != nil {
			switch input[field.Name].(type) {
			case int, int32, int64:
				// Value is a valid integer, continue checking the next value
				continue
			case string:
				if _, err := strconv.Atoi(input[field.Name].(string)); err != nil {
					return ValidationError{
						Field:   field.Label,
						Message: "field is not integer",
					}
				}
			default:
				return ValidationError{
					Field:   field.Label,
					Message: "field is not integer",
				}
			}
		}
		if field.Type.Name == "double" && input[field.Name] != "" && input[field.Name] != nil {
			switch input[field.Name].(type) {
			case float32, float64:
				// Value is a valid integer, continue checking the next value
				continue
			case string:
				if _, err := strconv.ParseFloat(input[field.Name].(string), 64); err != nil {
					return ValidationError{
						Field:   field.Label,
						Message: "field is not double",
					}
				}
			default:
				return ValidationError{
					Field:   field.Label,
					Message: "field is not double",
				}
			}
		}
		if field.Type.Name == "currency" && input[field.Name] != "" && input[field.Name] != nil {
			switch input[field.Name].(type) {
			case float32, float64:
				// Value is a valid integer, continue checking the next value
				continue
			case string:
				if _, err := strconv.ParseFloat(input[field.Name].(string), 64); err != nil {
					return ValidationError{
						Field:   field.Label,
						Message: "field is not currency",
					}
				}
			default:
				return ValidationError{
					Field:   field.Label,
					Message: "field is not currency",
				}
			}
		}
	}
	return nil
}
