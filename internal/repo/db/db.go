package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JMURv/org-struct/internal/config"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

func New(conf config.Config) *Repository {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Database,
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("failed to connect to database", zap.Error(err))
	}

	sqlDB, err := conn.DB()
	if err != nil {
		zap.L().Fatal("failed to get sql.DB", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err = sqlDB.Ping(); err != nil {
		zap.L().Fatal("failed to ping database", zap.Error(err))
	}

	if err = applyMigrations(sqlDB); err != nil {
		zap.L().Fatal("failed to apply migrations", zap.Error(err))
	}

	return &Repository{
		conn: conn,
	}
}

func (r *Repository) Close(ctx context.Context) error {
	done := make(chan error, 1)
	go func() {
		d, err := r.conn.DB()
		if err != nil {
			zap.L().Error("failed to get database obj", zap.Error(err))

			done <- err
		}

		done <- d.Close()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func applyMigrations(db *sql.DB) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		zap.L().Error("failed to set postgres dialect", zap.Error(err))
		return err
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		return err
	}

	zap.L().Info("migrations applied")
	return nil
}
