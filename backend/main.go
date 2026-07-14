package main

import (
	"log"
	"net/http"

	"istream/backend/internal/config"
	"istream/backend/internal/database"
	"istream/backend/internal/handlers"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	srv := handlers.NewServer(db, cfg)
	log.Printf("istream backend listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, srv); err != nil {
		log.Fatal(err)
	}
}
