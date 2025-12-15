package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 存储应用程序配置
type Config struct {
	DBHost            string
	DBPort            int
	DBUser            string
	DBPassword        string
	DBName            string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime int
	DBLogMode         bool
	DBLogLevel        string
}

func LoadConfig() (*Config, error) {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("加载.env文件失败: %v", err)
	}

	// 解析端口
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "3306"))
	if err != nil {
		return nil, fmt.Errorf("解析DB_PORT失败: %v", err)
	}

	// 解析连接池配置
	maxOpenConns, err := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "100"))
	if err != nil {
		return nil, fmt.Errorf("解析DB_MAX_OPEN_CONNS失败: %v", err)
	}

	maxIdleConns, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	if err != nil {
		return nil, fmt.Errorf("解析DB_MAX_IDLE_CONNS失败: %v", err)
	}

	connMaxLifetime, err := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "300"))
	if err != nil {
		return nil, fmt.Errorf("解析DB_CONN_MAX_LIFETIME失败: %v", err)
	}

	// 解析日志模式
	logMode := getEnv("DB_LOG_MODE", "true") == "true"

	return &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            dbPort,
		DBUser:            getEnv("DB_USER", "root"),
		DBPassword:        getEnv("DB_PASSWORD", "password"),
		DBName:            getEnv("DB_NAME", "gorm_learning_db"),
		DBMaxOpenConns:    maxOpenConns,
		DBMaxIdleConns:    maxIdleConns,
		DBConnMaxLifetime: connMaxLifetime,
		DBLogMode:         logMode,
		DBLogLevel:        getEnv("DB_LOG_LEVEL", "info"),
	}, nil
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
