package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"s/service"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	httpSwagger "github.com/swaggo/http-swagger"

	"s/config"
	_ "s/docs"
	"s/handler"
	"s/repository"
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
	outboxRepo := repository.NewOutbox(db)
	serviceSub := service.NewSubscriptionSerivce(db, repo, outboxRepo)

	hand := handler.NewSubscriptionHandler(serviceSub)
	http.Handle("/swagger/", httpSwagger.WrapHandler)
	handler.RegisterRoutes(hand)
	print("server start")
	handler.ServeStart()
}
