package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sfreiberg/gotwilio"
	"gopkg.in/mgo.v2/bson"
)

// Message type is the message type that we need to send
type Message struct {
	To      string `json:"to" xml:"to" form:"to" query:"to"`
	Message string `json:"message" xml:"message" form:"message" query:"message"`
}

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	client, err := mongo.Connect(context.TODO(), "mongodb://admin:admin123@localhost:27017")

	if err != nil {
		e.Logger.Debug(err)
		e.Logger.Fatal("Problem connecting to mongo database- localhost, 27017")
	} else {
		e.Logger.Info("Connected to MongoDB - localhost, 27017")
	}

	// Connect to collection
	collection := client.Database("users").Collection("message")

	// Here we define all the routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hi, Here you can send messages to you favorite ones.")
	})

	e.POST("/new-message", func(c echo.Context) error {
		to := c.FormValue("to")
		message := c.FormValue("message")

		err := godotenv.Load()

		if err != nil {
			e.Logger.Info("Environment variables not set perfectly")
		}

		newMessage := Message{to, message}

		if to == "" || message == "" || err != nil {
			e.Logger.Info("Not all values are set")
			return c.String(http.StatusBadRequest, "Not all values are set")
		}

		accountSid := os.Getenv("TWILIOACCOUNTSID")
		authToken := os.Getenv("TWILIOAUTHTOKEN")
		twilio := gotwilio.NewTwilioClient(accountSid, authToken)

		from := os.Getenv("FROM")
		res, excep, err := twilio.SendSMS(from, to, message, "", "")

		if err != nil || excep != nil {
			e.Logger.Info("Error in sending sms via twilio")
		} else {
			e.Logger.Info(res)
		}

		insertResult, err := collection.InsertOne(context.TODO(), newMessage)
		if err != nil {
			e.Logger.Info("Could not insert value into the database")
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
		return c.JSON(http.StatusCreated, res)
	})

	e.GET("/messages", func(c echo.Context) error {
		var messages []*Message
		ctx := c.Request().Context()
		if ctx == nil {
			ctx = context.Background()
		}

		size, err := collection.Count(ctx, bson.M{})
		if size == 0 {
			return c.String(http.StatusNoContent, "Either there are no values or there is some internal error")
			e.Logger.Info("Some error in finding the messages")
		}

		cur, err := collection.Find(ctx, bson.M{})
		if err != nil {
			return c.String(http.StatusInternalServerError, "There is some internal error")
			e.Logger.Info("Some error in finding the messages")
		}

		for cur.Next(context.TODO()) {
			// create a value into which the single document can be decoded
			var elem Message
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			messages = append(messages, &elem)
			fmt.Printf("Found multiple documents (array of pointers): %+v\n", &elem)
		}

		if err := cur.Err(); err != nil {
			e.Logger.Info("Some error in finding the messages")
		}

		// Close the cursor once finished
		cur.Close(context.TODO())

		fmt.Printf("Found multiple messages : %+v\n", messages)
		return c.JSONPretty(http.StatusOK, messages, " ")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
