package main

import (
	"log"

	"jmeter-test-api/cache"
	"jmeter-test-api/config"
	"jmeter-test-api/database"
	"jmeter-test-api/router"
)

func main() {
	cfg := config.Load()

	db, err := database.InitDB(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	itemCountCache := &cache.ItemCountCache{}
	r := router.Setup(db, itemCountCache)

	addr := ":" + cfg.Port
	log.Printf("Server starting on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
