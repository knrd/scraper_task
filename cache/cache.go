package cache

import "sync"

type cacheStorage struct {
	mu   sync.RWMutex
	data map[string]bool
}

// Try do add value to cache in non-blocking manner. If value alread exists return false, else return true
func (c *cacheStorage) Add(value string) bool {
	c.mu.RLock()
	if _, ok := c.data[value]; ok {
		c.mu.RUnlock()
		return false
	}
	c.mu.RUnlock()

	c.mu.Lock()
	c.data[value] = true
	c.mu.Unlock()
	return true
}

func (c *cacheStorage) Clear() {
	clear(c.data)
}

func New() *cacheStorage {
	return &cacheStorage{data: map[string]bool{}}
}
