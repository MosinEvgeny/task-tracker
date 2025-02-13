package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(databaseURL string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ошибка при проверке подключения к базе данных: %w", err)
	}

	log.Println("Connected to PostgreSQL database")

	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() error {
	log.Println("Closing database connection...")
	if err := p.DB.Close(); err != nil {
		return fmt.Errorf("ошибка при закрытии подключения к базе данных: %w", err)
	}
	return nil
}
