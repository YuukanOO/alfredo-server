package handlers

import (
	"net/http"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
)

func getAdapters(c *gin.Context) {
	db := middlewares.GetDB(c)

	var adapters []domain.Adapter

	err := db.C(env.AdaptersCollection).Find(domain.All()).All(&adapters)

	if err != nil {
		c.Error(err)
	} else {
		if adapters == nil {
			adapters = []domain.Adapter{}
		}

		c.JSON(http.StatusOK, adapters)
	}
}
