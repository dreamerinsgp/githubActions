package cache

import (
	"sync"
	"time"
)

const (
	itemCountCacheThreshold = 100000 // 超过此数量时启用缓存
	itemCountCacheTTL       = 60     // 缓存 TTL，秒
)

// ItemCountCache 用于 items 总数缓存（当 total > 10 万时启用）
type ItemCountCache struct {
	mu      sync.RWMutex
	count   int64
	expires time.Time
}

// Get 返回缓存的 count 及是否命中
func (c *ItemCountCache) Get() (int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if time.Now().Before(c.expires) {
		return c.count, true
	}
	return 0, false
}

// Set 设置缓存
func (c *ItemCountCache) Set(count int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count = count
	c.expires = time.Now().Add(itemCountCacheTTL * time.Second)
}

// Invalidate 使缓存失效（写操作后调用）
func (c *ItemCountCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.expires = time.Time{}
}

// Threshold 返回启用缓存的阈值
func (ItemCountCache) Threshold() int64 {
	return itemCountCacheThreshold
}
