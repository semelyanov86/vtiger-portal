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

func (m CustomModuleCrm) GetAll(ctx context.Context, filter PaginationQueryFilter, custom vtiger.Module, fields []string) ([]map[string]any, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM " + custom.Name + " WHERE "
	sort := filter.Sort
	accountField := m.findAccountField(custom)
	contactField := m.findContactField(custom)
	if filter.Search != "" {
		for i, field := range fields {
			if i == 0 {
				continue
			}
			query += field + " LIKE '%" + filter.Search + "%' OR"
		}
		query = strings.Trim(query, "OR")
		query += " "
	} else {
		if accountField != nil {
			query += accountField.Name + " = " + filter.Client + " OR "
		}
		if contactField != nil {
			query += contactField.Name + " = " + filter.Client + " "
		}
		query = strings.Trim(query, "OR")
		query += " "
	}
	query += GenerateOrderByClause(sort)
	query += " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"

	records := make([]map[string]any, 0)

	result, err := m.vtiger.Query(ctx, query)
	if err != nil {
		return nil, e.Wrap("can not execute query "+query+", got error", err)
	}
	fieldsInfo := m.getFieldTypesInfo(custom)
	for _, data := range result.Result {
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
	accountField := m.findAccountField(custom)
	contactField := m.findContactField(custom)
	if accountField != nil && result[accountField.Name].(ReferenceField).Id == user.AccountId {
		return result, nil
	}
	if contactField != nil && result[contactField.Name].(ReferenceField).Id == user.Crmid {
		return result, nil
	}
	return nil, ErrRecordNotFound
}

func (m CustomModuleCrm) Count(ctx context.Context, client string, contact string, custom vtiger.Module) (int, error) {
	accountField := m.findAccountField(custom)
	contactField := m.findContactField(custom)
	body := make(map[string]string)
	if accountField != nil && client != "" {
		body[accountField.Name] = client
	}
	if contactField != nil && contact != "" {
		body[contactField.Name] = client
	}
	return m.vtiger.Count(ctx, custom.Name, body)
}

func (m CustomModuleCrm) findAccountField(module vtiger.Module) *vtiger.ModuleField {
	return utils.FindFieldByRefers(module, "Accounts")
}

func (m CustomModuleCrm) findContactField(module vtiger.Module) *vtiger.ModuleField {
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
