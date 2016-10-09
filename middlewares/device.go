package middlewares

import (
	"net/http"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DeviceKey in the context
const DeviceKey = "device"

// Device ensure a valid device has been given in the request.
func Device() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := bson.ObjectIdHex(c.Param("device_id"))

		if !id.Valid() {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			db := c.MustGet(DBKey).(*mgo.Database)

			var device domain.Device

			err := domain.Find(db.C(env.DevicesCollection))(domain.ByID(id)).One(&device)

			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				c.Set(DeviceKey, &device)
				c.Next()
			}
		}
	}
}
