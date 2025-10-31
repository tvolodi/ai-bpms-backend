package main

import (
	"flag"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tvolodi/ai-bpms-backend/shared/common/config"
	"github.com/tvolodi/ai-bpms-backend/shared/database/migration"
)

func main() {
	var (
		rollback = flag.Bool("rollback", false, "Rollback the last migration")
		help     = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create migrator
	migrator := migration.NewMigrator(db)

	if *rollback {
		// Rollback last migration
		log.Println("Rolling back last migration...")
		if err := migrator.Rollback(); err != nil {
			log.Fatalf("Migration rollback failed: %v", err)
		}
		log.Println("Migration rollback completed successfully")
	} else {
		// Run migrations
		log.Println("Running database migrations...")
		if err := migrator.Run(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("All migrations completed successfully")
	}
}

func connectDB(cfg *config.Config) (*gorm.DB, error) {
	// Configure GORM logger
	logLevel := logger.Silent
	if cfg.Logging.Level == "debug" {
		logLevel = logger.Info
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.MaxLifetime)

	return db, nil
}
