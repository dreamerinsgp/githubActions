package cache

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
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

	hitsTotal         prometheus.Counter
	missesTotal       prometheus.Counter
	invalidationsTotal prometheus.Counter
}

// NewItemCountCache 创建带 Prometheus 指标的缓存实例
func NewItemCountCache() *ItemCountCache {
	c := &ItemCountCache{
		hitsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "jmeter_api_item_count_cache_hits_total",
			Help: "Total number of item count cache hits",
		}),
		missesTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "jmeter_api_item_count_cache_misses_total",
			Help: "Total number of item count cache misses",
		}),
		invalidationsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "jmeter_api_item_count_cache_invalidations_total",
			Help: "Total number of item count cache invalidations",
		}),
	}
	prometheus.MustRegister(c.hitsTotal, c.missesTotal, c.invalidationsTotal)
	return c
}

// Get 返回缓存的 count 及是否命中
func (c *ItemCountCache) Get() (int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if time.Now().Before(c.expires) {
		c.hitsTotal.Inc()
		return c.count, true
	}
	c.missesTotal.Inc()
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
	c.invalidationsTotal.Inc()
	c.expires = time.Time{}
}

// Threshold 返回启用缓存的阈值
func (ItemCountCache) Threshold() int64 {
	return itemCountCacheThreshold
}
