package repository

import (
	"WBTech_Level0/models"
	"github.com/jmoiron/sqlx"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (op *OrderPostgres) GetOrderById(uid string) (models.Order, error) {
	return models.Order{}, nil
}

func (op *OrderPostgres) CreateOrder(order models.Order) (string, error) {
	return "", nil
}
