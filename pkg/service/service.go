package service

import (
	"WBTech_Level0/models"
	"WBTech_Level0/pkg/repository"
)

type OrdersDB interface {
	GetOrderById(uid string) (models.Order, error)
	CreateOrder(order models.Order) error
}

type OrderCache interface {
	PutOrder(order models.Order) error
	GetOrder(uid string) (models.Order, error)
	GetAllOrders() ([]models.Order, error)
}

type Service struct {
	OrdersDB
	OrderCache
}

func NewService(repo *repository.Repository) *Service {
	return &Service{OrdersDB: repo.PostgresRepository, OrderCache: repo.CacheRepository}
}
