# AI-OCS Go版本

AI题库系统Go语言重构版本，支持多种AI平台API调用。

版本：v1.4.0

## 功能特性

- 题目答案自动生成
- 支持多种AI平台（硅基流动、阿里云百炼、智普AI、Ollama、DeepSeek、ChatGPT、Gemini）
- 本地数据库缓存答案（支持MySQL和SQLite）
- RESTful API接口
- 数据库自由切换（MySQL/SQLite）
- 管理后台界面（带登录验证）
- API密钥验证系统
- API密钥管理界面

## 技术栈

- Go 1.21+
- Gin Web框架
- MySQL 数据库 / SQLite 数据库
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
│   ├── models/          # 数据模型
│   └── tools/           # 工具脚本
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
       "database_type": "mysql",
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
       "admin": {
           "username": "admin",
           "password": "$2a$10$XovBFarUSNyp/Ux.DYOwqu/zGKyU3XbVEM6qKhS2U9Nq9WxxNpgk6"
       },
       "mysql": {
           "host": "127.0.0.1",
           "port": 3306,
           "user": "root",
           "password": "your_mysql_password",
           "database": "question_bank"
       },
       "sqlite": {
           "path": "question_bank.db"
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
GET /api/query?title=问题内容[&options=选项内容][&type=问题类型][&api-key=API密钥]
```

示例：
```
curl "http://127.0.0.1:8000/api/query?title=中国的首都是哪里%3F&options=北京###上海###广州###深圳&type=选择题&api-key=生成api-key"
```

注意：**API密钥是必需的，必须提供有效的密钥才能访问API。**

## 管理后台

项目包含一个Web管理后台，可用于查看题目统计、搜索题目等。

访问地址：`http://服务器地址:端口/admin/`

默认管理员账户：
- 用户名：`admin`
- 密码：`admin123`

注意：为安全起见，建议在生产环境中修改默认密码。

### 管理后台功能

- 系统统计信息展示
- 题目列表查看（支持分页）
- 关键词搜索题目
- 会话管理（登录/登出）
- API密钥管理（创建、查看、删除API密钥）

## API密钥系统

为了提高API安全性，系统引入了API密钥验证机制：

1. 系统启动时会自动生成一个默认的API密钥并存储在数据库中
2. 在启动日志中会打印包含API密钥的完整配置信息
3. 调用API时必须使用`api-key`参数传递API密钥进行验证
4. 系统会验证API密钥的有效性；无效的密钥将被拒绝访问
5. 可以通过管理后台界面创建、查看和删除API密钥

## 测试工具

项目提供了一个Python测试工具，用于交互式测试生成的API配置信息是否能正常使用：

### 使用方法

1. 确保已安装Python 3和requests库：
   ```bash
   pip install requests
   ```

2. 启动服务并获取API配置信息：
   ```bash
   go run cmd/main.go
   ```

3. 从启动日志中复制API配置信息

4. 将配置信息保存为 `api_config.json` 文件

5. 运行Python测试工具：
   ```bash
   python test_api_config.py api_config.json
   ```

### 功能特点

- 交互式界面，易于使用
- 支持从文件加载或手动输入API配置
- 可自定义测试参数（题目、选项、题型、API密钥）
- 显示详细的请求和响应信息
- 支持重新加载配置和修改测试参数
- 提供响应数据解析功能

### 命令行参数

运行测试工具时可以指定配置文件路径：
```bash
python test_api_config.py [配置文件路径]
```

## 配置说明

- `host`: 服务器监听地址
- `port`: 服务器监听端口
- `platform`: 使用的AI平台
- `database_type`: 数据库类型（mysql 或 sqlite）
- `api_keys`: 各平台的API密钥
- `models`: 各平台使用的模型
- `admin`: 管理员账户配置
- `mysql`: MySQL数据库配置
- `sqlite`: SQLite数据库配置

## 平台切换

在配置文件中修改 `platform` 字段：
- `siliconflow`: 使用硅基流动平台
- `aliyun`: 使用阿里云百炼平台
- `zhipu`: 使用智普AI平台
- `ollama`: 使用Ollama本地模型
- `deepseek`: 使用DeepSeek官方API
- `chatgpt`: 使用ChatGPT API
- `gemini`: 使用Gemini API

## 数据库切换

在配置文件中修改 `database_type` 字段：
- `mysql`: 使用MySQL数据库
- `sqlite`: 使用SQLite数据库

SQLite数据库无需额外安装，文件会自动创建在配置指定的路径。

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
- 选择数据库类型（MySQL/SQLite）
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

## 安全说明

1. 管理员账户的密码使用bcrypt加密存储，即使在配置文件中也不会明文保存
2. 会话管理采用安全策略，包括HttpOnly标志、SameSite策略等
3. 未登录用户访问管理后台会自动重定向到登录页面
4. API密钥存储在数据库中，提供访问控制机制
5. 建议在生产环境中使用HTTPS协议，并设置相应的安全选项

## 版本历史

详细版本更新信息请查看 [RELEASE_NOTES.md](file:///f:/WEB-PR/Go-ocsBase/RELEASE_NOTES.md) 文件。