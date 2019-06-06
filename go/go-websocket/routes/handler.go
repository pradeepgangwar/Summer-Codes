package route

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/pradeepgangwar/go-websocket/controller"
	"github.com/pradeepgangwar/go-websocket/user"
)

type handler struct {
	Controller *controller.Controller
}

var (
	upgrader = websocket.Upgrader{}
)

// NewHandler : This returns new handler object
func NewHandler(g *echo.Group, controller *controller.Controller) {
	handler := &handler{
		Controller: controller,
	}
	g.Static("/", "public/index.html")
	g.GET("/add", handler.UsersGet)
}

func (h *handler) UsersGet(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	allusers := h.Controller.AllUser(ctx)

	// Write
	err = ws.WriteJSON(allusers)
	if err != nil {
		c.Logger().Error(err)
	}

	for {
		var newUser user.Model

		// Read
		err = ws.ReadJSON(&newUser)
		if err != nil {
			c.Logger().Error(err)
		}

		err = h.Controller.NewUser(ctx, newUser)
		if err != nil {
			c.Logger().Error(err)
		}

		// Track changes through Mongo Change Stream
		changeDoc := struct {
			User user.Model `bson:"fullDocument"`
		}{}
		var users []*user.Model

		for h.Controller.Repo.ChangeStream.Next(&changeDoc) {
			users = append(users, &changeDoc.User)
		}

		err = ws.WriteJSON(users)
		if err != nil {
			c.Logger().Error(err)
		}
		// Writing through mongo change stream ends here
	}
}
