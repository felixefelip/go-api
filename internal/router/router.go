// Package router concentra o wiring das camadas e o registro das rotas HTTP.
package router

import (
	"net/http"

	"go-api/internal/controller"
	"go-api/internal/repository"
	"go-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// New devolve o servidor da aplicacao, ja com as rotas registradas.
func New(connection *gorm.DB) *gin.Engine {
	server := gin.Default()

	Register(server, connection)

	return server
}

// Register monta o stack de produtos e registra as rotas no servidor recebido.
// Os testes chamam esta mesma funcao, entao exercitam exatamente as rotas que
// a aplicacao expoe — elas nao tem como divergir.
func Register(server *gin.Engine, connection *gorm.DB) {
	productRepository := repository.NewProductRepository(connection)
	productUsecase := usecase.NewProductUsecase(productRepository)
	productController := controller.NewProductController(productUsecase)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	server.GET("/products", productController.GetProducts)
	server.GET("/products/:id", productController.GetProductByID)
	server.POST("/products", productController.CreateProduct)
}
