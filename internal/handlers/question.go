package handlers

import (
	"ai-ocs/internal/ai"
	"ai-ocs/internal/database"
	"ai-ocs/internal/models"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth API密钥验证中间件
func APIKeyAuth(config *models.Config) gin.HandlerFunc {
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
		valid := false
		for _, key := range config.APIKeys {
			if key == apiKey {
				valid = true
				break
			}
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

// GenerateAPIKey 生成新的API密钥
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SearchAnswer 处理查询答案的请求
func SearchAnswer(config *models.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取查询参数
		title := strings.TrimSpace(c.Query("title"))
		options := c.Query("options")
		questionType := c.Query("type")

		// 参数校验
		if title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "题目不能为空"})
			return
		}

		// 先从数据库查询答案
		answer, err := database.GetAnswer(title)
		if err != nil {
			// 如果数据库查询出错，记录日志但不中断流程
			log.Printf("数据库查询失败: %v", err)
		}

		// 如果数据库中有答案，直接返回
		if answer != "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "获取成功",
				"data": gin.H{
					"data": answer,
				},
			})
			return
		}

		// 如果数据库中没有答案，调用AI模型获取答案
		answer, err = ai.QueryLargeModel(
			title,
			options,
			questionType,
			config.Platform,
			config.APIKeys,
			config.Models,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "AI模型调用失败: " + err.Error()})
			return
		}

		// 将答案存入数据库
		err = database.SaveAnswer(title, answer)
		if err != nil {
			// 如果数据库保存出错，记录日志但不中断流程
			log.Printf("数据库保存失败: %v", err)
		}

		// 去除可能的markdown代码块标记
		answer = strings.TrimPrefix(answer, "```json")
		answer = strings.TrimSuffix(answer, "```")
		answer = strings.TrimSpace(answer)

		// 返回结果
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "获取成功",
			"data": gin.H{
				"data": answer,
			},
		})
	}
}