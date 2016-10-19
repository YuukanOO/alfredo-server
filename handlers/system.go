package handlers

import (
	"net/http"

	"github.com/YuukanOO/alfredo/env"
	"github.com/gin-gonic/gin"
)

func getAlfredoSystemInfos(c *gin.Context) {

	curEnv := env.Current()
	protocol := "http://"
	remote := curEnv.Server.Remote

	if curEnv.Server.UseHTTPS {
		protocol = "https://"
	}

	if remote != "" {
		remote = protocol + remote
	}

	c.JSON(http.StatusOK, gin.H{
		"version": env.Version,
		"local":   protocol + curEnv.Server.Listen,
		"remote":  remote,
	})
}
