package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func main() {
	e := echo.New()

	type User struct {
		Name        string `json:"name" xml:"name" form:"name" query:"name"`
		Email       string `json:"email" xml:"email" form:"email" query:"email"`
		Description string `json:"description" xml:"description" form:"description" query:"description"`
	}
	client, err := mongo.Connect(context.TODO(), "mongodb://admin:admin123@localhost:27017")

	if err != nil {
		log.Fatal(err)
	}

	// Connect to collection
	collection := client.Database("users").Collection("users")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/users", func(c echo.Context) error {
		var users []*User
		// findOptions := options.Find()
		// findOptions.SetLimit(2)
		cur, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Fatal(err)
		}

		for cur.Next(context.TODO()) {
			// create a value into which the single document can be decoded
			var elem User
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, &elem)
			fmt.Printf("Found multiple documents (array of pointers): %+v\n", &elem)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		// Close the cursor once finished
		cur.Close(context.TODO())

		fmt.Printf("Found multiple documents (array of pointers): %+v\n", users)
		return c.JSONPretty(http.StatusOK, users, " ")
	})

	e.GET("/users/:name", func(c echo.Context) error {
		name := c.Param("name")
		filter := bson.M{"name": name}
		// create a value into which the result can be decoded
		var result User
		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Found a single document: %+v\n", result)
		return c.JSONPretty(http.StatusOK, result, " ")
	})

	e.POST("/users", func(c echo.Context) error {
		u := new(User)

		if err := c.Bind(u); err != nil {
			return err
		}

		insertResult, err := collection.InsertOne(context.TODO(), u)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
		return c.JSON(http.StatusCreated, u)
	})

	e.DELETE("/users/:name", func(c echo.Context) error {
		name := c.Param("name")
		filter := bson.M{"name": name}
		deleteResult, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
		return c.String(http.StatusOK, "Deleted 1 entry with Name "+name)
	})

	e.PUT("/users/:name", func(c echo.Context) error {
		name := c.Param("name")
		newName := c.FormValue("name")
		newEmail := c.FormValue("email")
		newDes := c.FormValue("description")

		filter := bson.M{"name": name}

		update := bson.M{"$set": bson.M{
			"name":        newName,
			"email":       newEmail,
			"description": newDes},
		}

		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
		return c.JSONPretty(http.StatusOK, updateResult, " ")
	})

	e.Logger.Fatal(e.Start(":8000"))
}
