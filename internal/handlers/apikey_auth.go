package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// APISQLDB API密钥管理系统的数据库连接
var APISQLDB *sql.DB

// InitAPIDB 初始化API密钥数据库连接
func InitAPIDB() error {
	// 从环境变量获取数据库配置
	dbHost := getEnv("API_KEY_DB_HOST", "localhost")
	dbPort := getEnv("API_KEY_DB_PORT", "3306")
	dbUser := getEnv("API_KEY_DB_USER", "root")
	dbPassword := getEnv("API_KEY_DB_PASSWORD", "")
	dbName := getEnv("API_KEY_DB_NAME", "api_key_manager")

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// 连接数据库
	var err error
	APISQLDB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("无法连接到API密钥数据库: %v", err)
	}

	// 测试连接
	if err := APISQLDB.Ping(); err != nil {
		return fmt.Errorf("无法ping通API密钥数据库: %v", err)
	}

	log.Println("成功连接到API密钥管理系统数据库")
	return nil
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// APIKeyAuthWithDB 使用数据库验证API密钥的中间件
func APIKeyAuthWithDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取API密钥
		apiKey := c.GetHeader("API-Key")
		if apiKey == "" {
			// 如果请求头中没有，从查询参数中获取
			apiKey = c.Query("api_key")
		}

		// 如果没有提供API密钥
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "缺少API密钥",
			})
			c.Abort()
			return
		}

		// 验证API密钥是否有效
		valid, err := validateAPIKey(apiKey)
		if err != nil {
			log.Printf("验证API密钥时出错: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 1,
				"msg":  "服务器内部错误",
			})
			c.Abort()
			return
		}

		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "无效的API密钥",
			})
			c.Abort()
			return
		}

		// 验证通过，继续处理请求
		c.Next()
	}
}

// validateAPIKey 验证API密钥
func validateAPIKey(apiKey string) (bool, error) {
	// 查询所有密钥哈希值
	rows, err := APISQLDB.Query("SELECT key_value FROM api_keys")
	if err != nil {
		return false, fmt.Errorf("查询API密钥失败: %v", err)
	}
	defer rows.Close()

	// 验证密钥
	for rows.Next() {
		var hashedKey string
		if err := rows.Scan(&hashedKey); err != nil {
			return false, fmt.Errorf("扫描API密钥失败: %v", err)
		}

		// 使用bcrypt比较密钥
		if err := bcrypt.CompareHashAndPassword([]byte(hashedKey), []byte(apiKey)); err == nil {
			return true, nil
		}
	}

	if err = rows.Err(); err != nil {
		return false, fmt.Errorf("遍历API密钥时出错: %v", err)
	}

	return false, nil
}