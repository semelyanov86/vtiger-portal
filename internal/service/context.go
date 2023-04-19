package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
)

type ContextService struct {
}

type contextKey string

const userContextKey = contextKey("user")

func NewContextService() ContextService {
	return ContextService{}
}

func (cs ContextService) ContextSetUser(c *gin.Context, user *domain.User) *gin.Context {
	var ctx = context.WithValue(c.Request.Context(), userContextKey, user)
	c.Request = c.Request.WithContext(ctx)
	return c
}

func (cs ContextService) ContextGetUser(c *gin.Context) *domain.User {
	user, ok := c.Request.Context().Value(userContextKey).(*domain.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}

type MockedContextService struct {
	MockedUser *domain.User
}

func (cs MockedContextService) ContextSetUser(c *gin.Context, user *domain.User) *gin.Context {
	return c
}

func (cs MockedContextService) ContextGetUser(c *gin.Context) *domain.User {
	return cs.MockedUser
}
