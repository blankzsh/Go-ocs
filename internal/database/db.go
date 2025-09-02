package database

import (
	"database/sql"
	"fmt"
	"log"
	"ai-ocs/internal/models"

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

	log.Println("数据库表初始化成功")
	return nil
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