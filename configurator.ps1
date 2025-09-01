# Go-ocsBase AI Platform Configuration Tool

# Function to show main menu
function Show-MainMenu {
    Write-Host "==========================================" -ForegroundColor Green
    Write-Host "   Go-ocsBase AI Platform Configuration Tool" -ForegroundColor Green
    Write-Host "==========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Please select an operation:" -ForegroundColor Cyan
    Write-Host "1. Configure API Platform" -ForegroundColor White
    Write-Host "2. Configure API Keys" -ForegroundColor White
    Write-Host "3. Configure Models" -ForegroundColor White
    Write-Host "4. Configure MySQL Database" -ForegroundColor White
    Write-Host "5. View Current Configuration" -ForegroundColor White
    Write-Host "6. Exit" -ForegroundColor White
    Write-Host ""
}

# Function to configure platform
function Configure-Platform {
    Write-Host ""
    Write-Host "Available AI Platforms:" -ForegroundColor Cyan
    Write-Host "1. siliconflow (Silicon Flow)" -ForegroundColor White
    Write-Host "2. aliyun (Aliyun Bailian)" -ForegroundColor White
    Write-Host "3. zhipu (Zhipu AI)" -ForegroundColor White
    Write-Host "4. ollama (Ollama Local Model)" -ForegroundColor White
    Write-Host "5. deepseek (DeepSeek Official API)" -ForegroundColor White
    Write-Host "6. chatgpt (ChatGPT API)" -ForegroundColor White
    Write-Host "7. gemini (Gemini API)" -ForegroundColor White
    Write-Host ""

    $platformChoice = Read-Host "Please select a platform (1-7)"

    switch ($platformChoice) {
        "1" { $platform = "siliconflow" }
        "2" { $platform = "aliyun" }
        "3" { $platform = "zhipu" }
        "4" { $platform = "ollama" }
        "5" { $platform = "deepseek" }
        "6" { $platform = "chatgpt" }
        "7" { $platform = "gemini" }
        default { 
            Write-Host "Invalid option" -ForegroundColor Red
            return
        }
    }

    # Read config file and update platform
    $config = Get-Content "configs\config.json" -Raw | ConvertFrom-Json
    $config.platform = $platform
    $config | ConvertTo-Json -Depth 10 | Set-Content "configs\config.json" -Encoding UTF8

    Write-Host ""
    Write-Host "Platform has been set to: $platform" -ForegroundColor Green
}

# Function to configure API keys
function Configure-ApiKeys {
    Write-Host ""
    Write-Host "Please select an API key to configure:" -ForegroundColor Cyan
    Write-Host "1. siliconflow" -ForegroundColor White
    Write-Host "2. aliyun" -ForegroundColor White
    Write-Host "3. zhipu" -ForegroundColor White
    Write-Host "4. deepseek" -ForegroundColor White
    Write-Host "5. chatgpt" -ForegroundColor White
    Write-Host "6. gemini" -ForegroundColor White
    Write-Host ""

    $keyChoice = Read-Host "Please select (1-6)"

    switch ($keyChoice) {
        "1" { $keyName = "siliconflow" }
        "2" { $keyName = "aliyun" }
        "3" { $keyName = "zhipu" }
        "4" { $keyName = "deepseek" }
        "5" { $keyName = "chatgpt" }
        "6" { $keyName = "gemini" }
        default { 
            Write-Host "Invalid option" -ForegroundColor Red
            return
        }
    }

    $apiKey = Read-Host "Please enter the API key for $keyName"
    if ([string]::IsNullOrEmpty($apiKey)) {
        Write-Host "API key cannot be empty" -ForegroundColor Red
        return
    }

    # Read config file and update API key
    $config = Get-Content "configs\config.json" -Raw | ConvertFrom-Json
    $config.api_keys.$keyName = $apiKey
    $config | ConvertTo-Json -Depth 10 | Set-Content "configs\config.json" -Encoding UTF8

    Write-Host ""
    Write-Host "API key for $keyName has been updated" -ForegroundColor Green
}

# Function to configure models
function Configure-Models {
    Write-Host ""
    Write-Host "Please select a model to configure:" -ForegroundColor Cyan
    Write-Host "1. siliconflow (Current: deepseek-ai/DeepSeek-R1)" -ForegroundColor White
    Write-Host "2. aliyun (Current: qwen-plus-latest)" -ForegroundColor White
    Write-Host "3. zhipu (Current: glm-4)" -ForegroundColor White
    Write-Host "4. ollama (Current: llama3)" -ForegroundColor White
    Write-Host "5. deepseek (Current: deepseek-chat)" -ForegroundColor White
    Write-Host "6. chatgpt (Current: gpt-3.5-turbo)" -ForegroundColor White
    Write-Host "7. gemini (Current: gemini-pro)" -ForegroundColor White
    Write-Host ""

    $modelChoice = Read-Host "Please select (1-7)"

    switch ($modelChoice) {
        "1" { 
            $modelName = "siliconflow"
            Write-Host "Current recommended model: deepseek-ai/DeepSeek-R1" -ForegroundColor Yellow
        }
        "2" { 
            $modelName = "aliyun"
            Write-Host "Current recommended model: qwen-plus-latest" -ForegroundColor Yellow
        }
        "3" { 
            $modelName = "zhipu"
            Write-Host "Current recommended model: glm-4" -ForegroundColor Yellow
        }
        "4" { 
            $modelName = "ollama"
            Write-Host "Current recommended model: llama3" -ForegroundColor Yellow
        }
        "5" { 
            $modelName = "deepseek"
            Write-Host "Current recommended model: deepseek-chat" -ForegroundColor Yellow
        }
        "6" { 
            $modelName = "chatgpt"
            Write-Host "Current recommended model: gpt-3.5-turbo" -ForegroundColor Yellow
        }
        "7" { 
            $modelName = "gemini"
            Write-Host "Current recommended model: gemini-pro" -ForegroundColor Yellow
        }
        default { 
            Write-Host "Invalid option" -ForegroundColor Red
            return
        }
    }

    $modelValue = Read-Host "Please enter the model name"
    if ([string]::IsNullOrEmpty($modelValue)) {
        Write-Host "Model name cannot be empty" -ForegroundColor Red
        return
    }

    # Read config file and update model
    $config = Get-Content "configs\config.json" -Raw | ConvertFrom-Json
    $config.models.$modelName = $modelValue
    $config | ConvertTo-Json -Depth 10 | Set-Content "configs\config.json" -Encoding UTF8

    Write-Host ""
    Write-Host "Model for $modelName has been updated to: $modelValue" -ForegroundColor Green
}

# Function to configure MySQL
function Configure-MySQL {
    Write-Host ""
    Write-Host "MySQL Database Configuration:" -ForegroundColor Cyan

    $mysqlHost = Read-Host "Host (Default: 127.0.0.1)"
    if ([string]::IsNullOrEmpty($mysqlHost)) { $mysqlHost = "127.0.0.1" }

    $mysqlPort = Read-Host "Port (Default: 3306)"
    if ([string]::IsNullOrEmpty($mysqlPort)) { $mysqlPort = 3306 }

    $mysqlUser = Read-Host "User (Default: root)"
    if ([string]::IsNullOrEmpty($mysqlUser)) { $mysqlUser = "root" }

    $mysqlPassword = Read-Host "Password"
    if ([string]::IsNullOrEmpty($mysqlPassword)) {
        Write-Host "Password cannot be empty" -ForegroundColor Red
        return
    }

    $mysqlDatabase = Read-Host "Database (Default: question_bank)"
    if ([string]::IsNullOrEmpty($mysqlDatabase)) { $mysqlDatabase = "question_bank" }

    # Read config file and update MySQL configuration
    $config = Get-Content "configs\config.json" -Raw | ConvertFrom-Json
    $config.mysql.host = $mysqlHost
    $config.mysql.port = [int]$mysqlPort
    $config.mysql.user = $mysqlUser
    $config.mysql.password = $mysqlPassword
    $config.mysql.database = $mysqlDatabase
    $config | ConvertTo-Json -Depth 10 | Set-Content "configs\config.json" -Encoding UTF8

    Write-Host ""
    Write-Host "MySQL configuration has been updated" -ForegroundColor Green
}

# Function to view configuration
function View-Config {
    Write-Host ""
    Write-Host "Current Configuration:" -ForegroundColor Cyan
    Write-Host "=================" -ForegroundColor Cyan
    Get-Content "configs\config.json"
    Write-Host ""
    Write-Host "=================" -ForegroundColor Cyan
}

# Main program
Write-Host "==========================================" -ForegroundColor Green
Write-Host "   Go-ocsBase AI Platform Configuration Tool" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green

# Check if config file exists
if (-not (Test-Path "configs\config.json")) {
    Write-Host ""
    Write-Host "Config file not found, copying from example config..." -ForegroundColor Yellow
    Copy-Item "configs\config.example.json" "configs\config.json"
    Write-Host "Config file created: configs/config.json" -ForegroundColor Green
}

# Main loop
while ($true) {
    Show-MainMenu
    $choice = Read-Host "Please enter your choice (1-6)"

    switch ($choice) {
        "1" { Configure-Platform }
        "2" { Configure-ApiKeys }
        "3" { Configure-Models }
        "4" { Configure-MySQL }
        "5" { View-Config }
        "6" { 
            Write-Host ""
            Write-Host "Thank you for using Go-ocsBase Configuration Tool!" -ForegroundColor Green
            exit 0
        }
        default { Write-Host "Invalid option, please try again" -ForegroundColor Red }
    }

    Write-Host ""
    $continue = Read-Host "Continue configuration? (y/n)"
    if ($continue -ne "y" -and $continue -ne "Y") {
        Write-Host ""
        Write-Host "Thank you for using Go-ocsBase Configuration Tool!" -ForegroundColor Green
        exit 0
    }
    Clear-Host
}