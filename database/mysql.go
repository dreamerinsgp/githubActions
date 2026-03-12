package database

import (
	"jmeter-test-api/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

// InitDB initializes the MySQL connection and runs AutoMigrate
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true, // 复用预编译语句，减少解析开销
	})
	if err != nil {
		return nil, err
	}

	// 连接池：高并发下避免连接耗尽
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(20)   // 空闲连接数
	sqlDB.SetMaxOpenConns(100)  // 最大打开连接数
	sqlDB.SetConnMaxLifetime(0) // 连接不复用时长限制

	// 注册 DB 指标到 Prometheus（连接池状态等）
	if err := db.Use(prometheus.New(prometheus.Config{
		DBName:          "jmeter_test",
		RefreshInterval: 15,
		StartServer:     false, // 不单独起服务，由应用 /metrics 统一暴露
	})); err != nil {
		return nil, err
	}

	// Auto migrate the Item model
	if err := db.AutoMigrate(&models.Item{}); err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
