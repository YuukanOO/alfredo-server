package handlers

import (
	"net/http"

	"github.com/YuukanOO/alfredo/env"
	"github.com/gin-gonic/gin"
)

func getAdapters(c *gin.Context) {
	c.JSON(http.StatusOK, env.Current().Adapters)
}
