package repository

import (
	"fmt"
	"os"
	"testing"

	"go-api/db"
	"go-api/model"

	"gorm.io/gorm"
)

var testConnection *gorm.DB

// TestMain prepara o banco do ambiente de teste uma unica vez para todo o pacote.
func TestMain(m *testing.M) {
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "test")
	}

	// Trava de seguranca: os testes truncam tabelas, entao so podem rodar
	// contra o banco de teste.
	if db.Env() != "test" {
		fmt.Fprintf(os.Stderr, "recusando rodar testes com GO_ENV=%q; use GO_ENV=test\n", db.Env())
		os.Exit(1)
	}

	if err := db.EnsureDatabase(); err != nil {
		fmt.Fprintln(os.Stderr, "erro criando banco de teste:", err)
		os.Exit(1)
	}

	connection, err := db.ConnectDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro conectando no banco de teste:", err)
		os.Exit(1)
	}

	if err := connection.AutoMigrate(&model.Product{}); err != nil {
		fmt.Fprintln(os.Stderr, "erro migrando banco de teste:", err)
		os.Exit(1)
	}

	testConnection = connection

	os.Exit(m.Run())
}

// newRepository devolve um repositorio sobre uma tabela product vazia.
func newRepository(t *testing.T) ProductRepository {
	t.Helper()

	if err := testConnection.Exec("TRUNCATE TABLE product RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatalf("limpando product: %v", err)
	}

	return NewProductRepository(testConnection)
}

func TestCreateProduct(t *testing.T) {
	repository := newRepository(t)

	id, err := repository.CreateProduct(model.Product{Name: "Camiseta", Price: 30.99})
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}

	if id == 0 {
		t.Error("esperava um id gerado pelo banco, veio 0")
	}

	var saved model.Product
	if err := testConnection.First(&saved, id).Error; err != nil {
		t.Fatalf("relendo o produto salvo: %v", err)
	}

	if saved.Name != "Camiseta" {
		t.Errorf("Name = %q, esperado %q", saved.Name, "Camiseta")
	}

	if saved.Price != 30.99 {
		t.Errorf("Price = %v, esperado %v", saved.Price, 30.99)
	}
}

func TestGetProductsQuandoVazio(t *testing.T) {
	repository := newRepository(t)

	products, err := repository.GetProducts()
	if err != nil {
		t.Fatalf("GetProducts: %v", err)
	}

	if len(products) != 0 {
		t.Errorf("esperava lista vazia, veio %d produto(s)", len(products))
	}
}

func TestGetProductsRetornaOsProdutosCriados(t *testing.T) {
	repository := newRepository(t)

	criados := []model.Product{
		{Name: "Camiseta", Price: 30.99},
		{Name: "Calca Jeans", Price: 89.99},
	}

	for _, product := range criados {
		if _, err := repository.CreateProduct(product); err != nil {
			t.Fatalf("CreateProduct(%q): %v", product.Name, err)
		}
	}

	products, err := repository.GetProducts()
	if err != nil {
		t.Fatalf("GetProducts: %v", err)
	}

	if len(products) != len(criados) {
		t.Fatalf("esperava %d produtos, veio %d", len(criados), len(products))
	}

	for i, esperado := range criados {
		if products[i].Name != esperado.Name {
			t.Errorf("produto %d: Name = %q, esperado %q", i, products[i].Name, esperado.Name)
		}

		if products[i].Price != esperado.Price {
			t.Errorf("produto %d: Price = %v, esperado %v", i, products[i].Price, esperado.Price)
		}
	}
}

// TestIsolamentoEntreTestes garante que a limpeza de estado esta funcionando:
// o produto criado no teste anterior nao pode vazar para este.
func TestIsolamentoEntreTestes(t *testing.T) {
	repository := newRepository(t)

	products, err := repository.GetProducts()
	if err != nil {
		t.Fatalf("GetProducts: %v", err)
	}

	if len(products) != 0 {
		t.Errorf("estado vazou de outro teste: %d produto(s) na tabela", len(products))
	}
}
