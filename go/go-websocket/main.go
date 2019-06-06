package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/pradeepgangwar/go-websocket/boot"
	"github.com/pradeepgangwar/go-websocket/controller"
	"github.com/pradeepgangwar/go-websocket/middleware"
	"github.com/pradeepgangwar/go-websocket/repo"
	route "github.com/pradeepgangwar/go-websocket/routes"
)

func main() {

	// Setup echo
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	// Setup the config
	config, err := boot.SetupConfig()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Connect to Mongo
	mongoClient, err := boot.SetupMongo(config)
	if err != nil {
		e.Logger.Debug(err)
		e.Logger.Fatal("Error connecting to mongo at host " + config.MongoHost + " port " + config.MongoPort)
	}

	// Setup middleware
	boot.SetupMiddleware(e)

	e.Use(middleware.JWT([]byte(config.JwtSecret)))

	GroupBaseURL := "/users"

	// Get new mongo structure that contains collection and connection info
	repository := repo.NewRepository(mongoClient)
	controller := controller.NewController(config, repository)

	// Setup the routes
	usersGroup := e.Group(GroupBaseURL)
	route.NewHandler(usersGroup, controller)

	// changeDoc := struct {
	// 	User user.Model `bson:"fullDocument"`
	// }{}
	// var users []*user.Model

	// for repository.ChangeStream.Next(&changeDoc) {
	// 	users = append(users, &changeDoc.User)
	// 	fmt.Println(users)
	// }

	e.Logger.Fatal(e.Start(config.ServiceHost))
}
