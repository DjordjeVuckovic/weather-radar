package cache

import (
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type EvictPolicy string

const (
	EvictNO                EvictPolicy = "EvictNO"
	EvictLRU               EvictPolicy = "EvictLRU"
	defaultEvictInterval               = 1 * time.Minute
	defaultMemoryThreshold             = 0.8 // 80% of memory usage
)

type InMemItem struct {
	Value      []byte
	TTL        int64
	AccessTime int64
}

type InMemCache struct {
	items sync.Map

	stopCh  chan struct{}
	stopped int32

	evictPolicy     EvictPolicy
	evictInterval   time.Duration
	memoryThreshold float64
}

func NewInMemCache(cleanupInterval time.Duration, evictPolicy EvictPolicy) *InMemCache {
	c := &InMemCache{
		stopCh:          make(chan struct{}),
		evictPolicy:     evictPolicy,
		evictInterval:   defaultEvictInterval,
		memoryThreshold: defaultMemoryThreshold,
	}
	go c.startCleanupLoop(cleanupInterval)
	if evictPolicy != EvictNO {
		go c.startEvictionLoop(c.evictInterval)
	}
	return c
}

func (c *InMemCache) Set(key string, value []byte, ttl time.Duration) {
	expiration := time.Now().Add(ttl).UnixNano()

	c.items.Store(key, InMemItem{
		Value:      value,
		TTL:        expiration,
		AccessTime: time.Now().UnixNano(),
	})

}

func (c *InMemCache) Get(key string) ([]byte, bool) {
	itemAny, found := c.items.Load(key)
	if !found {
		return nil, false
	}

	item := itemAny.(InMemItem)
	if time.Now().UnixNano() > item.TTL {
		c.Delete(key)
		return nil, false
	}

	if c.evictPolicy == EvictLRU {
		item.AccessTime = time.Now().UnixNano()
		c.items.Store(key, item)
	}

	return item.Value, true
}

func (c *InMemCache) Delete(key string) {
	c.items.Delete(key)
}

func (c *InMemCache) Stop() {
	if atomic.CompareAndSwapInt32(&c.stopped, 0, 1) {
		close(c.stopCh)
	}
}

func (c *InMemCache) startCleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCh:
			return
		}
	}
}

func (c *InMemCache) cleanup() {
	slog.Debug("Cache cleaning up started ...")

	now := time.Now().UnixNano()
	c.items.Range(func(key, value interface{}) bool {
		item := value.(InMemItem)
		if item.TTL > 0 && now > item.TTL {
			slog.Debug("Deleting item due to TTL expiration", slog.String("cache_key", key.(string)))
			c.Delete(key.(string))
		}
		return true
	})
}

func (c *InMemCache) startEvictionLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if c.memoryUsageExceedsThreshold() {
				c.evict()
			} else {
				slog.Debug("Memory usage is below threshold. No eviction needed.")
			}
		case <-c.stopCh:
			return
		}
	}
}

func (c *InMemCache) evict() {
	switch c.evictPolicy {
	case EvictNO:
		slog.Info("Eviction policy set to EvictNO. No eviction will be performed.")
	case EvictLRU:
		c.evictLRU()
	default:
		slog.Info("No eviction policy set or policy not recognized.")
	}
}

func (c *InMemCache) evictLRU() {
	slog.Debug("Started evicting LRU...")
	oldestKey := ""
	oldestAccessTime := time.Now().UnixNano()

	c.items.Range(func(key, value interface{}) bool {
		item, ok := value.(InMemItem)
		if !ok {
			return true
		}
		if item.AccessTime < oldestAccessTime {
			oldestAccessTime = item.AccessTime
			oldestKey = key.(string)
		}
		return true
	})

	if oldestKey != "" {
		slog.Info("Evicting item with key due to LRU policy", slog.String("cache_key", oldestKey))
		c.Delete(oldestKey)
	} else {
		slog.Info("LRU eviction found no item to evict")
	}
}

func (c *InMemCache) memoryUsageExceedsThreshold() bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	allocMB := float64(m.Alloc) / (1024 * 1024)

	return allocMB > c.memoryThreshold*float64(m.Sys)
}
