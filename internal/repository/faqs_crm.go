package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strconv"
)

type FaqsCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewFaqsCrm(config config.Config, cache cache.Cache) FaqsCrm {
	return FaqsCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (f FaqsCrm) GetAllFaqs(ctx context.Context, filter PaginationQueryFilter) ([]domain.Faq, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM Faq WHERE faqstatus = 'Published' LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"
	faqs := make([]domain.Faq, 0)
	result, err := f.vtiger.Query(ctx, query)
	if err != nil {
		return faqs, e.Wrap("can not execute query "+query+", got error: "+result.Error.Message, err)
	}
	for _, data := range result.Result {
		faq, err := domain.ConvertMapToFaq(data)
		if err != nil {
			return faqs, e.Wrap("can not convert map to faq", err)
		}
		faqs = append(faqs, faq)
	}
	return faqs, nil
}

func (f FaqsCrm) Count(ctx context.Context, _ string) (int, error) {
	body := make(map[string]string)
	body["faqstatus"] = "Published"
	return f.vtiger.Count(ctx, "Faq", body)
}
