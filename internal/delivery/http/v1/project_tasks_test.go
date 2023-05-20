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

func TestHandler_getProjectTasks(t *testing.T) {
	type mockRepositoryProject func(r *mock_repository.MockProject)
	type mockRepositoryProjectTasks func(r *mock_repository.MockProjectTask)

	tests := []struct {
		name            string
		id              string
		mockProject     mockRepositoryProject
		mockProjectTask mockRepositoryProjectTasks
		userModel       *domain.User
		statusCode      int
		responseBody    string
	}{
		{
			name: "Project Tasks received",
			id:   "29x54",
			mockProject: func(r *mock_repository.MockProject) {
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(domain.MockedProject, nil)
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {
				r.EXPECT().GetFromProject(context.Background(), repository.PaginationQueryFilter{
					Page:     1,
					PageSize: 20,
					Client:   "11x1",
					Contact:  "12x11",
					Parent:   "29x54",
				}).Return([]domain.ProjectTask{domain.MockedProjectTask}, nil)
				r.EXPECT().Count(context.Background(), "29x54").Return(1, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"projecttaskname":"Install hosting"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "17x1",
			mockProject: func(r *mock_repository.MockProject) {
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		},
		{
			name: "Not owned ticket",
			id:   "29x54",
			mockProject: func(r *mock_repository.MockProject) {
				project := domain.MockedProject
				project.Linktoaccountscontacts = "11x16"
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(project, nil)
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {
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

			rm := mock_repository.NewMockProject(c)
			tt.mockProject(rm)

			rt := mock_repository.NewMockProjectTask(c)
			tt.mockProjectTask(rt)

			projectsService := service.NewProjectsService(rm, cache.NewMemoryCache(), service.Comments{}, mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{})

			tasksService := service.NewProjectTasksService(rt, cache.NewMemoryCache(), &mock_service.MockCommentServiceInterface{}, mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{}, projectsService)

			services := &service.Services{Projects: projectsService, ProjectTasks: tasksService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultPagination: 20}}}}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/projects/:id/tasks", func(c *gin.Context) {

			}, handler.getAllProjectTasks)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/projects/"+tt.id+"/tasks",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getRelatedCommentsFromProjectTask(t *testing.T) {
	type mockRepositoryProject func(r *mock_repository.MockProject)
	type mockRepositoryComment func(r *mock_repository.MockComment)
	type mockRepositoryProjectTasks func(r *mock_repository.MockProjectTask)

	tests := []struct {
		name            string
		id              string
		taskId          string
		mockProject     mockRepositoryProject
		mockProjectTask mockRepositoryProjectTasks
		mockComment     mockRepositoryComment
		userModel       *domain.User
		statusCode      int
		responseBody    string
	}{
		{
			name:   "Comments received",
			id:     "29x54",
			taskId: "28x56",
			mockProject: func(r *mock_repository.MockProject) {
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(domain.MockedProject, nil)
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {
				r.EXPECT().RetrieveById(context.Background(), "28x56").Return(domain.MockedProjectTask, nil)
			},
			mockComment: func(r *mock_repository.MockComment) {
				r.EXPECT().RetrieveFromModule(context.Background(), "28x56").Return([]domain.Comment{domain.MockedComment}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"commentcontent":"This is a test comment."`,
			userModel:    &repository.MockedUser,
		},
		{
			name:   "Anonymous Access",
			id:     "29x5",
			taskId: "28x56",
			mockProject: func(r *mock_repository.MockProject) {
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {

			},
			mockComment: func(r *mock_repository.MockComment) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name:   "Wrong ID",
			id:     "29",
			taskId: "28x56",
			mockProject: func(r *mock_repository.MockProject) {
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {

			},
			mockComment: func(r *mock_repository.MockComment) {
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `wrong id`,
			userModel:    &repository.MockedUser,
		},
		{
			name:   "Not owned ticket",
			id:     "29x54",
			taskId: "28x56",
			mockProject: func(r *mock_repository.MockProject) {
				project := domain.MockedProject
				project.Linktoaccountscontacts = "11x16"
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(project, nil)
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {
				r.EXPECT().RetrieveById(context.Background(), "28x56").Return(domain.MockedProjectTask, nil)
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

			rm := mock_repository.NewMockProject(c)
			rc := mock_repository.NewMockComment(c)
			tt.mockProject(rm)
			tt.mockComment(rc)
			rt := mock_repository.NewMockProjectTask(c)
			tt.mockProjectTask(rt)

			commentService := service.NewComments(rc, cache.NewMemoryCache(), config.Config{})

			projectsService := service.NewProjectsService(rm, cache.NewMemoryCache(), commentService, mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{})
			projectTaskService := service.NewProjectTasksService(rt, cache.NewMemoryCache(), commentService, mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{}, projectsService)

			services := &service.Services{Projects: projectsService, Comments: commentService, Context: service.MockedContextService{MockedUser: tt.userModel}, ProjectTasks: projectTaskService}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/projects/:id/tasks/:task/comments", func(c *gin.Context) {

			}, handler.getProjectTaskComments)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/projects/"+tt.id+"/tasks/"+tt.taskId+"/comments",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getRelatedDocumentsFromProjectTask(t *testing.T) {
	type mockRepositoryProject func(r *mock_repository.MockProject)
	type mockRepositoryDocument func(r *mock_repository.MockDocument)
	type mockRepositoryProjectTasks func(r *mock_repository.MockProjectTask)

	tests := []struct {
		name            string
		id              string
		mockProject     mockRepositoryProject
		mockProjectTask mockRepositoryProjectTasks
		mockDocument    mockRepositoryDocument
		userModel       *domain.User
		statusCode      int
		responseBody    string
	}{
		{
			name: "Document received",
			id:   "29x54",
			mockProject: func(r *mock_repository.MockProject) {
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(domain.MockedProject, nil)
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {
				r.EXPECT().RetrieveById(context.Background(), "28x56").Return(domain.MockedProjectTask, nil)
			},
			mockDocument: func(r *mock_repository.MockDocument) {
				r.EXPECT().RetrieveFromModule(context.Background(), "28x56").Return([]domain.Document{domain.MockedDocument}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"notes_title":"customer-portal"`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "17x1",
			mockProject: func(r *mock_repository.MockProject) {
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {
			},
			mockDocument: func(r *mock_repository.MockDocument) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name: "Wrong ID",
			id:   "17",
			mockProject: func(r *mock_repository.MockProject) {
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {
			},
			mockDocument: func(r *mock_repository.MockDocument) {
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `wrong id`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Not owned ticket",
			id:   "29x54",
			mockProject: func(r *mock_repository.MockProject) {
				project := domain.MockedProject
				project.Linktoaccountscontacts = "11x16"
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(project, nil)
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {
				r.EXPECT().RetrieveById(context.Background(), "28x56").Return(domain.MockedProjectTask, nil)
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

			rm := mock_repository.NewMockProject(c)
			rt := mock_repository.NewMockProjectTask(c)
			rc := mock_repository.NewMockComment(c)
			rd := mock_repository.NewMockDocument(c)
			tt.mockProject(rm)
			tt.mockDocument(rd)
			tt.mockProjectTask(rt)

			commentService := service.NewComments(rc, cache.NewMemoryCache(), config.Config{})
			documentService := service.NewDocuments(rd, cache.NewMemoryCache())

			projectsService := service.NewProjectsService(rm, cache.NewMemoryCache(), commentService, documentService, service.ModulesService{}, config.Config{})
			projectTasksService := service.NewProjectTasksService(rt, cache.NewMemoryCache(), commentService, documentService, service.ModulesService{}, config.Config{}, projectsService)

			services := &service.Services{Projects: projectsService, Comments: commentService, Documents: documentService, Context: service.MockedContextService{MockedUser: tt.userModel}, ProjectTasks: projectTasksService}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/projects/:id/tasks/:task/documents", func(c *gin.Context) {

			}, handler.getProjectTaskDocuments)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/projects/"+tt.id+"/tasks/28x56/documents",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
