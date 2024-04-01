package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:12345@tcp(localhost:3306)/go_database")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database is connected successfully")
}

func main() {
	initDB()

	defer db.Close()

	app := fiber.New()

	app.Post("/users", createUser)
	app.Get("/users",getUsers)
	app.Put("/users/:id", updateUser)
	app.Delete("/delete/:id",deleteUser)

	log.Fatal(app.Listen(":3000"))
}

func createUser(c *fiber.Ctx) error {
	user := new(User)

	log.Println("Request Body: ", c.Body())
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing JSON data")
	}

	_, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
func getUsers(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT name, age FROM users")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Name, &user.Age); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		users = append(users, user)
	}

	return c.JSON(users)
}
func updateUser(c *fiber.Ctx) error {
	user := new(User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing JSON data")
	}
	name := c.Query("name")

	userID := c.Query("id")

	_, err := db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", name, user.Age, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("User updated successfully")
}
func deleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	_, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("User deleted successfully")
}
