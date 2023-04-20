package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/pkg/logger"
	"net/http"
)

/*type dataResponse struct {
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}*/

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
