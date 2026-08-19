package repository

import (
	"testing"

	"go-api/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	repository := newRepository(t)

	id, err := repository.CreateProduct(model.Product{Name: "Camiseta", Price: 30.99})
	require.NoError(t, err)
	assert.NotZero(t, id, "o banco deveria ter gerado um id")

	var saved model.Product
	require.NoError(t, testConnection.First(&saved, id).Error)

	assert.Equal(t, "Camiseta", saved.Name)
	assert.Equal(t, 30.99, saved.Price)
}

func TestGetProductsQuandoVazio(t *testing.T) {
	repository := newRepository(t)

	products, err := repository.GetProducts()

	require.NoError(t, err)
	assert.Empty(t, products)
}

func TestGetProductsRetornaOsProdutosCriados(t *testing.T) {
	repository := newRepository(t)

	criados := []model.Product{
		{Name: "Camiseta", Price: 30.99},
		{Name: "Calca Jeans", Price: 89.99},
	}

	for _, product := range criados {
		_, err := repository.CreateProduct(product)
		require.NoError(t, err)
	}

	products, err := repository.GetProducts()
	require.NoError(t, err)
	require.Len(t, products, len(criados))

	for i, esperado := range criados {
		assert.Equal(t, esperado.Name, products[i].Name)
		assert.Equal(t, esperado.Price, products[i].Price)
	}
}

func TestIsolamentoEntreTestes(t *testing.T) {
	repository := newRepository(t)

	products, err := repository.GetProducts()

	require.NoError(t, err)
	assert.Empty(t, products, "estado vazou de outro teste")
}
