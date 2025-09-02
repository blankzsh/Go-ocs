package database

import (
	"database/sql"
	"fmt"
	"log"
	"ai-ocs/internal/models"
	"crypto/rand"
	"encoding/hex"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbType string

// InitDB 初始化数据库连接
func InitDB(config *models.Config) error {
	var err error
	
	// 根据配置选择数据库类型
	switch config.DatabaseType {
	case "sqlite":
		db, err = initSQLite(config.SQLiteConfig)
		dbType = "sqlite"
	case "mysql":
		fallthrough
	default:
		db, err = initMySQL(config.MySQLConfig)
		dbType = "mysql"
	}
	
	if err != nil {
		return fmt.Errorf("无法初始化数据库: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("无法连接到数据库: %v", err)
	}

	// 初始化表结构
	if err := initSchema(); err != nil {
		return fmt.Errorf("初始化数据库表失败: %v", err)
	}
	
	// 检查并修复api_keys表结构
	if err := checkAndFixAPIKeyTable(); err != nil {
		return fmt.Errorf("检查并修复api_keys表失败: %v", err)
	}

	// 确保至少有一个API密钥存在
	if err := ensureAPIKey(); err != nil {
		return fmt.Errorf("初始化API密钥失败: %v", err)
	}

	log.Printf("数据库连接成功，使用 %s 数据库", dbType)
	return nil
}

// initMySQL 初始化MySQL数据库连接
func initMySQL(config models.MySQLConfig) (*sql.DB, error) {
	// 构建MySQL连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("无法打开MySQL数据库: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * 60)

	return db, nil
}

// initSQLite 初始化SQLite数据库连接
func initSQLite(config models.SQLiteConfig) (*sql.DB, error) {
	// 构建SQLite连接字符串
	dsn := config.Path

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("无法打开SQLite数据库: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * 60)

	return db, nil
}

// initSchema 初始化数据库表结构
func initSchema() error {
	var createTableSQL string
	var createIndexSQL string
	var createAPIKeyTableSQL string
	var createAPIKeyUsageTableSQL string
	
	// 根据数据库类型选择合适的SQL语法
	if dbType == "sqlite" {
		// SQLite语法
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS question_answer (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			question TEXT NOT NULL,
			answer TEXT NOT NULL,
			options TEXT,
			type TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
		
		createIndexSQL = `CREATE INDEX IF NOT EXISTS idx_question ON question_answer(question);`
		
		createAPIKeyTableSQL = `
		CREATE TABLE IF NOT EXISTS api_keys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			api_key TEXT NOT NULL UNIQUE,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
		
		createAPIKeyUsageTableSQL = `
		CREATE TABLE IF NOT EXISTS api_key_usage (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			api_key_id INTEGER NOT NULL,
			call_count INTEGER DEFAULT 0,
			last_used_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (api_key_id) REFERENCES api_keys(id) ON DELETE CASCADE
		);`
	} else {
		// MySQL语法
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS question_answer (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
			question TEXT NOT NULL,
			answer TEXT NOT NULL,
			options TEXT,
			type TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
		
		// 为MySQL的TEXT字段指定索引长度
		createIndexSQL = `CREATE INDEX idx_question ON question_answer(question(255));`
		
		createAPIKeyTableSQL = `
		CREATE TABLE IF NOT EXISTS api_keys (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
			api_key VARCHAR(64) NOT NULL UNIQUE,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
		
		createAPIKeyUsageTableSQL = `
		CREATE TABLE IF NOT EXISTS api_key_usage (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
			api_key_id INTEGER NOT NULL,
			call_count INTEGER DEFAULT 0,
			last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (api_key_id) REFERENCES api_keys(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	}

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	_, err = db.Exec(createIndexSQL)
	if err != nil {
		// 如果索引已存在，忽略错误
		log.Printf("创建索引时出现警告（可能已存在）: %v", err)
	}
	
	// 创建API密钥表
	_, err = db.Exec(createAPIKeyTableSQL)
	if err != nil {
		return err
	}
	
	// 创建API密钥使用统计表
	_, err = db.Exec(createAPIKeyUsageTableSQL)
	if err != nil {
		return err
	}

	log.Println("数据库表初始化成功")
	return nil
}

// checkAndFixAPIKeyTable 检查并修复api_keys表结构
func checkAndFixAPIKeyTable() error {
	// 检查api_keys表是否存在api_key字段
	var count int
	query := ""
	
	if dbType == "sqlite" {
		query = "SELECT COUNT(*) FROM pragma_table_info('api_keys') WHERE name='api_key'"
	} else {
		// MySQL
		query = "SELECT COUNT(*) FROM information_schema.columns WHERE table_name='api_keys' AND column_name='api_key'"
	}
	
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	
	// 如果api_key字段不存在，则尝试修复表结构
	if count == 0 {
		log.Println("检测到api_keys表结构不正确，尝试修复...")
		
		if dbType == "sqlite" {
			// SQLite修复方法
			_, err = db.Exec("DROP TABLE IF EXISTS api_keys")
			if err != nil {
				return err
			}
			
			_, err = db.Exec(`
				CREATE TABLE api_keys (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					api_key TEXT NOT NULL UNIQUE,
					description TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);`)
			if err != nil {
				return err
			}
		} else {
			// MySQL修复方法
			_, err = db.Exec("DROP TABLE IF EXISTS api_keys")
			if err != nil {
				return err
			}
			
			_, err = db.Exec(`
				CREATE TABLE api_keys (
					id INTEGER PRIMARY KEY AUTO_INCREMENT,
					api_key VARCHAR(64) NOT NULL UNIQUE,
					description TEXT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)
			if err != nil {
				return err
			}
		}
		
		log.Println("api_keys表结构已修复")
	}
	
	return nil
}

// ensureAPIKey 确保至少有一个API密钥存在
func ensureAPIKey() error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM api_keys").Scan(&count)
	if err != nil {
		return err
	}

	// 如果没有API密钥，则生成一个默认的
	if count == 0 {
		apiKey, err := generateAPIKey()
		if err != nil {
			return err
		}
		
		_, err = db.Exec("INSERT INTO api_keys (api_key, description) VALUES (?, ?)", apiKey, "默认API密钥")
		if err != nil {
			return err
		}
		
		log.Printf("已生成默认API密钥: %s", apiKey)
	}
	
	return nil
}

// generateAPIKey 生成一个新的API密钥
func generateAPIKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetDB 获取数据库连接实例
func GetDB() *sql.DB {
	return db
}

// GetDBType 获取当前数据库类型
func GetDBType() string {
	return dbType
}

// GetAnswer 根据问题查询答案
func GetAnswer(question string) (string, error) {
	var answer string
	err := db.QueryRow("SELECT answer FROM question_answer WHERE question = ? LIMIT 1", question).Scan(&answer)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // 没有找到答案，返回空字符串而不是错误
		}
		return "", err
	}
	return answer, nil
}

// SaveAnswer 保存问题和答案到数据库
func SaveAnswer(question, answer string) error {
	// 检查问题是否已存在
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM question_answer WHERE question = ?", question).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// 如果问题已存在，则更新答案
		_, err = db.Exec("UPDATE question_answer SET answer = ? WHERE question = ?", answer, question)
	} else {
		// 如果问题不存在，则插入新记录
		_, err = db.Exec("INSERT INTO question_answer (question, answer) VALUES (?, ?)", question, answer)
	}
	
	return err
}

// ValidateAPIKey 验证API密钥是否有效
func ValidateAPIKey(apiKey string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM api_keys WHERE api_key = ?", apiKey).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetDefaultAPIKey 获取默认API密钥
func GetDefaultAPIKey() (string, error) {
	var apiKey string
	err := db.QueryRow("SELECT api_key FROM api_keys ORDER BY created_at ASC LIMIT 1").Scan(&apiKey)
	if err != nil {
		return "", err
	}
	return apiKey, nil
}

// GetAllAPIKeys 获取所有API密钥
func GetAllAPIKeys() ([]*models.APIKey, error) {
	rows, err := db.Query(`
		SELECT 
			k.id, 
			k.api_key, 
			k.description, 
			k.created_at,
			COALESCE(u.call_count, 0) as call_count,
			u.last_used_at
		FROM api_keys k
		LEFT JOIN api_key_usage u ON k.id = u.api_key_id
		ORDER BY k.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apiKeys []*models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		var lastUsedAt sql.NullTime
		err := rows.Scan(&apiKey.ID, &apiKey.APIKey, &apiKey.Description, &apiKey.CreatedAt, &apiKey.CallCount, &lastUsedAt)
		if err != nil {
			return nil, err
		}
		
		if lastUsedAt.Valid {
			apiKey.LastUsedAt = lastUsedAt.Time
		} else {
			apiKey.LastUsedAt = apiKey.CreatedAt
		}
		
		apiKeys = append(apiKeys, &apiKey)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return apiKeys, nil
}

// CreateAPIKey 创建新的API密钥
func CreateAPIKey(description string) (*models.APIKey, error) {
	// 生成新的API密钥
	key, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	// 插入到数据库
	result, err := db.Exec("INSERT INTO api_keys (api_key, description) VALUES (?, ?)", key, description)
	if err != nil {
		return nil, err
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 返回新创建的API密钥信息
	return &models.APIKey{
		ID:          id,
		APIKey:      key,
		Description: description,
		CreatedAt:   time.Now(),
	}, nil
}

// DeleteAPIKey 删除指定的API密钥
func DeleteAPIKey(id int64) error {
	// 检查是否是最后一个API密钥
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM api_keys").Scan(&count)
	if err != nil {
		return err
	}

	// 如果只有一个API密钥，不允许删除
	if count <= 1 {
		return fmt.Errorf("不能删除最后一个API密钥")
	}

	// 删除API密钥
	_, err = db.Exec("DELETE FROM api_keys WHERE id = ?", id)
	return err
}

// IncrementAPIKeyUsage 增加API密钥调用次数
func IncrementAPIKeyUsage(apiKey string) error {
	// 获取API密钥ID
	var keyID int64
	err := db.QueryRow("SELECT id FROM api_keys WHERE api_key = ?", apiKey).Scan(&keyID)
	if err != nil {
		return err
	}

	// 检查是否已存在使用记录
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM api_key_usage WHERE api_key_id = ?", keyID).Scan(&count)
	if err != nil {
		return err
	}

	// 根据数据库类型选择合适的时间函数
	timeFunc := "NOW()"
	if dbType == "sqlite" {
		timeFunc = "DATETIME('now')"
	}

	if count > 0 {
		// 如果已存在记录，则更新调用次数和最后使用时间
		_, err = db.Exec(`
			UPDATE api_key_usage 
			SET call_count = call_count + 1, last_used_at = `+timeFunc+`
			WHERE api_key_id = ?
		`, keyID)
	} else {
		// 如果不存在记录，则插入新记录
		_, err = db.Exec(`
			INSERT INTO api_key_usage (api_key_id, call_count, last_used_at)
			VALUES (?, 1, `+timeFunc+`)
		`, keyID)
	}

	return err
}
