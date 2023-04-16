package http

import (
	"net/http"

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
