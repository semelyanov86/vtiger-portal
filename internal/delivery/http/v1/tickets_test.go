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

			managerService := service.NewHelpDeskService(rm, cache.NewMemoryCache(), &mock_service.MockCommentServiceInterface{})

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

func TestHandler_getRelatedComments(t *testing.T) {
	type mockRepositoryTicket func(r *mock_repository.MockHelpDesk)
	type mockRepositoryComment func(r *mock_repository.MockComment)

	tests := []struct {
		name         string
		id           string
		mockTicket   mockRepositoryTicket
		mockComment  mockRepositoryComment
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
			mockComment: func(r *mock_repository.MockComment) {
				r.EXPECT().RetrieveFromModule(context.Background(), "17x16").Return([]domain.Comment{domain.MockedComment}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"commentcontent":"This is a test comment."`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "17x1",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
			},
			mockComment: func(r *mock_repository.MockComment) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name: "Wrong ID",
			id:   "17",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
			},
			mockComment: func(r *mock_repository.MockComment) {
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
			mockComment: func(r *mock_repository.MockComment) {
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
			rc := mock_repository.NewMockComment(c)
			tt.mockTicket(rm)
			tt.mockComment(rc)

			commentService := service.NewComments(rc, cache.NewMemoryCache())

			helpDeskService := service.NewHelpDeskService(rm, cache.NewMemoryCache(), commentService)

			services := &service.Services{HelpDesk: helpDeskService, Comments: commentService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/tickets/:id/comments", func(c *gin.Context) {

			}, handler.getComments)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/tickets/"+tt.id+"/comments",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getAllTickets(t *testing.T) {
	type mockRepositoryTicket func(r *mock_repository.MockHelpDesk)
	wrongIdUser := repository.MockedUser
	wrongIdUser.AccountId = ""

	tests := []struct {
		name         string
		postfix      string
		mockTicket   mockRepositoryTicket
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name:    "Tickets received",
			postfix: "?page=1&size=20",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
				r.EXPECT().GetAll(context.Background(), repository.TicketsQueryFilter{
					Page:     1,
					PageSize: 20,
					Client:   "11x1",
				}).Return([]domain.HelpDesk{domain.MockedHelpDesk}, nil)
				r.EXPECT().Count(context.Background(), "11x1").Return(1, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"description":"They are not attached to client"`,
			userModel:    &repository.MockedUser,
		},
		{
			name:    "Anonymous Access",
			postfix: "?page=1&size=20",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name:    "Wrong ID",
			postfix: "?page=1&size=20",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
			},
			statusCode:   http.StatusForbidden,
			responseBody: `Access Not Permitted`,
			userModel:    &wrongIdUser,
		}, {
			name:    "Wrong Pagination",
			postfix: "?page=notknown&size=smth",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
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

			rm := mock_repository.NewMockHelpDesk(c)
			rc := mock_repository.NewMockComment(c)
			tt.mockTicket(rm)

			commentService := service.NewComments(rc, cache.NewMemoryCache())

			helpDeskService := service.NewHelpDeskService(rm, cache.NewMemoryCache(), commentService)

			services := &service.Services{HelpDesk: helpDeskService, Comments: commentService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultPagination: 20}}}}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/tickets", func(c *gin.Context) {

			}, handler.getAllTickets)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/tickets"+tt.postfix,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}