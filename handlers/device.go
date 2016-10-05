package handlers

import (
	"net/http"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type createDeviceParams struct {
	Name    string
	Adapter string
	Config  map[string]interface{}
	RoomID  bson.ObjectId `json:"room_id"`
}

func createDevice(c *gin.Context) {
	var params createDeviceParams

	if err := c.BindJSON(&params); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		db := c.MustGet(middlewares.DBKey).(*mgo.Database)
		room, err := domain.RoomByID(db)(params.RoomID)
		adapter := env.Current().Adapters[params.Adapter]

		if err != nil || adapter == nil {
			c.AbortWithError(http.StatusNotFound, err)
		} else {
			device, err := room.RegisterDevice(domain.DeviceByName(db), params.Name, adapter, params.Config)

			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
			} else {
				if err = domain.DeviceUpsert(db)(device); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
				} else {
					c.JSON(http.StatusOK, device)
				}
			}
		}
	}
}

func deviceExecuteCommand(c *gin.Context) {
	var commandParameters map[string]interface{}

	if err := c.BindJSON(&commandParameters); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		device := c.MustGet(middlewares.DeviceKey).(*domain.Device)
		adapter := env.Current().Adapters[device.Adapter]
		command := c.Param("device_command")

		if adapter == nil {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			err = device.Execute(adapter, command, commandParameters)

			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
			} else {
				c.AbortWithStatus(http.StatusOK)
			}
		}
	}
}
