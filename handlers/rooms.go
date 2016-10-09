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
	Name string
}

func createRoom(c *gin.Context) {
	var params createRoomParams

	if err := c.BindJSON(&params); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		controller := c.MustGet(middlewares.ControllerKey).(domain.Controller)
		db := c.MustGet(middlewares.DBKey).(*mgo.Database)
		roomsCollection := db.C(env.RoomsCollection)

		room, err := controller.CreateRoom(domain.Find(roomsCollection), params.Name)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		} else {
			if err = domain.Insert(roomsCollection)(room); err != nil {
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
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		db := c.MustGet(middlewares.DBKey).(*mgo.Database)
		room := c.MustGet(middlewares.RoomKey).(*domain.Room)
		roomsCollection := db.C(env.RoomsCollection)

		err := room.Rename(domain.Find(roomsCollection), params.Name)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		} else {
			if err := domain.UpdateWithDoc(roomsCollection)([]bson.M{domain.ByID(room.ID)}, room); err != nil {
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

	err := domain.Find(db.C(env.RoomsCollection))().All(&rooms)

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

	if err := domain.Remove(db.C(env.RoomsCollection))(domain.ByID(room.ID)); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
