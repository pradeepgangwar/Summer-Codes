package repo

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pradeepgangwar/go-websocket/boot"
	"github.com/pradeepgangwar/go-websocket/user"
)

// MongoRepo contains the connection and collection information
type MongoRepo struct {
	Connection   *mgo.Session
	Collection   *mgo.Collection
	ChangeStream *mgo.ChangeStream
}

// NewRepository : This makes the ew object of the MongoRepo struct
func NewRepository(c *mgo.Session) (m *MongoRepo) {
	pipeline := []bson.M{}
	config := boot.GetConfig()
	coll := c.DB(config.MongoDatabaseName).C(config.MongoCollectionName)
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{})
	if err != nil {
		fmt.Println(err)
	}
	return &MongoRepo{c, coll, changeStream}
}

// New : Writes new object to the mongo database
func (r *MongoRepo) New(ctx context.Context, model user.Model) (string, error) {
	err := r.Collection.Insert(model)
	if err != nil {
		return "", err
	}
	return model.GetId().Hex(), nil
}

// FindAll : Writes new object to the mongo database
func (r *MongoRepo) FindAll(ctx context.Context) []*user.Model {

	var users []*user.Model
	r.Collection.Find(nil).All(&users)

	return users
}

// TrackChange : Implemented through the mongo changestreams, tracks the changes in the collection
// func (r *MongoRepo) TrackChange(ctx context.Context) ([]*user.Model, error) {

// 	var users []*user.Model

// 	pipeline := []bson.M{}
// 	changeStream, err := r.Collection.Watch(pipeline, mgo.ChangeStreamOptions{})
// 	var user user.Model

// 	if err != nil {
// 		return nil, err
// 	}

// 	for changeStream.Next(&user) {
// 		fmt.Printf("Change: %v\n", user)
// 		users = append(users, &user)
// 	}

// 	return users, nil
// }
