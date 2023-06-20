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

func TestHandler_getAllFaqs(t *testing.T) {
	type mockRepositoryFaq func(r *mock_repository.MockFaq)

	tests := []struct {
		name         string
		mockFaq      mockRepositoryFaq
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "faqs received",
			mockFaq: func(r *mock_repository.MockFaq) {
				r.EXPECT().GetAllFaqs(context.Background(), vtiger.PaginationQueryFilter{
					Page:     1,
					PageSize: 20,
					Client:   "11x1",
				}).Return([]domain.Faq{domain.MockedFaq}, nil)
				r.EXPECT().Count(context.Background(), "11x1").Return(1, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"faq_answer":"Just write it and that is it"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			mockFaq: func(r *mock_repository.MockFaq) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rm := mock_repository.NewMockFaq(c)
			tt.mockFaq(rm)

			managerService := service.NewFaqsService(rm, cache.NewMemoryCache(), service.ModulesService{}, config.Config{})

			services := &service.Services{Faqs: managerService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultPagination: 20}}}}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/faqs", func(c *gin.Context) {

			}, handler.getAllFaqs)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/faqs",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
