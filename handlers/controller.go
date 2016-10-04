package handlers

import (
	"net/http"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
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
		session, db := env.Current().Database.GetSession()
		defer session.Close()

		controller, err := domain.RegisterController(domain.ControllerByUID(db), []byte(env.Current().Security.Secret), params.UID)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		} else {
			if err = domain.ControllerUpsert(db)(controller); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusOK, controller.Token)
			}
		}
	}
}
