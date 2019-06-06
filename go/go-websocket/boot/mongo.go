package boot

import (
	"fmt"

	"github.com/globalsign/mgo"
)

// SetupMongo Sets up the mongo connection
func SetupMongo(c *Config) (*mgo.Session, error) {
	//////////////////////////////////////////
	/** Initialize the database connections **/
	//////////////////////////////////////////

	// DialInfo := &mgo.DialInfo{
	// 	AppName:  "go-websocket",
	// 	Database: c.MongoDatabaseName,
	// 	Addrs:    []string{c.MongoHost},
	// 	Timeout:  time.Duration(c.MongoTimeout) * time.Second,
	// 	Direct:   false,
	// 	// Username: c.MongoUserName,
	// 	// Password: c.MongoPassword,
	// }

	connection, err := mgo.Dial("mongodb://localhost:27019/?replicaSet=rs0")
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to mongo")
	return connection, nil
}
