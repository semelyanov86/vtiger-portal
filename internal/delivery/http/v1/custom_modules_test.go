package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	mock_service "github.com/semelyanov86/vtiger-portal/internal/service/mocks"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_getCustomModuleById(t *testing.T) {
	var notOwnedUser = repository.MockedUser
	notOwnedUser.Crmid = "22x44"

	tests := []struct {
		name         string
		id           string
		userModel    *domain.User
		statusCode   int
		responseBody string
		module       string
	}{
		{
			name:         "Custom module received",
			id:           "23x42343",
			statusCode:   http.StatusOK,
			responseBody: `"description":"Some description for mocked entity"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
		},
		{
			name:         "Anonymous Access",
			id:           "23x42343",
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
			module:       "Assets",
		}, {
			name:         "Wrong ID",
			id:           "17",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `wrong id`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
		},
		{
			name:         "Not owned entity",
			id:           "17x16",
			statusCode:   http.StatusForbidden,
			responseBody: `"message":"You are not allowed to view this record"`,
			userModel:    &notOwnedUser,
			module:       "Assets",
		},
		{
			name:         "Module not supported",
			id:           "23x42343",
			statusCode:   http.StatusBadRequest,
			responseBody: `module not supported`,
			userModel:    &repository.MockedUser,
			module:       "TestModule",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rm := repository.NewCustomModuleConcrete(config.Config{}, vtiger.NewMockedVtigerConnector())
			rt := repository.NewModulesCrmConcrete(config.Config{}, vtiger.NewMockedVtigerConnector())

			moduleService := service.NewModulesService(rt, cache.NewMemoryCache())
			customModuleService := service.NewCustomModuleService(rm, cache.NewMemoryCache(), &mock_service.MockCommentServiceInterface{}, mock_service.NewMockDocumentServiceInterface(c), moduleService, config.Config{
				Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{CustomModules: map[string][]string{tt.module: {"Documents"}}}},
			})

			services := &service.Services{CustomModules: customModuleService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/custom-modules/:module/:id", func(c *gin.Context) {

			}, handler.getEntityById)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/custom-modules/Assets/"+tt.id,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getRelatedDocumentsFromCustomModule(t *testing.T) {
	var notOwnedUser = repository.MockedUser
	notOwnedUser.Crmid = "22x44"

	tests := []struct {
		name         string
		id           string
		userModel    *domain.User
		statusCode   int
		responseBody string
		module       string
	}{
		{
			name:         "Related documents received",
			id:           "23x42343",
			userModel:    &repository.MockedUser,
			statusCode:   http.StatusOK,
			responseBody: `"tags":["tag1","tag2"]`,
			module:       "Assets",
		},
		{
			name:         "Module Not Supported",
			id:           "23x42343",
			userModel:    &repository.MockedUser,
			statusCode:   http.StatusBadRequest,
			responseBody: `module not supported`,
			module:       "TestModule",
		},
		{
			name:         "Anonymous Access",
			id:           "17x1",
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
			module:       "Assets",
		},
		{
			name:         "Wrong ID",
			id:           "17",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `wrong id`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
		},
		{
			name:         "Not owned module",
			id:           "17x16",
			statusCode:   http.StatusForbidden,
			responseBody: `"message":"You are not allowed to view this record"`,
			userModel:    &notOwnedUser,
			module:       "Assets",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rm := repository.NewCustomModuleConcrete(config.Config{}, vtiger.NewMockedVtigerConnector())
			rd := repository.NewDocumentConcrete(config.Config{}, vtiger.NewMockedVtigerConnector())
			rt := repository.NewModulesCrmConcrete(config.Config{}, vtiger.NewMockedVtigerConnector())

			moduleService := service.NewModulesService(rt, cache.NewMemoryCache())

			documentService := service.NewDocuments(rd, cache.NewMemoryCache(), config.Config{})

			customModuleService := service.NewCustomModuleService(rm, cache.NewMemoryCache(), service.Comments{}, documentService, moduleService, config.Config{
				Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{CustomModules: map[string][]string{tt.module: {"Documents"}}}},
			})

			services := &service.Services{CustomModules: customModuleService, Documents: documentService, Context: service.MockedContextService{MockedUser: tt.userModel}, Modules: moduleService}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/custom-modules/:module/:id/documents", func(c *gin.Context) {

			}, handler.getCustomDocuments)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/custom-modules/Assets/"+tt.id+"/documents",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
