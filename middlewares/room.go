package middlewares

import (
	"net/http"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
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

			var room domain.Room

			err := domain.Find(db.C(env.RoomsCollection))(domain.ByID(id)).One(&room)

			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				c.Set(RoomKey, &room)
				c.Next()
			}
		}
	}
}
