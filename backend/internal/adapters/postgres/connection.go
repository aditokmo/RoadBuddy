package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBService interface {
	Health() map[string]string
	Close() error
	Pool() *pgxpool.Pool
}

type service struct {
	db     *pgxpool.Pool
	dbName string
}

func New(connString string, dbName string) DBService {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return &service{
		db:     pool,
		dbName: dbName,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"

	dbStats := s.db.Stat()
	stats["total_connections"] = fmt.Sprintf("%d", dbStats.TotalConns())
	stats["idle_connections"] = fmt.Sprintf("%d", dbStats.IdleConns())
	stats["acquired_connections"] = fmt.Sprintf("%d", dbStats.AcquiredConns())

	return stats
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", s.dbName)
	s.db.Close()
	return nil
}

func (s *service) Pool() *pgxpool.Pool {
	return s.db
}
