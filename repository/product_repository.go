package repository

import (
	"go-api/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	connection *gorm.DB
}

func NewProductRepository(connection *gorm.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	var productList []model.Product

	err := pr.connection.Find(&productList).Error
	if err != nil {
		return []model.Product{}, err
	}

	return productList, nil
}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {
	err := pr.connection.Create(&product).Error
	if err != nil {
		return 0, err
	}

	return product.ID, nil
}
