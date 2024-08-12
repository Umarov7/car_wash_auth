package postgres

import (
	"auth-service/config"
	"auth-service/storage"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) User() storage.IUserStorage {
	return NewUserRepo(s.db)
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		log.Printf("error while closing the database connection: %v", err)
	}
}

func ConnectDB(cfg *config.Config) (storage.IStorage, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_NAME, cfg.DB_PASSWORD,
	)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}
