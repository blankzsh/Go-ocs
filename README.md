# AI-OCS Go版本

AI题库系统Go语言重构版本，支持多种AI平台API调用。

## 功能特性

- 题目答案自动生成
- 支持多种AI平台（硅基流动、阿里云百炼、智普AI、Ollama、DeepSeek、ChatGPT、Gemini）
- 本地数据库缓存答案
- RESTful API接口

## 技术栈

- Go 1.21+
- Gin Web框架
- MySQL 数据库
- 多种AI平台API：
  - SiliconFlow API
  - 阿里云百炼平台API
  - 智普AI API
  - Ollama本地模型API
  - DeepSeek官方API
  - ChatGPT API
  - Gemini API

## 目录结构

```
Go-ocs/
├── cmd/                 # 主程序入口
├── configs/             # 配置文件
│   ├── config.example.json  # 示例配置文件
│   └── config.json      # 实际配置文件（需手动创建，不会被版本控制）
├── internal/            # 内部模块
│   ├── ai/              # AI服务相关
│   ├── database/        # 数据库相关
│   ├── handlers/        # HTTP处理器
│   └── models/          # 数据模型
├── .gitignore           # Git忽略文件
├── go.mod              # Go模块定义
├── go.sum              # Go模块校验和
└── README.md           # 项目说明
```

## 安装与运行

1. 确保已安装Go 1.21+

2. 克隆项目或复制文件到本地

3. 安装依赖：
   ```bash
   cd Go-ocs
   go mod tidy
   ```

4. 配置文件设置：
   - 复制 [configs/config.example.json](file:///f:/WEB-PR/Go-ocsBase/configs/config.example.json) 为 `configs/config.json`
   - 修改 `configs/config.json` 中的配置项：
   ```json
   {
       "host": "127.0.0.1",
       "port": 8000,
       "platform": "siliconflow",
       "api_keys": {
           "aliyun": "your_aliyun_api_key_here",
           "siliconflow": "your_siliconflow_api_key_here",
           "zhipu": "your_zhipu_api_key_here",
           "deepseek": "your_deepseek_api_key_here",
           "chatgpt": "your_chatgpt_api_key_here",
           "gemini": "your_gemini_api_key_here"
       },
       "models": {
           "aliyun": "qwen-plus-latest",
           "siliconflow": "deepseek-ai/DeepSeek-R1",
           "zhipu": "glm-4",
           "ollama": "llama3",
           "deepseek": "deepseek-chat",
           "chatgpt": "gpt-3.5-turbo",
           "gemini": "gemini-pro"
       },
       "mysql": {
           "host": "127.0.0.1",
           "port": 3306,
           "user": "root",
           "password": "your_mysql_password",
           "database": "question_bank"
       }
   }
   ```

5. 运行程序：
   ```bash
   go run cmd/main.go
   ```

## API使用

### 查询题目答案

```
GET /api/query?title=问题内容[&options=选项内容][&type=问题类型]
```

示例：
```
curl "http://127.0.0.1:8000/api/query?title=中国的首都是哪里%3F&options=北京###上海###广州###深圳&type=选择题"
```

## 配置说明

- `host`: 服务器监听地址
- `port`: 服务器监听端口
- `platform`: 使用的AI平台
- `api_keys`: 各平台的API密钥
- `models`: 各平台使用的模型
- `mysql`: MySQL数据库配置

## 平台切换

在配置文件中修改 `platform` 字段：
- `siliconflow`: 使用硅基流动平台
- `aliyun`: 使用阿里云百炼平台
- `zhipu`: 使用智普AI平台
- `ollama`: 使用Ollama本地模型
- `deepseek`: 使用DeepSeek官方API
- `chatgpt`: 使用ChatGPT API
- `gemini`: 使用Gemini API

## 配置工具

项目提供了命令行配置工具，可以方便地配置各平台参数：

### Windows系统
由于编码兼容性问题，Windows用户建议使用PowerShell版本的配置工具：

英文版本：
```cmd
powershell -ExecutionPolicy Bypass -File .\configurator.ps1
```

中文版本：
```cmd
powershell -ExecutionPolicy Bypass -File .\configurator_zh.ps1
```

### Linux/Mac系统
```bash
./configurator.sh
```

通过配置工具可以：
- 选择AI平台
- 配置各平台API密钥
- 设置各平台使用的模型
- 配置MySQL数据库连接
- 查看当前配置

## 项目索引和版本控制

### .gitignore文件
项目包含一个 [.gitignore](file:///f:/WEB-PR/Go-ocsBase/.gitignore) 文件，用于排除以下内容：
- 编译生成的二进制文件和构建产物
- IDE配置文件
- 操作系统生成的文件
- 本地配置文件（如 `configs/config.json`）
- 日志文件
- 临时文件
- 环境变量文件

### 配置文件管理
- [configs/config.example.json](file:///f:/WEB-PR/Go-ocsBase/configs/config.example.json)：示例配置文件，包含所有配置项的模板
- `configs/config.json`：实际配置文件，不会被版本控制，需要用户根据示例文件手动创建