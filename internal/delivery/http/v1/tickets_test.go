package v1

import (
	"bytes"
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

			managerService := service.NewHelpDeskService(rm, cache.NewMemoryCache(), &mock_service.MockCommentServiceInterface{}, mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{})

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

			helpDeskService := service.NewHelpDeskService(rm, cache.NewMemoryCache(), commentService, mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{})

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

			helpDeskService := service.NewHelpDeskService(rm, cache.NewMemoryCache(), commentService, mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{})

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

func TestHandler_getRelatedDocuments(t *testing.T) {
	type mockRepositoryTicket func(r *mock_repository.MockHelpDesk)
	type mockRepositoryDocument func(r *mock_repository.MockDocument)

	tests := []struct {
		name         string
		id           string
		mockTicket   mockRepositoryTicket
		mockDocument mockRepositoryDocument
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Document received",
			id:   "17x16",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
				r.EXPECT().RetrieveById(context.Background(), "17x16").Return(domain.MockedHelpDesk, nil)
			},
			mockDocument: func(r *mock_repository.MockDocument) {
				r.EXPECT().RetrieveFromModule(context.Background(), "17x16").Return([]domain.Document{domain.MockedDocument}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"notes_title":"customer-portal"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "17x1",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
			},
			mockDocument: func(r *mock_repository.MockDocument) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name: "Wrong ID",
			id:   "17",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
			},
			mockDocument: func(r *mock_repository.MockDocument) {
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
			mockDocument: func(r *mock_repository.MockDocument) {
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
			rd := mock_repository.NewMockDocument(c)
			tt.mockTicket(rm)
			tt.mockDocument(rd)

			commentService := service.NewComments(rc, cache.NewMemoryCache())
			documentService := service.NewDocuments(rd, cache.NewMemoryCache())

			helpDeskService := service.NewHelpDeskService(rm, cache.NewMemoryCache(), commentService, documentService, service.ModulesService{}, config.Config{})

			services := &service.Services{HelpDesk: helpDeskService, Comments: commentService, Documents: documentService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/tickets/:id/documents", func(c *gin.Context) {

			}, handler.getDocuments)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/tickets/"+tt.id+"/documents",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getFile(t *testing.T) {
	type mockRepositoryTicket func(r *mock_repository.MockHelpDesk)
	type mockRepositoryDocument func(r *mock_repository.MockDocument)

	tests := []struct {
		name         string
		id           string
		fileId       string
		mockTicket   mockRepositoryTicket
		mockDocument mockRepositoryDocument
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name:   "File received",
			id:     "17x16",
			fileId: "15x42",
			mockTicket: func(r *mock_repository.MockHelpDesk) {

			},
			mockDocument: func(r *mock_repository.MockDocument) {
				r.EXPECT().RetrieveFile(context.Background(), "15x42").Return(domain.MockedFile, nil)
				r.EXPECT().RetrieveFromModule(context.Background(), "17x16").Return([]domain.Document{domain.MockedDocument}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"filecontents":"iVBORw0KGgoAAAANSUhEUgAAAnAAAAD1CAIAAA"`,
			userModel:    &repository.MockedUser,
		},
		{
			name:   "Anonymous Access",
			id:     "17x1",
			fileId: "15x42",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
			},
			mockDocument: func(r *mock_repository.MockDocument) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name:   "Wrong ID",
			id:     "17",
			fileId: "15x42",
			mockTicket: func(r *mock_repository.MockHelpDesk) {
			},
			mockDocument: func(r *mock_repository.MockDocument) {
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `wrong id`,
			userModel:    &repository.MockedUser,
		},
		{
			name:   "Not owned file",
			id:     "17x16",
			fileId: "15x42",
			mockTicket: func(r *mock_repository.MockHelpDesk) {

			},
			mockDocument: func(r *mock_repository.MockDocument) {
				document := domain.MockedDocument
				document.Imageattachmentids = "15x15"
				r.EXPECT().RetrieveFromModule(context.Background(), "17x16").Return([]domain.Document{document}, nil)
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
			rd := mock_repository.NewMockDocument(c)
			tt.mockTicket(rm)
			tt.mockDocument(rd)

			commentService := service.NewComments(rc, cache.NewMemoryCache())
			documentService := service.NewDocuments(rd, cache.NewMemoryCache())

			helpDeskService := service.NewHelpDeskService(rm, cache.NewMemoryCache(), commentService, documentService, service.ModulesService{}, config.Config{})

			services := &service.Services{HelpDesk: helpDeskService, Comments: commentService, Documents: documentService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/tickets/:id/file/:file", func(c *gin.Context) {

			}, handler.getFile)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/tickets/"+tt.id+"/file/"+tt.fileId,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_createTicket(t *testing.T) {
	type mockRepositoryModule func(r *mock_repository.MockModules)

	tests := []struct {
		name         string
		mockModule   mockRepositoryModule
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Ticket created",
			mockModule: func(r *mock_repository.MockModules) {
				r.EXPECT().GetModuleInfo(context.Background(), "HelpDesk").Return(vtiger.MockedModule, nil)
			},
			statusCode:   http.StatusCreated,
			responseBody: `"ticket_no":"TICKET_28"`,
			userModel:    &repository.MockedUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rc := mock_repository.NewMockComment(c)
			rd := mock_repository.NewMockDocument(c)
			rmm := mock_repository.NewMockModules(c)
			tt.mockModule(rmm)

			commentService := service.NewComments(rc, cache.NewMemoryCache())
			documentService := service.NewDocuments(rd, cache.NewMemoryCache())

			helpDeskService := service.NewHelpDeskService(repository.HelpDeskMockRepository{}, cache.NewMemoryCache(), commentService, documentService, service.NewModulesService(rmm, cache.NewMemoryCache()), config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultUser: "19x1"}}})

			services := &service.Services{HelpDesk: helpDeskService, Comments: commentService, Documents: documentService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/api/v1/tickets", func(c *gin.Context) {

			}, handler.createTicket)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/tickets",
				bytes.NewBufferString(`{
  "ticket_title": "Problem with internet",
  "ticketpriorities": "Normal",
  "ticketseverities": "Minor",
  "ticketcategories": "Big Problem",
  "description": "There are no internet in my appartment."
}`))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_updateTicket(t *testing.T) {
	type mockRepositoryModule func(r *mock_repository.MockModules)

	otherUser := repository.MockedUser
	otherUser.AccountId = "11x223"

	tests := []struct {
		name         string
		mockModule   mockRepositoryModule
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Ticket updated",
			mockModule: func(r *mock_repository.MockModules) {
				r.EXPECT().GetModuleInfo(context.Background(), "HelpDesk").Return(vtiger.MockedModule, nil)
			},
			statusCode:   http.StatusAccepted,
			responseBody: `"ticket_no":"TICKET_28"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Update not permitted",
			mockModule: func(r *mock_repository.MockModules) {
			},
			statusCode:   http.StatusForbidden,
			responseBody: `"error":"Access Not Permitted"`,
			userModel:    &otherUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rc := mock_repository.NewMockComment(c)
			rd := mock_repository.NewMockDocument(c)
			rmm := mock_repository.NewMockModules(c)
			tt.mockModule(rmm)

			commentService := service.NewComments(rc, cache.NewMemoryCache())
			documentService := service.NewDocuments(rd, cache.NewMemoryCache())

			helpDeskService := service.NewHelpDeskService(repository.HelpDeskMockRepository{}, cache.NewMemoryCache(), commentService, documentService, service.NewModulesService(rmm, cache.NewMemoryCache()), config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultUser: "19x1"}}})

			services := &service.Services{HelpDesk: helpDeskService, Comments: commentService, Documents: documentService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.PUT("/api/v1/tickets/:id", func(c *gin.Context) {

			}, handler.updateTicket)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/v1/tickets/17x28",
				bytes.NewBufferString(`{
  "ticket_title": "Problem with internet",
  "ticketpriorities": "Normal",
  "ticketseverities": "Minor",
  "ticketcategories": "Big Problem",
  "description": "There are no internet in my appartment."
}`))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
