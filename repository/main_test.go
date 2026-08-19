package repository

import (
	"fmt"
	"os"
	"testing"

	"go-api/db"
	"go-api/model"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// testConnection e a conexao compartilhada por todos os testes do pacote.
var testConnection *gorm.DB

// TestMain prepara o banco do ambiente de teste uma unica vez para todo o pacote.
// Aqui nao ha *testing.T, entao o testify nao se aplica: uma falha de setup
// acontece antes de existir qualquer teste, e a saida correta e encerrar com
// codigo de erro.
func TestMain(m *testing.M) {
	if err := setupTestDB(); err != nil {
		fmt.Fprintln(os.Stderr, "setup do banco de teste:", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// setupTestDB garante que o banco de teste existe, esta migrado e conectado.
func setupTestDB() error {
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "test")
	}

	// Trava de seguranca: os testes truncam tabelas, entao so podem rodar
	// contra o banco de teste.
	if db.Env() != "test" {
		return fmt.Errorf("recusando rodar com GO_ENV=%q; use GO_ENV=test", db.Env())
	}

	if err := db.EnsureDatabase(); err != nil {
		return fmt.Errorf("criando o banco: %w", err)
	}

	connection, err := db.ConnectDB()
	if err != nil {
		return fmt.Errorf("conectando no banco: %w", err)
	}

	if err := connection.AutoMigrate(&model.Product{}); err != nil {
		return fmt.Errorf("migrando o banco: %w", err)
	}

	testConnection = connection

	return nil
}

// newRepository devolve um repositorio sobre uma tabela product vazia.
func newRepository(t *testing.T) ProductRepository {
	t.Helper()

	err := testConnection.Exec("TRUNCATE TABLE product RESTART IDENTITY CASCADE").Error
	require.NoError(t, err, "limpando a tabela product")

	return NewProductRepository(testConnection)
}
