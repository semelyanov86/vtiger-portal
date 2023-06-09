package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/pkg/logger"
	"net/http"
)

type DataResponse[T any] struct {
	Data  []T `json:"data"`
	Count int `json:"count"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

type AloneDataResponse[T any] struct {
	Data T `json:"data"`
}

type response struct {
	Message string `json:"message"`
}

type validationResponse struct {
	Error   string `json:"error"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	logger.Error(logger.GenerateErrorMessageFromString(message))
	c.AbortWithStatusJSON(statusCode, response{message})
}

func anonymousResponse(c *gin.Context) {
	message := validationResponse{
		Error:   "Anonymous Access",
		Field:   "crmid",
		Message: "Got anonymous user from token or user without crmid",
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, message)
}

func notPermittedResponse(c *gin.Context) {
	message := validationResponse{
		Error:   "Access Not Permitted",
		Field:   "crmid",
		Message: "You are not allowed to view this record",
	}
	c.AbortWithStatusJSON(http.StatusForbidden, message)
}

func moduleNotSupportedResponse(c *gin.Context) {
	message := validationResponse{
		Error:   "module not supported",
		Field:   "module",
		Message: "Module you selected is not supported",
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, message)
}
