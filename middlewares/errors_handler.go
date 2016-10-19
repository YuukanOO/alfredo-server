package middlewares

import "github.com/gin-gonic/gin"

// AbortWithError formats the error and deals with the gin context.
// This is not a middleware in a gin sense and you should call it in handlers.
func AbortWithError(c *gin.Context, status int, err error) {
	c.Error(err)
	c.JSON(status, gin.H{"error": err.Error()})
	c.Abort()
}
