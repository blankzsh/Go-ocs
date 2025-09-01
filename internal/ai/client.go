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

// ZhipuAIResponse 智普AI模型响应结构
type ZhipuAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// OllamaResponse Ollama模型响应结构
type OllamaResponse struct {
	Model   string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

// GeminiResponse Gemini模型响应结构
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
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

// OllamaRequest Ollama请求结构
type OllamaRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Format  string `json:"format"`
}

// ZhipuAIRequest 智普AI请求结构
type ZhipuAIRequest struct {
	Model       string `json:"model"`
	Messages    []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	MaxTokens   int     `json:"max_tokens"`
}

// GeminiRequest Gemini请求结构
type GeminiRequest struct {
	Contents []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
	GenerationConfig struct {
		MaxOutputTokens   int     `json:"maxOutputTokens"`
		Temperature       float64 `json:"temperature"`
		TopP              float64 `json:"topP"`
	} `json:"generationConfig"`
}

// QueryLargeModel 调用AI模型获取问题答案
func QueryLargeModel(title, options, questionType, platform string, apiKeys map[string]string, models map[string]string) (string, error) {
	switch platform {
	case "siliconflow":
		return querySiliconFlow(title, options, questionType, apiKeys["siliconflow"], models["siliconflow"])
	case "aliyun":
		return queryAliyunBailian(title, options, questionType, apiKeys["aliyun"], models["aliyun"])
	case "zhipu":
		return queryZhipuAI(title, options, questionType, apiKeys["zhipu"], models["zhipu"])
	case "ollama":
		return queryOllama(title, options, questionType, models["ollama"])
	case "deepseek":
		return queryDeepSeek(title, options, questionType, apiKeys["deepseek"], models["deepseek"])
	case "chatgpt":
		return queryChatGPT(title, options, questionType, apiKeys["chatgpt"], models["chatgpt"])
	case "gemini":
		return queryGemini(title, options, questionType, apiKeys["gemini"], models["gemini"])
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

// queryZhipuAI 调用智普AI平台API获取问题答案
func queryZhipuAI(title, options, questionType, apiKey, model string) (string, error) {
	url := "https://open.bigmodel.cn/api/paas/v4/chat/completions"

	// 构建简化的提问内容，减少token数量
	content := fmt.Sprintf(`你是题库接口，根据问题和选项提供答案。选择题返回选项内容；多选题用###连接；判断题返回"对"或"错"；填空题用###连接多个空。格式：{"anwser":"答案"}。只返回json格式。
{
	"问题": "%s",
	"选项": "%s",
	"类型": "%s"
}`, title, options, questionType)

	// 构建请求体
	requestBody := ZhipuAIRequest{
		Model:       model,
		Stream:      false,
		Temperature: 0.05,
		TopP:        0.95,
		MaxTokens:   256,
	}
	requestBody.Messages = append(requestBody.Messages, struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{Role: "user", Content: content})

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
	var aiResp ZhipuAIResponse
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

// queryOllama 调用Ollama本地模型获取问题答案
func queryOllama(title, options, questionType, model string) (string, error) {
	url := "http://localhost:11434/api/generate"

	// 构建简化的提问内容，减少token数量
	content := fmt.Sprintf(`你是题库接口，根据问题和选项提供答案。选择题返回选项内容；多选题用###连接；判断题返回"对"或"错"；填空题用###连接多个空。格式：{"anwser":"答案"}。只返回json格式。
{
	"问题": "%s",
	"选项": "%s",
	"类型": "%s"
}`, title, options, questionType)

	// 构建请求体
	requestBody := OllamaRequest{
		Model:  model,
		Prompt: content,
		Stream: false,
		Format: "json",
	}

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
	var aiResp OllamaResponse
	err = json.Unmarshal(body, &aiResp)
	if err != nil {
		return "无法解析API响应: " + err.Error(), err
	}

	// 提取答案
	return aiResp.Message.Content, nil
}

// queryDeepSeek 调用DeepSeek官方API获取问题答案
func queryDeepSeek(title, options, questionType, apiKey, model string) (string, error) {
	url := "https://api.deepseek.com/chat/completions"

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
		return "无法解析API响应: " + err.Error(), err
	}

	// 提取答案
	if len(aiResp.Choices) > 0 {
		return aiResp.Choices[0].Message.Content, nil
	}

	return "无法从API获取答案", nil
}

// queryChatGPT 调用ChatGPT API获取问题答案
func queryChatGPT(title, options, questionType, apiKey, model string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

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
		return "无法解析API响应: " + err.Error(), err
	}

	// 提取答案
	if len(aiResp.Choices) > 0 {
		return aiResp.Choices[0].Message.Content, nil
	}

	return "无法从API获取答案", nil
}

// queryGemini 调用Gemini API获取问题答案
func queryGemini(title, options, questionType, apiKey, model string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)

	// 构建简化的提问内容，减少token数量
	content := fmt.Sprintf(`你是题库接口，根据问题和选项提供答案。选择题返回选项内容；多选题用###连接；判断题返回"对"或"错"；填空题用###连接多个空。格式：{"anwser":"答案"}。只返回json格式。
{
	"问题": "%s",
	"选项": "%s",
	"类型": "%s"
}`, title, options, questionType)

	// 构建请求体
	requestBody := GeminiRequest{}
	requestBody.Contents = append(requestBody.Contents, struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{Role: "user"})
	requestBody.Contents[0].Parts = append(requestBody.Contents[0].Parts, struct {
		Text string `json:"text"`
	}{Text: content})
	
	requestBody.GenerationConfig.MaxOutputTokens = 256
	requestBody.GenerationConfig.Temperature = 0.05
	requestBody.GenerationConfig.TopP = 0.95

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
	var aiResp GeminiResponse
	err = json.Unmarshal(body, &aiResp)
	if err != nil {
		return "无法解析API响应: " + err.Error(), err
	}

	// 提取答案
	if len(aiResp.Candidates) > 0 && len(aiResp.Candidates[0].Content.Parts) > 0 {
		return aiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "无法从API获取答案", nil
}