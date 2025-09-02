package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func main() {
	password := "admin123"
	
	// 生成加密后的密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成密码哈希失败:", err)
	}
	
	fmt.Printf("明文密码: %s\n", password)
	fmt.Printf("加密后的密码哈希: %s\n", string(hashedPassword))
	fmt.Println("\n请将此哈希值添加到您的config.json配置文件中的admin.password字段")
}