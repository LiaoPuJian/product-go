package services

import (
	"product-go/models"
	"product-go/repositories"
)

type IProductService interface {
	GetProductById(int64) (*models.Product, error)
	GetAllProduct() ([]*models.Product, error)
	DeleteProductById(int64) bool
	InsertProduct(*models.Product) (int64, error)
	UpdateProduct(product *models.Product) error
}

type ProductService struct {
	Repository repositories.IProduct
}

func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{Repository: repository}
}

func (p *ProductService) GetProductById(id int64) (*models.Product, error) {
	return p.Repository.SelectByKey(id)
}

func (p *ProductService) GetAllProduct() ([]*models.Product, error) {
	return p.Repository.SelectAll()
}

func (p *ProductService) DeleteProductById(id int64) bool {
	return p.Repository.Delete(id)
}

func (p *ProductService) InsertProduct(product *models.Product) (int64, error) {
	return p.Repository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *models.Product) error {
	return p.Repository.Update(product)
}
