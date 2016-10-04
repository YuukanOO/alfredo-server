package handlers

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/env"
	"github.com/gin-gonic/gin"
)

type createRoomParams struct {
	Name string
}

func createRoom(c *gin.Context) {
	var params createRoomParams

	if err := c.BindJSON(&params); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		controller := c.MustGet("controller").(domain.Controller)
		session, db := env.Current().Database.GetSession()
		defer session.Close()

		room, err := controller.CreateRoom(domain.RoomByName(db), params.Name)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		} else {
			if err = domain.RoomUpsert(db)(room); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusOK, room)
			}
		}
	}
}

func getAllRooms(c *gin.Context) {
	session, db := env.Current().Database.GetSession()
	defer session.Close()

	rooms, err := domain.RoomsAll(db)()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, rooms)
	}
}

func removeRoom(c *gin.Context) {
	id := bson.ObjectIdHex(c.Param("id"))

	session, db := env.Current().Database.GetSession()
	defer session.Close()

	if err := domain.RoomRemove(db)(id); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
