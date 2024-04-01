package Fiber

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func CreateOne() {
	c := fiber.New()
	c.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello wold")
	})

	if err := c.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v",err)
	}

}
