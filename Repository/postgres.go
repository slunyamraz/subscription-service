package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	config "s/Config"
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.SQL.Host, cfg.SQL.Port, cfg.SQL.Username, cfg.SQL.Password, cfg.SQL.Database)
	print(dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("db open error", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		slog.Error("db ping error", err)
		return nil, err
	}
	slog.Info("db open success")
	return db, nil
}
