package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
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

func (f FaqsCrm) GetAllFaqs(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.Faq, error) {
	items, err := f.vtiger.GetByWhereClause(ctx, filter, "faqstatus", "Published", "Faq")
	if err != nil {
		return nil, err
	}

	faqs := make([]domain.Faq, 0, len(items))
	for _, data := range items {
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
