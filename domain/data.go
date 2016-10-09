package domain

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// TODO: If multiple keys are found, apply an AND
func applySelectors(selectors ...bson.M) bson.M {
	query := bson.M{}

	for _, s := range selectors {
		for k, v := range s {
			query[k] = v
		}
	}

	return query
}

// FindFunc is the delegate type used to query the database via selectors.
type FindFunc func(selectors ...bson.M) *mgo.Query

// UpdateAllFunc delegate used when performing find and updates for many documents.
type UpdateAllFunc func(selectors []bson.M, setters []bson.M) (*mgo.ChangeInfo, error)

// UpdateAllWithDocFunc delegate used when performing updates with a doc.
type UpdateAllWithDocFunc func(selectors []bson.M, doc interface{}) (*mgo.ChangeInfo, error)

// UpdateFunc delegate used when performing find and update for one document.
type UpdateFunc func(selectors []bson.M, setters []bson.M) error

// UpdateWithDocFunc delegate used when performing an update with a doc.
type UpdateWithDocFunc func(selectors []bson.M, doc interface{}) error

// UpsertFunc delegate used when upserting.
type UpsertFunc func(selectors []bson.M, setters []bson.M) (*mgo.ChangeInfo, error)

// UpsertWithDocFunc delegate used when upserting with a doc.
type UpsertWithDocFunc func(selectors []bson.M, doc interface{}) (*mgo.ChangeInfo, error)

// InsertFunc is the delegate type used to insert documents in the database.
type InsertFunc func(docs ...interface{}) error

// RemoveFunc delegate used when removing one matched document.
type RemoveFunc func(selectors ...bson.M) error

// RemoveAllFunc delegate used when removing all matched documents.
type RemoveAllFunc func(selectors ...bson.M) (*mgo.ChangeInfo, error)

// Find query the given collection with given selectors.
func Find(collection *mgo.Collection) FindFunc {
	return func(selectors ...bson.M) *mgo.Query {
		return collection.Find(applySelectors(selectors...))
	}
}

// Insert documents into the given collection.
func Insert(collection *mgo.Collection) InsertFunc {
	return func(docs ...interface{}) error {
		return collection.Insert(docs...)
	}
}

// UpdateAll all matched documents with given setters.
func UpdateAll(collection *mgo.Collection) UpdateAllFunc {
	return func(selectors []bson.M, setters []bson.M) (*mgo.ChangeInfo, error) {
		return collection.UpdateAll(applySelectors(selectors...), applySelectors(setters...))
	}
}

// UpdateAllWithDoc all matched documents with given document.
func UpdateAllWithDoc(collection *mgo.Collection) UpdateAllWithDocFunc {
	return func(selectors []bson.M, doc interface{}) (*mgo.ChangeInfo, error) {
		return collection.UpdateAll(applySelectors(selectors...), doc)
	}
}

// Update one matched document with given setters.
func Update(collection *mgo.Collection) UpdateFunc {
	return func(selectors []bson.M, setters []bson.M) error {
		return collection.Update(applySelectors(selectors...), applySelectors(setters...))
	}
}

// Upsert matched documents.
func Upsert(collection *mgo.Collection) UpsertFunc {
	return func(selectors []bson.M, setters []bson.M) (*mgo.ChangeInfo, error) {
		return collection.Upsert(applySelectors(selectors...), applySelectors(setters...))
	}
}

// UpsertWithDoc upserts matched documents with the given doc.
func UpsertWithDoc(collection *mgo.Collection) UpsertWithDocFunc {
	return func(selectors []bson.M, doc interface{}) (*mgo.ChangeInfo, error) {
		return collection.Upsert(applySelectors(selectors...), doc)
	}
}

// UpdateWithDoc one matched document with given document.
func UpdateWithDoc(collection *mgo.Collection) UpdateWithDocFunc {
	return func(selectors []bson.M, doc interface{}) error {
		return collection.Update(applySelectors(selectors...), doc)
	}
}

// Remove one matched document from the database.
func Remove(collection *mgo.Collection) RemoveFunc {
	return func(selectors ...bson.M) error {
		return collection.Remove(applySelectors(selectors...))
	}
}

// RemoveAll removes all matched documents.
func RemoveAll(collection *mgo.Collection) RemoveAllFunc {
	return func(selectors ...bson.M) (*mgo.ChangeInfo, error) {
		return collection.RemoveAll(applySelectors(selectors...))
	}
}
