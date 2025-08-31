package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AIResponse AI模型响应结构
type AIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// QueryRequest 查询请求结构
type QueryRequest struct {
	Model       string `json:"model"`
	Messages    []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream        bool `json:"stream"`
	MaxTokens     int  `json:"max_tokens"`
	Temperature   float64 `json:"temperature"`
	TopP          float64 `json:"top_p"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	N             int  `json:"n"`
	ResponseFormat struct {
		Type string `json:"type"`
	} `json:"response_format"`
}

// QueryLargeModel 调用AI模型获取问题答案
func QueryLargeModel(title, options, questionType, platform string, apiKeys map[string]string, models map[string]string) (string, error) {
	switch platform {
	case "siliconflow":
		return querySiliconFlow(title, options, questionType, apiKeys["siliconflow"], models["siliconflow"])
	case "aliyun":
		return queryAliyunBailian(title, options, questionType, apiKeys["aliyun"], models["aliyun"])
	default:
		return querySiliconFlow(title, options, questionType, apiKeys["siliconflow"], models["siliconflow"])
	}
}

// querySiliconFlow 调用SiliconFlow API获取问题答案
func querySiliconFlow(title, options, questionType, apiKey, model string) (string, error) {
	url := "https://api.siliconflow.cn/v1/chat/completions"

	// 构建简化的提问内容，减少token数量
	content := fmt.Sprintf(`你是题库接口，根据问题和选项提供答案。选择题返回选项内容；多选题用###连接；判断题返回"对"或"错"；填空题用###连接多个空。格式：{"anwser":"答案"}。只返回json格式。
{
	"问题": "%s",
	"选项": "%s",
	"类型": "%s"
}`, title, options, questionType)

	// 构建请求体
	requestBody := QueryRequest{
		Model:     model,
		Stream:    false,
		MaxTokens: 256,
		Temperature: 0.05,
		TopP:      0.95,
		FrequencyPenalty: 0.0,
		N:         1,
	}
	requestBody.Messages = append(requestBody.Messages, struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{Role: "user", Content: content})
	requestBody.ResponseFormat.Type = "json_object"

	// 转换为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "API调用失败", err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "API调用失败", err
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 创建HTTP客户端并发送请求
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "API调用失败", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "API调用失败", err
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("API调用失败，状态码: %d", resp.StatusCode), nil
	}

	// 解析响应
	var aiResp AIResponse
	err = json.Unmarshal(body, &aiResp)
	if err != nil {
		return "无法解析API响应", err
	}

	// 提取答案
	if len(aiResp.Choices) > 0 {
		return aiResp.Choices[0].Message.Content, nil
	}

	return "无法从API获取答案", nil
}

// queryAliyunBailian 调用阿里云百炼平台API获取问题答案
func queryAliyunBailian(title, options, questionType, apiKey, model string) (string, error) {
	url := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

	// 构建简化的提问内容，减少token数量
	content := fmt.Sprintf(`你是题库接口，根据问题和选项提供答案。选择题返回选项内容；多选题用###连接；判断题返回"对"或"错"；填空题用###连接多个空。格式：{"anwser":"答案"}。只返回json格式。
{
	"问题": "%s",
	"选项": "%s",
	"类型": "%s"
}`, title, options, questionType)

	// 构建请求体
	requestBody := QueryRequest{
		Model:     model,
		Stream:    false,
		MaxTokens: 256,
		Temperature: 0.05,
		TopP:      0.95,
		N:         1,
	}
	requestBody.Messages = append(requestBody.Messages, struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{Role: "user", Content: content})
	requestBody.ResponseFormat.Type = "json_object"

	// 转换为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "API调用失败", err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "API调用失败", err
	}

	// 设置请求头，使用阿里云的API密钥
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-SSE", "enable")

	// 创建HTTP客户端并发送请求
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "API调用失败", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "API调用失败", err
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("API调用失败，状态码: %d", resp.StatusCode), nil
	}

	// 解析响应
	var aiResp AIResponse
	err = json.Unmarshal(body, &aiResp)
	if err != nil {
		return "无法解析API响应: " + err.Error(), err
	}

	// 提取答案
	if len(aiResp.Choices) > 0 {
		return aiResp.Choices[0].Message.Content, nil
	}

	return "无法从API获取答案", nil
}