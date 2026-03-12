package main

import (
	"errors"
	"log"

	"jmeter-test-api/cache"
	"jmeter-test-api/config"
	"jmeter-test-api/database"
	"jmeter-test-api/router"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	var db *gorm.DB
	if cfg.SkipDB {
		log.Println("SKIP_DB=1: 跳过 MySQL，仅提供 /metrics、/health、/api/items/slow")
		db = nil
	} else {
		var err error
		db, err = database.InitDB(cfg.MySQLDSN)
		if err != nil {
			log.Printf("MySQL 连接失败: %v，自动降级为仅指标模式（/metrics、/health、/items/slow）", err)
			db = nil
		}
	}

	// 注册进程指标（CPU、内存、文件句柄等）；若已注册则忽略（避免与其它包冲突）
	if err := prometheus.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		var already prometheus.AlreadyRegisteredError
		if !errors.As(err, &already) {
			log.Fatalf("注册 ProcessCollector 失败: %v", err)
		}
	}

	itemCountCache := cache.NewItemCountCache()
	r := router.Setup(db, itemCountCache)

	addr := ":" + cfg.Port
	log.Printf("Server starting on http://localhost%s (SKIP_DB=%v)", addr, cfg.SkipDB)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
