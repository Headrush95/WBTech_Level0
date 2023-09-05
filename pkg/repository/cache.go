package repository

import (
	"WBTech_Level0/models"
	"sync"
)

// TODO поменять на REDIS, если будет время (или целесообразность)

type Cache struct {
	Data map[string]models.Order
	Mx   *sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{Data: make(map[string]models.Order, 1000), Mx: &sync.RWMutex{}}
}

func (c *Cache) PutOrder(order models.Order) error {
	return nil
}

func (c *Cache) GetOrder(uid string) (models.Order, error) {
	return models.Order{}, nil
}
