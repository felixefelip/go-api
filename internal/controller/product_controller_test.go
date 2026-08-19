package controller

import (
	"encoding/json"
	"net/http"
	"testing"

	"go-api/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProductRetorna201ComOProdutoCriado(t *testing.T) {
	server := newServer(t)

	response := post(t, server, "/product", `{"name":"Camiseta","price":30.99,"stock":12}`)

	require.Equal(t, http.StatusCreated, response.Code)

	var criado model.Product
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &criado))

	assert.NotZero(t, criado.ID, "a resposta deveria trazer o id gerado")
	assert.Equal(t, "Camiseta", criado.Name)
	assert.Equal(t, 30.99, criado.Price)
	assert.Equal(t, 12, criado.Stock)
}

func TestCreateProductSemStockRetornaZero(t *testing.T) {
	server := newServer(t)

	response := post(t, server, "/product", `{"name":"Camiseta","price":30.99}`)

	require.Equal(t, http.StatusCreated, response.Code)

	var criado model.Product
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &criado))

	assert.Zero(t, criado.Stock)
}

func TestCreateProductPersisteNoBanco(t *testing.T) {
	server := newServer(t)

	response := post(t, server, "/product", `{"name":"Camiseta","price":30.99,"stock":12}`)
	require.Equal(t, http.StatusCreated, response.Code)

	var salvos []model.Product
	require.NoError(t, testConnection.Find(&salvos).Error)

	require.Len(t, salvos, 1)
	assert.Equal(t, "Camiseta", salvos[0].Name)
	assert.Equal(t, 30.99, salvos[0].Price)
	assert.Equal(t, 12, salvos[0].Stock)
}

func TestCreateProductComJSONInvalidoRetorna400(t *testing.T) {
	server := newServer(t)

	response := post(t, server, "/product", `{"name":`)

	assert.Equal(t, http.StatusBadRequest, response.Code)

	var salvos []model.Product
	require.NoError(t, testConnection.Find(&salvos).Error)
	assert.Empty(t, salvos, "nada deveria ter sido gravado")
}

func TestCreateProductComPrecoDeTipoErradoRetorna400(t *testing.T) {
	server := newServer(t)

	response := post(t, server, "/product", `{"name":"Camiseta","price":"muito caro"}`)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCreateProductComStockDeTipoErradoRetorna400(t *testing.T) {
	server := newServer(t)

	response := post(t, server, "/product", `{"name":"Camiseta","price":30.99,"stock":"muitos"}`)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestGetProductsRetorna200ComOsProdutos(t *testing.T) {
	server := newServer(t)

	require.Equal(t, http.StatusCreated, post(t, server, "/product", `{"name":"Camiseta","price":30.99,"stock":12}`).Code)
	require.Equal(t, http.StatusCreated, post(t, server, "/product", `{"name":"Calca Jeans","price":89.99,"stock":3}`).Code)

	response := get(t, server, "/products")

	require.Equal(t, http.StatusOK, response.Code)

	var produtos []model.Product
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &produtos))

	require.Len(t, produtos, 2)
	assert.Equal(t, "Camiseta", produtos[0].Name)
	assert.Equal(t, 30.99, produtos[0].Price)
	assert.Equal(t, 12, produtos[0].Stock)
	assert.Equal(t, "Calca Jeans", produtos[1].Name)
	assert.Equal(t, 89.99, produtos[1].Price)
	assert.Equal(t, 3, produtos[1].Stock)
}

func TestGetProductsQuandoNaoHaProdutos(t *testing.T) {
	server := newServer(t)

	response := get(t, server, "/products")

	require.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, `[]`, response.Body.String())
}

// TestGetProductsQuandoOBancoFalhaRetorna500 cobre o caminho de erro do handler:
// a resposta precisa ser 500 e ter um corpo JSON unico. Antes do `return` na
// linha do erro, o handler seguia adiante e escrevia um segundo corpo, gerando
// um payload concatenado e impossivel de parsear.
func TestGetProductsQuandoOBancoFalhaRetorna500(t *testing.T) {
	server := newServerComBancoIndisponivel(t)

	response := get(t, server, "/products")

	require.Equal(t, http.StatusInternalServerError, response.Code)

	var corpo any
	assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &corpo),
		"o corpo precisa ser um JSON unico e valido, nao dois concatenados")
}
