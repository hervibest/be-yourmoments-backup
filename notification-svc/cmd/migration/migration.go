package migration

import (
	"database/sql"
	"log"
	"time"

	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/utils"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Run() {
	time.Sleep(10 * time.Second)
	dbURL := utils.GetEnv("NOTIFICATION_DB_URL")
	if dbURL == "" {
		log.Fatal("NOTIFICATION_DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	migrationPath := "/db/migrations"

	if err := goose.Up(db, migrationPath); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("✅ Migration completed successfully")
}
