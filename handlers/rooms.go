package handlers

import (
	"net/http"

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

	if c.BindJSON(&params) == nil {
		controller := middlewares.GetController(c)
		db := middlewares.GetDB(c)
		roomsCollection := db.C(env.RoomsCollection)

		room, err := controller.CreateRoom(roomsCollection.Find, params.Name)

		if err != nil {
			c.Error(err)
		} else {
			if err = roomsCollection.Insert(room); err != nil {
				c.Error(err)
			} else {
				c.JSON(http.StatusOK, room)
			}
		}
	}
}

func updateRoom(c *gin.Context) {
	var params createRoomParams

	if c.BindJSON(&params) == nil {
		db := middlewares.GetDB(c)
		room := middlewares.GetRoom(c)
		roomsCollection := db.C(env.RoomsCollection)

		if err := room.Rename(roomsCollection.Find, params.Name); err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			if err := roomsCollection.UpdateId(room.ID, room); err != nil {
				c.Error(err)
			} else {
				c.JSON(http.StatusOK, room)
			}
		}
	}
}

func getAllRooms(c *gin.Context) {
	db := middlewares.GetDB(c)

	var rooms []domain.Room

	err := db.C(env.RoomsCollection).Find(domain.All()).All(&rooms)

	if err != nil {
		c.Error(err)
	} else {
		if rooms == nil {
			rooms = []domain.Room{}
		}

		c.JSON(http.StatusOK, rooms)
	}
}

func removeRoom(c *gin.Context) {
	db := middlewares.GetDB(c)
	room := middlewares.GetRoom(c)

	if err := db.C(env.RoomsCollection).RemoveId(room.ID); err != nil {
		c.Error(err)
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
