package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type DocumentCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewDocumentCrm(config config.Config, cache cache.Cache) DocumentCrm {
	return DocumentCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (d DocumentCrm) RetrieveFromModule(ctx context.Context, id string) ([]domain.Document, error) {
	result, err := d.vtiger.RetrieveRelated(ctx, id, "Documents")
	if err != nil {
		return []domain.Document{}, e.Wrap("can not retrieve related documents from id "+id, err)
	}
	documents := make([]domain.Document, 0)
	for _, m := range result.Result {
		document := domain.ConvertMapToDocument(m)
		documents = append(documents, document)
	}
	return documents, nil
}

func (d DocumentCrm) RetrieveFile(ctx context.Context, id string) (vtiger.File, error) {
	result, err := d.vtiger.RetrieveFiles(ctx, id)
	if err != nil {
		return vtiger.File{}, e.Wrap("can not get file with id "+id, err)
	}
	if len(result.Result) == 0 {
		return vtiger.File{}, nil
	}
	return result.Result[0], nil
}

func (d DocumentCrm) AttachFile(ctx context.Context, doc domain.Document, parent string) (domain.Document, error) {
	docMap, err := doc.ConvertToMap()
	if err != nil {
		return doc, err
	}
	resp, err := d.vtiger.Create(ctx, "Documents", docMap)
	if err != nil {
		return doc, err
	}
	newDoc := domain.ConvertMapToDocument(resp.Result)

	d.vtiger.AddRelated(ctx, parent, newDoc.Id, "Documents")

	return newDoc, nil
}

func (d DocumentCrm) DeleteFile(ctx context.Context, id string) error {
	return d.vtiger.Delete(ctx, id)
}
