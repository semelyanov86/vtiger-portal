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
	mock_service "github.com/semelyanov86/vtiger-portal/internal/service/mocks"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_getServiceContractById(t *testing.T) {
	type mockRepositorySc func(r *mock_repository.MockServiceContract)

	tests := []struct {
		name         string
		id           string
		mockSc       mockRepositorySc
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Sc received",
			id:   "24x59",
			mockSc: func(r *mock_repository.MockServiceContract) {
				r.EXPECT().RetrieveById(context.Background(), "24x59").Return(domain.MockedServiceContract, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"subject":"Mocked Service Contract"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "17x1",
			mockSc: func(r *mock_repository.MockServiceContract) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name: "Wrong ID",
			id:   "17",
			mockSc: func(r *mock_repository.MockServiceContract) {
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `wrong id`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Not owned ticket",
			id:   "17x16",
			mockSc: func(r *mock_repository.MockServiceContract) {
				ticket := domain.MockedServiceContract
				ticket.ScRelatedTo = "11x16"
				r.EXPECT().RetrieveById(context.Background(), "17x16").Return(ticket, nil)
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

			rm := mock_repository.NewMockServiceContract(c)
			tt.mockSc(rm)

			managerService := service.NewServiceContractsService(rm, cache.NewMemoryCache(), mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{})

			services := &service.Services{ServiceContracts: managerService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/service-contracts/:id", func(c *gin.Context) {

			}, handler.getServiceContract)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/service-contracts/"+tt.id,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
