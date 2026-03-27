package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tugkanmeral/the-host-go/internal/cli"
	"github.com/tugkanmeral/the-host-go/internal/config"
	"github.com/tugkanmeral/the-host-go/internal/database"
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

	if err := cli.Run(authSvc, noteSvc); err != nil {
		fmt.Fprintf(os.Stderr, "cli: %v\n", err)
		os.Exit(1)
	}
}
