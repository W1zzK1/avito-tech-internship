package main

import (
	"avito-tech-internship/server"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	dbURL := "postgres://user:password@localhost:5432/reviewer_db?sslmode=disable"
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	s := server.NewServer(db)
	s.Start()
}
