package handlers

import (
	"ai-ocs/internal/ai"
	"ai-ocs/internal/models"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// TestHandler 测试接口处理器
func TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "测试接口正常工作",
		"data": gin.H{
			"timestamp": 0,
		},
	})
}

// TestAnswerHandler 测试答题功能的接口处理器
func TestAnswerHandler(c *gin.Context) {
	// 使用相对路径获取配置文件路径
	configPath := filepath.Join("configs", "config.json")

	// 加载配置
	config, err := models.LoadConfig(configPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "配置加载失败: " + err.Error(),
		})
		return
	}

	// 获取查询参数
	title := c.Query("title")
	options := c.Query("options")
	questionType := c.Query("type")

	// 参数校验
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "题目不能为空",
		})
		return
	}

	// 调用AI模型获取答案
	answer, err := ai.QueryLargeModel(
		title,
		options,
		questionType,
		config.Platform,
		config.APIKeys,
		config.Models,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "AI模型调用失败: " + err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "测试答题成功",
		"data": gin.H{
			"answer": answer,
		},
	})
}