package repository

import (
	"WBTech_Level0/models"
	"errors"
	"github.com/spf13/viper"
	"sync"
)

var CacheContainsValueError = errors.New("error: the cache already has this order")
var CacheHasNoValue = errors.New("error: there is no such order uid in the cache")

// TODO поменять на REDIS, если будет время (или целесообразность)

type Cache struct {
	Data map[string]models.Order
	Mx   *sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{Data: make(map[string]models.Order, viper.GetInt("cache.initialSize")), Mx: &sync.RWMutex{}}
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
