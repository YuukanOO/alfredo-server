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

type createRoomParams struct {
	Name string `binding:"required"`
}

func createRoom(c *gin.Context) {
	var params createRoomParams

	if err := c.BindJSON(&params); err != nil {
		middlewares.AbortWithError(c, http.StatusBadRequest, err)
	} else {
		controller := c.MustGet(middlewares.ControllerKey).(domain.Controller)
		db := c.MustGet(middlewares.DBKey).(*mgo.Database)
		roomsCollection := db.C(env.RoomsCollection)

		room, err := controller.CreateRoom(roomsCollection.Find, params.Name)

		if err != nil {
			middlewares.AbortWithError(c, http.StatusBadRequest, err)
		} else {
			if err = roomsCollection.Insert(room); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusOK, room)
			}
		}
	}
}

func updateRoom(c *gin.Context) {
	var params createRoomParams

	if err := c.BindJSON(&params); err != nil {
		middlewares.AbortWithError(c, http.StatusBadRequest, err)
	} else {
		db := c.MustGet(middlewares.DBKey).(*mgo.Database)
		room := c.MustGet(middlewares.RoomKey).(*domain.Room)
		roomsCollection := db.C(env.RoomsCollection)

		err := room.Rename(roomsCollection.Find, params.Name)

		if err != nil {
			middlewares.AbortWithError(c, http.StatusBadRequest, err)
		} else {
			if err := roomsCollection.UpdateId(room.ID, room); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusOK, room)
			}
		}
	}
}

func getAllRooms(c *gin.Context) {
	db := c.MustGet(middlewares.DBKey).(*mgo.Database)

	var rooms []domain.Room

	err := db.C(env.RoomsCollection).Find(bson.M{}).All(&rooms)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		if rooms == nil {
			rooms = []domain.Room{}
		}

		c.JSON(http.StatusOK, rooms)
	}
}

func removeRoom(c *gin.Context) {
	db := c.MustGet(middlewares.DBKey).(*mgo.Database)
	room := c.MustGet(middlewares.RoomKey).(*domain.Room)

	if err := db.C(env.RoomsCollection).RemoveId(room.ID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
