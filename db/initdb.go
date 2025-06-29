package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

// NewDBConfig создает конфиг с значениями по умолчанию
func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "1234",
		DBName:   "shift_board",
		SSLMode:  "disable",
	}
}

// InitDB инициализирует подключение к PostgreSQL
func InitDB(cfg *DBConfig) (*Database, error) {
	if cfg == nil {
		cfg = NewDBConfig()
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка ping БД: %w", err)
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Успешное подключение к PostgreSQL")
	return &Database{db}, nil
}

// Close закрывает подключение к БД
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}
