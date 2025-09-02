package handlers

import (
	"ai-ocs/internal/database"
	"ai-ocs/internal/models"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// 使用更强的会话密钥（实际应用中应从环境变量或配置文件读取）
var store = sessions.NewCookieStore([]byte("this-is-a-very-strong-session-key-with-32-bytes-long"))

// AdminStats 管理后台统计数据
type AdminStats struct {
	TotalQuestions int64     `json:"total_questions"`
	LastUpdated    time.Time `json:"last_updated"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录处理
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 从配置文件读取管理员账户信息
	configPath := filepath.Join("configs", "config.json")
	config, err := models.LoadConfig(configPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法加载配置文件"})
		return
	}

	// 验证用户名
	if req.Username != config.Admin.Username {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
			"code":  1,
		})
		return
	}

	// 验证密码（使用bcrypt验证）
	err = bcrypt.CompareHashAndPassword([]byte(config.Admin.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
			"code":  1,
		})
		return
	}

	// 创建会话
	session, _ := store.Get(c.Request, "admin-session")
	session.Values["authenticated"] = true
	session.Values["username"] = req.Username
	
	// 生成随机session ID
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成会话"})
		return
	}
	session.Values["session_id"] = hex.EncodeToString(bytes)
	
	// 实施更严格的会话管理策略
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1小时
		HttpOnly: true,
		Secure:   false, // 在生产环境中应设为true（需要HTTPS）
		SameSite: http.SameSiteStrictMode,
	}
	
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存会话"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"code":    0,
	})
}

// Logout 登出处理
func Logout(c *gin.Context) {
	session, _ := store.Get(c.Request, "admin-session")
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)
	
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// RequireAuth 中间件，验证是否已登录
func RequireAuth(c *gin.Context) {
	session, _ := store.Get(c.Request, "admin-session")
	
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// 如果没有登录，重定向到登录页面
		c.Redirect(http.StatusFound, "/admin/login")
		c.Abort()
		return
	}
	
	c.Next()
}

// GetStats 获取统计数据
func GetStats(c *gin.Context) {
	db := database.GetDB()
	
	var count int64
	var lastUpdated time.Time
	
	// 查询总题目数
	err := db.QueryRow("SELECT COUNT(*), MAX(created_at) FROM question_answer").Scan(&count, &lastUpdated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法获取统计数据: " + err.Error(),
		})
		return
	}
	
	stats := AdminStats{
		TotalQuestions: count,
		LastUpdated:    lastUpdated,
	}
	
	c.JSON(http.StatusOK, stats)
}

// GetQuestions 获取所有题目和答案（支持分页）
func GetQuestions(c *gin.Context) {
	db := database.GetDB()
	
	// 获取分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}
	
	offset := (page - 1) * limit
	
	// 查询总数
	var total int64
	err = db.QueryRow("SELECT COUNT(*) FROM question_answer").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法获取题目总数: " + err.Error(),
		})
		return
	}
	
	// 查询题目和答案（分页）
	rows, err := db.Query("SELECT id, question, answer, options, type, created_at FROM question_answer ORDER BY created_at DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法获取题目数据: " + err.Error(),
		})
		return
	}
	defer rows.Close()
	
	type QuestionAnswer struct {
		ID        int64     `json:"id"`
		Question  string    `json:"question"`
		Answer    string    `json:"answer"`
		Options   *string   `json:"options,omitempty"`
		Type      *string   `json:"type,omitempty"`
		CreatedAt time.Time `json:"created_at"`
	}
	
	var results []QuestionAnswer
	
	for rows.Next() {
		var qa QuestionAnswer
		var options, qtype *string
		
		err := rows.Scan(&qa.ID, &qa.Question, &qa.Answer, &options, &qtype, &qa.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "扫描数据时出错: " + err.Error(),
			})
			return
		}
		
		qa.Options = options
		qa.Type = qtype
		results = append(results, qa)
	}
	
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "遍历数据时出错: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// SearchQuestion 搜索特定题目
func SearchQuestion(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请提供搜索关键词",
		})
		return
	}
	
	db := database.GetDB()
	
	// 获取分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}
	
	offset := (page - 1) * limit
	
	// 模糊搜索题目
	rows, err := db.Query("SELECT id, question, answer, options, type, created_at FROM question_answer WHERE question LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?", "%"+keyword+"%", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "搜索时出错: " + err.Error(),
		})
		return
	}
	defer rows.Close()
	
	type QuestionAnswer struct {
		ID        int64     `json:"id"`
		Question  string    `json:"question"`
		Answer    string    `json:"answer"`
		Options   *string   `json:"options,omitempty"`
		Type      *string   `json:"type,omitempty"`
		CreatedAt time.Time `json:"created_at"`
	}
	
	var results []QuestionAnswer
	
	for rows.Next() {
		var qa QuestionAnswer
		var options, qtype *string
		
		err := rows.Scan(&qa.ID, &qa.Question, &qa.Answer, &options, &qtype, &qa.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "扫描数据时出错: " + err.Error(),
			})
			return
		}
		
		qa.Options = options
		qa.Type = qtype
		results = append(results, qa)
	}
	
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "遍历数据时出错: " + err.Error(),
		})
		return
	}
	
	// 获取搜索结果总数
	var total int64
	err = db.QueryRow("SELECT COUNT(*) FROM question_answer WHERE question LIKE ?", "%"+keyword+"%").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法获取搜索结果总数: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// LoginPage 返回管理后台登录页面
func LoginPage(c *gin.Context) {
	// 返回内嵌的管理登录页面
	html := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OCS网课助手AI题库 - 管理登录</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            margin: 0;
            padding: 0;
            height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
        }
        .login-container {
            background: white;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
            width: 100%;
            max-width: 400px;
        }
        .login-header {
            text-align: center;
            margin-bottom: 30px;
        }
        .login-header h1 {
            color: #333;
            margin-bottom: 10px;
        }
        .login-header p {
            color: #666;
            margin: 0;
        }
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            color: #333;
            font-weight: 500;
        }
        .form-group input {
            width: 100%;
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 5px;
            font-size: 16px;
            box-sizing: border-box;
        }
        .form-group input:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
        }
        .btn {
            width: 100%;
            padding: 12px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 5px;
            font-size: 16px;
            cursor: pointer;
            transition: transform 0.2s;
        }
        .btn:hover {
            transform: translateY(-2px);
        }
        .btn:active {
            transform: translateY(0);
        }
        .error-message {
            color: #e74c3c;
            background-color: #fdf2f2;
            padding: 10px;
            border-radius: 5px;
            margin-bottom: 20px;
            display: none;
        }
        .footer {
            text-align: center;
            margin-top: 20px;
            color: #666;
            font-size: 14px;
        }
        .footer a {
            color: #667eea;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <div class="login-header">
            <h1>OCS网课助手AI题库</h1>
            <p>管理后台登录</p>
        </div>
        <div class="error-message" id="errorMessage"></div>
        <form id="loginForm">
            <div class="form-group">
                <label for="username">用户名</label>
                <input type="text" id="username" name="username" required>
            </div>
            <div class="form-group">
                <label for="password">密码</label>
                <input type="password" id="password" name="password" required>
            </div>
            <button type="submit" class="btn">登录</button>
        </form>
        <div class="footer">
            <p>默认账号: admin / admin123</p>
        </div>
    </div>

    <script>
        document.getElementById('loginForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const errorMessage = document.getElementById('errorMessage');
            
            // 隐藏错误消息
            errorMessage.style.display = 'none';
            
            // 发送登录请求
            fetch('/admin/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    password: password
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.code === 0) {
                    // 登录成功，跳转到管理页面
                    window.location.href = '/admin/';
                } else {
                    // 显示错误消息
                    errorMessage.textContent = data.error;
                    errorMessage.style.display = 'block';
                }
            })
            .catch(error => {
                errorMessage.textContent = '登录请求失败，请稍后重试';
                errorMessage.style.display = 'block';
            });
        });
    </script>
</body>
</html>
`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AdminPage 返回管理后台页面
func AdminPage(c *gin.Context) {
	// 返回内嵌的管理页面
	html := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OCS网课助手AI题库 - 管理后台</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 20px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .logout-btn {
            background: rgba(255, 255, 255, 0.2);
            color: white;
            border: 1px solid rgba(255, 255, 255, 0.3);
            padding: 8px 16px;
            border-radius: 4px;
            cursor: pointer;
            transition: background 0.3s;
        }
        .logout-btn:hover {
            background: rgba(255, 255, 255, 0.3);
        }
        .container {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        .stats-container {
            display: flex;
            justify-content: space-between;
            flex-wrap: wrap;
            gap: 20px;
            margin-bottom: 20px;
        }
        .stat-card {
            background: #f8f9fa;
            border-radius: 8px;
            padding: 20px;
            flex: 1;
            min-width: 200px;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        .stat-number {
            font-size: 2em;
            font-weight: bold;
            color: #667eea;
        }
        .stat-label {
            color: #6c757d;
            margin-top: 5px;
        }
        .search-container {
            margin: 20px 0;
        }
        .search-box {
            display: flex;
            gap: 10px;
        }
        input[type="text"] {
            flex: 1;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
        }
        button {
            background: #667eea;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            transition: background 0.3s;
        }
        button:hover {
            background: #5a6fd8;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #f8f9fa;
            font-weight: bold;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .question-text {
            font-weight: bold;
        }
        .answer-text {
            color: #666;
        }
        .loading {
            text-align: center;
            padding: 20px;
            color: #666;
        }
        .error {
            background-color: #f8d7da;
            color: #721c24;
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .hidden {
            display: none;
        }
        .pagination {
            display: flex;
            justify-content: center;
            gap: 10px;
            margin-top: 20px;
        }
        .pagination button {
            padding: 8px 12px;
        }
        .pagination .current {
            background: #5a6fd8;
        }
        @media (max-width: 768px) {
            .stats-container {
                flex-direction: column;
            }
            .search-box {
                flex-direction: column;
            }
            th, td {
                padding: 8px;
                font-size: 14px;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <div>
            <h1>OCS网课助手AI题库</h1>
            <p>管理后台 - 题目统计与查询</p>
        </div>
        <button class="logout-btn" onclick="logout()">退出登录</button>
    </div>

    <div class="container">
        <h2>系统统计</h2>
        <div class="stats-container">
            <div class="stat-card">
                <div class="stat-number" id="totalQuestions">0</div>
                <div class="stat-label">总题目数</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="lastUpdated">-</div>
                <div class="stat-label">最后更新</div>
            </div>
        </div>
    </div>

    <div class="container">
        <h2>题目查询</h2>
        <div class="search-container">
            <div class="search-box">
                <input type="text" id="searchKeyword" placeholder="输入关键词搜索题目...">
                <button onclick="searchQuestions()">搜索</button>
                <button onclick="loadAllQuestions()">显示全部</button>
            </div>
        </div>

        <div id="loading" class="loading hidden">加载中...</div>
        <div id="error" class="error hidden"></div>
        
        <table id="questionsTable">
            <thead>
                <tr>
                    <th>ID</th>
                    <th>题目</th>
                    <th>答案</th>
                    <th>创建时间</th>
                </tr>
            </thead>
            <tbody id="questionsBody">
                <!-- 数据将通过JavaScript动态加载 -->
            </tbody>
        </table>
        
        <div class="pagination" id="pagination">
            <!-- 分页控件将通过JavaScript动态加载 -->
        </div>
    </div>

    <script>
        let currentPage = 1;
        let currentKeyword = '';
        const limit = 50;

        // 页面加载完成后获取统计数据
        document.addEventListener('DOMContentLoaded', function() {
            loadStats();
            loadAllQuestions();
        });

        // 获取统计数据
        function loadStats() {
            fetch('/admin/stats')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('totalQuestions').textContent = data.total_questions;
                    if (data.last_updated) {
                        const date = new Date(data.last_updated);
                        document.getElementById('lastUpdated').textContent = date.toLocaleString('zh-CN');
                    } else {
                        document.getElementById('lastUpdated').textContent = '暂无数据';
                    }
                })
                .catch(error => {
                    console.error('获取统计数据失败:', error);
                });
        }

        // 加载所有题目
        function loadAllQuestions(page = 1) {
            currentKeyword = '';
            currentPage = page;
            showLoading();
            hideError();
            
            fetch('/admin/questions?page=' + page + '&limit=' + limit)
                .then(response => response.json())
                .then(data => {
                    hideLoading();
                    renderQuestions(data.data);
                    renderPagination(data.total, data.page, data.limit);
                })
                .catch(error => {
                    hideLoading();
                    showError('加载题目失败: ' + error.message);
                });
        }

        // 搜索题目
        function searchQuestions(page = 1) {
            const keyword = document.getElementById('searchKeyword').value.trim();
            if (!keyword) {
                loadAllQuestions(page);
                return;
            }
            
            currentKeyword = keyword;
            currentPage = page;
            showLoading();
            hideError();
            
            fetch('/admin/search?keyword=' + encodeURIComponent(keyword) + '&page=' + page + '&limit=' + limit)
                .then(response => response.json())
                .then(data => {
                    hideLoading();
                    renderQuestions(data.data);
                    renderPagination(data.total, data.page, data.limit);
                })
                .catch(error => {
                    hideLoading();
                    showError('搜索失败: ' + error.message);
                });
        }

        // 渲染题目列表
        function renderQuestions(questions) {
            const tbody = document.getElementById('questionsBody');
            tbody.innerHTML = '';
            
            if (questions.length === 0) {
                const row = tbody.insertRow();
                const cell = row.insertCell(0);
                cell.colSpan = 4;
                cell.textContent = '未找到相关题目';
                cell.style.textAlign = 'center';
                cell.style.padding = '20px';
                return;
            }
            
            questions.forEach(question => {
                const row = tbody.insertRow();
                
                const idCell = row.insertCell(0);
                idCell.textContent = question.id;
                
                const questionCell = row.insertCell(1);
                questionCell.innerHTML = '<div class="question-text">' + escapeHtml(question.question) + '</div>';
                if (question.options) {
                    questionCell.innerHTML += '<div style="margin-top: 5px; font-size: 0.9em; color: #666;">选项: ' + escapeHtml(question.options) + '</div>';
                }
                if (question.type) {
                    questionCell.innerHTML += '<div style="margin-top: 3px; font-size: 0.8em; color: #999;">类型: ' + escapeHtml(question.type) + '</div>';
                }
                
                const answerCell = row.insertCell(2);
                answerCell.innerHTML = '<div class="answer-text">' + escapeHtml(question.answer) + '</div>';
                
                const dateCell = row.insertCell(3);
                if (question.created_at) {
                    const date = new Date(question.created_at);
                    dateCell.textContent = date.toLocaleString('zh-CN');
                } else {
                    dateCell.textContent = '-';
                }
            });
        }

        // 渲染分页控件
        function renderPagination(total, currentPage, limit) {
            const totalPages = Math.ceil(total / limit);
            const pagination = document.getElementById('pagination');
            pagination.innerHTML = '';
            
            if (totalPages <= 1) {
                return;
            }
            
            // 上一页按钮
            if (currentPage > 1) {
                const prevButton = document.createElement('button');
                prevButton.textContent = '上一页';
                prevButton.onclick = () => {
                    if (currentKeyword) {
                        searchQuestions(currentPage - 1);
                    } else {
                        loadAllQuestions(currentPage - 1);
                    }
                };
                pagination.appendChild(prevButton);
            }
            
            // 页码按钮
            const startPage = Math.max(1, currentPage - 2);
            const endPage = Math.min(totalPages, currentPage + 2);
            
            for (let i = startPage; i <= endPage; i++) {
                const pageButton = document.createElement('button');
                pageButton.textContent = i;
                if (i === currentPage) {
                    pageButton.classList.add('current');
                }
                pageButton.onclick = () => {
                    if (currentKeyword) {
                        searchQuestions(i);
                    } else {
                        loadAllQuestions(i);
                    }
                };
                pagination.appendChild(pageButton);
            }
            
            // 下一页按钮
            if (currentPage < totalPages) {
                const nextButton = document.createElement('button');
                nextButton.textContent = '下一页';
                nextButton.onclick = () => {
                    if (currentKeyword) {
                        searchQuestions(currentPage + 1);
                    } else {
                        loadAllQuestions(currentPage + 1);
                    }
                };
                pagination.appendChild(nextButton);
            }
        }

        // 显示加载状态
        function showLoading() {
            document.getElementById('loading').classList.remove('hidden');
        }

        // 隐藏加载状态
        function hideLoading() {
            document.getElementById('loading').classList.add('hidden');
        }

        // 显示错误信息
        function showError(message) {
            const errorElement = document.getElementById('error');
            errorElement.textContent = message;
            errorElement.classList.remove('hidden');
        }

        // 隐藏错误信息
        function hideError() {
            document.getElementById('error').classList.add('hidden');
        }

        // 转义HTML特殊字符
        function escapeHtml(text) {
            if (!text) return '';
            const map = {
                '&': '&amp;',
                '<': '&lt;',
                '>': '&gt;',
                '"': '&quot;',
                "'": '&#039;'
            };
            return text.replace(/[&<>"']/g, function(m) { return map[m]; });
        }

        // 回车键触发搜索
        document.getElementById('searchKeyword').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                searchQuestions(1);
            }
        });

        // 退出登录
        function logout() {
            fetch('/admin/logout', {
                method: 'POST'
            })
            .then(() => {
                window.location.href = '/admin/login';
            })
            .catch(error => {
                console.error('登出失败:', error);
                window.location.href = '/admin/login';
            });
        }
    </script>
</body>
</html>
`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}