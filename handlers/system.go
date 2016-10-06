package handlers

import (
	"net/http"

	"github.com/YuukanOO/alfredo/env"
	"github.com/gin-gonic/gin"
)

func getAlfredoSystemInfos(c *gin.Context) {

	curEnv := env.Current()

	c.JSON(http.StatusOK, gin.H{
		"version": env.Version,
		"local":   curEnv.Server.Listen,
		"remote":  curEnv.Server.Remote,
	})
}
