package websocket

import (
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/pradeepgangwar/go-websocket/repo"
	"github.com/pradeepgangwar/go-websocket/user"
)

// DocumentID :
type DocumentID struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
}

// ChangeDoc :
type ChangeDoc struct {
	ID        DocumentID `bson:"documentKey" json:"id"`
	User      user.Model `bson:"fullDocument" json:"user"`
	Operation string     `bson:"operationType" json:"operation"`
}

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
		var changedoc *ChangeDoc
		for repo.ChangeStream.Next(&changedoc) {
			w.WebSocket.WriteJSON(changedoc)
		}
	}
}
