package route

import (
	"context"
	"fmt"

	"github.com/atlanhq/sc-go-edge/module/user"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/pradeepgangwar/go-websocket/controller"
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

// func (h *handler) UserPost(c echo.Context) error {

// 	var model POSTRequest

// 	// Binding
// 	if err := c.Bind(&model); err != nil {
// 		c.Logger().Error(err.Error())
// 		return utils.ThrowError(c, common.BIND_ERROR, "")
// 	}

// 	if ok, err := model.IsPostRequestValid(); !ok {
// 		c.Logger().Error(err.Error())
// 	}

// 	// Get Context
// 	ctx := c.Request().Context()
// 	if ctx == nil {
// 		ctx = context.Background()
// 	}

// 	userModel := user.NewModel(model.Name, model.Email, model.About)

// 	_, err := h.Controller.New(ctx, userModel)
// 	if err != nil {
// 		c.Logger().Error(err.Error())
// 		return utils.ThrowError(c, common.SMS_CLIENT_ERROR, "")
// 	}

// 	return c.JSON(http.StatusOK, userModel)
// }

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

	var user user.Model

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		msg := ws.ReadJSON(&user)
		fmt.Println(msg)
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
