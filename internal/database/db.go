package database

import (
	"fmt"
	"log/slog"

	"github.com/knadh/koanf"
	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(k *koanf.Koanf, logger *slog.Logger) {
	con := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", k.String("db.host"), k.String("db.user"), k.String("db.password"), k.String("db.dbname"), k.String("db.port"))
	db, err := gorm.Open(postgres.Open(con), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database")
	}

	DB = db

	logger.Info("Connection Opened to Database")
	db.AutoMigrate(
		&User{},           // User must be created first
		&Admin{},          // Admin depends on User
		&Member{},         // Member depends on User
		&CoreMember{},     // CoreMember depends on User
		&Staff{},          // Staff depends on User
		&Event{},          // Event does not depend on any other table
		&GalleryImage{},   // GalleryImage depends on Event and User
		&Attendee{},       // Attendee depends on User and Event
		&Project{},        // Project does not depend on any other table
		&ProjectMember{},  // ProjectMember depends on User and Project
		&ProjectRequest{}, // ProjectRequest depends on User and Project
		&Notice{},         // Notice does not depend on any other table
	)
	logger.Info("Database Migrated")

}
