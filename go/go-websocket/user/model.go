package user

import (
	"time"

	"github.com/atlanhq/go-bongo"
	"github.com/globalsign/mgo/bson"
)

// Model Represents user model
type Model struct {
	bongo.DocumentBase `bson:"inline"`
	ID                 bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name               string        `json:"name" bson:"name"`
	Email              string        `json:"email" bson:"email"`
	About              string        `json:"about" bson:"about"`
}

// NewModel : returns new Model
func NewModel(name string, email string, about string) *Model {
	id := bson.NewObjectId()

	model := Model{
		Name:  name,
		Email: email,
		About: about,
	}

	now := time.Now()
	model.SetId(id)
	model.SetCreated(now)
	model.SetModified(now)
	return &model
}
