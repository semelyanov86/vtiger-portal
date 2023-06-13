package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type SearchCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewSearchCrm(config config.Config, cache cache.Cache) SearchCrm {
	return SearchCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (s SearchCrm) SearchFaqs(ctx context.Context, query string) ([]domain.Search, error) {
	sql := "SELECT id, question FROM Faq WHERE faqstatus = 'Published' AND question LIKE '%" + query + "%';"
	searches := make([]domain.Search, 0)
	result, err := s.vtiger.Query(ctx, sql)
	if err != nil {
		return searches, e.Wrap("can not execute query "+sql+", got error: "+result.Error.Message, err)
	}
	for _, data := range result.Result {
		search := domain.ConvertMapToSearch(data)
		search.Module = "Faq"
		searches = append(searches, search)
	}
	return searches, nil
}
