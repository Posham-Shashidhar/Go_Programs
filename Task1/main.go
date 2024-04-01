package main
 
import (
    "context"
    "log"
    "net/http"
    //"sync"
 
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
 
type Person struct {
    Key         int    `json:"Key" bson:"Key"`
    Name        string `json:"name,omitempty" bson:"name,omitempty"`
    RollNo      int    `json:"RollNo,omitempty" bson:"RollNo,omitempty"`
    Address     string `json:"address,omitempty" bson:"address,omitempty"`
}
 
var (
    client *mongo.Client
    //mux    sync.Mutex
)
 
func main() {
    app := fiber.New()
 
    connectToMongoDB()
 
    app.Post("/person", createPerson)
    app.Get("/persons", getPersons)
 
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
 
    key, err := getNextKey()
    if err != nil {
        return err
    }
 
    person.Key = key
 
    collection := client.Database("test").Collection("people")
    _, err = collection.InsertOne(context.Background(), person)
    if err != nil {
        return err
    }
 
    return c.Status(http.StatusCreated).JSON(person.Key)
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
 
func getNextKey() (int, error) {
    collection := client.Database("test").Collection("people")
 
    var highestKey Person
    err := collection.FindOne(context.Background(), bson.D{}, options.FindOne().SetSort(bson.D{{"Key", -1}})).Decode(&highestKey)
    if err != nil && err != mongo.ErrNoDocuments {
        return 0, err
    }
 
    nextKey := highestKey.Key + 1
    return nextKey, nil
}