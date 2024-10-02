package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	connect()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	}))
	app.Static("/", "./ui/dist")

	app.Get("", func(c *fiber.Ctx) error {
    return c.Render("/ui/dist/index", fiber.Map{})
	})

	app.Get("/api/records", func(c *fiber.Ctx) error {
		time.Sleep(200 * time.Millisecond)
		records, err := getRecords()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(records)
	})

	app.Post("/api/records", func(c *fiber.Ctx) error {
		time.Sleep(200 * time.Millisecond)
		var payload struct {
			Name string `json:"name"`
		}
		err := c.BodyParser(&payload)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		if payload.Name == "" {
			return c.Status(400).SendString("Name is required")
		}
		if len(payload.Name) > 55 {
			return c.Status(400).SendString("Name cannot be longer than 55 characters")
		}
    newRecord, err := createRecord(payload.Name)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
	  return c.Status(201).JSON(newRecord)
	})

	app.Put("/api/records", func(c *fiber.Ctx) error {
    time.Sleep(400 * time.Millisecond)
		var payload struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}
		err := c.BodyParser(&payload)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		if payload.ID == 0 {
			return c.Status(400).SendString("Id is required")
		}
		if payload.Name == "" {
			return c.Status(400).SendString("Name is required")
		}
		if len(payload.Name) > 55 {
			return c.Status(400).SendString("Name cannot be longer than 55 characters")
		}
		err = editRecord(payload.ID, payload.Name)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
	  return c.Status(200).JSON(payload)
	})

	app.Delete("/api/records/:id", func(c *fiber.Ctx) error {
    time.Sleep(400 * time.Millisecond)
		err := deleteRecord(c.Params("id"))
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendStatus(204)
	})

	log.Fatal(app.Listen(":8081"))
}
