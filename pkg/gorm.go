package pkg

import (
	"context"
	"fmt"
	"time"

	"main/internal/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormDB struct {
	*gorm.DB
}

func NewGormDatabase(logger Logger, env config.Env) GormDB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.PGHost,
		env.PGPort,
		env.PGUser,
		env.PGPass,
		env.PGName,
	)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connStr,
	}), &gorm.Config{
		Logger: NewGormLogger(logger, gormlogger.Info),
	})
	if err != nil {
		logger.Fatalf("failed to initialize GORM: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Fatalf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Fatalf("failed to ping PostgreSQL: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL")

	if err := RunPostgresMigrations(sqlDB, env.MigrationPath); err != nil {
		logger.Fatalf("migration failed: %w", err)
	}
	logger.Info("Successfully applied database migrations")

	return GormDB{gormDB}
}

func (db GormDB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db GormDB) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		return fn(tx.WithContext(ctx))
	})
}

func (db GormDB) Paginate(query *gorm.DB, hasNext *bool, hasPrev *bool) (tx *gorm.DB) {

	return nil
}
