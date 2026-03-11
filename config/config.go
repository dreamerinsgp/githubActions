package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	MySQLDSN string
	Port     string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		user := getEnv("MYSQL_USER", "root")
		pass := getEnv("MYSQL_PASS", "root")
		host := getEnv("MYSQL_HOST", "localhost")
		port := getEnv("MYSQL_PORT", "3306")
		db := getEnv("MYSQL_DB", "jmeter_test")
		mysqlDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", user, pass, host, port, db)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		MySQLDSN: mysqlDSN,
		Port:     port,
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
