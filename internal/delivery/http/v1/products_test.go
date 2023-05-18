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
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_getProductById(t *testing.T) {
	type mockRepositoryProduct func(r *mock_repository.MockProduct)
	type mockRepositoryCurrency func(r *mock_repository.MockCurrency)
	type mockRepositoryDocument func(r *mock_repository.MockDocument)

	tests := []struct {
		name         string
		id           string
		mockProduct  mockRepositoryProduct
		mockCurrency mockRepositoryCurrency
		mockDocument mockRepositoryDocument
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Product received",
			id:   "14x9",
			mockProduct: func(r *mock_repository.MockProduct) {
				r.EXPECT().RetrieveById(context.Background(), "14x9").Return(domain.MockedProduct, nil)
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
				r.EXPECT().RetrieveById(context.Background(), "21x11").Return(domain.MockedCurrency, nil)
			},
			mockDocument: func(r *mock_repository.MockDocument) {
				r.EXPECT().RetrieveFile(context.Background(), "14x62").Return(vtiger.File{}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"productname":"Keyboard Logitech"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "17x1",
			mockProduct: func(r *mock_repository.MockProduct) {
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			mockDocument: func(r *mock_repository.MockDocument) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name: "Wrong ID",
			id:   "17",
			mockProduct: func(r *mock_repository.MockProduct) {
			},
			mockCurrency: func(r *mock_repository.MockCurrency) {
			},
			mockDocument: func(r *mock_repository.MockDocument) {
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

			rm := mock_repository.NewMockProduct(c)
			tt.mockProduct(rm)

			rc := mock_repository.NewMockCurrency(c)
			tt.mockCurrency(rc)

			rd := mock_repository.NewMockDocument(c)
			tt.mockDocument(rd)

			productService := service.NewProductService(rm, cache.NewMemoryCache(), service.NewCurrencyService(rc, cache.NewMemoryCache()), rd, service.ModulesService{}, config.Config{})

			services := &service.Services{Products: productService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/products/:id", func(c *gin.Context) {

			}, handler.getProduct)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/products/"+tt.id,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
