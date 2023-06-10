package v1

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	mock_repository "github.com/semelyanov86/vtiger-portal/internal/repository/mocks"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

var mockedUser = []domain.User{
	{
		Id:             0,
		Crmid:          "12x11",
		FirstName:      "Sergei",
		LastName:       "Emelianov",
		Description:    "Test Description",
		AccountId:      "11x6",
		AccountName:    "",
		Title:          "Manager",
		Department:     "IT",
		Email:          "emelyanov86@km.ru",
		Password:       domain.Password{},
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
		IsActive:       false,
		MailingCity:    "Cheboksary",
		MailingStreet:  "",
		MailingCountry: "",
		OtherCountry:   "",
		MailingState:   "",
		MailingPoBox:   "",
		OtherCity:      "",
		OtherState:     "",
		MailingZip:     "",
		OtherZip:       "",
		OtherStreet:    "",
		OtherPoBox:     "",
		Image:          "",
		Version:        0,
		Code:           "1234",
	},
}

func TestHandler_createUser(t *testing.T) {
	type mockRepositoryCrm func(r *mock_repository.MockUsersCrm)
	type mockRepositoryAccount func(r *mock_repository.MockAccount)

	const email = "emelyanov86@Km.ru"
	tests := []struct {
		name         string
		body         string
		mockCrm      mockRepositoryCrm
		mockAccount  mockRepositoryAccount
		statusCode   int
		responseBody string
	}{
		{
			name: "Correct user",
			body: fmt.Sprintf(`{
				"email": "%s",
				"code": "1234",
				"password": "ThisPasswordCool"
			}`, email),
			mockCrm: func(r *mock_repository.MockUsersCrm) {
				r.EXPECT().FindByEmail(context.Background(), email).Return(mockedUser, nil)
			},
			mockAccount: func(r *mock_repository.MockAccount) {
				r.EXPECT().RetrieveById(context.Background(), "11x6").Return(domain.Account{}, nil)
			},
			statusCode:   201,
			responseBody: `"id":1,"crmid":"12x11",`,
		}, {
			name: "wrong email",
			body: fmt.Sprintf(`{
				"email": "%s",
				"code": "1234",
				"password": "ThisPasswordCool"
			}`, "wront-email"),
			mockCrm: func(r *mock_repository.MockUsersCrm) {
				//r.EXPECT().FindByEmail(context.Background(), email).Return(mockedUser, nil)
			},
			mockAccount: func(r *mock_repository.MockAccount) {
			},
			statusCode:   400,
			responseBody: `Key: 'UserSignUpInput.Email' Error:Field validation for 'Email' failed on the 'email' tag`,
		}, {
			name: "wrong code",
			body: fmt.Sprintf(`{
				"email": "%s",
				"code": "",
				"password": "ThisPasswordCool"
			}`, email),
			mockCrm: func(r *mock_repository.MockUsersCrm) {
				//r.EXPECT().FindByEmail(context.Background(), email).Return(mockedUser, nil)
			},
			mockAccount: func(r *mock_repository.MockAccount) {
			},
			statusCode:   400,
			responseBody: `Key: 'UserSignUpInput.Code' Error:Field validation for 'Code' failed on the 'required' tag`,
		}, {
			name: "wrong password",
			body: fmt.Sprintf(`{
				"email": "%s",
				"code": "1143",
				"password": "123"
			}`, email),
			mockCrm: func(r *mock_repository.MockUsersCrm) {
				//r.EXPECT().FindByEmail(context.Background(), email).Return(mockedUser, nil)
			},
			mockAccount: func(r *mock_repository.MockAccount) {
			},
			statusCode:   400,
			responseBody: `Key: 'UserSignUpInput.Password' Error:Field validation for 'Password' failed on the 'min' tag`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rdb := repository.NewUsersMock()

			rcrm := mock_repository.NewMockUsersCrm(c)
			ra := mock_repository.NewMockAccount(c)
			tt.mockCrm(rcrm)
			tt.mockAccount(ra)

			companyService := service.NewCompanyService(repository.NewCompanyMock(), cache.NewMemoryCache())
			usersService := service.NewUsersService(rdb, rcrm, &wg, service.NewMockEmailService(), companyService, mock_repository.NewMockTokens(c), mock_repository.NewMockDocument(c), cache.NewMemoryCache(), service.NewAccountService(ra, cache.NewMemoryCache()))

			services := &service.Services{Users: usersService}
			handler := Handler{services: services, config: &config.Config{}}

			// Init Endpoint
			r := gin.New()
			r.POST("/api/v1/users", func(c *gin.Context) {

			}, handler.userSignUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/users",
				bytes.NewBufferString(tt.body))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_Login(t *testing.T) {
	type mockRepositoryUser func(r *mock_repository.MockUsers)
	type mockRepositoryToken func(r *mock_repository.MockTokens)
	mockedUserModel := repository.MockedUser

	mockedToken := &domain.Token{
		ID:        1,
		Plaintext: "SOME_TEXT",
		Hash:      "SOME_HASH",
		UserId:    1,
		Expiry:    time.Time{},
		Scope:     domain.ScopeAuthentication,
	}

	tests := []struct {
		name         string
		email        string
		password     string
		mockUser     mockRepositoryUser
		mockToken    mockRepositoryToken
		statusCode   int
		responseBody string
	}{
		{
			name:     "Login Successfull",
			email:    "emelyanov86@km.ru",
			password: "GoodPasswordHele",
			mockUser: func(r *mock_repository.MockUsers) {
				var pass domain.Password
				pass.Set("GoodPasswordHele")
				mockedUserModel.Password = pass
				r.EXPECT().GetByEmail(context.Background(), "emelyanov86@km.ru").Return(mockedUserModel, nil)
			},
			mockToken: func(r *mock_repository.MockTokens) {
				r.EXPECT().New(context.Background(), int64(1), 24*time.Hour*90, domain.ScopeAuthentication).Return(mockedToken, nil)
			},
			statusCode:   201,
			responseBody: `"id":1,"token":"SOME_TEXT"`,
		}, {
			name:     "Wrong Email",
			email:    "emelyanov8611@km.ru",
			password: "GoodPasswordHele",
			mockUser: func(r *mock_repository.MockUsers) {
				var pass domain.Password
				pass.Set("GoodPasswordHele")
				mockedUserModel.Password = pass
				r.EXPECT().GetByEmail(context.Background(), "emelyanov8611@km.ru").Return(mockedUserModel, repository.ErrRecordNotFound)
			},
			mockToken: func(r *mock_repository.MockTokens) {
				//r.EXPECT().New(context.Background(), int64(1), 24*time.Hour*90, domain.ScopeAuthentication).Return(mockedToken, nil)
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `User with this email not found`,
		}, {
			name:     "Wrong Password",
			email:    "emelyanov86@km.ru",
			password: "PasswordWrong",
			mockUser: func(r *mock_repository.MockUsers) {
				r.EXPECT().GetByEmail(context.Background(), "emelyanov86@km.ru").Return(mockedUserModel, service.ErrPasswordDoesNotMatch)
			},
			mockToken: func(r *mock_repository.MockTokens) {
				//r.EXPECT().New(context.Background(), int64(1), 24*time.Hour*90, domain.ScopeAuthentication).Return(mockedToken, nil)
			},
			statusCode:   http.StatusUnprocessableEntity,
			responseBody: `Password you passed to us is incorrect`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rt := mock_repository.NewMockTokens(c)

			ru := mock_repository.NewMockUsers(c)
			tt.mockUser(ru)
			tt.mockToken(rt)

			companyService := service.NewCompanyService(repository.NewCompanyMock(), cache.NewMemoryCache())
			tokensService := service.NewTokensService(rt, ru, service.NewMockEmailService(), config.Config{}, companyService)

			services := &service.Services{Tokens: tokensService}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/api/v1/users/login", func(c *gin.Context) {

			}, handler.signIn)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/users/login",
				bytes.NewBufferString(fmt.Sprintf(`{
				  "email": "%s",
				  "password": "%s"
				}`, tt.email, tt.password)))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_userInfo(t *testing.T) {
	type mockRepositoryAccount func(r *mock_repository.MockAccount)

	tests := []struct {
		name         string
		userModel    *domain.User
		mockAccount  mockRepositoryAccount
		statusCode   int
		responseBody string
	}{
		{
			name:         "Get user info",
			statusCode:   http.StatusOK,
			userModel:    &repository.MockedUser,
			responseBody: `"email":"emelyanov86@km.ru",`,
			mockAccount: func(r *mock_repository.MockAccount) {
				r.EXPECT().RetrieveById(context.Background(), "11x1").Return(domain.Account{}, nil)
			},
		}, {
			name:         "Get user if anonymous",
			statusCode:   http.StatusUnauthorized,
			userModel:    domain.AnonymousUser,
			responseBody: `"error":"Anonymous Access",`,
			mockAccount: func(r *mock_repository.MockAccount) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rc := repository.NewUsersCrmMock(repository.MockedUser)

			rd := repository.NewUsersMock()
			ra := mock_repository.NewMockAccount(c)
			tt.mockAccount(ra)

			companyService := service.NewCompanyService(repository.NewCompanyMock(), cache.NewMemoryCache())
			usersService := service.NewUsersService(rd, rc, &wg, service.NewMockEmailService(), companyService, mock_repository.NewMockTokens(c), mock_repository.NewMockDocument(c), cache.NewMemoryCache(), service.NewAccountService(ra, cache.NewMemoryCache()))

			services := &service.Services{Users: usersService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()

			r.GET("/api/v1/users/my", func(c *gin.Context) {

			}, handler.getUserInfo)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/users/my", nil)

			// Make Request
			r.ServeHTTP(w, req)

			wg.Wait()
			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_restorePassword(t *testing.T) {
	type mockRepositoryUser func(r *mock_repository.MockUsers)
	type mockRepositoryToken func(r *mock_repository.MockTokens)

	mockedUserModel := repository.MockedUser

	mockedToken := &domain.Token{
		ID:        2,
		Plaintext: "PASS_TEXT",
		Hash:      "PASS_HASH",
		UserId:    1,
		Expiry:    time.Time{},
		Scope:     domain.ScopePasswordReset,
	}

	tests := []struct {
		name         string
		email        string
		mockUser     mockRepositoryUser
		mockToken    mockRepositoryToken
		statusCode   int
		responseBody string
	}{
		{
			name:  "Token sent",
			email: "emelyanov86@km.ru",
			mockUser: func(r *mock_repository.MockUsers) {
				r.EXPECT().GetByEmail(context.Background(), "emelyanov86@km.ru").Return(mockedUserModel, nil)
			},
			mockToken: func(r *mock_repository.MockTokens) {
				r.EXPECT().New(context.Background(), int64(1), 45*time.Minute, domain.ScopePasswordReset).Return(mockedToken, nil)
			},
			statusCode:   201,
			responseBody: `"message":"Token successfully created, please check an email"`,
		},
		{
			name:  "Token sent",
			email: "wrong-email",
			mockUser: func(r *mock_repository.MockUsers) {
			},
			mockToken: func(r *mock_repository.MockTokens) {
			},
			statusCode:   400,
			responseBody: `Error:Field validation for 'Email' failed on the 'email' tag`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rt := mock_repository.NewMockTokens(c)

			ru := mock_repository.NewMockUsers(c)
			tt.mockUser(ru)
			tt.mockToken(rt)

			companyService := service.NewCompanyService(repository.NewCompanyMock(), cache.NewMemoryCache())
			tokensService := service.NewTokensService(rt, ru, service.NewMockEmailService(), config.Config{}, companyService)

			services := &service.Services{Tokens: tokensService}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/api/v1/users/restore", func(c *gin.Context) {

			}, handler.sendRestoreToken)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/users/restore",
				bytes.NewBufferString(fmt.Sprintf(`{
				  "email": "%s"
				}`, tt.email)))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_resetPassword(t *testing.T) {
	type mockRepositoryToken func(r *mock_repository.MockTokens)

	tests := []struct {
		name         string
		token        string
		password     string
		mockToken    mockRepositoryToken
		statusCode   int
		responseBody string
	}{
		{
			name:     "Password updated",
			token:    "ISBTSR6CWE7ZNQGWN2QIHRFGNA",
			password: "NewPasswordHere",
			mockToken: func(r *mock_repository.MockTokens) {
				r.EXPECT().DeleteAllForUser(context.Background(), domain.ScopePasswordReset, int64(1)).Return(nil)
			},
			statusCode:   http.StatusAccepted,
			responseBody: `"email":"emelyanov86@km.ru"`,
		}, {
			name:     "Wrong Token",
			token:    "",
			password: "NewPasswordHere",
			mockToken: func(r *mock_repository.MockTokens) {
				//r.EXPECT().DeleteAllForUser(context.Background(), domain.ScopePasswordReset, int64(1)).Return(nil)
			},
			statusCode:   http.StatusBadRequest,
			responseBody: `"error":"Validation Error",`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rt := mock_repository.NewMockTokens(c)

			tt.mockToken(rt)

			companyService := service.NewCompanyService(repository.NewCompanyMock(), cache.NewMemoryCache())
			rc := repository.NewUsersCrmMock(repository.MockedUser)
			usersService := service.NewUsersService(repository.NewUsersMock(), rc, &wg, service.NewMockEmailService(), companyService, rt, mock_repository.NewMockDocument(c), cache.NewMemoryCache(), service.AccountService{})

			services := &service.Services{Users: usersService}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.PUT("/api/v1/users/password", func(c *gin.Context) {

			}, handler.resetPassword)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/v1/users/password",
				bytes.NewBufferString(fmt.Sprintf(`{
				  "token": "%s", "password": "%s"
				}`, tt.token, tt.password)))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}

func TestHandler_allUsers(t *testing.T) {
	tests := []struct {
		name         string
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name:         "Get all users",
			statusCode:   http.StatusOK,
			userModel:    &repository.MockedUser,
			responseBody: `"email":"emelyanov86@km.ru",`,
		}, {
			name:         "Get user if anonymous",
			statusCode:   http.StatusUnauthorized,
			userModel:    domain.AnonymousUser,
			responseBody: `"error":"Anonymous Access",`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rc := repository.NewUsersCrmMock(repository.MockedUser)

			rd := repository.NewUsersMock()

			companyService := service.NewCompanyService(repository.NewCompanyMock(), cache.NewMemoryCache())
			usersService := service.NewUsersService(rd, rc, &wg, service.NewMockEmailService(), companyService, mock_repository.NewMockTokens(c), mock_repository.NewMockDocument(c), cache.NewMemoryCache(), service.AccountService{})

			services := &service.Services{Users: usersService, Context: service.MockedContextService{MockedUser: tt.userModel}}
			handler := Handler{services: services, config: &config.Config{Vtiger: config.VtigerConfig{Business: config.VtigerBusinessConfig{DefaultPagination: 20}}}}

			// Init Endpoint
			r := gin.New()

			r.GET("/api/v1/users/all", func(c *gin.Context) {

			}, handler.usersFromAccount)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/users/all", nil)

			// Make Request
			r.ServeHTTP(w, req)

			wg.Wait()
			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.True(t, strings.Contains(w.Body.String(), tt.responseBody), "response body does not match, expected "+w.Body.String()+" has a string "+tt.responseBody)
		})
	}
}
