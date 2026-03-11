package router

import (
	"jmeter-test-api/cache"
	"jmeter-test-api/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup configures and returns the Gin router
func Setup(db *gorm.DB, itemCountCache *cache.ItemCountCache) *gin.Engine {
	r := gin.Default()

	itemHandler := handlers.NewItemHandler(db, itemCountCache)

	api := r.Group("/api")
	{
		api.GET("/health", itemHandler.Health)
		api.GET("/items/slow", itemHandler.Slow)
		api.GET("/items", itemHandler.ListItems)
		api.GET("/items/:id", itemHandler.GetItem)
		api.POST("/items", itemHandler.CreateItem)
		api.PUT("/items/:id", itemHandler.UpdateItem)
		api.DELETE("/items/:id", itemHandler.DeleteItem)
	}

	return r
}
