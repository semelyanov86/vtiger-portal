package v1

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Connector struct {
}

func (c *Connector) Create(ctx context.Context, element string, data map[string]any) (*vtiger.VtigerResponse[map[string]any], error) {
	return &vtiger.VtigerResponse[map[string]any]{
		Success: true,
		Result:  map[string]any{"email": "test@mail.ru"},
		Error:   vtiger.ErrorData{},
	}, nil
}

func TestHandler_createLead(t *testing.T) {
	tests := []struct {
		name         string
		userModel    *domain.User
		body         string
		statusCode   int
		responseBody string
	}{
		{
			name: "lead created",
			body: `{
			  "lastname": "Boris Johns",
			  "email": "test@mail.ru",
			  "phone": "+79578889654"
			}`,
			statusCode:   http.StatusCreated,
			responseBody: `"email":"test@mail.ru"`,
			userModel:    &repository.MockedUser,
		},
		{
			name:         "Anonymous Access",
			statusCode:   http.StatusUnauthorized,
			body:         "",
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rm := repository.NewLeadCrmConcrete(config.Config{}, &Connector{})

			leadsService := service.NewLeads(rm, config.Config{})

			services := &service.Services{Leads: leadsService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultPagination: 20}}}}

			// Init Endpoint
			r := gin.New()
			r.POST("/api/v1/leads", func(c *gin.Context) {

			}, handler.createLead)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/leads",
				bytes.NewBufferString(tt.body))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
