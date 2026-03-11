package handlers

import (
	"net/http"
	"strconv"
	"time"

	"jmeter-test-api/cache"
	"jmeter-test-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ItemHandler handles item-related HTTP requests
type ItemHandler struct {
	DB    *gorm.DB
	Cache *cache.ItemCountCache
}

// NewItemHandler creates a new ItemHandler
func NewItemHandler(db *gorm.DB, itemCache *cache.ItemCountCache) *ItemHandler {
	return &ItemHandler{DB: db, Cache: itemCache}
}

// Health returns a simple health check (no DB)
func (h *ItemHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "service is healthy",
	})
}

// ListItems returns a paginated list of items
func (h *ItemHandler) ListItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	// 限制最大 offset，防止 OFFSET 过大导致全表扫描变慢
	const maxOffset = 10000
	if offset > maxOffset {
		offset = maxOffset
		page = maxOffset/pageSize + 1
	}

	var items []models.Item
	var total int64

	// 优先使用内存缓存（当 total > 10 万时，避免频繁 COUNT）
	if cached, ok := h.Cache.Get(); ok {
		total = cached
	} else {
		if err := h.DB.Model(&models.Item{}).Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if total > h.Cache.Threshold() {
			h.Cache.Set(total)
		}
	}

	// 使用主键排序，利用索引避免 filesort
	if err := h.DB.Order("id ASC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// GetItem returns a single item by ID
func (h *ItemHandler) GetItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var item models.Item
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// CreateItem creates a new item
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := models.Item{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Cache.Invalidate() // 写操作后使 count 缓存失效
	c.JSON(http.StatusCreated, item)
}

// UpdateItem updates an existing item
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构造更新字段，仅更新传入的非空字段（避免 SELECT+Save 两次查询）
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	result := h.DB.Model(&models.Item{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	h.Cache.Invalidate() // 写操作后使 count 缓存失效

	// 返回更新后的数据（仅需一次查询）
	var item models.Item
	if err := h.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteItem deletes an item
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result := h.DB.Delete(&models.Item{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	h.Cache.Invalidate() // 写操作后使 count 缓存失效
	c.JSON(http.StatusOK, gin.H{"message": "item deleted"})
}

// Slow simulates a delayed response for response time testing
func (h *ItemHandler) Slow(c *gin.Context) {
	ms, _ := strconv.Atoi(c.DefaultQuery("ms", "100"))
	if ms < 0 {
		ms = 0
	}
	if ms > 10000 {
		ms = 10000
	}

	time.Sleep(time.Duration(ms) * time.Millisecond)

	c.JSON(http.StatusOK, gin.H{
		"message": "delayed response",
		"delay_ms": ms,
	})
}
