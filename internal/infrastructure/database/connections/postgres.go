package db_connections

import (
	"context"
	"fmt"
	"golang/internal/infrastructure/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	
)


func runMigrations(url string) error {
	migration, err := migrate.New("file://migrations", url)

	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run up migrations")
	}

	defer migration.Close()
	return nil
}


func NewPostgresConnection() (*pgxpool.Pool, error) {
	cfg, err := config.LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}
	
	connUrl := cfg.ConnectionString()
	
	err = runMigrations(connUrl)
	if err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
        pool.Close()
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
	return pool, nil
}

