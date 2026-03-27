package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tugkanmeral/the-host-go/internal/config"
	"github.com/tugkanmeral/the-host-go/internal/database"
	"github.com/tugkanmeral/the-host-go/internal/handlers"
)

func main() {
	cfg := config.LoadConfig()

	if err := database.ConnectDB(cfg.MongoURI, cfg.DBName); err != nil {
		log.Fatalf("MongoDB bağlantı hatası: %v", err)
	}
	defer func() {
		if err := database.DisconnectDB(); err != nil {
			log.Printf("MongoDB bağlantısı kapatılırken hata: %v", err)
		}
	}()

	setupRoutes()
}

func setupRoutes() {
	app := fiber.New()

	// Auth
	auth := app.Group("api/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/register", handlers.Register)

	// Note
	note := app.Group("/api/note")
	note.Post("/", handlers.Add)
	note.Get("/", handlers.GetList)
	note.Get("/:id", handlers.Get)
	note.Put("/:id", handlers.Update)

	log.Fatal(app.Listen(":8000"))
}
