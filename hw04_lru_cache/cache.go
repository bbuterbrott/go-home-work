package hw04_lru_cache //nolint:golint,stylecheck

// Cache is interface for cache implementations which use string as a key.
type Cache interface {
	// Adds cache value by key. Returns if key was already present in cache
	Set(key string, value interface{}) bool
	// Gets cache value by key. Also returns if key was in cache
	Get(key string) (interface{}, bool)
	// Clears cache
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	cache    map[string]*listItem
}

type cacheItem struct {
	Key   string
	Value interface{}
}

func (c lruCache) Set(key string, value interface{}) bool {
	li := c.cache[key]
	if li != nil {
		li.Value = cacheItem{key, value}
		c.queue.MoveToFront(li)

		return true
	}

	nli := c.queue.PushFront(cacheItem{key, value})
	if c.queue.Len() > c.capacity {
		oli := c.queue.Back()
		oci := oli.Value.(cacheItem)
		delete(c.cache, oci.Key)
		c.queue.Remove(oli)
	}
	c.cache[key] = nli
	return false
}

func (c lruCache) Get(key string) (interface{}, bool) {
	li := c.cache[key]
	if li == nil {
		return nil, false
	}

	c.queue.MoveToFront(li)
	ci := li.Value.(cacheItem)
	return ci.Value, true
}

func (c lruCache) Clear() {
	c.queue.Clear()
	c.cache = make(map[string]*listItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{capacity: capacity, queue: NewList(), cache: make(map[string]*listItem)}
}
