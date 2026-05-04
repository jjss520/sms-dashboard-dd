package database

import (
	"log"
	"sms-dashboard/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Optimize SQLite performance
	DB.Exec("PRAGMA journal_mode=WAL")          // Write-Ahead Logging
	DB.Exec("PRAGMA synchronous=NORMAL")         // Balance safety and speed
	DB.Exec("PRAGMA cache_size=-64000")          // 64MB cache
	DB.Exec("PRAGMA temp_store=MEMORY")          // Store temp tables in memory
	DB.Exec("PRAGMA mmap_size=268435456")        // 256MB memory-mapped I/O
	DB.Exec("PRAGMA page_size=4096")             // Optimal page size

	// Migrate the schema
	err = DB.AutoMigrate(&model.SMS{}, &model.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
