package controller

import (
	"context"
	"time"

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

// NewUser : Create a user
func (c *Controller) NewUser(ctx context.Context, user user.Model) error {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	_, err := c.Repo.New(ctx, user)

	return err
}

// AllUser : Returns all users
func (c *Controller) AllUser(ctx context.Context) []*user.Model {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	users := c.Repo.FindAll(ctx)

	return users
}

// GetChange : Chacks for change in the collection
// func (c *Controller) GetChange(ctx context.Context) ([]*user.Model, error) {
// 	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
// 	defer cancel()

// 	users, err := c.Repo.TrackChange(ctx)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return users, nil
// }
