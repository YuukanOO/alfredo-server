package handlers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
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

	if c.BindJSON(&params) == nil {
		db := middlewares.GetDB(c)
		var room domain.Room
		var adapter domain.Adapter

		if err := waterfall(
			db.C(env.RoomsCollection).FindId(params.RoomID).One(&room),
			db.C(env.AdaptersCollection).FindId(params.Adapter).One(&adapter),
		); err != nil {
			c.Error(err)
		} else {
			devicesCollection := db.C(env.DevicesCollection)
			device, err := room.RegisterDevice(devicesCollection.Find, params.Name, &adapter, params.Config)

			if err != nil {
				c.JSON(http.StatusBadRequest, err)
			} else {
				// Upon device creation, try to execute a command named status
				if outRes, _ := adapter.Execute(env.Current().Server.ShellCommand, "status", *device, map[string]interface{}{}); outRes != nil {
					device.UpdateStatus(outRes)
				}

				if err = devicesCollection.Insert(device); err != nil {
					c.Error(err)
				} else {
					c.JSON(http.StatusOK, device)
				}
			}
		}
	}
}

type updateDeviceParams struct {
	Name   string
	Config map[string]interface{}
}

func updateDevice(c *gin.Context) {
	var params updateDeviceParams

	if c.BindJSON(&params) == nil {
		db := middlewares.GetDB(c)
		device := middlewares.GetDevice(c)
		deviceCollection := db.C(env.DevicesCollection)
		var adapter domain.Adapter

		if err := db.C(env.AdaptersCollection).FindId(device.Adapter).One(&adapter); err != nil {
			c.Error(err)
		} else {

			if err := waterfall(
				device.Rename(deviceCollection.Find, params.Name),
				device.UpdateConfig(&adapter, params.Config),
			); err != nil {
				c.JSON(http.StatusBadRequest, err)
			} else {
				if err := deviceCollection.UpdateId(device.ID, device); err != nil {
					c.Error(err)
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

	if err := db.C(env.DevicesCollection).Find(domain.All()).All(&devices); err != nil {
		c.Error(err)
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

	if err := db.C(env.DevicesCollection).Find(domain.ByRoomID(room.ID)).All(&devices); err != nil {
		c.Error(err)
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
		c.Error(err)
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func deviceExecuteCommand(c *gin.Context) {
	var commandParameters map[string]interface{}

	if c.BindJSON(&commandParameters) == nil {
		db := middlewares.GetDB(c)
		device := middlewares.GetDevice(c)
		var adapter domain.Adapter
		command := c.Param("device_command")

		if err := db.C(env.AdaptersCollection).FindId(device.Adapter).One(&adapter); err != nil {
			c.Error(err)
		} else {
			// Executes the command
			res, err := adapter.Execute(env.Current().Server.ShellCommand, command, *device, commandParameters)

			if err == nil {
				logrus.WithField("out", res.Out).Debug("Command success")

				// If everything is good, update the device status given the execution result
				// and returns the new device status.
				device.UpdateStatus(res)

				if err := db.C(env.DevicesCollection).UpdateId(device.ID, device); err != nil {
					c.Error(err)
				} else {
					c.JSON(http.StatusOK, device.Status)
				}
			} else {
				// If something goes wrong, it will print the execution result to ease the debugging
				if res == nil {
					logrus.WithError(err).Error("Command execution failed")

					c.JSON(http.StatusBadRequest, err)
				} else {
					logrus.WithField("error", res.Err).Warn("Command execution failed")

					c.JSON(http.StatusBadRequest, res)
				}
			}
		}
	}
}
