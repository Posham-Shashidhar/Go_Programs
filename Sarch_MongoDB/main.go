package main
 
import (
    "context"
    // "encoding/json"
    "log"
    //"net/http"
 
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
 
 
func getPersons(c *fiber.Ctx) error {
    queryParam := c.Query("q")
 
    filter := bson.M{
        "$or": []bson.M{
            {"_id": bson.M{"$regex": queryParam, "$options": "i"}},
            {"firstname": bson.M{"$regex": queryParam, "$options": "i"}},
            {"lastname": bson.M{"$regex": queryParam, "$options": "i"}},
            {"email": bson.M{"$regex": queryParam, "$options": "i"}},
        },
    }
 
    collection := client.Database("test").Collection("people")
 
    cursor, err := collection.Find(context.Background(), filter)
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