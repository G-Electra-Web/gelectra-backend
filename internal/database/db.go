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
	db.AutoMigrate(&User{}, &Member{}, &CoreMember{}, &Staff{}, &Staff{}, &Event{}, &Attendee{}, &GalleryImage{}, &ProjectRequest{}, &ProjectMember{}, &Notice{}, &Admin{})
	logger.Info("Database Migrated")

}
