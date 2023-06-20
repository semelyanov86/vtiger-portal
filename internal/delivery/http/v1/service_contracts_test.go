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
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
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

func TestHandler_getAllServiceContracts(t *testing.T) {
	type mockRepositorySc func(r *mock_repository.MockServiceContract)
	wrongIdUser := repository.MockedUser
	wrongIdUser.AccountId = ""

	tests := []struct {
		name         string
		postfix      string
		mockSc       mockRepositorySc
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name:    "Service Contracts received",
			postfix: "?page=1&size=20",
			mockSc: func(r *mock_repository.MockServiceContract) {
				r.EXPECT().GetAll(context.Background(), vtiger.PaginationQueryFilter{
					Page:     1,
					PageSize: 20,
					Client:   "11x1",
					Contact:  "12x11",
				}).Return([]domain.ServiceContract{domain.MockedServiceContract}, nil)
				r.EXPECT().Count(context.Background(), "11x1", "12x11").Return(1, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"subject":"Mocked Service Contract"`,
			userModel:    &repository.MockedUser,
		},
		{
			name:    "Anonymous Access",
			postfix: "?page=1&size=20",
			mockSc: func(r *mock_repository.MockServiceContract) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name:    "Wrong ID",
			postfix: "?page=1&size=20",
			mockSc: func(r *mock_repository.MockServiceContract) {
			},
			statusCode:   http.StatusForbidden,
			responseBody: `Access Not Permitted`,
			userModel:    &wrongIdUser,
		}, {
			name:    "Wrong Pagination",
			postfix: "?page=notknown&size=smth",
			mockSc: func(r *mock_repository.MockServiceContract) {
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

			rm := mock_repository.NewMockServiceContract(c)
			tt.mockSc(rm)

			scService := service.NewServiceContractsService(rm, cache.NewMemoryCache(), mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{})

			services := &service.Services{ServiceContracts: scService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultPagination: 20}}}}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/service-contracts", func(c *gin.Context) {

			}, handler.getAllServiceContracts)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/service-contracts"+tt.postfix,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
