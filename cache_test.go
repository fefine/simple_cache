package cache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	interval = 1 * time.Second
	keyPattern = "key-%d"
	valuePattern = "value-%d"
)

func key(i int) string {
	return fmt.Sprintf(keyPattern, i)
}

func value(i int) string {
	return fmt.Sprintf(valuePattern, i)
}

func TestGetAndSet(t *testing.T)  {
	cache := NewCache(interval)

	for i := 0; i < 1000; i++ {
		cache.Set(key(i), value(i), interval)
	}

	for i := 0; i < 1000; i++ {
		v, ok := cache.Get(key(i))
		if assert.True(t, ok) {
			val := v.(string)
			assert.Equal(t, value(i), val)
		}
	}
}

func TestGetExpired(t *testing.T) {
	cache := NewCache(interval)

	cache.Set(key(10), value(20), interval)

	time.Sleep(2 * time.Second)

	_, ok := cache.Get(key(10))
	assert.False(t, ok)
}

func TestGetRefresh(t *testing.T) {
	cache := NewCache(interval)

	cache.Set(key(10), value(20), interval)

	// not expired
	time.Sleep(500 * time.Millisecond)

	// refresh time
	v, ok := cache.Get(key(10))
	if assert.True(t, ok) {
		val := v.(string)
		assert.Equal(t, value(20), val)
	}

	// not expired
	time.Sleep(800 * time.Millisecond)
	v, ok = cache.Get(key(10))
	if assert.True(t, ok) {
		val := v.(string)
		assert.Equal(t, value(20), val)
	}

	// expired
	time.Sleep(2 * time.Second)
	_, ok = cache.Get(key(10))
	assert.False(t, ok)
}
