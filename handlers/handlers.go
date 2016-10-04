package handlers

import (
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
)

// Register register handlers for the entire application.
func Register(eng *gin.Engine) {
	eng.POST("/controller", registerController)

	// Require privileges
	auth := eng.Group("", middlewares.Auth())

	auth.GET("/", getAlfredoSystemInfos)

	// Adapters
	auth.GET("/adapters", getAdapters)

	// Rooms
	auth.GET("/rooms", getAllRooms)
	auth.POST("/rooms", createRoom)
	auth.DELETE("/rooms/:id", removeRoom)
}
