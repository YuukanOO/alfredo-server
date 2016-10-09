package handlers

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
)

type registerControllerParams struct {
	UID string
}

func registerController(c *gin.Context) {
	var params registerControllerParams

	if err := c.BindJSON(&params); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		db := c.MustGet(middlewares.DBKey).(*mgo.Database)
		controllersCollection := db.C(env.ControllersCollection)

		controller, err := domain.RegisterController(controllersCollection.Find, []byte(env.Current().Security.Secret), params.UID)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		} else {
			if _, err = controllersCollection.UpsertId(controller.ID, controller); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusOK, controller.Token)
			}
		}
	}
}
