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
	Name    string `binding:"required"`
	Adapter string `binding:"required"`
	Config  map[string]interface{}
	RoomID  bson.ObjectId `json:"room_id" binding:"required"`
}

func createDevice(c *gin.Context) {
	var params createDeviceParams

	if err := c.BindJSON(&params); err != nil {
		middlewares.AbortWithError(c, http.StatusBadRequest, err)
	} else {
		db := c.MustGet(middlewares.DBKey).(*mgo.Database)
		var room domain.Room

		err := db.C(env.RoomsCollection).FindId(params.RoomID).One(&room)

		adapter := env.Current().GetAdapter(params.Adapter)

		if err != nil || adapter == nil {
			c.AbortWithError(http.StatusNotFound, err)
		} else {
			devicesCollection := db.C(env.DevicesCollection)
			device, err := room.RegisterDevice(devicesCollection.Find, params.Name, adapter, params.Config)

			if err != nil {
				middlewares.AbortWithError(c, http.StatusBadRequest, err)
			} else {
				if err = devicesCollection.Insert(device); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
				} else {
					c.JSON(http.StatusOK, device)
				}
			}
		}
	}
}

func getAllDevices(c *gin.Context) {
	db := c.MustGet(middlewares.DBKey).(*mgo.Database)

	var devices []domain.Device

	err := db.C(env.DevicesCollection).Find(bson.M{}).All(&devices)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		if devices == nil {
			devices = []domain.Device{}
		}

		c.JSON(http.StatusOK, devices)
	}
}

func getRoomDevices(c *gin.Context) {
	db := c.MustGet(middlewares.DBKey).(*mgo.Database)
	room := c.MustGet(middlewares.RoomKey).(*domain.Room)

	var devices []domain.Device

	err := db.C(env.DevicesCollection).Find(domain.ByRoomID(room.ID)).All(&devices)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		if devices == nil {
			devices = []domain.Device{}
		}

		c.JSON(http.StatusOK, devices)
	}
}

func removeDevice(c *gin.Context) {
	db := c.MustGet(middlewares.DBKey).(*mgo.Database)
	device := c.MustGet(middlewares.DeviceKey).(*domain.Device)

	if err := db.C(env.DevicesCollection).RemoveId(device.ID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func deviceExecuteCommand(c *gin.Context) {
	var commandParameters map[string]interface{}

	if err := c.BindJSON(&commandParameters); err != nil {
		middlewares.AbortWithError(c, http.StatusBadRequest, err)
	} else {
		device := c.MustGet(middlewares.DeviceKey).(*domain.Device)
		adapter := env.Current().GetAdapter(device.Adapter)
		command := c.Param("device_command")

		if adapter == nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			res, err := adapter.Execute(env.Current().Server.ShellCommand, command, device, commandParameters)

			if err == nil {
				c.JSON(http.StatusOK, res)
			} else {
				if res == nil {
					middlewares.AbortWithError(c, http.StatusBadRequest, err)
				} else {
					c.JSON(http.StatusBadRequest, res)
				}
			}
		}
	}
}
