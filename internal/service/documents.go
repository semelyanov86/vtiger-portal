package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type Documents struct {
	repository repository.Document
	cache      cache.Cache
}

const CacheDocuments = "documents-"

const CacheDocumentTtl = 500

func NewDocuments(repository repository.Document, cache cache.Cache) Documents {
	return Documents{
		repository: repository,
		cache:      cache,
	}
}

func (d Documents) GetRelated(ctx context.Context, id string) ([]domain.Document, error) {
	documents := &[]domain.Document{}
	err := GetFromCache[*[]domain.Document](CacheDocuments+id, documents, d.cache)
	if err == nil {
		return *documents, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		documentsData, err := d.repository.RetrieveFromModule(ctx, id)
		if err != nil {
			return documentsData, e.Wrap("can not get a documents", err)
		}
		err = StoreInCache[*[]domain.Document](CacheDocuments+id, &documentsData, CacheDocumentTtl, d.cache)
		if err != nil {
			return documentsData, err
		}
		return documentsData, nil
	} else {
		return *documents, e.Wrap("can not convert caches data to documents", err)
	}
}

func (d Documents) GetFile(ctx context.Context, id string, relatedId string) (vtiger.File, error) {
	documents, err := d.GetRelated(ctx, relatedId)
	if err != nil {
		return vtiger.File{}, e.Wrap("can not get related documents", err)
	}
	for _, document := range documents {
		if document.Imageattachmentids == id {
			return d.repository.RetrieveFile(ctx, id)
		}
	}
	return vtiger.File{}, ErrOperationNotPermitted
}
