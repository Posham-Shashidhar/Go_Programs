package main

import (
	"context"
	//"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type en struct {
	Tags []string `bson:"tags"`
}

var (
	client *mongo.Client
)

func main() {

	app := fiber.New()

	connectToDB()

	app.Get("/result", getList)

	log.Fatal(app.Listen(":3000"))

}
func connectToDB() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error

	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
}

func getList(c *fiber.Ctx) error {

    var tags []string

    collection := client.Database("test").Collection("livetv")

    
    cursor, err := collection.Find(context.Background(), bson.M{"en": bson.M{"tags": tags , "key" : "1"}})
    if err != nil {
        return err
    }
    defer cursor.Close(context.Background())

  
    //var allTags []string

    
    for cursor.Next(context.Background()) {
       // var en en
        if err := cursor.Decode(&tags); err != nil {
            return err
        }
    
        //allTags = append(allTags, en.Tags...)
    }

 
    if err := cursor.Err(); err != nil {
        return err
    }

  
    return c.JSON(tags)
}


