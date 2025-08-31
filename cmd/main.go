package main

import (
	"ai-ocs/internal/database"
	"ai-ocs/internal/handlers"
	"ai-ocs/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，将使用环境变量或默认值")
	}

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

	// 初始化API密钥数据库
	if err := handlers.InitAPIDB(); err != nil {
		log.Fatalf("API密钥数据库初始化失败: %v", err)
	}

	// 设置Gin为发布模式（生产环境）
	// gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	r.GET("/api/query", handlers.APIKeyAuthWithDB(), handlers.SearchAnswer(config))

	// 添加一个生成API密钥的管理端点（仅用于生成密钥，实际使用中应该有额外的保护）
	r.POST("/api/generate-key", func(c *gin.Context) {
		key, err := handlers.GenerateAPIKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 1,
				"msg":  "生成API密钥失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "生成API密钥成功",
			"data": gin.H{
				"api_key": key,
			},
		})
	})

	// 打印API配置信息
	printAPIConfig(config)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Printf("启动Gin服务器，地址: %s:%d", config.Host, config.Port)
	log.Fatal(r.Run(addr))
}

// printAPIConfig 打印API配置信息
func printAPIConfig(config *models.Config) {
	apiConfig := map[string]interface{}{
		"name":        "完美题库",
		"url":         fmt.Sprintf("http://%s:%d/api/query", config.Host, config.Port),
		"method":      "get",
		"type":        "GM_xmlhttpRequest",
		"homepage":    "https://currso.com/",
		"contentType": "json",
		"data": map[string]string{
			"title":   "${title}",
			"options": "${options}",
			"type":    "${type}",
			"api-key": "YOUR_API_KEY_HERE", // 添加API密钥占位符
		},
		"handler": "return (res)=>res.code === 0 ? [undefined, undefined] : [undefined,res.data.data]",
	}

	fmt.Println("\nAPI配置信息:")
	jsonData, _ := json.MarshalIndent(apiConfig, "", "  ")
	fmt.Println(string(jsonData))
}
