package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tugkanmeral/the-host-go/internal/config"
	"github.com/tugkanmeral/the-host-go/internal/database"
	"github.com/tugkanmeral/the-host-go/internal/handlers"
	"github.com/tugkanmeral/the-host-go/internal/service"
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

	db := database.GetDB()
	authSvc := service.NewAuthService(db)
	noteSvc := service.NewNoteService(db)
	authHandler := handlers.NewAuthHandler(authSvc)
	noteHandler := handlers.NewNoteHandler(noteSvc)

	addr := cfg.Port
	if addr == "" {
		addr = ":8000"
	} else if addr[0] != ':' {
		addr = ":" + addr
	}

	app := fiber.New()

	auth := app.Group("/api/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)

	note := app.Group("/api/note")
	note.Post("/", noteHandler.Add)
	note.Get("/", noteHandler.GetList)
	note.Get("/:id", noteHandler.Get)
	note.Put("/:id", noteHandler.Update)
	note.Delete("/:id", noteHandler.Delete)

	log.Fatal(app.Listen(addr))
}
