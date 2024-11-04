package cache

import (
	"testing"
	"time"
)

func TestInMemCacheSetAndGet(t *testing.T) {
	cache := NewInMemCache(1*time.Second, EvictNO)
	defer cache.Stop()

	key := "testKey"
	value := []byte("testValue")
	cache.Set(key, value, 5*time.Minute)

	storedValue, found := cache.Get(key)
	if !found {
		t.Errorf("Expected to find value for key %s, but it was not found", key)
	}
	if string(storedValue) != "testValue" {
		t.Errorf("Expected value %s, got %s", "testValue", string(storedValue))
	}
}

func TestInMemCacheTTLExpiration(t *testing.T) {
	cache := NewInMemCache(1*time.Minute, EvictNO)
	defer cache.Stop()

	key := "testKey"
	value := []byte("testValue")
	cache.Set(key, value, 1*time.Second)

	_, found := cache.Get(key)
	if !found {
		t.Errorf("Expected to find value for key %s, but it was not found", key)
	}

	time.Sleep(2 * time.Second)

	_, found = cache.Get(key)
	if found {
		t.Errorf("Expected key %s to be expired, but it was found", key)
	}
}

func TestInMemCacheLRUEviction(t *testing.T) {
	cache := NewInMemCache(1*time.Minute, EvictLRU, WithEvictInterval(10*time.Second))
	defer cache.Stop()

	cache.Set("key1", []byte("value1"), 5*time.Minute)
	cache.Set("key2", []byte("value2"), 5*time.Minute)
	cache.Get("key1") // Access key1 to update its access time

	cache.evictLRU()

	_, ok := cache.items.Load("key1")
	if !ok {
		t.Error("Expected key1 to be found after LRU eviction, but it was not found")
	}
}

func TestInMemCacheCleanup(t *testing.T) {
	cache := NewInMemCache(500*time.Millisecond, EvictNO)
	defer cache.Stop()

	cache.Set("key1", []byte("value1"), 1*time.Second)

	time.Sleep(2 * time.Second)

	_, found := cache.Get("key1")
	if found {
		t.Error("Expected key1 to be cleaned up due to TTL expiration, but it was found")
	}
}

func TestInMemCacheStop(t *testing.T) {
	cache := NewInMemCache(1*time.Minute, EvictNO)
	cache.Stop()

	cache.Set("key1", []byte("value1"), 5*time.Second)
	time.Sleep(1 * time.Millisecond)

	_, found := cache.Get("key1")
	if !found {
		t.Error("Expected key1 to be accessible after cache stop, but it was not found")
	}
}
