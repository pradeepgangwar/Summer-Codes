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

// Socket : Create a new socke
func (c *Controller) Socket(ctx context.Context, user user.Model) error {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	_, err := c.Repo.New(ctx, user)

	return err
}
