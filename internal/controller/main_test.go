package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"go-api/internal/db"
	"go-api/internal/repository"
	"go-api/internal/testdb"
	"go-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// testConnection e a conexao compartilhada por todos os testes do pacote.
var testConnection *gorm.DB

func TestMain(m *testing.M) {
	connection, err := testdb.Setup()
	if err != nil {
		fmt.Fprintln(os.Stderr, "setup do banco de teste:", err)
		os.Exit(1)
	}

	testConnection = connection
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

// newServer monta o stack real (controller -> usecase -> repository -> Postgres)
// sobre uma tabela product vazia, com as mesmas rotas registradas no main.go.
func newServer(t *testing.T) *gin.Engine {
	t.Helper()

	testdb.Reset(t, testConnection)

	productRepository := repository.NewProductRepository(testConnection)
	productUsecase := usecase.NewProductUsecase(productRepository)
	productController := NewProductController(productUsecase)

	server := gin.New()
	server.GET("/products", productController.GetProducts)
	server.GET("/products/:id", productController.GetProductByID)
	server.POST("/products", productController.CreateProduct)

	return server
}

// newServerComBancoIndisponivel monta o mesmo stack sobre uma conexao propria
// que ja foi fechada, para exercitar o caminho de erro dos handlers. A conexao
// compartilhada pelos demais testes nao e afetada.
func newServerComBancoIndisponivel(t *testing.T) *gin.Engine {
	t.Helper()

	connection, err := db.ConnectDB()
	require.NoError(t, err)

	sqlDB, err := connection.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	productRepository := repository.NewProductRepository(connection)
	productUsecase := usecase.NewProductUsecase(productRepository)
	productController := NewProductController(productUsecase)

	server := gin.New()
	server.GET("/products", productController.GetProducts)
	server.GET("/products/:id", productController.GetProductByID)
	server.POST("/products", productController.CreateProduct)

	return server
}

// do dispara uma requisicao contra o servidor e devolve a resposta gravada.
func do(t *testing.T, server *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}

	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	return response
}

func get(t *testing.T, server *gin.Engine, path string) *httptest.ResponseRecorder {
	t.Helper()

	return do(t, server, http.MethodGet, path, "")
}

func post(t *testing.T, server *gin.Engine, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	return do(t, server, http.MethodPost, path, body)
}
