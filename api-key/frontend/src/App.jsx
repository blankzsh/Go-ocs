import React, { useState, useEffect } from 'react';
import './App.css';

const API_BASE_URL = 'http://localhost:3002/api';

function App() {
  const [keys, setKeys] = useState([]);
  const [newKeyName, setNewKeyName] = useState('');
  const [newKey, setNewKey] = useState('');
  const [showKey, setShowKey] = useState(false);
  const [validationKey, setValidationKey] = useState('');
  const [validationResult, setValidationResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // 获取所有API密钥
  const fetchKeys = async () => {
    try {
      setLoading(true);
      setError('');
      console.log('开始获取API密钥列表');
      const response = await fetch(`${API_BASE_URL}/keys`);
      console.log('API密钥列表响应状态:', response.status);
      
      const text = await response.text();
      console.log('API密钥列表响应内容:', text);
      
      let data;
      try {
        data = JSON.parse(text);
      } catch (parseError) {
        throw new Error(`响应不是有效的JSON格式: ${text}`);
      }
      
      if (!response.ok) {
        throw new Error(data.error || `HTTP错误: ${response.status}`);
      }
      
      setKeys(data.keys || []);
    } catch (err) {
      console.error('获取API密钥列表失败:', err);
      setError(`获取API密钥列表失败: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  // 生成新API密钥
  const generateKey = async () => {
    if (!newKeyName.trim()) {
      setError('请输入密钥名称');
      return;
    }

    try {
      setLoading(true);
      setError('');
      console.log('开始生成API密钥，名称:', newKeyName);
      
      const response = await fetch(`${API_BASE_URL}/keys`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: newKeyName }),
      });
      
      console.log('生成API密钥响应状态:', response.status);
      
      const text = await response.text();
      console.log('生成API密钥响应内容:', text);
      
      let data;
      try {
        data = JSON.parse(text);
      } catch (parseError) {
        throw new Error(`响应不是有效的JSON格式: ${text}`);
      }
      
      if (!response.ok) {
        throw new Error(data.error || `HTTP错误: ${response.status}`);
      }
      
      setNewKey(data.key);
      setShowKey(true);
      setNewKeyName('');
      fetchKeys(); // 刷新密钥列表
    } catch (err) {
      console.error('生成API密钥失败:', err);
      setError(`生成API密钥失败: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  // 验证API密钥
  const validateKey = async () => {
    if (!validationKey.trim()) {
      setError('请输入要验证的API密钥');
      return;
    }

    try {
      setLoading(true);
      setError('');
      console.log('开始验证API密钥');
      
      const response = await fetch(`${API_BASE_URL}/keys/validate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ key: validationKey }),
      });
      
      console.log('验证API密钥响应状态:', response.status);
      
      const text = await response.text();
      console.log('验证API密钥响应内容:', text);
      
      let data;
      try {
        data = JSON.parse(text);
      } catch (parseError) {
        throw new Error(`响应不是有效的JSON格式: ${text}`);
      }
      
      if (!response.ok) {
        throw new Error(data.error || `HTTP错误: ${response.status}`);
      }
      
      setValidationResult(data.valid);
    } catch (err) {
      console.error('验证API密钥失败:', err);
      setError(`验证API密钥失败: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  // 删除API密钥
  const deleteKey = async (id) => {
    if (!window.confirm('确定要删除这个API密钥吗？此操作不可恢复。')) {
      return;
    }

    try {
      setLoading(true);
      setError('');
      console.log('开始删除API密钥，ID:', id);
      
      const response = await fetch(`${API_BASE_URL}/keys/${id}`, {
        method: 'DELETE',
      });
      
      console.log('删除API密钥响应状态:', response.status);
      
      const text = await response.text();
      console.log('删除API密钥响应内容:', text);
      
      let data;
      try {
        data = JSON.parse(text);
      } catch (parseError) {
        throw new Error(`响应不是有效的JSON格式: ${text}`);
      }
      
      if (!response.ok) {
        throw new Error(data.error || `HTTP错误: ${response.status}`);
      }
      
      fetchKeys(); // 刷新密钥列表
    } catch (err) {
      console.error('删除API密钥失败:', err);
      setError(`删除API密钥失败: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  // 复制文本到剪贴板
  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    alert('已复制到剪贴板');
  };

  // 组件加载时获取密钥列表
  useEffect(() => {
    fetchKeys();
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <h1>API密钥管理系统</h1>
      </header>

      <main>
        {/* 生成新密钥 */}
        <section className="section">
          <h2>生成新API密钥</h2>
          <div className="form-group">
            <label htmlFor="keyName">密钥名称:</label>
            <input
              type="text"
              id="keyName"
              value={newKeyName}
              onChange={(e) => setNewKeyName(e.target.value)}
              placeholder="请输入密钥名称"
            />
            <button onClick={generateKey} disabled={loading}>
              {loading ? '生成中...' : '生成密钥'}
            </button>
          </div>

          {newKey && (
            <div className="new-key-result">
              <h3>新生成的API密钥:</h3>
              <div className="key-display">
                {showKey ? (
                  <>
                    <span>{newKey}</span>
                    <button onClick={() => copyToClipboard(newKey)}>复制</button>
                  </>
                ) : (
                  <>
                    <span>••••••••••••••••</span>
                    <button onClick={() => setShowKey(true)}>显示</button>
                  </>
                )}
                <button onClick={() => copyToClipboard(newKey)}>复制</button>
              </div>
              <p className="warning">
                请妥善保管此密钥，一旦关闭此窗口将无法再次查看！
              </p>
            </div>
          )}
        </section>

        {/* 验证密钥 */}
        <section className="section">
          <h2>验证API密钥</h2>
          <div className="form-group">
            <label htmlFor="validateKey">API密钥:</label>
            <input
              type="text"
              id="validateKey"
              value={validationKey}
              onChange={(e) => setValidationKey(e.target.value)}
              placeholder="请输入要验证的API密钥"
            />
            <button onClick={validateKey} disabled={loading}>
              {loading ? '验证中...' : '验证密钥'}
            </button>
          </div>

          {validationResult !== null && (
            <div className={`validation-result ${validationResult ? 'valid' : 'invalid'}`}>
              {validationResult ? '✓ 有效的API密钥' : '✗ 无效的API密钥'}
            </div>
          )}
        </section>

        {/* 密钥列表 */}
        <section className="section">
          <h2>API密钥列表</h2>
          {error && <div className="error">{error}</div>}
          {loading && <div className="loading">加载中...</div>}
          
          <div className="keys-list">
            {keys.length > 0 ? (
              keys.map((key) => (
                <div key={key.id} className="key-item">
                  <div className="key-info">
                    <h3>{key.name}</h3>
                    <p>创建时间: {new Date(key.created_at).toLocaleString()}</p>
                  </div>
                  <div className="key-actions">
                    <button 
                      onClick={() => deleteKey(key.id)}
                      className="delete-btn"
                    >
                      删除
                    </button>
                  </div>
                </div>
              ))
            ) : (
              <p>暂无API密钥</p>
            )}
          </div>
        </section>
      </main>
    </div>
  );
}

export default App;