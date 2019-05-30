package boot

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
)

// SetupMongo Sets up the mongo connection
func SetupMongo(c *Config) (*mgo.Session, error) {
	//////////////////////////////////////////
	/** Initialize the database connections **/
	//////////////////////////////////////////

	DialInfo := &mgo.DialInfo{
		AppName:  "go-websocket",
		Database: c.MongoDatabaseName,
		Addrs:    []string{c.MongoHost},
		Timeout:  time.Duration(c.MongoTimeout) * time.Second,
		Username: c.MongoUserName,
		Password: c.MongoPassword,
	}

	connection, err := mgo.DialWithInfo(DialInfo)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to mongo")
	return connection, nil
}
