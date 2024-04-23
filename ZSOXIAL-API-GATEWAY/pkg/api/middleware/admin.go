package middleware

import (
	"net/http"

	"github.com/akhi9550/api-gateway/pkg/helper"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := helper.GetTokenFromHeader(authHeader)
		if tokenString == "" {
			var err error
			tokenString, err = c.Cookie("Authorization")
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		tokenCliams, err := helper.ExtractAdminFromToken(tokenString)
		if err != nil {
			response := response.ClientResponse(http.StatusUnauthorized, "Invalid Token", nil, err.Error())
			c.JSON(http.StatusUnauthorized, response)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("tokenClaims", tokenCliams)
		c.Next()
	}
}
