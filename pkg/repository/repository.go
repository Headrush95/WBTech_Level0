package repository

import (
	"WBTech_Level0/models"
	"github.com/jmoiron/sqlx"
)

type PostgresRepository interface {
	CreateOrder(order models.Order) error
	GetOrderById(uid string) (models.Order, error)
	CloseStatements() error
}

type CacheRepository interface {
	PutOrder(order models.Order) error
	GetOrder(uid string) (models.Order, error)
	GetAllOrders() ([]models.Order, error)
}

type Repository struct {
	PostgresRepository
	CacheRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	orderPostgres := NewOrderPostgres(db)
	cache := NewCache(orderPostgres)

	return &Repository{
		PostgresRepository: orderPostgres,
		CacheRepository:    cache,
	}
}
