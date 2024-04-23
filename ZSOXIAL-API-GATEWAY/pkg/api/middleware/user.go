package middleware

import (
	"fmt"
	"net/http"

	"github.com/akhi9550/api-gateway/pkg/helper"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := helper.GetTokenFromHeader(authHeader)
		fmt.Println("token:-", tokenString)
		if tokenString == "" {
			var err error
			tokenString, err = c.Cookie("Authorization")
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		userID, userName, err := helper.ExtractUserIDFromToken(tokenString)
		if err != nil {
			response := response.ClientResponse(http.StatusUnauthorized, "Invalid Token", nil, err.Error())
			c.JSON(http.StatusUnauthorized, response)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user_id", userID)
		c.Set("user_name", userName)
		c.Next()
	}
}
