package router

import (
	"jmeter-test-api/cache"
	"jmeter-test-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zsais/go-gin-prometheus"
	"gorm.io/gorm"
)

// Setup configures and returns the Gin router
func Setup(db *gorm.DB, itemCountCache *cache.ItemCountCache) *gin.Engine {
	r := gin.Default()

	// Prometheus 指标：只加中间件采集 HTTP 请求，不调用 p.Use 避免重复注册 /metrics
	p := ginprometheus.NewPrometheus("jmeter_api")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		if path := c.FullPath(); path != "" {
			return path
		}
		return "unknown"
	}
	r.Use(p.HandlerFunc())
	// 统一用 promhttp 暴露全部指标（含 HTTP、process、cache、gorm）
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	itemHandler := handlers.NewItemHandler(db, itemCountCache)

	api := r.Group("/api")
	{
		api.GET("/health", itemHandler.Health)
		api.GET("/items/slow", itemHandler.Slow)
		if db != nil {
			api.GET("/items", itemHandler.ListItems)
			api.GET("/items/:id", itemHandler.GetItem)
			api.POST("/items", itemHandler.CreateItem)
			api.PUT("/items/:id", itemHandler.UpdateItem)
			api.DELETE("/items/:id", itemHandler.DeleteItem)
		}
	}

	return r
}
