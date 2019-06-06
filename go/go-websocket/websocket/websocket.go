package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/pradeepgangwar/go-websocket/repo"
	"github.com/pradeepgangwar/go-websocket/user"
)

var (
	upgrader = websocket.Upgrader{}
)

// WebSocket : contains the websocket
type WebSocket struct {
	WebSocket *websocket.Conn
}

// NewWebsocket : initializes the new websocket
func NewWebsocket(c echo.Context) *WebSocket {
	ws, _ := upgrader.Upgrade(c.Response(), c.Request(), nil)
	return &WebSocket{
		ws,
	}
}

// SendUsers : sends the users to the app
func (w *WebSocket) SendUsers(users []*user.Model) {
	w.WebSocket.WriteJSON(users)
}

// NewUser This function reads the changes in the socket and writes to mongo
func (w *WebSocket) NewUser() user.Model {
	var newUser user.Model
	w.WebSocket.ReadJSON(&newUser)
	return newUser
}

// Listen this listens for the change in the mongo
func (w *WebSocket) Listen(repo *repo.MongoRepo) {
	for {
		changeDoc := struct {
			User user.Model `bson:"fullDocument"`
		}{}
		var newusers []*user.Model
		for repo.ChangeStream.Next(&changeDoc) {
			newusers = append(newusers, &changeDoc.User)
		}
		if len(newusers) > 0 {
			w.WebSocket.WriteJSON(newusers)
		}
	}
}
