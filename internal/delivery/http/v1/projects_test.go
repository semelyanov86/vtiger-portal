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

func TestHandler_getProjectById(t *testing.T) {
	type mockRepositoryProject func(r *mock_repository.MockProject)
	type mockRepositoryTask func(r *mock_repository.MockProjectTask)

	tests := []struct {
		name            string
		id              string
		mockProject     mockRepositoryProject
		mockProjectTask mockRepositoryTask
		userModel       *domain.User
		statusCode      int
		responseBody    string
	}{
		{
			name: "Project received",
			id:   "29x54",
			mockProject: func(r *mock_repository.MockProject) {
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(domain.MockedProject, nil)
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {
				r.EXPECT().GetFromProject(context.Background(), vtiger.PaginationQueryFilter{
					Page:     1,
					PageSize: 100,
					Client:   "",
					Contact:  "",
					Parent:   "29x54",
					Sort:     "",
					Filters:  nil,
					Search:   "",
				}).Return([]domain.ProjectTask{domain.MockedProjectTask}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"projectname":"Website development"`,
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
		}, {
			name: "Wrong ID",
			id:   "17",
			mockProject: func(r *mock_repository.MockProject) {
			},
			mockProjectTask: func(r *mock_repository.MockProjectTask) {

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
			tt.mockProject(rm)
			tt.mockProjectTask(rt)

			managerService := service.NewProjectsService(rm, cache.NewMemoryCache(), &mock_service.MockCommentServiceInterface{}, mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{}, rt)

			services := &service.Services{Projects: managerService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/projects/:id", func(c *gin.Context) {

			}, handler.getProject)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/projects/"+tt.id,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getRelatedCommentsFromProject(t *testing.T) {
	type mockRepositoryProject func(r *mock_repository.MockProject)
	type mockRepositoryComment func(r *mock_repository.MockComment)
	type mockRepositoryManager func(r *mock_repository.MockManagers)

	tests := []struct {
		name         string
		id           string
		mockProject  mockRepositoryProject
		mockComment  mockRepositoryComment
		mockManager  mockRepositoryManager
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Project received",
			id:   "29x54",
			mockProject: func(r *mock_repository.MockProject) {
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(domain.MockedProject, nil)
			},
			mockComment: func(r *mock_repository.MockComment) {
				r.EXPECT().RetrieveFromModule(context.Background(), "29x54").Return([]domain.Comment{domain.MockedComment}, nil)
			},
			mockManager: func(r *mock_repository.MockManagers) {
				r.EXPECT().RetrieveById(context.Background(), "19x1").Return(domain.MockedManager, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"commentcontent":"This is a test comment."`,
			userModel:    &repository.MockedUser,
		},
		{
			name: "Anonymous Access",
			id:   "29x5",
			mockProject: func(r *mock_repository.MockProject) {
			},
			mockComment: func(r *mock_repository.MockComment) {
			},
			mockManager: func(r *mock_repository.MockManagers) {

			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name: "Wrong ID",
			id:   "29",
			mockProject: func(r *mock_repository.MockProject) {
			},
			mockComment: func(r *mock_repository.MockComment) {
			},
			mockManager: func(r *mock_repository.MockManagers) {

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
			mockComment: func(r *mock_repository.MockComment) {
			},
			mockManager: func(r *mock_repository.MockManagers) {

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
			rman := mock_repository.NewMockManagers(c)
			tt.mockProject(rm)
			tt.mockComment(rc)
			tt.mockManager(rman)

			commentService := service.NewComments(rc, cache.NewMemoryCache(), config.Config{}, service.UsersService{}, service.NewManagerService(rman, cache.NewMemoryCache()))

			projectsService := service.NewProjectsService(rm, cache.NewMemoryCache(), commentService, mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{}, repository.ProjectTaskCrm{})

			services := &service.Services{Projects: projectsService, Comments: commentService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/projects/:id/comments", func(c *gin.Context) {

			}, handler.getProjectComments)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/projects/"+tt.id+"/comments",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getAllProjects(t *testing.T) {
	type mockRepositoryProject func(r *mock_repository.MockProject)
	wrongIdUser := repository.MockedUser
	wrongIdUser.AccountId = ""

	tests := []struct {
		name         string
		postfix      string
		mockProject  mockRepositoryProject
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name:    "Projects received",
			postfix: "?page=1&size=20",
			mockProject: func(r *mock_repository.MockProject) {
				r.EXPECT().GetAll(context.Background(), vtiger.PaginationQueryFilter{
					Page:     1,
					PageSize: 20,
					Client:   "11x1",
					Contact:  "12x11",
					Sort:     "-project_no",
					Search:   "",
				}).Return([]domain.Project{domain.MockedProject}, nil)
				r.EXPECT().Count(context.Background(), "11x1", "12x11").Return(1, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"projectname":"Website development"`,
			userModel:    &repository.MockedUser,
		},
		{
			name:    "Anonymous Access",
			postfix: "?page=1&size=20",
			mockProject: func(r *mock_repository.MockProject) {
			},
			statusCode:   http.StatusUnauthorized,
			responseBody: `"error":"Anonymous Access",`,
			userModel:    domain.AnonymousUser,
		}, {
			name:    "Wrong ID",
			postfix: "?page=1&size=20",
			mockProject: func(r *mock_repository.MockProject) {
			},
			statusCode:   http.StatusForbidden,
			responseBody: `Access Not Permitted`,
			userModel:    &wrongIdUser,
		}, {
			name:    "Wrong Pagination",
			postfix: "?page=notknown&size=smth",
			mockProject: func(r *mock_repository.MockProject) {
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

			rm := mock_repository.NewMockProject(c)
			rc := mock_repository.NewMockComment(c)
			tt.mockProject(rm)

			commentService := service.NewComments(rc, cache.NewMemoryCache(), config.Config{}, service.UsersService{}, service.ManagerService{})

			projectsService := service.NewProjectsService(rm, cache.NewMemoryCache(), commentService, mock_service.NewMockDocumentServiceInterface(c), service.ModulesService{}, config.Config{}, repository.ProjectTaskCrm{})

			services := &service.Services{Projects: projectsService, Comments: commentService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultPagination: 20}}}}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/projects", func(c *gin.Context) {

			}, handler.getAllProjects)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/projects"+tt.postfix,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getRelatedDocumentsFromProject(t *testing.T) {
	type mockRepositoryProject func(r *mock_repository.MockProject)
	type mockRepositoryDocument func(r *mock_repository.MockDocument)

	tests := []struct {
		name         string
		id           string
		mockProject  mockRepositoryProject
		mockDocument mockRepositoryDocument
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name: "Document received",
			id:   "29x54",
			mockProject: func(r *mock_repository.MockProject) {
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(domain.MockedProject, nil)
			},
			mockDocument: func(r *mock_repository.MockDocument) {
				r.EXPECT().RetrieveFromModule(context.Background(), "29x54").Return([]domain.Document{domain.MockedDocument}, nil)
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
			rc := mock_repository.NewMockComment(c)
			rd := mock_repository.NewMockDocument(c)
			tt.mockProject(rm)
			tt.mockDocument(rd)

			commentService := service.NewComments(rc, cache.NewMemoryCache(), config.Config{}, service.UsersService{}, service.ManagerService{})
			documentService := service.NewDocuments(rd, cache.NewMemoryCache(), config.Config{})

			projectsService := service.NewProjectsService(rm, cache.NewMemoryCache(), commentService, documentService, service.ModulesService{}, config.Config{}, repository.ProjectTaskCrm{})

			services := &service.Services{Projects: projectsService, Comments: commentService, Documents: documentService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/projects/:id/documents", func(c *gin.Context) {

			}, handler.getProjectDocuments)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/projects/"+tt.id+"/documents",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_getFileFromProject(t *testing.T) {
	type mockRepositoryProject func(r *mock_repository.MockProject)
	type mockRepositoryDocument func(r *mock_repository.MockDocument)

	tests := []struct {
		name         string
		id           string
		fileId       string
		mockProject  mockRepositoryProject
		mockDocument mockRepositoryDocument
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name:   "File received",
			id:     "29x54",
			fileId: "15x42",
			mockProject: func(r *mock_repository.MockProject) {
				project := domain.MockedProject
				project.Linktoaccountscontacts = "12x11"
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(project, nil)
			},
			mockDocument: func(r *mock_repository.MockDocument) {
				r.EXPECT().RetrieveFile(context.Background(), "15x42").Return(domain.MockedFile, nil)
				r.EXPECT().RetrieveFromModule(context.Background(), "29x54").Return([]domain.Document{domain.MockedDocument}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: `"filecontents":"iVBORw0KGgoAAAANSUhEUgAAAnAAAAD1CAIAAA"`,
			userModel:    &repository.MockedUser,
		},
		{
			name:   "Anonymous Access",
			id:     "29x1",
			fileId: "15x42",
			mockProject: func(r *mock_repository.MockProject) {
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
			mockProject: func(r *mock_repository.MockProject) {
			},
			mockDocument: func(r *mock_repository.MockDocument) {
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `wrong id`,
			userModel:    &repository.MockedUser,
		},
		{
			name:   "Not owned file",
			id:     "29x54",
			fileId: "15x42",
			mockProject: func(r *mock_repository.MockProject) {
				project := domain.MockedProject
				project.Linktoaccountscontacts = "11x16"
				r.EXPECT().RetrieveById(context.Background(), "29x54").Return(project, nil)
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
			rc := mock_repository.NewMockComment(c)
			rd := mock_repository.NewMockDocument(c)
			tt.mockProject(rm)
			tt.mockDocument(rd)

			commentService := service.NewComments(rc, cache.NewMemoryCache(), config.Config{}, service.UsersService{}, service.ManagerService{})
			documentService := service.NewDocuments(rd, cache.NewMemoryCache(), config.Config{})

			projectsService := service.NewProjectsService(rm, cache.NewMemoryCache(), commentService, documentService, service.ModulesService{}, config.Config{}, repository.ProjectTaskCrm{})

			services := &service.Services{Projects: projectsService, Comments: commentService, Documents: documentService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/projects/:id/file/:file", func(c *gin.Context) {

			}, handler.getProjectFile)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/projects/"+tt.id+"/file/"+tt.fileId,
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
