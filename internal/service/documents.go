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
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Documents struct {
	repository repository.Document
	cache      cache.Cache
	config     config.Config
}

const CacheDocuments = "documents-"

const CacheDocumentTtl = 500

func NewDocuments(repository repository.Document, cache cache.Cache, config config.Config) Documents {
	return Documents{
		repository: repository,
		cache:      cache,
		config:     config,
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

func (d Documents) AttachFile(ctx context.Context, file multipart.File, id string, userModel domain.User, header *multipart.FileHeader) (domain.Document, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return domain.Document{}, err
	}

	destinationPath := filepath.Join("storage", userModel.Crmid, id, header.Filename)
	err = os.MkdirAll(filepath.Dir(destinationPath), 0755)
	if err != nil {
		return domain.Document{}, err
	}

	outFile, err := os.Create(destinationPath)
	if err != nil {
		return domain.Document{}, err
	}
	defer outFile.Close()

	_, err = outFile.Write(fileBytes)
	if err != nil {
		return domain.Document{}, err
	}
	link := d.config.HTTP.Host + ":" + strconv.Itoa(d.config.HTTP.Port) + "/" + destinationPath
	doc := domain.Document{
		NotesTitle:       header.Filename,
		Createdtime:      time.Now(),
		Modifiedtime:     time.Now(),
		Filename:         link,
		AssignedUserId:   d.config.Vtiger.Business.DefaultUser,
		Notecontent:      "Document uploaded from " + userModel.FirstName + " " + userModel.LastName,
		Filetype:         "",
		Filesize:         strconv.Itoa(int(header.Size)),
		Filelocationtype: "E",
		Fileversion:      "1",
		Filestatus:       "1",
	}

	return d.repository.AttachFile(ctx, doc, id)
}
