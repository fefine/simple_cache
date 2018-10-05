package cache

import (
	"runtime"
	"sync"
	"time"
)

type Cache struct {
	cache *cache
	janitor *janitor
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache:&cache{
			items:make(map[string]*cacheItem),
		},
		janitor:&janitor{
			interval:interval,
			stop:make(chan bool),
		},
	}
	runJanitor(c)
	runtime.SetFinalizer(c, stopJanitor)
	return c
}

func (c *Cache) Set(key string, value interface{}, expire time.Duration) {
	c.cache.Set(key, value, expire)
}

func (c *Cache) Get(key string) (value interface{}, ok bool) {
	return c.cache.Get(key)
}

func runJanitor(c *Cache) {
	go c.janitor.run(c.cache)
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

type cache struct {
	mux sync.Mutex
	items map[string]*cacheItem
}

func (c *cache) Set(key string, value interface{}, expire time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.items[key] = newCacheItem(value, expire)
}

func (c *cache) Get(key string) (interface{}, bool) {
	if res, ok := c.items[key]; ok {
		// may be conflict with set, but doesn't matter
		res.UpdateGetTime()
		return res.value, true
	}
	return nil, false
}

func (c *cache) clearExpired() {
	//c.mux.Lock()
	for key, item := range c.items {
		if item.Expired() {
			// double check
			c.mux.Lock()
			if item.Expired() {
				delete(c.items, key)
			}
			c.mux.Unlock()
		}
	}
}

type cacheItem struct {
	value interface{}
	expireTime time.Duration
	lastGetTime time.Time
}

func newCacheItem(value interface{}, expireTime time.Duration) *cacheItem {
	return &cacheItem{
		expireTime:expireTime,
		lastGetTime:time.Now(),
		value:value,
	}
}

func (item *cacheItem) Expired() bool {
	if item.expireTime == 0 {
		return false
	}
	return item.lastGetTime.Add(item.expireTime).Before(time.Now())
}

func (item *cacheItem) UpdateGetTime() {
	item.lastGetTime = time.Now()
}

type janitor struct {
	interval time.Duration
	stop chan bool
}

func (jan *janitor) run(c *cache) {
	ticker := time.NewTicker(jan.interval)
	for {
		select {
		case <-ticker.C:
			c.clearExpired()
		case <- jan.stop:
			return
		}
	}
}
