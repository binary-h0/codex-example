package main

import (
   "errors"
   "log"
   "os"

   "gorm.io/driver/postgres"
   "gorm.io/gorm"

   "codex-pipeline/internal/post"
   "codex-pipeline/internal/server"
)

func setupDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// Routes moved to internal/server.SetupRouter

// main starts the application, connecting to the database and serving HTTP routes
func main() {
   db, err := setupDatabase()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

   if err := db.AutoMigrate(&post.Post{}); err != nil {
       log.Fatalf("failed to migrate: %v", err)
   }

   r := server.SetupRouter(db)

	log.Println("Listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
