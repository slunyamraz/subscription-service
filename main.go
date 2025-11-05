package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	httpSwagger "github.com/swaggo/http-swagger"

	config "s/Config"
	handler "s/Handler"
	repository "s/Repository"
	_ "s/docs"
)

// @title Subscriptions API
// @version 1.0
// @description REST service for managing user subscriptions
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.NewConfig()
	db, err := repository.NewDB(cfg)
	if err != nil {
		slog.Error("bad db!", "error", err)
		return
	}
	defer db.Close()

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Println("Applying migrations...")

		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatal("Set dialect failed:", err)
		}

		if err := goose.Up(db, "migrations"); err != nil {
			log.Fatal("Migration failed:", err)
		}

		log.Println("Migrations applied successfully")
		return
	}
	repo := repository.NewSubscriptionRepo(db)
	hand := handler.NewSubscriptionHandler(repo)
	http.Handle("/swagger/", httpSwagger.WrapHandler)
	handler.RegisterRoutes(hand)
	handler.ServeStart()
}
