# API密钥管理系统

这是一个基于Node.js和React的API密钥管理系统，用于生成、存储和验证API密钥。密钥通过MySQL数据库进行存储，并使用bcrypt进行加密。

## 功能特性

- 生成新的API密钥
- 存储和管理API密钥（名称和创建时间）
- 验证API密钥有效性
- 删除API密钥
- 响应式用户界面

## 技术栈

### 后端
- Node.js
- Express.js
- MySQL2
- Bcrypt.js（用于加密）
- Dotenv（环境变量配置）

### 前端
- React（Vite）
- JavaScript
- CSS

## 目录结构

```
api-key/
├── backend/          # 后端服务
│   ├── server.js     # 服务器入口文件
│   ├── .env          # 环境变量配置
│   ├── package.json  # 后端依赖配置
│   └── ...           # 其他后端文件
└── frontend/         # 前端应用
    ├── src/          # 源代码
    ├── package.json  # 前端依赖配置
    └── ...           # 其他前端文件
```

## 安装和运行

### 数据库设置

1. 安装MySQL数据库
2. 创建一个新的数据库，例如：`api_key_manager`
3. 更新 [backend/.env](file:///f:/Project/Go-ocs/api-key/backend/.env) 文件中的数据库配置

### 后端服务

1. 进入后端目录：
   ```
   cd backend
   ```

2. 安装依赖：
   ```
   npm install
   ```

3. 启动服务：
   ```
   npm start
   ```
   
   或者开发模式启动：
   ```
   npm run dev
   ```

### 前端应用

1. 进入前端目录：
   ```
   cd frontend
   ```

2. 安装依赖：
   ```
   npm install
   ```

3. 启动开发服务器：
   ```
   npm run dev
   ```

## API接口

### 获取所有API密钥
- **URL**: `/api/keys`
- **方法**: `GET`
- **响应**: 
  ```json
  {
    "keys": [
      {
        "id": 1,
        "name": "测试密钥",
        "created_at": "2023-01-01T00:00:00.000Z"
      }
    ]
  }
  ```

### 生成新API密钥
- **URL**: `/api/keys`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "name": "密钥名称"
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "name": "密钥名称",
    "key": "生成的API密钥",
    "message": "API密钥创建成功"
  }
  ```

### 验证API密钥
- **URL**: `/api/keys/validate`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "key": "要验证的API密钥"
  }
  ```
- **响应**:
  ```json
  {
    "valid": true
  }
  ```

### 删除API密钥
- **URL**: `/api/keys/:id`
- **方法**: `DELETE`
- **响应**:
  ```json
  {
    "message": "API密钥删除成功"
  }
  ```

## 安全说明

1. API密钥在数据库中使用bcrypt进行加密存储
2. 前端只显示密钥名称和创建时间，不显示实际密钥值
3. 新生成的密钥只在创建时显示一次，刷新页面后无法再次查看
4. 建议使用HTTPS来保护API通信安全

## 注意事项

1. 请确保正确配置数据库连接信息
2. 生产环境中应使用强密码和适当的访问控制
3. 定期备份数据库以防止数据丢失