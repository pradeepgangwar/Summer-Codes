package repo

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/pradeepgangwar/go-websocket/boot"
	"github.com/pradeepgangwar/go-websocket/user"
)

// MongoRepo contains the connection and collection information
type MongoRepo struct {
	Connection *mgo.Session
	Collection *mgo.Collection
}

// NewRepository : This makes the ew object of the MongoRepo struct
func NewRepository(c *mgo.Session) (m *MongoRepo) {
	config := boot.GetConfig()
	coll := c.DB(config.MongoDatabaseName).C(config.MongoCollectionName)
	return &MongoRepo{c, coll}
}

// New : Writes new object to the mongo database
func (r *MongoRepo) New(ctx context.Context, model *user.Model) (string, error) {
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

	//

	// user1 := &user.Model{}

	// for results.Next(user1) {
	// 	users = append(users, user1)
	// 	user1 = &user.Model{}
	// }

	return users
}
