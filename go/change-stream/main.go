package main

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pradeepgangwar/change-stream/boot"
	"github.com/pradeepgangwar/change-stream/user"
)

func main() {

	// Setup the config
	config, err := boot.SetupConfig()
	if err != nil {
		fmt.Println(err)
	}

	//////////////////////////////////////////
	/** Initialize the database connections **/
	//////////////////////////////////////////

	DialInfo := &mgo.DialInfo{
		AppName:        "go-websocket",
		Database:       config.MongoDatabaseName,
		Addrs:          []string{config.MongoHost},
		Timeout:        time.Duration(config.MongoTimeout) * time.Second,
		Direct:         false,
		ReplicaSetName: "rs0",
		// Username: config.MongoUserName,
		// Password: config.MongoPassword,
	}

	// connection, err := mgo.Dial("mongodb://localhost:27019/?replicaSet=rs0")
	// if err != nil {
	// 	return nil, err
	// }
	mongoClient, err := mgo.DialWithInfo(DialInfo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully connected to mongo")

	defer mongoClient.Close()

	pipeline := []bson.M{}
	coll := mongoClient.DB(config.MongoDatabaseName).C(config.MongoCollectionName)
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{})
	if err != nil {
		fmt.Println(err)
	}

	// var users []*user.Model
	// coll.Find(nil).All(&users)
	// fmt.Println(users)

	defer changeStream.Close()

	changeDoc := struct {
		User user.Model `bson:"fullDocument"`
	}{}
	var users []*user.Model

	for {
		for changeStream.Next(&changeDoc) {
			users = append(users, &changeDoc.User)
			fmt.Println(users)
		}
	}

	// if err := changeStream.Close(); err != nil {
	// 	fmt.Println(err)
	// }

}
