package domain

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ControllerCollectionName is the mongo collection name
const ControllerCollectionName = "controllers"

// ControllerByUIDFunc callback used by ControllerByUID
type ControllerByUIDFunc func(uid string) (*Controller, error)

// ControllerByUID retrieves a controller by its UID
func ControllerByUID(db *mgo.Database) ControllerByUIDFunc {
	return func(uid string) (*Controller, error) {
		var result Controller

		err := db.C(ControllerCollectionName).Find(bson.M{
			"uid": uid,
		}).One(&result)

		return &result, err
	}
}

// ControllerByIDFunc callback used by ControllerByID
type ControllerByIDFunc func(id bson.ObjectId) (*Controller, error)

// ControllerByID retrieves a controller by its ID
func ControllerByID(db *mgo.Database) ControllerByIDFunc {
	return func(id bson.ObjectId) (*Controller, error) {
		var result Controller

		err := db.C(ControllerCollectionName).FindId(id).One(&result)

		return &result, err
	}
}

// ControllerFunc callback
type ControllerFunc func(controller *Controller) error

// ControllerUpsert creates or updates a controller
func ControllerUpsert(db *mgo.Database) ControllerFunc {
	return func(controller *Controller) error {
		_, err := db.C(ControllerCollectionName).UpsertId(controller.ID, controller)

		return err
	}
}
