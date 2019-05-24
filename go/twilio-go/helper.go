package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/sfreiberg/gotwilio"
)

// Message type is the message type that we need to send
// type Message struct {
// 	To      string `json:"email" xml:"email" form:"email" query:"email"`
// 	Message string `json:"description" xml:"description" form:"description" query:"description"`
// }

// Home :This is the function that handles home route
func Home(c echo.Context) error {
	return c.String(http.StatusOK, "Hi, Here you can send messages to you favorite ones.")
}

// NewMessage : This sends a new Message to user
func NewMessage(c echo.Context) error {
	to := c.FormValue("to")
	message := c.FormValue("message")

	err := godotenv.Load()

	if err != nil {
		c.Logger().Error(err.Error())
		// e.Logger.Info("Environment variables not set perfectly")
	}

	newMessage := new(Message)
	err = c.Bind(newMessage)

	if to == "" || message == "" || err != nil {
		fmt.Println("Not all values are set")
		return c.String(http.StatusBadRequest, "Not all values are se")
	}

	accountSid := os.Getenv("TWILIOACCOUNTSID")
	authToken := os.Getenv("TWILIOAUTHTOKEN")
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := os.Getenv("FROM")
	res, exp, err := twilio.SendSMS(from, to, message, "", "")

	if err != nil {
		fmt.Println("Error in sending sms via twilio")
	} else {
		fmt.Println(res)
	}

	insertResult, err := collection.InsertOne(context.TODO(), newMessage)
	if err != nil {
		fmt.Println("Could not insert value into the database")
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return c.JSON(http.StatusCreated, newMessage)
}
