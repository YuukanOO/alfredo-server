package middlewares

import (
	"github.com/YuukanOO/alfredo/env"
	"github.com/gin-gonic/gin"
)

// DBKey in the context
const DBKey = "db"

// DB middleware which open and close a DB session for each request
func DB() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, db := env.Current().Database.GetSession()
		defer session.Close()

		c.Set(DBKey, db)

		c.Next()
	}
}
