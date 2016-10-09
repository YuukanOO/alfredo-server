package domain

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// QueryFunc used to retrieve data when inside the domain.
type QueryFunc func(query interface{}) *mgo.Query

// And join given selectors with a logic and.
func And(selectors ...bson.M) bson.M {
	result := bson.M{}

	for _, s := range selectors {
		for k, v := range s {
			result[k] = v
		}
	}

	return result
}

// ByName finds an element by its name
func ByName(name string) bson.M {
	return bson.M{
		"name": name,
	}
}

// ByUID finds an element by its UID
func ByUID(uid string) bson.M {
	return bson.M{
		"uid": uid,
	}
}

// ByRoomID finds an element by its room ID
func ByRoomID(roomID bson.ObjectId) bson.M {
	return bson.M{
		"room_id": roomID,
	}
}
