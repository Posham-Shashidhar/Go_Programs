package main

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var db *sql.DB

func main() {
	fmt.Println("Hello World")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World")

	})

	app.Get("/hello",func(c *fiber.Ctx) error {
		return c.SendString("Hello fiber$$$$")
	})

	app.Listen(":8080")
}
