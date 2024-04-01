package main

import (
    "github.com/gofiber/fiber/v2"
    "log"
)

func main() {
    // Create a new Fiber instance
    app := fiber.New()

    // Define a GET route
    app.Get("/hello", func(c *fiber.Ctx) error {
        // Return a simple response
        return c.SendString("Hello, World!")
    })

    // Start the Fiber server on port 3000
    if err := app.Listen(":3000"); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}