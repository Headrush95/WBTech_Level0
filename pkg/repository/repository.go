package repository

import (
	"WBTech_Level0/models"
	"github.com/jmoiron/sqlx"
)

type PostgresRepository interface {
	GetOrderById(uid string) (models.Order, error)
	CreateOrder(order models.Order) error // возвращаем uid
}

type CacheRepository interface {
	PutOrder(order models.Order) error
	GetOrder(uid string) (models.Order, error)
}

type Repository struct {
	PostgresRepository
	CacheRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		PostgresRepository: NewOrderPostgres(db),
		CacheRepository:    NewCache(),
	}
}
