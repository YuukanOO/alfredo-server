package handlers

import (
	"net/http"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
)

type registerControllerParams struct {
	UID string
}

func registerController(c *gin.Context) {
	if !env.Current().Server.AllowRegistering {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		var params registerControllerParams

		if c.BindJSON(&params) == nil {
			db := middlewares.GetDB(c)
			controllersCollection := db.C(env.ControllersCollection)

			controller, err := domain.RegisterController(controllersCollection.Find, []byte(env.Current().Security.Secret), params.UID)

			if err != nil {
				middlewares.AbortWithError(c, http.StatusBadRequest, err)
			} else {
				if _, err = controllersCollection.UpsertId(controller.ID, controller); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
				} else {
					c.JSON(http.StatusOK, controller.Token)
				}
			}
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}
}
