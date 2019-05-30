package controller

import (
	"context"
	"time"

	bson "github.com/globalsign/mgo/bson"
	"github.com/pradeepgangwar/go-websocket/boot"
	"github.com/pradeepgangwar/go-websocket/repo"
	"github.com/pradeepgangwar/go-websocket/user"
)

// Controller contains the repo and timeout information
type Controller struct {
	Repo    *repo.MongoRepo
	Timeout time.Duration
}

// NewController makes new struct
func NewController(config *boot.Config, repo *repo.MongoRepo) *Controller {
	timeout := time.Duration(config.ContextTimeout) * time.Second
	return &Controller{
		repo, timeout,
	}
}

// New : This sends the new object to be written to the DB
func (c *Controller) New(ctx context.Context, m *user.Model) (*user.Model, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	id, err := c.Repo.New(ctx, m)
	if err != nil {
		return nil, err
	}
	m.Id = bson.ObjectIdHex(id)
	return m, nil
}

// FindAll : Returns the user set from DB
func (c *Controller) FindAll(ctx context.Context) ([]*user.Model, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	users := c.Repo.FindAll(ctx)

	return users, nil
}
