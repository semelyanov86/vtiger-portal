package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/utils"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strconv"
	"strings"
	"time"
)

type CustomModuleCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

type ReferenceField struct {
	Label string `json:"label"`
	Id    string `json:"id"`
}

var cachedLabels = make(map[string]string)

func NewCustomModuleCrm(config config.Config, cache cache.Cache) CustomModuleCrm {
	return CustomModuleCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (m CustomModuleCrm) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter, custom vtiger.Module, fields []string) ([]map[string]any, error) {
	fieldProps := vtiger.QueryFieldsProps{
		DefaultSort:  "id",
		SearchFields: fields,
		ClientField:  "",
		AccountField: "",
		TableName:    custom.Name,
	}
	accountField := m.FindAccountField(custom)
	contactField := m.FindContactField(custom)
	if accountField != nil {
		fieldProps.AccountField = accountField.Name
	}
	if contactField != nil {
		fieldProps.ClientField = contactField.Name
	}

	items, err := m.vtiger.GetAll(ctx, filter, fieldProps)
	if err != nil {
		return nil, err
	}
	fieldsInfo := m.getFieldTypesInfo(custom)
	records := make([]map[string]any, 0)
	for _, data := range items {
		data = m.convertData(ctx, fieldsInfo, data, false)
		if accountField != nil && data[accountField.Name].(ReferenceField).Id == filter.Client {
			records = append(records, data)
			continue
		}
		if contactField != nil && data[contactField.Name].(ReferenceField).Id == filter.Client {
			records = append(records, data)
		}
	}
	return records, nil
}

func (m CustomModuleCrm) GetById(ctx context.Context, id string, custom vtiger.Module, user domain.User) (map[string]any, error) {
	entity, err := m.vtiger.Retrieve(ctx, id)
	if err != nil {
		return nil, e.Wrap("can not retrieve entity "+id, err)
	}
	if entity.Result == nil {
		return nil, e.Wrap("can not retrieve entity "+id+" got unexpected error", err)
	}
	fieldsInfo := m.getFieldTypesInfo(custom)
	result := m.convertData(ctx, fieldsInfo, entity.Result, true)
	accountField := m.FindAccountField(custom)
	contactField := m.FindContactField(custom)
	if accountField != nil && result[accountField.Name].(ReferenceField).Id == user.AccountId {
		return result, nil
	}
	if contactField != nil && result[contactField.Name].(ReferenceField).Id == user.Crmid {
		return result, nil
	}
	return nil, ErrRecordNotFound
}

func (m CustomModuleCrm) Count(ctx context.Context, client string, contact string, custom vtiger.Module) (int, error) {
	accountField := m.FindAccountField(custom)
	contactField := m.FindContactField(custom)
	body := make(map[string]string)
	if accountField != nil && client != "" {
		body[accountField.Name] = client
	}
	if contactField != nil && contact != "" {
		body[contactField.Name] = contact
	}
	return m.vtiger.Count(ctx, custom.Name, body)
}

func (m CustomModuleCrm) Create(ctx context.Context, entity map[string]any, module vtiger.Module, user domain.User) (map[string]any, error) {
	accountField := m.FindAccountField(module)
	contactField := m.FindContactField(module)
	if accountField != nil {
		entity[accountField.Name] = user.AccountId
	}
	if contactField != nil {
		entity[contactField.Name] = user.Crmid
	}
	result, err := m.vtiger.Create(ctx, module.Name, entity)
	if err != nil {
		return nil, e.Wrap("can not create entity", err)
	}
	return result.Result, nil
}

func (m CustomModuleCrm) Update(ctx context.Context, entity map[string]any) (map[string]any, error) {
	result, err := m.vtiger.Update(ctx, entity)
	if err != nil {
		return nil, e.Wrap("can send update map to vtiger", err)
	}
	return result.Result, nil
}

func (m CustomModuleCrm) Revise(ctx context.Context, entity map[string]any) (map[string]any, error) {
	result, err := m.vtiger.Revise(ctx, entity)
	if err != nil {
		return nil, e.Wrap("can send update map to vtiger", err)
	}
	return result.Result, nil
}

func (m CustomModuleCrm) FindAccountField(module vtiger.Module) *vtiger.ModuleField {
	return utils.FindFieldByRefers(module, "Accounts")
}

func (m CustomModuleCrm) FindContactField(module vtiger.Module) *vtiger.ModuleField {
	return utils.FindFieldByRefers(module, "Contacts")
}

func (m CustomModuleCrm) getFieldTypesInfo(module vtiger.Module) map[string]string {
	res := make(map[string]string)
	for _, field := range module.Fields {
		res[field.Name] = field.Type.Name
	}
	return res
}

func (m CustomModuleCrm) convertData(ctx context.Context, fieldTypes map[string]string, entity map[string]any, retrieveRelated bool) map[string]any {
	result := make(map[string]any)
	for field, data := range entity {
		ftype, ok := fieldTypes[field]
		if !ok {
			result[field] = data
		}
		switch ftype {
		case "date":
			layout := "2006-01-02"
			parsedTime, err := time.Parse(layout, data.(string))
			if err != nil {
				result[field] = time.Time{}
			} else {
				result[field] = parsedTime
			}
		case "datetime":
			layout := "2006-01-02 15:04:05"
			parsedTime, err := time.Parse(layout, data.(string))
			if err != nil {
				result[field] = time.Time{}
			} else {
				result[field] = parsedTime
			}
		case "boolean":
			result[field] = data == "1" || data == 1
		case "reference":
			if retrieveRelated {
				cached, ok := cachedLabels[data.(string)]
				if ok {
					result[field] = ReferenceField{Id: data.(string), Label: cached}
				} else {
					retrieved, err := m.vtiger.Retrieve(ctx, data.(string))
					if err != nil {
						result[field] = ReferenceField{Id: data.(string)}
					} else if retrieved.Result != nil {
						label, ok := retrieved.Result["label"]
						if !ok {
							result[field] = ReferenceField{Id: data.(string)}
						} else {
							cachedLabels[field] = label.(string)
							result[field] = ReferenceField{Id: data.(string), Label: label.(string)}
						}
					} else {
						result[field] = ReferenceField{Id: data.(string)}
					}
				}

			} else {
				result[field] = ReferenceField{Id: data.(string)}
			}
		case "integer":
			converted, err := strconv.Atoi(data.(string))
			if err == nil {
				result[field] = converted
			} else {
				result[field] = 0
			}
		case "currency", "double":
			parsedFloat, err := strconv.ParseFloat(data.(string), 64)
			if err == nil {
				result[field] = parsedFloat
			} else {
				result[field] = 0
			}
		default:
			result[field] = data
		}
	}
	tags, ok := entity["tags"]
	if ok {
		result["tags"] = strings.Split(tags.(string), ",")
	}
	return result
}
