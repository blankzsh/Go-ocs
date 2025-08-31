const express = require('express');
const mysql = require('mysql2');
const bcrypt = require('bcryptjs');
const cors = require('cors');
require('dotenv').config();

const app = express();
const PORT = process.env.PORT || 3001;

// 中间件
app.use(cors());
app.use(express.json());

// 创建数据库连接
const db = mysql.createConnection({
  host: process.env.DB_HOST,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME
});

// 连接数据库
db.connect((err) => {
  if (err) {
    console.error('数据库连接失败:', err);
    return;
  }
  console.log('成功连接到MySQL数据库');
  
  // 创建api_keys表
  const createTableQuery = `
    CREATE TABLE IF NOT EXISTS api_keys (
      id INT AUTO_INCREMENT PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      key_value VARCHAR(255) NOT NULL UNIQUE,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
  `;
  
  db.query(createTableQuery, (err) => {
    if (err) {
      console.error('创建表失败:', err);
      return;
    }
    console.log('api_keys表已创建或已存在');
  });
});

// 添加请求日志中间件
app.use((req, res, next) => {
  console.log(`收到请求: ${req.method} ${req.path}`);
  next();
});

// API路由

// 获取所有API密钥（不返回实际密钥值）
app.get('/api/keys', (req, res) => {
  console.log('处理获取API密钥列表请求');
  const query = 'SELECT id, name, created_at FROM api_keys ORDER BY created_at DESC';
  
  db.query(query, (err, results) => {
    if (err) {
      console.error('查询API密钥失败:', err);
      return res.status(500).json({ error: '服务器内部错误' });
    }
    
    console.log('成功获取API密钥列表，数量:', results.length);
    res.json({ keys: results });
  });
});

// 生成新的API密钥
app.post('/api/keys', (req, res) => {
  console.log('处理生成新API密钥请求');
  const { name } = req.body;
  
  // 确保始终有响应
  try {
    if (!name) {
      console.log('密钥名称为空');
      return res.status(400).json({ error: '密钥名称不能为空' });
    }
    
    console.log('生成API密钥，名称:', name);
    // 生成随机API密钥
    const keyValue = 'key_' + Math.random().toString(36).substr(2, 16);
    console.log('生成的API密钥:', keyValue);
    
    // 加密API密钥
    const saltRounds = 10;
    bcrypt.hash(keyValue, saltRounds, (err, hashedKey) => {
      if (err) {
        console.error('加密API密钥失败:', err);
        // 确保响应只发送一次
        if (!res.headersSent) {
          return res.status(500).json({ error: '服务器内部错误' });
        }
        return;
      }
      
      console.log('API密钥加密完成');
      // 存储到数据库
      const query = 'INSERT INTO api_keys (name, key_value) VALUES (?, ?)';
      db.query(query, [name, hashedKey], (err, result) => {
        if (err) {
          console.error('存储API密钥失败:', err);
          // 确保响应只发送一次
          if (!res.headersSent) {
            return res.status(500).json({ error: '服务器内部错误' });
          }
          return;
        }
        
        console.log('API密钥存储成功，ID:', result.insertId);
        // 确保响应只发送一次
        if (!res.headersSent) {
          res.status(201).json({ 
            id: result.insertId, 
            name, 
            key: keyValue, 
            message: 'API密钥创建成功' 
          });
        }
      });
    });
  } catch (error) {
    console.error('生成API密钥时发生未处理的错误:', error);
    // 确保即使发生未捕获的异常也有响应
    if (!res.headersSent) {
      res.status(500).json({ error: '服务器内部错误' });
    }
  }
});

// 验证API密钥
app.post('/api/keys/validate', (req, res) => {
  console.log('处理API密钥验证请求');
  const { key } = req.body;
  
  try {
    if (!key) {
      console.log('API密钥为空');
      return res.status(400).json({ valid: false, error: 'API密钥不能为空' });
    }
    
    console.log('验证API密钥');
    // 查询所有密钥并逐一验证
    const query = 'SELECT key_value FROM api_keys';
    db.query(query, (err, results) => {
      if (err) {
        console.error('查询API密钥失败:', err);
        // 确保响应只发送一次
        if (!res.headersSent) {
          return res.status(500).json({ valid: false, error: '服务器内部错误' });
        }
        return;
      }
      
      console.log('查询到API密钥数量:', results.length);
      // 验证密钥
      let isValid = false;
      for (const row of results) {
        if (bcrypt.compareSync(key, row.key_value)) {
          isValid = true;
          break;
        }
      }
      
      console.log('API密钥验证结果:', isValid);
      // 确保响应只发送一次
      if (!res.headersSent) {
        res.json({ valid: isValid });
      }
    });
  } catch (error) {
    console.error('验证API密钥时发生未处理的错误:', error);
    // 确保即使发生未捕获的异常也有响应
    if (!res.headersSent) {
      res.status(500).json({ valid: false, error: '服务器内部错误' });
    }
  }
});

// 删除API密钥
app.delete('/api/keys/:id', (req, res) => {
  console.log('处理删除API密钥请求');
  const { id } = req.params;
  
  try {
    console.log('删除API密钥，ID:', id);
    const query = 'DELETE FROM api_keys WHERE id = ?';
    db.query(query, [id], (err, result) => {
      if (err) {
        console.error('删除API密钥失败:', err);
        // 确保响应只发送一次
        if (!res.headersSent) {
          return res.status(500).json({ error: '服务器内部错误' });
        }
        return;
      }
      
      if (result.affectedRows === 0) {
        console.log('未找到要删除的API密钥');
        // 确保响应只发送一次
        if (!res.headersSent) {
          return res.status(404).json({ error: 'API密钥未找到' });
        }
        return;
      }
      
      console.log('API密钥删除成功');
      // 确保响应只发送一次
      if (!res.headersSent) {
        res.json({ message: 'API密钥删除成功' });
      }
    });
  } catch (error) {
    console.error('删除API密钥时发生未处理的错误:', error);
    // 确保即使发生未捕获的异常也有响应
    if (!res.headersSent) {
      res.status(500).json({ error: '服务器内部错误' });
    }
  }
});

// 根节点路由，用于健康检查
app.get('/', (req, res) => {
  console.log('处理根路径请求');
  res.json({ message: 'API密钥管理服务正在运行' });
});

// 404处理
app.use((req, res) => {
  console.log('404 - 未找到路由:', req.path);
  res.status(404).json({ error: '请求的资源未找到' });
});

// 全局错误处理
app.use((err, req, res, next) => {
  console.error('未处理的错误:', err);
  res.status(500).json({ error: '服务器内部错误' });
});

// 启动服务器
const server = app.listen(PORT, () => {
  console.log(`API密钥管理服务器运行在端口 ${PORT}`);
});

// 优雅关闭
process.on('SIGINT', () => {
  console.log('正在关闭服务器...');
  server.close(() => {
    console.log('服务器已关闭');
    process.exit(0);
  });
});