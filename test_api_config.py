#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import json
import requests
import urllib.parse
import sys
import os

class APITestTool:
    def __init__(self, config_file=None):
        self.config = None
        self.config_file = config_file
        if config_file and os.path.exists(config_file):
            self.load_config(config_file)
    
    def load_config(self, config_file):
        """从文件加载API配置"""
        try:
            with open(config_file, 'r', encoding='utf-8') as f:
                # 读取文件内容并处理可能的格式问题
                content = f.read().strip()
                
                # 如果是多行内容，尝试提取JSON部分
                if '\n' in content:
                    lines = content.split('\n')
                    json_lines = []
                    in_json = False
                    
                    for line in lines:
                        line = line.strip()
                        if line.startswith('{'):
                            in_json = True
                            json_lines.append(line)
                        elif in_json and line:
                            json_lines.append(line)
                        elif line.endswith('}'):
                            json_lines.append(line)
                            break
                    
                    content = ' '.join(json_lines)
                
                # 尝试解析JSON
                self.config = json.loads(content)
                print("✓ 成功加载API配置文件")
                return True
        except json.JSONDecodeError as e:
            print(f"✗ JSON解析错误: {e}")
            print("请确保配置文件是有效的JSON格式")
            return False
        except Exception as e:
            print(f"✗ 加载配置文件时出错: {e}")
            return False
    
    def input_config_manually(self):
        """手动输入API配置"""
        print("\n请手动输入API配置信息:")
        
        self.config = {
            "name": input("名称 (如: 完美题库): ") or "完美题库",
            "homepage": input("主页 (如: https://currso.com/): ") or "https://currso.com/",
            "url": input("API地址 (如: http://127.0.0.1:8000/api/query): "),
            "method": input("请求方法 (GET/POST, 默认GET): ") or "GET",
            "type": input("类型 (如: GM_xmlhttpRequest): ") or "GM_xmlhttpRequest",
            "contentType": input("内容类型 (如: json): ") or "json",
            "data": {},
            "handler": input("处理函数 (默认值): ") or "return (res)=>res.code === 0 ? [undefined, undefined] : [undefined,res.data.data]"
        }
        
        print("\n请输入数据字段 (按回车结束):")
        while True:
            key = input("字段名 (如: title, options, type, api-key): ")
            if not key:
                break
            value = input(f"{key} 的值: ")
            self.config["data"][key] = value
        
        return True
    
    def display_config(self):
        """显示当前配置"""
        if not self.config:
            print("✗ 没有加载配置")
            return
        
        print("\n当前API配置:")
        print("-" * 50)
        print(f"名称: {self.config.get('name', 'N/A')}")
        print(f"主页: {self.config.get('homepage', 'N/A')}")
        print(f"API地址: {self.config.get('url', 'N/A')}")
        print(f"请求方法: {self.config.get('method', 'N/A')}")
        print(f"类型: {self.config.get('type', 'N/A')}")
        print(f"内容类型: {self.config.get('contentType', 'N/A')}")
        print("数据字段:")
        for key, value in self.config.get('data', {}).items():
            print(f"  {key}: {value}")
        print(f"处理函数: {self.config.get('handler', 'N/A')}")
        print("-" * 50)
    
    def build_request_url(self, title, options, question_type, api_key):
        """构建请求URL"""
        if not self.config:
            raise Exception("没有加载配置")
        
        base_url = self.config.get('url', '')
        if not base_url:
            raise Exception("配置中缺少API地址")
        
        # 构建查询参数
        params = {}
        for key, template in self.config.get('data', {}).items():
            # 替换模板变量
            value = template
            value = value.replace('${title}', urllib.parse.quote(title))
            value = value.replace('${options}', urllib.parse.quote(options))
            value = value.replace('${type}', urllib.parse.quote(question_type))
            if api_key:
                value = value.replace('生成api-key', api_key)
            
            params[key] = value
        
        # 构造完整URL
        if params:
            query_string = urllib.parse.urlencode(params)
            if '?' in base_url:
                url = f"{base_url}&{query_string}"
            else:
                url = f"{base_url}?{query_string}"
        else:
            url = base_url
        
        return url
    
    def send_test_request(self, url):
        """发送测试请求"""
        try:
            print(f"\n正在发送请求到: {url}")
            
            # 根据配置设置请求头
            headers = {}
            content_type = self.config.get('contentType', 'json')
            if content_type == 'json':
                headers['Content-Type'] = 'application/json'
            else:
                headers['Content-Type'] = 'application/x-www-form-urlencoded'
            
            # 发送请求
            method = self.config.get('method', 'GET').upper()
            response = requests.request(method, url, headers=headers, timeout=30)
            
            print(f"✓ 请求成功")
            print(f"状态码: {response.status_code}")
            print(f"响应头: {dict(response.headers)}")
            
            # 尝试解析JSON响应
            try:
                data = response.json()
                print("响应数据 (JSON格式):")
                print(json.dumps(data, ensure_ascii=False, indent=2))
                return data
            except:
                # 如果不是JSON，直接显示文本
                print("响应数据 (文本格式):")
                print(response.text)
                return response.text
                
        except requests.exceptions.Timeout:
            print("✗ 请求超时")
            return None
        except requests.exceptions.ConnectionError:
            print("✗ 连接错误，请检查API地址是否正确且服务正在运行")
            return None
        except Exception as e:
            print(f"✗ 请求失败: {e}")
            return None
    
    def parse_response(self, response_data):
        """解析响应数据"""
        if not response_data:
            return None
        
        try:
            # 如果是字符串，尝试解析为JSON
            if isinstance(response_data, str):
                response_data = json.loads(response_data)
            
            # 如果是字典，尝试解析
            if isinstance(response_data, dict):
                code = response_data.get('code', -1)
                if code == 0:
                    print("\n✓ 响应解析成功: 请求成功")
                    return {"result": "success", "data": None}
                else:
                    # 尝试获取 data.data
                    data_field = response_data.get('data')
                    if isinstance(data_field, dict):
                        actual_data = data_field.get('data')
                        print(f"\n✓ 响应解析成功: {actual_data}")
                        return {"result": "success", "data": actual_data}
                    else:
                        print(f"\n✓ 响应解析成功: {data_field}")
                        return {"result": "success", "data": data_field}
            else:
                print(f"\n✓ 响应数据: {response_data}")
                return {"result": "success", "data": response_data}
                
        except Exception as e:
            print(f"✗ 解析响应时出错: {e}")
            return None

def main():
    print("=" * 60)
    print("OCS答题系统API配置测试工具")
    print("=" * 60)
    
    # 创建测试工具实例
    tester = APITestTool()
    
    # 尝试从命令行参数获取配置文件
    config_file = None
    if len(sys.argv) > 1:
        config_file = sys.argv[1]
        print(f"尝试加载配置文件: {config_file}")
        if not tester.load_config(config_file):
            print("配置文件加载失败，将使用手动输入模式")
    
    # 如果没有配置文件或加载失败，提供选项
    while not tester.config:
        print("\n请选择操作:")
        print("1. 从文件加载配置")
        print("2. 手动输入配置")
        print("3. 退出")
        
        choice = input("请输入选择 (1-3): ").strip()
        
        if choice == '1':
            file_path = input("请输入配置文件路径: ").strip()
            if os.path.exists(file_path):
                tester.load_config(file_path)
            else:
                print("✗ 文件不存在")
        elif choice == '2':
            tester.input_config_manually()
        elif choice == '3':
            print("退出程序")
            return
        else:
            print("无效选择，请重新输入")
    
    # 显示配置
    tester.display_config()
    
    # 进入测试循环
    while True:
        print("\n" + "=" * 40)
        print("测试选项:")
        print("1. 发送测试请求")
        print("2. 重新输入测试参数")
        print("3. 重新加载配置")
        print("4. 显示当前配置")
        print("5. 退出")
        
        choice = input("请选择操作 (1-5): ").strip()
        
        if choice == '1':
            # 获取测试参数
            print("\n请输入测试参数:")
            title = input("题目 (默认: 中国的首都是哪里?): ") or "中国的首都是哪里?"
            options = input("选项 (默认: 北京###上海###广州###深圳): ") or "北京###上海###广州###深圳"
            question_type = input("题目类型 (默认: 选择题): ") or "选择题"
            api_key = input("API密钥 (可选): ") or ""
            
            try:
                # 构建请求URL
                url = tester.build_request_url(title, options, question_type, api_key)
                print(f"\n构建的请求URL: {url}")
                
                # 发送请求
                response = tester.send_test_request(url)
                
                # 解析响应
                if response:
                    tester.parse_response(response)
                    
            except Exception as e:
                print(f"✗ 构建请求时出错: {e}")
                
        elif choice == '2':
            # 参数已经包含在选项1中，这里可以留空或添加其他功能
            print("请使用选项1来输入测试参数并发送请求")
            
        elif choice == '3':
            file_path = input("请输入配置文件路径: ").strip()
            if os.path.exists(file_path):
                tester.load_config(file_path)
                tester.display_config()
            else:
                print("✗ 文件不存在")
                
        elif choice == '4':
            tester.display_config()
            
        elif choice == '5':
            print("退出程序")
            break
            
        else:
            print("无效选择，请重新输入")

if __name__ == "__main__":
    main()