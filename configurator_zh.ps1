# Go-ocsBase AI平台配置工具

# 显示主菜单函数
function Show-MainMenu {
    Write-Host "==========================================" -ForegroundColor Green
    Write-Host "   Go-ocsBase AI平台配置工具" -ForegroundColor Green
    Write-Host "==========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "请选择操作:" -ForegroundColor Cyan
    Write-Host "1. 配置API平台" -ForegroundColor White
    Write-Host "2. 配置API密钥" -ForegroundColor White
    Write-Host "3. 配置模型" -ForegroundColor White
    Write-Host "4. 配置数据库" -ForegroundColor White
    Write-Host "5. 配置MySQL数据库" -ForegroundColor White
    Write-Host "6. 查看当前配置" -ForegroundColor White
    Write-Host "7. 退出" -ForegroundColor White
    Write-Host ""
}

# 配置平台函数
function Configure-Platform {
    Write-Host ""
    Write-Host "可用的AI平台:" -ForegroundColor Cyan
    Write-Host "1. siliconflow (硅基流动)" -ForegroundColor White
    Write-Host "2. aliyun (阿里云百炼)" -ForegroundColor White
    Write-Host "3. zhipu (智普AI)" -ForegroundColor White
    Write-Host "4. ollama (Ollama本地模型)" -ForegroundColor White
    Write-Host "5. deepseek (DeepSeek官方API)" -ForegroundColor White
    Write-Host "6. chatgpt (ChatGPT API)" -ForegroundColor White
    Write-Host "7. gemini (Gemini API)" -ForegroundColor White
    Write-Host ""

    $platformChoice = Read-Host "请选择平台 (1-7)"

    switch ($platformChoice) {
        "1" { $platform = "siliconflow" }
        "2" { $platform = "aliyun" }
        "3" { $platform = "zhipu" }
        "4" { $platform = "ollama" }
        "5" { $platform = "deepseek" }
        "6" { $platform = "chatgpt" }
        "7" { $platform = "gemini" }
        default { 
            Write-Host "无效选项" -ForegroundColor Red
            return
        }
    }

    # 读取配置文件并更新平台
    $config = Get-Content "configs\config.json" -Raw | ConvertFrom-Json
    $config.platform = $platform
    $config | ConvertTo-Json -Depth 10 | Set-Content "configs\config.json" -Encoding UTF8

    Write-Host ""
    Write-Host "平台已设置为: $platform" -ForegroundColor Green
}

# 配置API密钥函数
function Configure-ApiKeys {
    Write-Host ""
    Write-Host "请选择要配置的API密钥:" -ForegroundColor Cyan
    Write-Host "1. siliconflow" -ForegroundColor White
    Write-Host "2. aliyun" -ForegroundColor White
    Write-Host "3. zhipu" -ForegroundColor White
    Write-Host "4. deepseek" -ForegroundColor White
    Write-Host "5. chatgpt" -ForegroundColor White
    Write-Host "6. gemini" -ForegroundColor White
    Write-Host ""

    $keyChoice = Read-Host "请选择 (1-6)"

    switch ($keyChoice) {
        "1" { $keyName = "siliconflow" }
        "2" { $keyName = "aliyun" }
        "3" { $keyName = "zhipu" }
        "4" { $keyName = "deepseek" }
        "5" { $keyName = "chatgpt" }
        "6" { $keyName = "gemini" }
        default { 
            Write-Host "无效选项" -ForegroundColor Red
            return
        }
    }

    $apiKey = Read-Host "请输入 $keyName 的API密钥"
    if ([string]::IsNullOrEmpty($apiKey)) {
        Write-Host "API密钥不能为空" -ForegroundColor Red
        return
    }

    # 读取配置文件并更新API密钥
    $config = Get-Content "configs\config.json" -Raw | ConvertFrom-Json
    $config.api_keys.$keyName = $apiKey
    $config | ConvertTo-Json -Depth 10 | Set-Content "configs\config.json" -Encoding UTF8

    Write-Host ""
    Write-Host "$keyName 的API密钥已更新" -ForegroundColor Green
}

# 配置模型函数
function Configure-Models {
    Write-Host ""
    Write-Host "请选择要配置的模型:" -ForegroundColor Cyan
    Write-Host "1. siliconflow (当前: deepseek-ai/DeepSeek-R1)" -ForegroundColor White
    Write-Host "2. aliyun (当前: qwen-plus-latest)" -ForegroundColor White
    Write-Host "3. zhipu (当前: glm-4)" -ForegroundColor White
    Write-Host "4. ollama (当前: llama3)" -ForegroundColor White
    Write-Host "5. deepseek (当前: deepseek-chat)" -ForegroundColor White
    Write-Host "6. chatgpt (当前: gpt-3.5-turbo)" -ForegroundColor White
    Write-Host "7. gemini (当前: gemini-pro)" -ForegroundColor White
    Write-Host ""

    $modelChoice = Read-Host "请选择 (1-7)"

    switch ($modelChoice) {
        "1" { 
            $modelName = "siliconflow"
            Write-Host "当前推荐模型: deepseek-ai/DeepSeek-R1" -ForegroundColor Yellow
        }
        "2" { 
            $modelName = "aliyun"
            Write-Host "当前推荐模型: qwen-plus-latest" -ForegroundColor Yellow
        }
        "3" { 
            $modelName = "zhipu"
            Write-Host "当前推荐模型: glm-4" -ForegroundColor Yellow
        }
        "4" { 
            $modelName = "ollama"
            Write-Host "当前推荐模型: llama3" -ForegroundColor Yellow
        }
        "5" { 
            $modelName = "deepseek"
            Write-Host "当前推荐模型: deepseek-chat" -ForegroundColor Yellow
        }
        "6" { 
            $modelName = "chatgpt"
            Write-Host "当前推荐模型: gpt-3.5-turbo" -ForegroundColor Yellow
        }
        "7" { 
            $modelName = "gemini"
            Write-Host "当前推荐模型: gemini-pro" -ForegroundColor Yellow
        }
        default { 
            Write-Host "无效选项" -ForegroundColor Red
            return
        }
    }

    $modelValue = Read-Host "请输入模型名称"
    if ([string]::IsNullOrEmpty($modelValue)) {
        Write-Host "模型名称不能为空" -ForegroundColor Red
        return
    }

    # 读取配置文件并更新模型
    $config = Get-Content "configs\config.json" -Raw | ConvertFrom-Json
    $config.models.$modelName = $modelValue
    $config | ConvertTo-Json -Depth 10 | Set-Content "configs\config.json" -Encoding UTF8

    Write-Host ""
    Write-Host "$modelName 的模型已更新为: $modelValue" -ForegroundColor Green
}

# 配置数据库类型函数
function Configure-Database {
    Write-Host ""
    Write-Host "请选择数据库类型:" -ForegroundColor Cyan
    Write-Host "1. MySQL" -ForegroundColor White
    Write-Host "2. SQLite" -ForegroundColor White
    Write-Host ""

    $dbChoice = Read-Host "请选择 (1-2)"

    switch ($dbChoice) {
        "1" { 
            $dbType = "mysql"
            Write-Host "数据库类型设置为 MySQL" -ForegroundColor Yellow
        }
        "2" { 
            $dbType = "sqlite"
            Write-Host "数据库类型设置为 SQLite" -ForegroundColor Yellow
        }
        default { 
            Write-Host "无效选项" -ForegroundColor Red
            return
        }
    }

    # 读取配置文件并更新数据库类型
    $config = Get-Content "configs\config.json" -Raw | ConvertFrom-Json
    $config.database_type = $dbType
    $config | ConvertTo-Json -Depth 10 | Set-Content "configs\config.json" -Encoding UTF8

    Write-Host ""
    Write-Host "数据库类型已更新为: $dbType" -ForegroundColor Green
}

# 配置MySQL函数
function Configure-MySQL {
    Write-Host ""
    Write-Host "MySQL数据库配置:" -ForegroundColor Cyan

    $mysqlHost = Read-Host "主机地址 (默认: 127.0.0.1)"
    if ([string]::IsNullOrEmpty($mysqlHost)) { $mysqlHost = "127.0.0.1" }

    $mysqlPort = Read-Host "端口 (默认: 3306)"
    if ([string]::IsNullOrEmpty($mysqlPort)) { $mysqlPort = 3306 }

    $mysqlUser = Read-Host "用户名 (默认: root)"
    if ([string]::IsNullOrEmpty($mysqlUser)) { $mysqlUser = "root" }

    $mysqlPassword = Read-Host "密码"
    if ([string]::IsNullOrEmpty($mysqlPassword)) {
        Write-Host "密码不能为空" -ForegroundColor Red
        return
    }

    $mysqlDatabase = Read-Host "数据库名 (默认: question_bank)"
    if ([string]::IsNullOrEmpty($mysqlDatabase)) { $mysqlDatabase = "question_bank" }

    # 读取配置文件并更新MySQL配置
    $config = Get-Content "configs\config.json" -Raw | ConvertFrom-Json
    $config.mysql.host = $mysqlHost
    $config.mysql.port = [int]$mysqlPort
    $config.mysql.user = $mysqlUser
    $config.mysql.password = $mysqlPassword
    $config.mysql.database = $mysqlDatabase
    $config | ConvertTo-Json -Depth 10 | Set-Content "configs\config.json" -Encoding UTF8

    Write-Host ""
    Write-Host "MySQL配置已更新" -ForegroundColor Green
}

# 查看配置函数
function View-Config {
    Write-Host ""
    Write-Host "当前配置:" -ForegroundColor Cyan
    Write-Host "=================" -ForegroundColor Cyan
    Get-Content "configs\config.json"
    Write-Host ""
    Write-Host "=================" -ForegroundColor Cyan
}

# 主程序
Write-Host "==========================================" -ForegroundColor Green
Write-Host "   Go-ocsBase AI平台配置工具" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green

# 检查是否存在配置文件
if (-not (Test-Path "configs\config.json")) {
    Write-Host ""
    Write-Host "检测到尚未创建配置文件，将从示例配置文件复制..." -ForegroundColor Yellow
    Copy-Item "configs\config.example.json" "configs\config.json"
    Write-Host "配置文件已创建: configs/config.json" -ForegroundColor Green
}

# 主循环
while ($true) {
    Show-MainMenu
    $choice = Read-Host "请输入选项 (1-7)"

    switch ($choice) {
        "1" { Configure-Platform }
        "2" { Configure-ApiKeys }
        "3" { Configure-Models }
        "4" { Configure-Database }
        "5" { Configure-MySQL }
        "6" { View-Config }
        "7" { 
            Write-Host ""
            Write-Host "感谢使用Go-ocsBase配置工具！" -ForegroundColor Green
            exit 0
        }
        default { Write-Host "无效选项，请重新选择" -ForegroundColor Red }
    }

    Write-Host ""
    $continue = Read-Host "是否继续配置? (y/n)"
    if ($continue -ne "y" -and $continue -ne "Y") {
        Write-Host ""
        Write-Host "感谢使用Go-ocsBase配置工具！" -ForegroundColor Green
        exit 0
    }
    Clear-Host
}