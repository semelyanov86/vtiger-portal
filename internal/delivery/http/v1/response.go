package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/pkg/logger"
)

/*type dataResponse struct {
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}*/

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	logger.Error(logger.GenerateErrorMessageFromString(message))
	c.AbortWithStatusJSON(statusCode, response{message})
}
