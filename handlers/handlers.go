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

	auth.GET("/", getAlfredoSystemInfos)

	// Adapters
	auth.GET("/adapters", getAdapters)

	// Rooms
	auth.GET("/rooms", getAllRooms)
	auth.POST("/rooms", createRoom)
	auth.DELETE("/rooms/:room_id", middlewares.Room(), removeRoom)

	// Devices
	auth.GET("/devices", getAllDevices)
	auth.POST("/devices", createDevice)
	auth.PUT("/devices/:device_id/:device_command", middlewares.Device(), deviceExecuteCommand)
}
