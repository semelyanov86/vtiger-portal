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

func TestHandler_receiveInvoiceById(t *testing.T) {
	type mockRepositoryInvoice func(r *mock_repository.MockInvoice)
	type mockRepositoryCurrency func(r *mock_repository.MockCurrency)

	tests := []struct {
		name         string
		id           string
		mockInvoice  mockRepositoryInvoice
		mockCurrency mockRepositoryCurrency
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Invoice received",
			id:   "2x53",
			mockInvoice: func(r *mock_repository.MockInvoice) {
				r.EXPECT().RetrieveById(context.Background(), "2x53").Return(domain.Invoice{
					Description: "This is test description",
					AccountID:   "11x1",
					CurrencyID:  "22x22",
				}, nil)
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
				r.EXPECT().RetrieveById(context.Background(), "22x22").Return(domain.Currency{}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"description":"This is test description"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "2x53",
			mockInvoice: func(r *mock_repository.MockInvoice) {
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name: "Wrong ID",
			id:   "17",
			mockInvoice: func(r *mock_repository.MockInvoice) {
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `wrong id`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Not owned invoice",
			id:   "2x53",
			mockInvoice: func(r *mock_repository.MockInvoice) {
				r.EXPECT().RetrieveById(context.Background(), "2x53").Return(domain.Invoice{
					Description: "This is test description",
					AccountID:   "12x44",
				}, nil)
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			statusCode:   http.StatusForbidden,
			responseBody: `"message":"You are not allowed to view this record"`,
			userModel:    &repository.MockedUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rm := mock_repository.NewMockInvoice(c)
			rc := mock_repository.NewMockCurrency(c)
			tt.mockInvoice(rm)
			tt.mockCurrency(rc)

			currencyService := service.NewCurrencyService(rc, cache.NewMemoryCache())

			invoiceService := service.NewInvoiceService(rm, cache.NewMemoryCache(), service.ModulesService{}, config.Config{}, currencyService)

			services := &service.Services{Invoices: invoiceService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/invoices/:id", func(c *gin.Context) {

			}, handler.getInvoice)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/invoices/"+tt.id,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getAllInvoices(t *testing.T) {
	type mockRepositoryInvoice func(r *mock_repository.MockInvoice)
	type mockRepositoryCurrency func(r *mock_repository.MockCurrency)

	tests := []struct {
		name         string
		postfix      string
		mockInvoice  mockRepositoryInvoice
		mockCurrency mockRepositoryCurrency
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Invoices received",
			mockInvoice: func(r *mock_repository.MockInvoice) {
				r.EXPECT().GetAll(context.Background(), repository.PaginationQueryFilter{
					Page:     1,
					PageSize: 20,
					Client:   "11x1",
					Sort:     "-invoice_no",
				}).Return([]domain.Invoice{
					{Description: "This is test description",
						AccountID: "11x1", CurrencyID: "22x22"},
				}, nil)
				r.EXPECT().Count(context.Background(), "11x1").Return(1, nil)
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
				r.EXPECT().RetrieveById(context.Background(), "22x22").Return(domain.Currency{}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"description":"This is test description"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			mockInvoice: func(r *mock_repository.MockInvoice) {
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		},
		{
			name:    "Wrong Pagination",
			postfix: "?page=notknown&size=smth",
			mockInvoice: func(r *mock_repository.MockInvoice) {
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `Invalid page number`,
			userModel:    &repository.MockedUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rm := mock_repository.NewMockInvoice(c)
			rc := mock_repository.NewMockCurrency(c)
			tt.mockInvoice(rm)
			tt.mockCurrency(rc)

			curencyService := service.NewCurrencyService(rc, cache.NewMemoryCache())

			invoiceService := service.NewInvoiceService(rm, cache.NewMemoryCache(), service.ModulesService{}, config.Config{}, curencyService)

			services := &service.Services{Invoices: invoiceService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultPagination: 20}}}}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/invoices", func(c *gin.Context) {

			}, handler.getAllInvoices)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/invoices"+tt.postfix,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
