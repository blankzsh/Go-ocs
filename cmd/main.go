package main

import (
	"ai-ocs/internal/database"
	"ai-ocs/internal/handlers"
	"ai-ocs/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"path/filepath"
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
	// gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	r.GET("/api/query", handlers.SearchAnswer(config))

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
		},
		"handler": "return (res)=>res.code === 0 ? [undefined, undefined] : [undefined,res.data.data]",
	}

	fmt.Println("\nAPI配置信息:")
	jsonData, _ := json.MarshalIndent(apiConfig, "", "  ")
	fmt.Println(string(jsonData))
}
