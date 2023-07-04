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
	if strings.HasSuffix(c.Request.RequestURI, ".js") {
		c.Header("Content-Type", "application/javascript; charset=utf-8")
	} else if strings.HasSuffix(c.Request.RequestURI, ".css") {
		c.Header("Content-Type", "text/css; charset=utf-8")
	} else if strings.HasSuffix(c.Request.RequestURI, ".png") {
		c.Header("Content-Type", "image/png")
	} else if strings.HasSuffix(c.Request.RequestURI, "swagger_spec") {
		c.Header("Content-Type", "application/x-yaml")
	} else if strings.HasPrefix(c.Request.RequestURI, "/swagger/") {
		c.Header("Content-Type", "text/html; charset=utf-8")
	} else {
		c.Header("Content-Type", "application/json")
	}

	if origin != "" {
		for i := range h.config.Cors.TrustedOrigins {
			if origin == h.config.Cors.TrustedOrigins[i] {
				c.Header("Access-Control-Allow-Origin", origin)
				if c.Request.Method == http.MethodOptions && c.Request.Header.Get("Access-Control-Request-Method") != "" {
					c.Header("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
					c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With")
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
	if c.FullPath() == "/api/v1/users/" || c.FullPath() == "/api/v1/users/restore" || c.FullPath() == "/api/v1/users/password" || c.FullPath() == "/api/v1/users/login" || c.FullPath() == "/api/v1/payments/webhook" {
		c.Next()
		return
	}
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
	if user != nil && !user.IsActive {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Not Active", "field": "is_active", "message": "Authorized user is deactivated!"})
		return
	}

	if !user.Otp_verified && user.Otp_enabled && c.FullPath() != "/api/v1/otp/validate" && c.FullPath() != "/api/v1/otp/verify" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "OTP Required", "field": "otp_verified", "message": "You need to pass otp verification to use this service"})
		return
	}
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
