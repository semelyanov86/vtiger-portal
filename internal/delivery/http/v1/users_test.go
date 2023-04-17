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
	"net/http/httptest"
	"strings"
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
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			rdb := repository.NewUsersMock()

			rcrm := mock_repository.NewMockUsersCrm(c)
			tt.mockCrm(rcrm)

			usersService := service.NewUsersService(rdb, rcrm)

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
