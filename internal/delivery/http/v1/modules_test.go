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
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestModule_describeOperation(t *testing.T) {
	type mockRepositoryModule func(r *mock_repository.MockModules)

	tests := []struct {
		name         string
		id           string
		mockModule   mockRepositoryModule
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Module received success",
			id:   "HelpDesk",
			mockModule: func(r *mock_repository.MockModules) {
				r.EXPECT().GetModuleInfo(context.Background(), "HelpDesk").Return(vtiger.MockedModule, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"label":"Assets"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "HelpDesk",
			mockModule: func(r *mock_repository.MockModules) {
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

			rm := mock_repository.NewMockModules(c)
			tt.mockModule(rm)

			modulesService := service.NewModulesService(rm, cache.NewMemoryCache())

			services := &service.Services{Modules: modulesService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/modules/:name", func(c *gin.Context) {

			}, handler.describeModule)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/modules/"+tt.id,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
