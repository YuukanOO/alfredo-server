package handlers

import (
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
)

// Register register handlers for the entire application.
func Register(eng *gin.Engine) {
	eng.POST("/controllers", registerController)

	// Require privileges
	auth := eng.Group("", middlewares.Auth())
	requireRoom := auth.Group("", middlewares.Room())
	requireDevice := auth.Group("", middlewares.Device())

	auth.GET("/", getAlfredoSystemInfos)

	// Adapters
	auth.GET("/adapters", getAdapters)

	// Rooms
	auth.GET("/rooms", getAllRooms)
	auth.POST("/rooms", createRoom)
	requireRoom.PUT("/rooms/:room_id", updateRoom)
	requireRoom.DELETE("/rooms/:room_id", removeRoom)

	// Devices
	auth.GET("/devices", getAllDevices)
	auth.POST("/devices", createDevice)
	requireDevice.PUT("/devices/:device_id/:device_command", deviceExecuteCommand)
	requireDevice.DELETE("/devices/:device_id", removeDevice)
}
