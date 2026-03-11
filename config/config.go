package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
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
		// Railway 使用 MYSQL_PRIVATE_URL，其他平台可能用 MYSQL_URL 或 DATABASE_URL
		mysqlDSN = parseMySQLURL(os.Getenv("MYSQL_PRIVATE_URL"))
	}
	if mysqlDSN == "" {
		mysqlDSN = parseMySQLURL(os.Getenv("MYSQL_URL"))
	}
	if mysqlDSN == "" {
		mysqlDSN = parseMySQLURL(os.Getenv("DATABASE_URL"))
	}
	if mysqlDSN == "" {
		user := getEnv("MYSQL_USER", "root")
		pass := getEnv("MYSQL_PASS", "root")
		host := getEnv("MYSQL_HOST", "127.0.0.1") // 用 127.0.0.1 替代 localhost，避免 IPv6 [::1] 连接失败
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

// parseMySQLURL 将 mysql://user:pass@host:port/db 转为 Go MySQL DSN
func parseMySQLURL(raw string) string {
	if raw == "" {
		return ""
	}
	if !strings.HasPrefix(raw, "mysql://") {
		raw = "mysql://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	user := "root"
	if u.User != nil {
		user = u.User.Username()
	}
	pass, _ := u.User.Password()
	host := u.Host
	if host == "" {
		return ""
	}
	db := strings.TrimPrefix(u.Path, "/")
	if db == "" {
		db = "jmeter_test"
	}
	// localhost 可能解析为 [::1]，导致连接失败，统一改为 127.0.0.1
	if host == "localhost" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True", user, pass, host, db)
}
