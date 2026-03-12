package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// MySQL 默认配置（无环境变量时使用）
const (
	mysqlUser = "root"
	mysqlPass = "root"
	mysqlHost = "localhost" // 与 mysql -u root -proot 的默认连接一致；若报错可改为 "127.0.0.1"
	mysqlPort = "3306"
	mysqlDB   = "jmeter_test"
)

// Config holds application configuration
type Config struct {
	MySQLDSN string
	Port     string
	SkipDB   bool // 为 true 时跳过 MySQL，仅启动 /metrics、/health、/items/slow（供 MySQL 不可用时使用）
}

// Load loads configuration
func Load() *Config {
	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		mysqlDSN = parseMySQLURL(os.Getenv("MYSQL_PRIVATE_URL"))
	}
	if mysqlDSN == "" {
		mysqlDSN = parseMySQLURL(os.Getenv("MYSQL_URL"))
	}
	if mysqlDSN == "" {
		mysqlDSN = parseMySQLURL(os.Getenv("DATABASE_URL"))
	}
	if mysqlDSN == "" {
		mysqlDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&allowNativePasswords=true",
			mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDB)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	skipDB := strings.ToLower(os.Getenv("SKIP_DB")) == "1" || strings.ToLower(os.Getenv("SKIP_DB")) == "true"

	return &Config{
		MySQLDSN: mysqlDSN,
		Port:     port,
		SkipDB:   skipDB,
	}
}

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
	user, pass := "root", "root"
	if u.User != nil {
		user = u.User.Username()
		pass, _ = u.User.Password()
	}
	host := u.Host
	if host == "" {
		return ""
	}
	// 保持 localhost 原样，与 mysql CLI 默认行为一致
	db := strings.TrimPrefix(u.Path, "/")
	if db == "" {
		db = "jmeter_test"
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&allowNativePasswords=true", user, pass, host, db)
}
