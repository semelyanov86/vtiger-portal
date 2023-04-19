package http

import (
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/validator"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h Handler) corsMiddleware(c *gin.Context) {
	c.Header("Vary", "Origin")
	origin := c.Request.Header.Get("Origin")
	c.Header("Content-Type", "application/json")

	if origin != "" {
		for i := range h.config.Cors.TrustedOrigins {
			if origin == h.config.Cors.TrustedOrigins[i] {
				c.Header("Access-Control-Allow-Origin", origin)
				if c.Request.Method == http.MethodOptions && c.Request.Header.Get("Access-Control-Request-Method") != "" {
					c.Header("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
					c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
					c.AbortWithStatus(http.StatusOK)
					return
				}
				break
			}
		}
	}

	c.Next()
}

func (h Handler) authenticate(c *gin.Context) {
	c.Header("Vary", "Authorization")
	var authorizationHeader = c.Request.Header.Get("Authorization")
	if authorizationHeader == "" {
		h.services.Context.ContextSetUser(c, domain.AnonymousUser)
		c.Next()
		return
	}

	var headerParts = strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Token Error", "field": "Bearer", "message": "Wrong Authorization Header"})
		return
	}

	var token = headerParts[1]
	v := validator.New()

	if domain.ValidateTokenPlaintext(v, token); !v.Valid() {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Token Error", "field": "Bearer", "message": "Wrong Token in authorization header"})
		return
	}

	user, err := h.services.Users.GetUserByToken(c.Request.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Token Error", "field": "Bearer", "message": "Passed token is not attached to user"})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Token Error", "field": "Bearer", "message": "There was an error during checking a token"})
		}
		return
	}
	h.services.Context.ContextSetUser(c, user)
	c.Next()
}
