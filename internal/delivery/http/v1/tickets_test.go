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

func TestHandler_getTicketById(t *testing.T) {
	type mockRepositoryTicket func(r *mock_repository.MockHelpDesk)

	tests := []struct {
		name         string
		id           string
		mockTicket   mockRepositoryTicket
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Ticket received",
			id:   "17x16",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
				r.EXPECT().RetrieveById(context.Background(), "17x16").Return(domain.MockedHelpDesk, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"description":"They are not attached to client"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "17x1",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name: "Wrong ID",
			id:   "17",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `wrong id`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Not owned ticket",
			id:   "17x16",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
				ticket := domain.MockedHelpDesk
				ticket.ParentID = "11x16"
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

			rm := mock_repository.NewMockHelpDesk(c)
			tt.mockTicket(rm)

			managerService := service.NewHelpDeskService(rm, cache.NewMemoryCache())

			services := &service.Services{HelpDesk: managerService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/tickets/:id", func(c *gin.Context) {

			}, handler.getTicket)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/tickets/"+tt.id,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
