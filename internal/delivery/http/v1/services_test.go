package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	mock_repository "github.com/semelyanov86/vtiger-portal/internal/repository/mocks"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_getServiceById(t *testing.T) {
	type mockRepositoryService func(r *mock_repository.MockService)
	type mockRepositoryCurrency func(r *mock_repository.MockCurrency)

	tests := []struct {
		name         string
		id           string
		mockService  mockRepositoryService
		mockCurrency mockRepositoryCurrency
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Service received",
			id:   "25x52",
			mockService: func(r *mock_repository.MockService) {
				r.EXPECT().RetrieveById(context.Background(), "25x52").Return(domain.MockedService, nil)
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
				r.EXPECT().RetrieveById(context.Background(), "21x1").Return(domain.MockedCurrency, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"servicename":"Cleaning laptop"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "17x1",
			mockService: func(r *mock_repository.MockService) {
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name: "Wrong ID",
			id:   "17",
			mockService: func(r *mock_repository.MockService) {
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `wrong id`,
			userModel:    &repository.MockedUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rm := mock_repository.NewMockService(c)
			tt.mockService(rm)

			rc := mock_repository.NewMockCurrency(c)
			tt.mockCurrency(rc)

			serviceService := service.NewServicesService(rm, cache.NewMemoryCache(), service.NewCurrencyService(rc, cache.NewMemoryCache()), service.ModulesService{}, config.Config{})

			services := &service.Services{Services: serviceService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/services/:id", func(c *gin.Context) {

			}, handler.getService)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/services/"+tt.id,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getAllServices(t *testing.T) {
	type mockRepositoryService func(r *mock_repository.MockService)
	type mockRepositoryCurrency func(r *mock_repository.MockCurrency)

	tests := []struct {
		name         string
		postfix      string
		mockService  mockRepositoryService
		mockCurrency mockRepositoryCurrency
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name:    "Services received",
			postfix: "?page=1&size=20",
			mockService: func(r *mock_repository.MockService) {
				r.EXPECT().GetAll(context.Background(), repository.PaginationQueryFilter{
					Page:     1,
					PageSize: 20,
					Client:   "11x1",
					Contact:  "12x11",
					Filters: map[string]any{
						"discontinued": true,
					},
				}).Return([]domain.Service{domain.MockedService}, nil)
				r.EXPECT().RetrieveById(context.Background(), "25x52").Return(domain.MockedService, nil)
				r.EXPECT().Count(context.Background(), map[string]any{
					"discontinued": true,
				}).Return(1, nil)
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
				r.EXPECT().RetrieveById(context.Background(), "21x1").Return(domain.MockedCurrency, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"servicename":"Cleaning laptop"`,
			userModel:    &repository.MockedUser,
		},
		{
			name:    "Anonymous Access",
			postfix: "?page=1&size=20",
			mockService: func(r *mock_repository.MockService) {
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name:    "Wrong Filter",
			postfix: "?filter[discontinued]=fsdgfsafgsddsa",
			mockService: func(r *mock_repository.MockService) {

			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			statusCode:   http.StatusBadRequest,
			responseBody: `Wrong filter value for discontinued, expected boolean`,
			userModel:    &repository.MockedUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rm := mock_repository.NewMockService(c)
			tt.mockService(rm)

			rc := mock_repository.NewMockCurrency(c)
			tt.mockCurrency(rc)

			servicesService := service.NewServicesService(rm, cache.NewMemoryCache(), service.NewCurrencyService(rc, cache.NewMemoryCache()), service.ModulesService{}, config.Config{})

			services := &service.Services{Services: servicesService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultPagination: 20}}}}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/services", func(c *gin.Context) {

			}, handler.getAllServices)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/services"+tt.postfix,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
