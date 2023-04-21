package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_getCompany(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
	}{
		{
			name:         "Company received",
			statusCode:   200,
			responseBody: `"organizationname":"Индивидуальный предприниматель Емельянов Сергей Петрович"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			companyService := service.NewCompanyService(repository.NewCompanyMock(), cache.NewMemoryCache())

			services := &service.Services{Company: companyService}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{CompanyId: "23x1"}}}}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/company", func(c *gin.Context) {

			}, handler.getCompany)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/company",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
