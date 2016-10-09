package domain

import "gopkg.in/mgo.v2/bson"

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

// ByID finds an element by its ID
func ByID(id bson.ObjectId) bson.M {
	return bson.M{
		"_id": id,
	}
}

// ByRoomID finds an element by its room ID
func ByRoomID(roomID bson.ObjectId) bson.M {
	return bson.M{
		"room_id": roomID,
	}
}
