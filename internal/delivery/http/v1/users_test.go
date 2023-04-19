package v1

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	mock_repository "github.com/semelyanov86/vtiger-portal/internal/repository/mocks"
	"github.com/semelyanov86/vtiger-portal/internal/service"
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

	const email = "emelyanov86@Km.ru"
	tests := []struct {
		name         string
		body         string
		mockCrm      mockRepositoryCrm
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
			tt.mockCrm(rcrm)

			usersService := service.NewUsersService(rdb, rcrm, &wg)

			services := &service.Services{Users: usersService}
			handler := Handler{services: services}

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

			tokensService := service.NewTokensService(rt, ru)

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
	tests := []struct {
		name         string
		userModel    *domain.User
		statusCode   int
		responseBody string
	}{
		{
			name:         "Get user info",
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

			usersService := service.NewUsersService(rd, rc, &wg)

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
