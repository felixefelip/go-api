package main

import (
	"go-api/internal/db"
	"go-api/internal/model"
	"go-api/internal/router"
)

func main() {
	if err := db.EnsureDatabase(); err != nil {
		panic(err)
	}

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	if err := dbConnection.AutoMigrate(&model.Product{}); err != nil {
		panic(err)
	}

	server := router.New(dbConnection)

	server.Run(":8000")
}
