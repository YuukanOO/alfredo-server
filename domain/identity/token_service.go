package identity

import "gopkg.in/mgo.v2/bson"

// TokenService manage access token stuff.
type TokenService interface {
	// GenerateToken generates a token for the given user id.
	GenerateToken(userID bson.ObjectId) (Token, error)
}
