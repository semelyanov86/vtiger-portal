package v1

import (
	"bytes"
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

func TestHandler_getAllFromCustomModule(t *testing.T) {
	var notOwnedUser = repository.MockedUser
	notOwnedUser.Crmid = "22x44"

	tests := []struct {
		name         string
		postfix      string
		userModel    *domain.User
		statusCode   int
		responseBody string
		module       string
	}{
		{
			name:         "Custom Modules received",
			postfix:      "?page=1&size=20",
			statusCode:   http.StatusOK,
			responseBody: `"description":"Some description for mocked entity"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
		},
		{
			name:         "Anonymous Access",
			postfix:      "?page=1&size=20",
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
			module:       "Assets",
		}, {
			name:         "Wrong ID",
			postfix:      "?page=1&size=20",
			statusCode:   http.StatusOK,
			responseBody: `"data":[]`,
			userModel:    &notOwnedUser,
			module:       "Assets",
		}, {
			name:         "Wrong Pagination",
			postfix:      "?page=notknown&size=smth",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `Invalid page number`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
		},
		{
			name:         "Module not supported",
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
				Vtiger: config.VtigerConfig{
					Business: config.VtigerBusinessConfig{CustomModules: map[string][]string{tt.module: {"Documents"}}, DefaultPagination: 20}},
			})

			services := &service.Services{CustomModules: customModuleService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultPagination: 20}}}}
			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/custom-modules/:module", func(c *gin.Context) {

			}, handler.getAllEntities)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/custom-modules/Assets"+tt.postfix,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_createCustomModule(t *testing.T) {
	tests := []struct {
		name         string
		userModel    *domain.User
		statusCode   int
		responseBody string
		requestBody  string
		module       string
	}{
		{
			name:         "Entity created",
			statusCode:   http.StatusCreated,
			responseBody: `"description":"Some description for mocked entity"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody: `{
  "title": "Some test problem",
  "description": "This is some description",
  "status": "New",
  "hours": 4,
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Module not supported",
			statusCode:   http.StatusBadRequest,
			responseBody: `module not supported`,
			userModel:    &repository.MockedUser,
			module:       "TestModule",
			requestBody: `{
  "title": "Some test problem",
  "description": "This is some description",
  "status": "New",
  "hours": 4,
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Anonymous Access",
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
			module:       "Assets",
			requestBody: `{
  "title": "Some test problem",
  "description": "This is some description",
  "status": "New",
  "hours": 4,
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Empty Body",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `"error":"Validation Error"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody:  `{}`,
		},
		{
			name:         "Entity with no title",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `"field":"Title"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody: `{
  "title": "",
  "description": "This is some description",
  "status": "New",
  "hours": 4,
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Entity with wrong status",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `"field":"Status"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody: `{
  "title": "Some correct title",
  "description": "This is some description",
  "status": "some wrong value",
  "hours": 4,
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Entity with number as string",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `"field":"Hours"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody: `{
  "title": "Some correct title",
  "description": "This is some description",
  "status": "New",
  "hours": "hours",
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Entity with wrong date",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `"field":"Date"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody: `{
  "title": "Some correct title",
  "description": "This is some description",
  "status": "New",
  "hours": 2,
  "date": "4324234"
}`,
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
			r.POST("/api/v1/custom-modules/:module", func(c *gin.Context) {

			}, handler.createCustomModule)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/custom-modules/Assets",
				bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_updateCustomModule(t *testing.T) {
	tests := []struct {
		name         string
		userModel    *domain.User
		statusCode   int
		responseBody string
		requestBody  string
		module       string
	}{
		{
			name:         "Entity updated",
			statusCode:   http.StatusAccepted,
			responseBody: `"description":"Some description for mocked entity"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody: `{
  "title": "Some test problem",
  "description": "This is some description",
  "status": "New",
  "hours": 4,
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Module not supported",
			statusCode:   http.StatusBadRequest,
			responseBody: `module not supported`,
			userModel:    &repository.MockedUser,
			module:       "TestModule",
			requestBody: `{
  "title": "Some test problem",
  "description": "This is some description",
  "status": "New",
  "hours": 4,
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Anonymous Access",
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
			module:       "Assets",
			requestBody: `{
  "title": "Some test problem",
  "description": "This is some description",
  "status": "New",
  "hours": 4,
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Empty Body",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `"error":"Validation Error"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody:  `{}`,
		},
		{
			name:         "Entity with no title",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `"field":"Title"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody: `{
  "title": "",
  "description": "This is some description",
  "status": "New",
  "hours": 4,
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Entity with wrong status",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `"field":"Status"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody: `{
  "title": "Some correct title",
  "description": "This is some description",
  "status": "some wrong value",
  "hours": 4,
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Entity with number as string",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `"field":"Hours"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody: `{
  "title": "Some correct title",
  "description": "This is some description",
  "status": "New",
  "hours": "hours",
  "date": "2023-05-06"
}`,
		},
		{
			name:         "Entity with wrong date",
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `"field":"Date"`,
			userModel:    &repository.MockedUser,
			module:       "Assets",
			requestBody: `{
  "title": "Some correct title",
  "description": "This is some description",
  "status": "New",
  "hours": 2,
  "date": "4324234"
}`,
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
			r.PUT("/api/v1/custom-modules/:module/:id", func(c *gin.Context) {

			}, handler.updateEntity)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/v1/custom-modules/Assets/22x45",
				bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
