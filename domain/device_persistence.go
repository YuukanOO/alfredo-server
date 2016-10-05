package domain

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DeviceCollectionName mongo collection name
const DeviceCollectionName = "devices"

// DeviceUpsertFunc callback
type DeviceUpsertFunc func(device *Device) error

// DeviceUpsert inserts or updates a device.
func DeviceUpsert(db *mgo.Database) DeviceUpsertFunc {
	return func(device *Device) error {
		_, err := db.C(DeviceCollectionName).UpsertId(device.ID, device)

		return err
	}
}

// DeviceByIDFunc callback
type DeviceByIDFunc func(id bson.ObjectId) (*Device, error)

// DeviceByID retrieves a device by its ID
func DeviceByID(db *mgo.Database) DeviceByIDFunc {
	return func(id bson.ObjectId) (*Device, error) {
		var result Device

		err := db.C(DeviceCollectionName).FindId(id).One(&result)

		return &result, err
	}
}

// DeviceByNameFunc callback
type DeviceByNameFunc func(name string) (*Device, error)

// DeviceByName retrieves a device by its name
func DeviceByName(db *mgo.Database) DeviceByNameFunc {
	return func(name string) (*Device, error) {
		var result Device

		err := db.C(DeviceCollectionName).Find(bson.M{
			"name": name,
		}).One(&result)

		return &result, err
	}
}
