package route

import (
	"context"

	"github.com/labstack/echo"
	"github.com/pradeepgangwar/go-websocket/controller"
	"github.com/pradeepgangwar/go-websocket/websocket"
)

type handler struct {
	Controller *controller.Controller
}

// var (
// 	upgrader = websocket.Upgrader{}
// )

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

	// ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	// if err != nil {
	// 	return err
	// }
	// defer ws.Close()

	ws := websocket.NewWebsocket(c)
	h.Controller.WebSocket = ws

	allusers := h.Controller.AllUser(ctx)

	// Write
	// err = ws.WriteJSON(allusers)
	// if err != nil {
	// 	c.Logger().Error(err)
	// }
	ws.SendUsers(allusers)
	go ws.Listen(h.Controller.Repo)

	for {
		// var newUser user.Model

		// // Read
		// user = ws.ReadJSON(&newUser)
		// if err != nil {
		// 	c.Logger().Error(err)
		// }
		// changeDoc := struct {
		// 	User user.Model `bson:"fullDocument"`
		// }{}
		// var newusers []*user.Model

		// for h.Controller.Repo.ChangeStream.Next(&changeDoc) {
		// 	fmt.Println("Hello")
		// 	newusers = append(newusers, &changeDoc.User)
		// }

		// ws.SendUsers(newusers)

		newUser := ws.NewUser()

		err := h.Controller.NewUser(ctx, newUser)
		if err != nil {
			c.Logger().Error(err)
		}

		// // Track changes through Mongo Change Stream

		// if err != nil {
		// 	c.Logger().Error(err)
		// }
		// // Writing through mongo change stream ends here
	}
}
