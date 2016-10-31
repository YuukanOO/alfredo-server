package middlewares

import (
	"net/http"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

const deviceKey = "device"

// GetDevice retrieves the device attached to this context
func GetDevice(c *gin.Context) *domain.Device {
	return c.MustGet(deviceKey).(*domain.Device)
}

// Device ensure a valid device has been given in the request.
func Device() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := bson.ObjectIdHex(c.Param("device_id"))

		if !id.Valid() {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			db := GetDB(c)
			var device domain.Device

			if err := db.C(env.DevicesCollection).FindId(id).One(&device); err != nil {
				AbortWithError(c, http.StatusBadRequest, err)
			} else {
				c.Set(deviceKey, &device)
				c.Next()
			}
		}
	}
}
