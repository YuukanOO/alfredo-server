package handlers

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/YuukanOO/alfredo/domain"
	"github.com/YuukanOO/alfredo/middlewares"
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
		controller := c.MustGet(middlewares.ControllerKey).(domain.Controller)
		db := c.MustGet(middlewares.DBKey).(*mgo.Database)

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
	db := c.MustGet(middlewares.DBKey).(*mgo.Database)

	rooms, err := domain.RoomsAll(db)()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, rooms)
	}
}

func removeRoom(c *gin.Context) {
	db := c.MustGet(middlewares.DBKey).(*mgo.Database)
	room := c.MustGet(middlewares.RoomKey).(*domain.Room)

	if err := domain.RoomRemove(db)(room.ID); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
