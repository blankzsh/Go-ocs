#!/bin/bash

# Go-ocsBase AI平台配置工具

# 配置平台函数
config_platform() {
    echo ""
    echo "可用的AI平台:"
    echo "1. siliconflow (硅基流动)"
    echo "2. aliyun (阿里云百炼)"
    echo "3. zhipu (智普AI)"
    echo "4. ollama (Ollama本地模型)"
    echo "5. deepseek (DeepSeek官方API)"
    echo "6. chatgpt (ChatGPT API)"
    echo "7. gemini (Gemini API)"
    echo ""

    read -p "请选择平台 (1-7): " platform_choice

    case $platform_choice in
        1)
            platform="siliconflow"
            ;;
        2)
            platform="aliyun"
            ;;
        3)
            platform="zhipu"
            ;;
        4)
            platform="ollama"
            ;;
        5)
            platform="deepseek"
            ;;
        6)
            platform="chatgpt"
            ;;
        7)
            platform="gemini"
            ;;
        *)
            echo "无效选项"
            return
            ;;
    esac

    # 使用sed更新平台配置
    sed -i.bak "s/\"platform\": \"[^\"]*\"/\"platform\": \"$platform\"/g" configs/config.json
    rm configs/config.json.bak 2>/dev/null || true
    echo ""
    echo "平台已设置为: $platform"
}

# 配置API密钥函数
config_api_keys() {
    echo ""
    echo "请选择要配置的API密钥:"
    echo "1. siliconflow"
    echo "2. aliyun"
    echo "3. zhipu"
    echo "4. deepseek"
    echo "5. chatgpt"
    echo "6. gemini"
    echo ""

    read -p "请选择 (1-6): " key_choice

    case $key_choice in
        1)
            key_name="siliconflow"
            ;;
        2)
            key_name="aliyun"
            ;;
        3)
            key_name="zhipu"
            ;;
        4)
            key_name="deepseek"
            ;;
        5)
            key_name="chatgpt"
            ;;
        6)
            key_name="gemini"
            ;;
        *)
            echo "无效选项"
            return
            ;;
    esac

    read -p "请输入 $key_name 的API密钥: " api_key
    if [ -z "$api_key" ]; then
        echo "API密钥不能为空"
        return
    fi

    # 使用sed更新API密钥
    sed -i.bak "s/\"$key_name\": \"[^\"]*\"/\"$key_name\": \"$api_key\"/g" configs/config.json
    rm configs/config.json.bak 2>/dev/null || true
    echo ""
    echo "$key_name 的API密钥已更新"
}

# 配置模型函数
config_models() {
    echo ""
    echo "请选择要配置的模型:"
    echo "1. siliconflow (当前: deepseek-ai/DeepSeek-R1)"
    echo "2. aliyun (当前: qwen-plus-latest)"
    echo "3. zhipu (当前: glm-4)"
    echo "4. ollama (当前: llama3)"
    echo "5. deepseek (当前: deepseek-chat)"
    echo "6. chatgpt (当前: gpt-3.5-turbo)"
    echo "7. gemini (当前: gemini-pro)"
    echo ""

    read -p "请选择 (1-7): " model_choice

    case $model_choice in
        1)
            model_name="siliconflow"
            echo "当前推荐模型: deepseek-ai/DeepSeek-R1"
            ;;
        2)
            model_name="aliyun"
            echo "当前推荐模型: qwen-plus-latest"
            ;;
        3)
            model_name="zhipu"
            echo "当前推荐模型: glm-4"
            ;;
        4)
            model_name="ollama"
            echo "当前推荐模型: llama3"
            ;;
        5)
            model_name="deepseek"
            echo "当前推荐模型: deepseek-chat"
            ;;
        6)
            model_name="chatgpt"
            echo "当前推荐模型: gpt-3.5-turbo"
            ;;
        7)
            model_name="gemini"
            echo "当前推荐模型: gemini-pro"
            ;;
        *)
            echo "无效选项"
            return
            ;;
    esac

    read -p "请输入模型名称: " model_value
    if [ -z "$model_value" ]; then
        echo "模型名称不能为空"
        return
    fi

    # 使用sed更新模型配置
    sed -i.bak "s/\"$model_name\": \"[^\"]*\"/\"$model_name\": \"$model_value\"/g" configs/config.json
    rm configs/config.json.bak 2>/dev/null || true
    echo ""
    echo "$model_name 的模型已更新为: $model_value"
}

# 配置数据库类型函数
config_database() {
    echo ""
    echo "请选择数据库类型:"
    echo "1. MySQL"
    echo "2. SQLite"
    echo ""

    read -p "请选择 (1-2): " db_choice

    case $db_choice in
        1)
            db_type="mysql"
            echo "数据库类型设置为 MySQL"
            ;;
        2)
            db_type="sqlite"
            echo "数据库类型设置为 SQLite"
            ;;
        *)
            echo "无效选项"
            return
            ;;
    esac

    # 使用sed更新数据库类型配置
    sed -i.bak "s/\"database_type\": \"[^\"]*\"/\"database_type\": \"$db_type\"/g" configs/config.json
    rm configs/config.json.bak 2>/dev/null || true
    echo ""
    echo "数据库类型已更新为: $db_type"
}

# 配置MySQL函数
config_mysql() {
    echo ""
    echo "MySQL数据库配置:"
    
    read -p "主机地址 (默认: 127.0.0.1): " mysql_host
    if [ -z "$mysql_host" ]; then
        mysql_host="127.0.0.1"
    fi

    read -p "端口 (默认: 3306): " mysql_port
    if [ -z "$mysql_port" ]; then
        mysql_port="3306"
    fi

    read -p "用户名 (默认: root): " mysql_user
    if [ -z "$mysql_user" ]; then
        mysql_user="root"
    fi

    read -s -p "密码: " mysql_password
    echo ""
    if [ -z "$mysql_password" ]; then
        echo "密码不能为空"
        return
    fi

    read -p "数据库名 (默认: question_bank): " mysql_database
    if [ -z "$mysql_database" ]; then
        mysql_database="question_bank"
    fi

    # 更新MySQL配置
    sed -i.bak "s/\"host\": \"[^\"]*\"/\"host\": \"$mysql_host\"/g" configs/config.json
    sed -i.bak "s/\"port\": [0-9]*/\"port\": $mysql_port/g" configs/config.json
    sed -i.bak "s/\"user\": \"[^\"]*\"/\"user\": \"$mysql_user\"/g" configs/config.json
    sed -i.bak "s/\"password\": \"[^\"]*\"/\"password\": \"$mysql_password\"/g" configs/config.json
    sed -i.bak "s/\"database\": \"[^\"]*\"/\"database\": \"$mysql_database\"/g" configs/config.json
    rm configs/config.json.bak 2>/dev/null || true

    echo ""
    echo "MySQL配置已更新"
}

# 查看配置函数
view_config() {
    echo ""
    echo "当前配置:"
    echo "================="
    cat configs/config.json
    echo ""
    echo "================="
}

# 主程序开始
echo "=========================================="
echo "   Go-ocsBase AI平台配置工具"
echo "=========================================="

# 检查是否存在配置文件
if [ ! -f "configs/config.json" ]; then
    echo ""
    echo "检测到尚未创建配置文件，将从示例配置文件复制..."
    cp "configs/config.example.json" "configs/config.json"
    echo "配置文件已创建: configs/config.json"
fi

while true; do
    echo ""
    echo "请选择操作:"
    echo "1. 配置API平台"
    echo "2. 配置API密钥"
    echo "3. 配置模型"
    echo "4. 配置数据库"
    echo "5. 配置MySQL数据库"
    echo "6. 查看当前配置"
    echo "7. 退出"
    echo ""

    read -p "请输入选项 (1-7): " choice

    case $choice in
        1)
            config_platform
            ;;
        2)
            config_api_keys
            ;;
        3)
            config_models
            ;;
        4)
            config_database
            ;;
        5)
            config_mysql
            ;;
        6)
            view_config
            ;;
        7)
            echo "感谢使用Go-ocsBase配置工具！"
            exit 0
            ;;
        *)
            echo "无效选项，请重新选择"
            ;;
    esac

    echo ""
    read -p "按回车键继续..." dummy
    clear
done