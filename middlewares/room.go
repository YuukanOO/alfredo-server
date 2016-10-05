package middlewares

import (
	"net/http"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RoomKey in the context
const RoomKey = "room"

// Room middleware used to ensure a valid room_id has been given
func Room() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := bson.ObjectIdHex(c.Param("room_id"))

		if !id.Valid() {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			db := c.MustGet(DBKey).(*mgo.Database)
			room, err := domain.RoomByID(db)(id)

			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				c.Set(RoomKey, room)
				c.Next()
			}
		}
	}
}
