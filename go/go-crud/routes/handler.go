package route

import (
	"context"
	"log"
	"net/http"

	common "github.com/atlanhq/sc-go-common"
	"github.com/atlanhq/sc-go-sms/utils"
	"github.com/labstack/echo"
	"github.com/pradeepgangwar/go-websocket/controller"
	"github.com/pradeepgangwar/go-websocket/user"
)

type handler struct {
	Controller *controller.Controller
}

// NewHandler : This returns new handler object
func NewHandler(g *echo.Group, controller *controller.Controller) {
	handler := &handler{
		Controller: controller,
	}
	g.POST("", handler.UserPost)
	// g.POST("/callback", handler.SMSCallback)
	g.GET("", handler.UsersGet)
}

func (h *handler) UserPost(c echo.Context) error {

	var model POSTRequest

	// Binding
	if err := c.Bind(&model); err != nil {
		c.Logger().Error(err.Error())
		return utils.ThrowError(c, common.BIND_ERROR, "")
	}

	if ok, err := model.IsPostRequestValid(); !ok {
		c.Logger().Error(err.Error())
	}

	// Get Context
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userModel := user.NewModel(model.Name, model.Email, model.About)

	_, err := h.Controller.New(ctx, userModel)
	if err != nil {
		c.Logger().Error(err.Error())
		return utils.ThrowError(c, common.SMS_CLIENT_ERROR, "")
	}

	return c.JSON(http.StatusOK, userModel)
}

// func (h *handler) SMSCallback(c echo.Context) error {

// 	var model SMSCallback

// 	// Get Context
// 	ctx := c.Request().Context()
// 	if ctx == nil {
// 		ctx = context.Background()
// 	}

// 	// Binding
// 	if err := c.Bind(&model); err != nil {
// 		c.Logger().Error(err.Error())
// 		fmt.Println(err)
// 		return utils.ThrowError(c, common.BIND_ERROR, "")
// 	}
// 	_, err := h.controller.SMSController.Update(ctx, model.SmsSID, model.SMSStatus)
// 	if err != nil {
// 		c.Logger().Error(err.Error())
// 		return utils.ThrowError(c, common.SMS_CLIENT_ERROR, "")
// 	}
// 	return c.String(http.StatusOK, "Received the callback parameters")
// }

func (h *handler) UsersGet(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	users, err := h.Controller.FindAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSONPretty(http.StatusOK, users, " ")
}
