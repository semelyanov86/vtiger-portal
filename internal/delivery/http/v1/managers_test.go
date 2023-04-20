package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
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

func TestHandler_getUserById(t *testing.T) {
	type mockRepositoryManager func(r *mock_repository.MockManagers)

	tests := []struct {
		name         string
		id           string
		mockManager  mockRepositoryManager
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "User received success",
			id:   "19x1",
			mockManager: func(r *mock_repository.MockManagers) {
				r.EXPECT().RetrieveById(context.Background(), "19x1").Return(domain.MockedManager, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"description":"This is description for administrator user"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "19x1",
			mockManager: func(r *mock_repository.MockManagers) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name: "Wrong ID",
			id:   "19",
			mockManager: func(r *mock_repository.MockManagers) {
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

			rm := mock_repository.NewMockManagers(c)
			tt.mockManager(rm)

			managerService := service.NewManagerService(rm, cache.NewMemoryCache())

			services := &service.Services{Managers: managerService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/managers/:id", func(c *gin.Context) {

			}, handler.getById)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/managers/"+tt.id,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
