package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	host     = "go_db"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "go_api"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to " + dbname)

	return db, nil
}
