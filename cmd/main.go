package main

import (
	"go-api/internal/controller"
	"go-api/internal/db"
	"go-api/internal/model"
	"go-api/internal/repository"
	"go-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

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

	ProductRepository := repository.NewProductRepository(dbConnection)

	ProductUsecase := usecase.NewProductUsecase(ProductRepository)

	ProductController := controller.NewProductController(ProductUsecase)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.GET("/products", ProductController.GetProducts)
	server.POST("/product", ProductController.CreateProduct)

	server.Run(":8000")
}
