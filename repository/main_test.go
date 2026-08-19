package repository

import (
	"fmt"
	"os"
	"testing"

	"go-api/internal/testdb"

	"gorm.io/gorm"
)

// testConnection e a conexao compartilhada por todos os testes do pacote.
var testConnection *gorm.DB

// TestMain prepara o banco do ambiente de teste uma unica vez para todo o pacote.
// Aqui nao ha *testing.T, entao o testify nao se aplica: uma falha de setup
// acontece antes de existir qualquer teste, e a saida correta e encerrar com
// codigo de erro.
func TestMain(m *testing.M) {
	connection, err := testdb.Setup()
	if err != nil {
		fmt.Fprintln(os.Stderr, "setup do banco de teste:", err)
		os.Exit(1)
	}

	testConnection = connection

	os.Exit(m.Run())
}

// newRepository devolve um repositorio sobre uma tabela product vazia.
func newRepository(t *testing.T) ProductRepository {
	t.Helper()

	testdb.Reset(t, testConnection)

	return NewProductRepository(testConnection)
}
