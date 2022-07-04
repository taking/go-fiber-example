package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Test struct {
	Name string `json:"ConnectionName"`
}

func main() {
	// Custom config
	app := fiber.New(fiber.Config{
		Prefork:       false, // 멀티 go 프로세스
		CaseSensitive: true,  // 대소문자 구분 /Foo != /foo
		StrictRouting: true,  // 엄격한 라우팅 /foo !== /foo/
		ServerHeader:  "Fiber",
		AppName:       "Test App v1.0.1",
	})

	app.Use(
		logger.New(),
	)

	app.Get("/", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		log.Println("hostname : ", c.Hostname())
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/json", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name": "taking",
			"age":  "secret",
		})
	})

	app.Get("/test3", func(c *fiber.Ctx) error {
		fmt.Println("1st route!")
		return c.Next()
	})

	app.Get("*", func(c *fiber.Ctx) error {
		fmt.Println("2nd route!")
		return c.Next()
	})

	app.Get("/test3", func(c *fiber.Ctx) error {
		fmt.Println("3rd route!")
		return c.SendString("Hello, World!")
	})

	app.Get("/api/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	}).Name("api")

	app.Get("/user/:name/*", func(c *fiber.Ctx) error {
		log.Println(c.AllParams())    // "{"name": "fenny"}"
		log.Println(c.Params("name")) // "fenny"
		log.Println(c.Params("*"))    // "fenny/*"
		log.Println(c.Query("*"))     // ""

		return c.JSON(fiber.Map{
			"name": c.Params("name"),
		})
	})

	app.Post("/user", func(c *fiber.Ctx) error {
		p := new(Test)
		if err := c.BodyParser(p); err != nil {
			return err
		}
		log.Println("Cofigname :", p.Name)
		return c.Send(c.Body())
	})

	users := [...]string{"Alice", "Bob", "Charlie", "David"}
	fmt.Println("users : ", users)

	app.Get("/error", func(c *fiber.Ctx) error {
		return c.JSON(fiber.NewError(782, "Custom error messenger"))
	})

	// GET /dictionary.txt
	app.Get("/:file.:ext", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("📃 %s.%s", c.Params("file"), c.Params("ext"))
		return c.SendString(msg) // => 📃 dictionary.txt
	})

	data, _ := json.MarshalIndent(app.GetRoute("api"), "", "  ")
	fmt.Print(string(data))

	log.Fatal(app.Listen(":52529"))
}
