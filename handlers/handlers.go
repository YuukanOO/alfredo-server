package handlers

import "github.com/gin-gonic/gin"

// Register register handlers for the entire application.
func Register(eng *gin.Engine) {
	eng.GET("/", getAlfredoSystemInfos)
	eng.GET("/adapters", getAdapters)
}
