package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

import "github.com/YuukanOO/alfredo/domain"

// FormatErrors handles errors responses.
func FormatErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if len(c.Errors) > 0 {
				var err interface{} = c.Errors[0].Err

				switch err.(type) {
				case *domain.Error:
					c.JSON(http.StatusBadRequest, err)
				default:
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()

		c.Next()
	}
}
