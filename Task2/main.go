package main

import (
	"context"
	"log"

	// "net/http"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Page struct {
	Key      int        `json:"key" bson:"key"`
	Name     string     `json:"name" bson:"name"`
	Type     string     `json:"type" bson:"type"`
	Playlist []Playlist `json:"playlist" bson:"playlist"`
}

type Playlist struct {
	Key  int    `json:"key" bson:"key"`
	Name string `json:"name" bson:"name"`
}

var client *mongo.Client
var keyMutex sync.Mutex
var pageKey, playlistKey int

const (
	dbName             = "test"
	pageCollection     = "Page"
	playlistCollection = "Playlist"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	app := fiber.New()

	app.Post("/playlist", AddPlaylist)
	app.Post("/page", AddPage)
	app.Get("/playlist", GetPlaylist)
	app.Get("/page", GetPage)

	log.Fatal(app.Listen(":3000"))
}

func UpdateAllPagesWithPlaylists(playlist map[string]interface{}) error {
	pageCollection := client.Database(dbName).Collection(pageCollection)

	update := bson.M{"$push": bson.M{"playlist": playlist}}
	_, err := pageCollection.UpdateMany(context.Background(), bson.M{}, update)
	if err != nil {
		return err
	}

	return nil
}

func AddPlaylist(c *fiber.Ctx) error {
	var playlist Playlist
	if err := c.BodyParser(&playlist); err != nil {
		return err
	}

	playlistKey, err := getNextPlaylistKey()
	if err != nil {
		return err
	}
	playlist.Key = playlistKey

	playlistCollection := client.Database(dbName).Collection(playlistCollection)
	_, err = playlistCollection.InsertOne(context.Background(), playlist)
	if err != nil {
		return err
	}
	pagePlaylist := make(map[string]interface{})
	pagePlaylist["key"] = playlistKey
	pagePlaylist["name"] = playlist.Name

	err = UpdateAllPagesWithPlaylists(pagePlaylist)
	if err != nil {
		return err
	}

	return c.SendString("Playlist added successfully")
}

func AddPage(c *fiber.Ctx) error {
	var page Page
	if err := c.BodyParser(&page); err != nil {
		return err
	}

	pageKey, err := getNextPageKey()
	if err != nil {
		return err
	}
	page.Key = pageKey

	pageCollection := client.Database(dbName).Collection(pageCollection)
	_, err = pageCollection.InsertOne(context.Background(), page)
	if err != nil {
		return err
	}

	playlistCollection := client.Database(dbName).Collection(playlistCollection)
	cursor, err := playlistCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	var playlists []Playlist
	if err := cursor.All(context.Background(), &playlists); err != nil {
		return err
	}

	page.Playlist = playlists

	_, err = pageCollection.UpdateOne(context.Background(), bson.M{"key": pageKey}, bson.M{"$set": bson.M{"playlist": playlists}})
	if err != nil {
		return err
	}

	return c.SendString("Page added successfully")
}

func getNextPlaylistKey() (int, error) {
	keyMutex.Lock()
	defer keyMutex.Unlock()

	playlistCollection := client.Database(dbName).Collection(playlistCollection)

	var highestPlaylist Playlist
	err := playlistCollection.FindOne(context.Background(), bson.D{}, options.FindOne().SetSort(bson.D{{"key", -1}})).Decode(&highestPlaylist)
	if err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	}

	nextPlaylistKey := highestPlaylist.Key + 1
	return nextPlaylistKey, nil
}

func getNextPageKey() (int, error) {
	keyMutex.Lock()
	defer keyMutex.Unlock()

	pageCollection := client.Database(dbName).Collection(pageCollection)

	var highestPage Page
	err := pageCollection.FindOne(context.Background(), bson.D{}, options.FindOne().SetSort(bson.D{{"key", -1}})).Decode(&highestPage)
	if err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	}

	nextPageKey := highestPage.Key + 1
	return nextPageKey, nil
}

func GetPlaylist(c *fiber.Ctx) error {
	key := c.Query("key")

	playlistKey, err := strconv.Atoi(key)
	if err != nil {
		return err
	}

	playlistCollection := client.Database(dbName).Collection(playlistCollection)
	var playlist Playlist
	err = playlistCollection.FindOne(context.Background(), bson.M{"key": playlistKey}).Decode(&playlist)
	if err != nil {
		return err
	}

	return c.JSON(playlist)
}

func GetPage(c *fiber.Ctx) error {
	key := c.Query("key")

	pageKey, err := strconv.Atoi(key)
	if err != nil {
		return err
	}

	pageCollection := client.Database(dbName).Collection(pageCollection)
	var page Page
	err = pageCollection.FindOne(context.Background(), bson.M{"key": pageKey}).Decode(&page)
	if err != nil {
		return err
	}

	return c.JSON(page)
}
