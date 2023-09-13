package repository

import (
	"WBTech_Level0/models"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

var CacheContainsValueError = errors.New("error: the cache already has this order")
var CacheHasNoValue = errors.New("error: there is no such order uid in the cache")
var CacheIsEmpty = errors.New("error: the cache is empty")

// TODO поменять на REDIS, если будет время (или целесообразность)

type Cache struct {
	Data map[string]models.Order
	Mx   *sync.RWMutex
}

func NewCache(ps *OrderPostgres) *Cache {
	data, err := ps.getAllOrders()
	if err != nil {
		logrus.Panicf("error occurred during cache initialization: %v", err)
	}

	if len(data) == 0 {
		return &Cache{Data: make(map[string]models.Order, viper.GetInt("cache.initialSize")), Mx: &sync.RWMutex{}}
	}
	return &Cache{Data: data, Mx: &sync.RWMutex{}}
}

func (c *Cache) PutOrder(order models.Order) error {
	c.Mx.Lock()
	defer c.Mx.Unlock()
	if _, ok := c.Data[order.Uid]; ok {
		return CacheContainsValueError
	}

	c.Data[order.Uid] = order
	return nil
}

func (c *Cache) GetOrder(uid string) (models.Order, error) {
	c.Mx.RLock()
	defer c.Mx.RUnlock()
	if _, ok := c.Data[uid]; !ok {

		return models.Order{}, CacheHasNoValue
	}
	return c.Data[uid], nil
}

func (c *Cache) GetAllOrders() ([]models.Order, error) {
	c.Mx.RLock()
	defer c.Mx.RUnlock()
	if len(c.Data) == 0 {
		return nil, CacheIsEmpty
	}

	orders := make([]models.Order, 0, len(c.Data))
	for _, order := range c.Data {
		orders = append(orders, order)
	}
	return orders, nil
}
