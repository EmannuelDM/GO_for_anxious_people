package cache

import (
	"sync"
)

// Cache is a thread-safe, in-memory key-value store.
//
// Key Encapsulation Lesson:
// Both 'mu' and 'items' are unexported.
// If 'items' or 'mu' were exported:
// 1. External callers might read/write 'items' without locking 'mu', causing race conditions or crashes.
// 2. External callers might lock 'mu' and forget to unlock it (deadlock).
// By encapsulating both behind exported methods, concurrency safety is guaranteed.
type Cache struct {
	mu    sync.RWMutex
	items map[string]string
}

// New constructs and initializes a new thread-safe Cache.
func New() *Cache {
	return &Cache{
		items: make(map[string]string),
	}
}

// Set stores a key-value pair under an exclusive write lock.
func (c *Cache) Set(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = val
}

// Get retrieves a value by key under a shared read lock.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.items[key]
	return val, ok
}

// Delete removes a key under an exclusive write lock.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Len returns the number of cached entries.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Keys returns a snapshot slice of all keys currently in the cache.
// Encapsulation note: Returns a brand-new slice so the caller cannot
// observe or interfere with concurrent cache updates.
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}
