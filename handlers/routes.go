package handlers

import (
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
)

// Routes register handlers for the entire application.
func Routes(eng *gin.Engine) {
	eng.POST("/controllers", registerController)

	// Require privileges
	auth := eng.Group("", middlewares.RequireAuth())
	requireRoom := auth.Group("", middlewares.RequireRoom())
	requireDevice := auth.Group("", middlewares.RequireDevice())

	auth.GET("/", getAlfredoSystemInfos)

	// Adapters
	auth.GET("/adapters", getAdapters)

	// Rooms
	auth.GET("/rooms", getAllRooms)
	auth.POST("/rooms", createRoom)
	requireRoom.PUT("/rooms/:room_id", updateRoom)
	requireRoom.DELETE("/rooms/:room_id", removeRoom)
	requireRoom.GET("/rooms/:room_id/devices", getRoomDevices)

	// Devices
	auth.GET("/devices", getAllDevices)
	auth.POST("/devices", createDevice)
	requireDevice.PUT("/devices/:device_id", updateDevice)
	requireDevice.PUT("/devices/:device_id/:device_command", deviceExecuteCommand)
	requireDevice.DELETE("/devices/:device_id", removeDevice)
}

func waterfall(errors ...error) error {
	for _, v := range errors {
		if v != nil {
			return v
		}
	}

	return nil
}
