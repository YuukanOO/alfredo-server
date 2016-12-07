package middlewares

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

import "github.com/YuukanOO/alfredo/domain"

// FormatErrors handles errors responses.
func FormatErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if len(c.Errors) > 0 {
				err := c.Errors[0].Err

				switch err.(type) {
				case *domain.Error:
					c.JSON(http.StatusBadRequest, err)
				default:
					logrus.WithError(err).Error("Internal error")
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()

		c.Next()
	}
}
