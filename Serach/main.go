package main
 
import (
    "database/sql"
    "fmt"
    "log"
    "strconv"
 
    "github.com/gofiber/fiber/v2"
    _ "github.com/go-sql-driver/mysql"
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
}
 
func main() {
    initDB()
    fmt.Println("Database successfully connected")
    defer db.Close()
 
    app := fiber.New()
 
    app.Get("/users", getUsers)
    app.Get("/user",getUsersCase)
 
    log.Fatal(app.Listen(":3000"))
}
 
func getUsers(c *fiber.Ctx) error {
    name :=c.Query("name")
    rows, err := db.Query("SELECT name, age FROM users WHERE name LIKE ?","%"+name+"%")
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
func getUsersCase(c *fiber.Ctx) error {
    queryValue := c.Query("q")
 
    if queryValue == "" {
        return c.Status(fiber.StatusBadRequest).SendString("Please provide the value ")
    }
 
    query := "SELECT name, age FROM users WHERE"
    args := make([]interface{}, 0)
 
    if age, err := strconv.Atoi(queryValue); err == nil {
        query += " age = ?"
        args = append(args, age)
    } else {
        query += " name LIKE ?"
        args = append(args, "%"+queryValue+"%")
    }
 
    rows, err := db.Query(query, args...)
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