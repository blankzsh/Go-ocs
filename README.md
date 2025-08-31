# AI-OCS Go版本

AI题库系统Go语言重构版本，支持多种AI平台API调用。

## 功能特性

- 题目答案自动生成
- 支持多种AI平台（硅基流动、阿里云百炼）
- 本地数据库缓存答案
- RESTful API接口

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
- `platform`: 使用的AI平台 (siliconflow 或 aliyun)
- `api_keys`: 各平台的API密钥
- `models`: 各平台使用的模型
- `mysql`: MySQL数据库配置

## 平台切换

在配置文件中修改 `platform` 字段：
- `siliconflow`: 使用硅基流动平台
- `aliyun`: 使用阿里云百炼平台