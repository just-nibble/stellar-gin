package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Model interface {
	Save()
}

var DB *gorm.DB

func Connect() {
	var err error
	db_user := os.Getenv("DB_USER")
	db_host := os.Getenv("DB_HOST")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	db_password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", db_host, db_user, db_password, db_name, db_port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		QueryFields:                              true,
		PrepareStmt:                              true,
	})

	if err != nil {
		panic(err)
	} else {
		log.Println("Successfully connected to the database")
	}
}
