# AI-OCS Go版本

AI题库系统Go语言重构版本，支持多种AI平台API调用。

## 功能特性

- 题库查询API服务
- 支持多种AI平台（硅基流动、阿里云百炼）
- 本地数据库缓存
- **API密钥验证机制**

## 技术栈

- Go 1.21+
- Gin Web框架
- MySQL 数据库
- SiliconFlow API
- 阿里云百炼平台API

## 目录结构

```
Go-ocs/
├── cmd/                 # 主程序入口
├── configs/             # 配置文件
│   ├── config.example.json  # 示例配置文件
│   └── config.json      # 实际配置文件（需手动创建）
├── internal/            # 内部模块
│   ├── ai/              # AI服务相关
│   ├── database/        # 数据库相关
│   ├── handlers/        # HTTP处理器
│   └── models/          # 数据模型
├── go.mod              # Go模块定义
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
   - 复制 [configs/config.example.json](file:///f:/Project/Go-ocs/configs/config.example.json) 为 `configs/config.json`
   - 修改 `configs/config.json` 中的配置项：
   ```json
   {
       "host": "127.0.0.1",
       "port": 8000,
       "platform": "siliconflow",
       "api_keys": {
           "aliyun": "your_aliyun_api_key_here",
           "siliconflow": "your_siliconflow_api_key_here"
       },
       "models": {
           "aliyun": "qwen-plus-latest",
           "siliconflow": "deepseek-ai/DeepSeek-R1"
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

## API接口

### 查询接口

- 接口地址：`/api/query`
- 请求方法：`GET`
- 请求参数：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 是 | 题目内容 |
| options | string | 否 | 选项内容 |
| type | string | 否 | 题目类型 |
| api_key | string | 是 | API密钥（可通过请求头或查询参数传递） |

- 请求示例：

```bash
# 通过请求头传递API密钥
curl -H "API-Key: YOUR_API_KEY" "http://localhost:8000/api/query?title=题目内容"

# 通过查询参数传递API密钥
curl "http://localhost:8000/api/query?api_key=YOUR_API_KEY&title=题目内容"
```

## 配置说明

- `host`: 服务器监听地址
- `port`: 服务器监听端口
- `platform`: 使用的AI平台 (siliconflow 或 aliyun)
- `api_keys`: 各平台的API密钥
- `models`: 各平台使用的模型
- `mysql`: MySQL数据库配置

## 平台切换

在配置文件中修改 `platform` 字段：
- `siliconflow`: 使用硅基流动平台
- `aliyun`: 使用阿里云百炼平台

## 生成API密钥

- 接口地址：`/api/generate-key`
- 请求方法：`POST`
- 请求示例：

```bash
curl -X POST http://localhost:8000/api/generate-key
```

## 安全说明

API接口现在需要通过API密钥进行身份验证。请确保：

1. 为每个客户端生成唯一的API密钥
2. 通过请求头或查询参数传递API密钥
3. 妥善保管API密钥，避免泄露
4. 定期更换API密钥以确保安全