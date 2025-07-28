package main

import (
	"database/sql"
	"log"

	"github.com/artyomkorchagin/storeyourimages/internal/services/content"
	"github.com/artyomkorchagin/storeyourimages/internal/services/users"
	pgxcontent "github.com/artyomkorchagin/storeyourimages/internal/storage/pgx/content"
	pgxusers "github.com/artyomkorchagin/storeyourimages/internal/storage/pgx/users"
	bot "github.com/artyomkorchagin/storeyourimages/internal/tg-bot"
	"github.com/artyomkorchagin/storeyourimages/pkg/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	db, err := sql.Open("pgx", config.GetDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	repoContent := pgxcontent.NewRepository(db)
	svcContent := content.NewService(repoContent)

	repoUsers := pgxusers.NewRepository(db)
	svcUsers := users.NewService(repoUsers)

	svcs := bot.NewAllServices(svcUsers, svcContent)
	log.Println("Connected to database")
	bot.InitBot(svcs)
}
