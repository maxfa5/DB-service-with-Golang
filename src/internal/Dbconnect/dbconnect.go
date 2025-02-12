package dbconnect

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

// InitDB инициализирует пул подключений к PostgreSQL.
func InitDB(logger *slog.Logger, connectionString string) error {
	var err error
	dbPool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		logger.Error("Ошибка подключения к Kafka", slog.String("error", err.Error()))
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	query := "INSERT INTO public.threads(name, id_theme) VALUES ('test1', 1);"

	_, err = dbPool.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Error inserting thread: %v", err)
		return err
	}
	fmt.Println("Thread created successfully")
	return nil
}

// GetDB возвращает пул подключений к базе данных.
func GetDB() *pgxpool.Pool {
	return dbPool
}

// CloseDB закрывает пул подключений.
func CloseDB() {
	dbPool.Close()
}
