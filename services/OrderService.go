package services

import (
	"product-go/models"
	"product-go/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*models.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*models.Order) error
	InsertOrder(*models.Order) (int64, error)
	GetAllOrder() ([]*models.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
}

type OrderService struct {
	repository repositories.IOrder
}

func NewOrderService(r repositories.IOrder) *OrderService {
	return &OrderService{repository: r}
}

func (o *OrderService) GetOrderByID(id int64) (*models.Order, error) {
	return o.repository.SelectByKey(id)
}

func (o *OrderService) DeleteOrderByID(id int64) bool {
	return o.repository.Delete(id)
}

func (o *OrderService) UpdateOrder(order *models.Order) error {
	return o.repository.Update(order)
}

func (o *OrderService) InsertOrder(order *models.Order) (int64, error) {
	return o.repository.Insert(order)
}

func (o *OrderService) GetAllOrder() ([]*models.Order, error) {
	return o.repository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.repository.SelectAllWithInfo()
}
