package models

import (
	"encoding/json"
	"os"
)

// Config 应用配置结构体
type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Platform string `json:"platform"`
	// API密钥配置
	APIKeys map[string]string `json:"api_keys"`
	// 模型配置
	Models map[string]string `json:"models"`
	// MySQL配置
	MySQLConfig MySQLConfig `json:"mysql"`
}

// MySQLConfig MySQL数据库配置
type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// LoadConfig 从文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// 设置默认值
	if config.Platform == "" {
		config.Platform = "siliconflow"
	}

	// 初始化APIKeys映射
	if config.APIKeys == nil {
		config.APIKeys = make(map[string]string)
	}

	// 设置默认API密钥
	if config.APIKeys["siliconflow"] == "" {
		config.APIKeys["siliconflow"] = ""
	}

	if config.APIKeys["aliyun"] == "" {
		config.APIKeys["aliyun"] = ""
	}
	
	if config.APIKeys["zhipu"] == "" {
		config.APIKeys["zhipu"] = ""
	}
	
	if config.APIKeys["deepseek"] == "" {
		config.APIKeys["deepseek"] = ""
	}
	
	if config.APIKeys["chatgpt"] == "" {
		config.APIKeys["chatgpt"] = ""
	}
	
	if config.APIKeys["gemini"] == "" {
		config.APIKeys["gemini"] = ""
	}

	// 初始化Models映射
	if config.Models == nil {
		config.Models = make(map[string]string)
	}

	// 设置默认模型
	if config.Models["siliconflow"] == "" {
		config.Models["siliconflow"] = "deepseek-ai/DeepSeek-R1"
	}

	if config.Models["aliyun"] == "" {
		config.Models["aliyun"] = "qwen-plus"
	}
	
	if config.Models["zhipu"] == "" {
		config.Models["zhipu"] = "glm-4"
	}
	
	if config.Models["ollama"] == "" {
		config.Models["ollama"] = "llama3"
	}
	
	if config.Models["deepseek"] == "" {
		config.Models["deepseek"] = "deepseek-chat"
	}
	
	if config.Models["chatgpt"] == "" {
		config.Models["chatgpt"] = "gpt-3.5-turbo"
	}
	
	if config.Models["gemini"] == "" {
		config.Models["gemini"] = "gemini-pro"
	}

	// 设置MySQL默认值
	if config.MySQLConfig.Host == "" {
		config.MySQLConfig.Host = "localhost"
	}

	if config.MySQLConfig.Port == 0 {
		config.MySQLConfig.Port = 3306
	}

	if config.MySQLConfig.User == "" {
		config.MySQLConfig.User = "root"
	}

	if config.MySQLConfig.Database == "" {
		config.MySQLConfig.Database = "question_bank"
	}

	return &config, nil
}