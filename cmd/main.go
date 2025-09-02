package main

import (
	"ai-ocs/internal/database"
	"ai-ocs/internal/handlers"
	"ai-ocs/internal/models"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 使用相对路径获取配置文件路径
	configPath := filepath.Join("configs", "config.json")

	// 加载配置
	config, err := models.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(config); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置Gin为发布模式（生产环境）
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.Default()

	// 注册API路由
	r.GET("/api/query", handlers.SearchAnswer(config))
	r.GET("/api/test", handlers.TestHandler)              // 添加测试接口
	r.GET("/api/test-answer", handlers.TestAnswerHandler) // 添加测试答题接口

	// 注册管理后台路由
	admin := r.Group("/admin")
	{
		admin.GET("/login", handlers.LoginPage)
		admin.POST("/login", handlers.Login)
		admin.POST("/logout", handlers.Logout)
		admin.GET("/", handlers.RequireAuth, handlers.AdminPage)
		admin.GET("/stats", handlers.RequireAuth, handlers.GetStats)
		admin.GET("/questions", handlers.RequireAuth, handlers.GetQuestions)
		admin.GET("/search", handlers.RequireAuth, handlers.SearchQuestion)
	}

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	actualPort := config.Port

	// 检查端口是否可用
	if !isPortAvailable(config.Host, config.Port) {
		// 尝试寻找可用端口
		newPort := findAvailablePort(config.Host, config.Port)
		if newPort != -1 {
			log.Printf("端口 %d 已被占用，切换到可用端口 %d", config.Port, newPort)
			addr = fmt.Sprintf("%s:%d", config.Host, newPort)
			actualPort = newPort
		} else {
			log.Fatalf("无法找到可用端口")
		}
	}

	// 打印API配置信息
	printAPIConfig(config, actualPort)

	log.Printf("启动Gin服务器，地址: %s", addr)
	log.Fatal(r.Run(addr))
}

// printAPIConfig 打印API配置信息
func printAPIConfig(config *models.Config, actualPort int) {
	// 获取实际监听的IP和端口
	host := config.Host
	if host == "" {
		host = "127.0.0.1"
	}

	// 构造API配置信息
	url := fmt.Sprintf("http://%s:%d/api/query", host, actualPort)
	jsonStr := fmt.Sprintf("{\n  \"name\": \"完美题库\",\n  \"homepage\": \"https://currso.com/\",\n  \"url\": \"%s\",\n  \"method\": \"get\",\n  \"type\": \"GM_xmlhttpRequest\",\n  \"contentType\": \"json\",\n  \"data\": {\n    \"title\": \"${title}\",\n    \"options\": \"${options}\",\n    \"type\": \"${type}\"\n  },\n  \"handler\": \"return (res)=>res.code === 0 ? [undefined, undefined] : [undefined,res.data.data]\"\n}", url)

	log.Printf("API配置信息:\n%s", jsonStr)
}

// isPortAvailable 检查端口是否可用
func isPortAvailable(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return true
	}
	conn.Close()
	return false
}

// findAvailablePort 寻找可用端口
func findAvailablePort(host string, startPort int) int {
	for port := startPort + 1; port <= startPort+100; port++ {
		if isPortAvailable(host, port) {
			return port
		}
	}
	return -1
}
