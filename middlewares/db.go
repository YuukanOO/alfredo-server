package middlewares

import (
	"github.com/YuukanOO/alfredo/env"
	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
)

const dbKey = "db"

// GetDB retrieves the db attached to this context.
func GetDB(c *gin.Context) *mgo.Database {
	return c.MustGet(dbKey).(*mgo.Database)
}

// OpenDBHandle middleware which open and close a DB session for each request
func OpenDBHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, db := env.Current().Database.GetSession()
		defer session.Close()

		c.Set(dbKey, db)
		c.Next()
	}
}
