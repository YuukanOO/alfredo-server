package handlers

import (
	"net/http"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
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
		db := middlewares.GetDB(c)
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
				// Upon device creation, try to execute a command named status
				if outRes, _ := adapter.Execute(env.Current().Server.ShellCommand, "status", device, map[string]interface{}{}); outRes != nil {
					device.UpdateStatus(outRes)
				}

				if err = devicesCollection.Insert(device); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
				} else {
					c.JSON(http.StatusOK, device)
				}
			}
		}
	}
}

type updateDeviceParams struct {
	Name   string `binding:"required"`
	Config map[string]interface{}
}

func updateDevice(c *gin.Context) {
	var params updateDeviceParams

	if err := c.BindJSON(&params); err != nil {
		middlewares.AbortWithError(c, http.StatusBadRequest, err)
	} else {
		db := middlewares.GetDB(c)
		device := middlewares.GetDevice(c)
		adapter := env.Current().GetAdapter(device.Adapter)
		deviceCollection := db.C(env.DevicesCollection)

		if err := device.Rename(deviceCollection.Find, params.Name); err != nil {
			middlewares.AbortWithError(c, http.StatusBadRequest, err)
		} else {
			if err := device.UpdateConfig(adapter, params.Config); err != nil {
				middlewares.AbortWithError(c, http.StatusBadRequest, err)
			} else {
				if err := deviceCollection.UpdateId(device.ID, device); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
				} else {
					c.JSON(http.StatusOK, device)
				}
			}
		}
	}
}

func getAllDevices(c *gin.Context) {
	db := middlewares.GetDB(c)

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
	db := middlewares.GetDB(c)
	room := middlewares.GetRoom(c)

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
	db := middlewares.GetDB(c)
	device := middlewares.GetDevice(c)

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
		db := middlewares.GetDB(c)
		device := middlewares.GetDevice(c)
		adapter := env.Current().GetAdapter(device.Adapter)
		command := c.Param("device_command")

		if adapter == nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			// Executes the command
			res, err := adapter.Execute(env.Current().Server.ShellCommand, command, device, commandParameters)

			if err == nil {
				// If everything is good, update the device status given the execution result
				// and returns the new device status.
				device.UpdateStatus(res)

				if err := db.C(env.DevicesCollection).UpdateId(device.ID, device); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
				} else {
					c.JSON(http.StatusOK, device.Status)
				}
			} else {
				// If something goes wrong, it will print the execution result to ease the debugging
				if res == nil {
					middlewares.AbortWithError(c, http.StatusBadRequest, err)
				} else {
					c.JSON(http.StatusBadRequest, res)
				}
			}
		}
	}
}
