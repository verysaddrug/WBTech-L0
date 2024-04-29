package cache

import (
	"sync"
	"wbTechL0/internal/storage/pgsql"
	"wbTechL0/models"
)

type Cache struct {
	mu    sync.RWMutex
	cache sync.Map
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) AddOrder(orderID int, order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Store(orderID, order)
}

func (c *Cache) GetOrder(orderID int) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.cache.Load(orderID)
	if !ok {
		return models.Order{}, false
	}

	return val.(models.Order), true
}

func (c *Cache) GetAllOrders() []models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var orders []models.Order

	c.cache.Range(func(_, value interface{}) bool {
		orders = append(orders, value.(models.Order))
		return true
	})

	return orders
}

func (cache *Cache) AddAllOrdersToCache(d *pgsql.Database) error {
	orders, err := d.GetAllOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		cache.AddOrder(order.OrderID, order)
	}

	return nil
}
