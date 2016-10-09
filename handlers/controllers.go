package handlers

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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

		controller, err := domain.RegisterController(domain.Find(controllersCollection), []byte(env.Current().Security.Secret), params.UID)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		} else {
			if _, err = domain.UpsertWithDoc(controllersCollection)([]bson.M{domain.ByID(controller.ID)}, controller); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusOK, controller.Token)
			}
		}
	}
}
