package main
 
import (
    "context"
    "log"
    "sort"
    "strconv"
 
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
 
type livetv struct {
    En en `bson:"en" json:"en"`
}
 
type en struct {
    Key  int      `bson:"key" json:"key"`
    Tags []string `bson:"tags" json:"tags"`
}
 
var (
    client *mongo.Client
)
 
func main() {
    app := fiber.New()
 
    connectToMongoDB()
 
    app.Get("/details", getList)
    app.Get("/count",getCountList)
 
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
 
func getList(c *fiber.Ctx) error {
    keyStr := c.Query("key")
    key, err := strconv.Atoi(keyStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid key")
    }
 
    collection := client.Database("test").Collection("livetv")
 
    var result livetv
    err = collection.FindOne(context.Background(), bson.M{"en.key": key}).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return c.Status(fiber.StatusNotFound).SendString("Record not found")
        }
        return err
    }
 
 
    tags := result.En.Tags
 
    cursor, err := collection.Find(context.Background(), bson.M{"en.tags": bson.M{"$in": tags}})
    if err != nil {
        return err
    }
    defer cursor.Close(context.Background())
 
    var matchingRecords []livetv
 
 
    for cursor.Next(context.Background()) {
        var record livetv
        if err := cursor.Decode(&record); err != nil {
            return err
        }
        matchingRecords = append(matchingRecords, record)
    }
 
    if err := cursor.Err(); err != nil {
        return err
    }
 
    return c.JSON(matchingRecords)
}
 
func getCountList(c *fiber.Ctx) error {
    keyStr := c.Query("key")
    key, err := strconv.Atoi(keyStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid key")
    }
 
    collection := client.Database("test").Collection("livetv")
 
   
    var result livetv
    err = collection.FindOne(context.Background(), bson.M{"en.key": key}).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return c.Status(fiber.StatusNotFound).SendString("Record not found")
        }
        return err
    }
 
    tags := result.En.Tags
 
 
    cursor, err := collection.Find(context.Background(), bson.M{"en.tags": bson.M{"$in": tags}})
    if err != nil {
        return err
    }
    defer cursor.Close(context.Background())
 
    tagCounts := make(map[string]int)
 
   
    for cursor.Next(context.Background()) {
        var record livetv
        if err := cursor.Decode(&record); err != nil {
            return err
        }
 
        for _, tag := range record.En.Tags {
            tagCounts[tag]++
        }
    }
 
    if err := cursor.Err(); err != nil {
        return err
    }
 
 
    var response []map[string]interface{}
    for tag, count := range tagCounts {
        response = append(response, map[string]interface{}{
            "tag":   tag,
            "count": count,
        })
    }
 
 
    sort.Slice(response, func(i, j int) bool {
        return response[i]["count"].(int) > response[j]["count"].(int)
    })
 
    return c.JSON(response)
}