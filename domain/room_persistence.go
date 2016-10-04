package domain

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RoomCollectionName is the mongodb collection name.
const RoomCollectionName = "rooms"

// RoomByIDFunc callback.
type RoomByIDFunc func(id bson.ObjectId) (*Room, error)

// RoomByID finds a room by its id.
func RoomByID(db *mgo.Database) RoomByIDFunc {
	return func(id bson.ObjectId) (*Room, error) {
		var result Room

		err := db.C(RoomCollectionName).FindId(id).One(&result)

		return &result, err
	}
}

// RoomByNameFunc callback.
type RoomByNameFunc func(name string) (*Room, error)

// RoomByName finds a room by its name.
func RoomByName(db *mgo.Database) RoomByNameFunc {
	return func(name string) (*Room, error) {
		var result Room

		err := db.C(RoomCollectionName).Find(bson.M{
			"name": name,
		}).One(&result)

		return &result, err
	}
}

// RoomFunc callback.
type RoomFunc func(room *Room) error

// RoomUpsert adds or updates a room.
func RoomUpsert(db *mgo.Database) RoomFunc {
	return func(room *Room) error {
		_, err := db.C(RoomCollectionName).UpsertId(room.ID, room)

		return err
	}
}

// RoomRemoveFunc callback
type RoomRemoveFunc func(id bson.ObjectId) error

// RoomRemove removes a room from the db.
func RoomRemove(db *mgo.Database) RoomRemoveFunc {
	return func(id bson.ObjectId) error {
		return db.C(RoomCollectionName).RemoveId(id)
	}
}

// RoomsAllFunc callback
type RoomsAllFunc func() ([]Room, error)

// RoomsAll retrieves all rooms.
func RoomsAll(db *mgo.Database) RoomsAllFunc {
	return func() ([]Room, error) {
		var result []Room

		err := db.C(RoomCollectionName).Find(bson.M{}).All(&result)

		if result == nil {
			result = []Room{}
		}

		return result, err
	}
}
