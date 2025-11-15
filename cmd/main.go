package main

import (
	"avito-tech-internship/internal/server"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		slog.Error("could not connect to database: %v", err)
	}
	defer db.Close()

	s := server.NewServer(db)
	s.Start()
}
