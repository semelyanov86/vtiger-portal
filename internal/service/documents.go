package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	cachedDocumentData, err := d.cache.Get(CacheDocuments + id)
	if errors.Is(cache.ErrItemNotFound, err) || cachedDocumentData == nil {
		documentData, err := d.repository.RetrieveFromModule(ctx, id)
		if err != nil {
			return documentData, e.Wrap("can not get a documents", err)
		}
		cachedValue, err := json.Marshal(documentData)
		if err != nil {
			return documentData, err
		}
		err = d.cache.Set(CacheDocuments+id, cachedValue, CacheDocumentTtl)
		if err != nil {
			return documentData, err
		}
		return documentData, nil
	} else {
		decodedDocuments := &[]domain.Document{}
		err = json.Unmarshal(cachedDocumentData, decodedDocuments)
		if err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := cachedDocumentData[jsonErr.Offset-10 : jsonErr.Offset+10]

				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
			}
			return *decodedDocuments, e.Wrap("can not convert caches data to documents", err)
		}
		return *decodedDocuments, nil
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
