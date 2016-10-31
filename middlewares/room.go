package middlewares

import (
	"net/http"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

const roomKey = "room"

// GetRoom retrieves the room attached to this context
func GetRoom(c *gin.Context) *domain.Room {
	return c.MustGet(roomKey).(*domain.Room)
}

// Room middleware used to ensure a valid room_id has been given
func Room() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := bson.ObjectIdHex(c.Param("room_id"))

		if !id.Valid() {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			db := GetDB(c)
			var room domain.Room

			if err := db.C(env.RoomsCollection).FindId(id).One(&room); err != nil {
				AbortWithError(c, http.StatusNotFound, err)
			} else {
				c.Set(roomKey, &room)
				c.Next()
			}
		}
	}
}
