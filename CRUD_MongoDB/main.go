package main
 
import (
    "context"
    // "encoding/json"
    "log"
    "net/http"
 
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
 
type Person struct {
    ID        string `json:"id,omitempty" bson:"_id,omitempty"`
    FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
    LastName  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
    Email     string `json:"email,omitempty" bson:"email,omitempty"`
}
 
var client *mongo.Client
 
func main() {
    app := fiber.New()
 
    connectToMongoDB()
 
    app.Post("/person", createPerson)
 
    app.Get("/persons", getPersons)
 
    app.Put("/person/:id", updatePerson)
 
    app.Delete("/persons/:id", deletePerson)
 
 
    log.Fatal(app.Listen(":3000"))
}
 
func connectToMongoDB() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    var err error
    client, err = mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
}
 
func createPerson(c *fiber.Ctx) error {
    var person Person
    if err := c.BodyParser(&person); err != nil {
        return err
    }
 
    collection := client.Database("test").Collection("people")
    result, err := collection.InsertOne(context.Background(), person)
    if err != nil {
        return err
    }
 
    return c.Status(http.StatusCreated).JSON(result.InsertedID)
}
func getPersons(c *fiber.Ctx) error {
    collection := client.Database("test").Collection("people")
    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        return err
    }
    defer cursor.Close(context.Background())
 
    var persons []Person
    if err := cursor.All(context.Background(), &persons); err != nil {
        return err
    }
 
    return c.JSON(persons)
}
func updatePerson(c *fiber.Ctx) error {
    id := c.Params("id")
    var updatedPerson Person
    if err := c.BodyParser(&updatedPerson); err != nil {
        return err
    }
 
    collection := client.Database("test").Collection("people")
    filter := bson.M{"_id": id}
    update := bson.M{"$set": updatedPerson}
    _, err := collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return err
    }
 
    return c.SendString("Person successfully updated")
}
func deletePerson(c *fiber.Ctx) error {
    id := c.Params("id")
 
    collection := client.Database("test").Collection("people")
    _, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
    if err != nil {
        return err
    }
 
    return c.SendString("Person successfully deleted")
}