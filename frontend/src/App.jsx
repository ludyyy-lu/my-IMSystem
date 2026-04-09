import { Component, useState, useEffect, useRef, useCallback } from 'react';
import './App.css';
import { api } from './api';

// WS_URL: empty string uses the Vite proxy /ws → ws://localhost:8888.
// In production, set VITE_WS_URL to the ws-gateway WebSocket endpoint.
const WS_BASE = import.meta.env.VITE_WS_URL ||
  (typeof window !== 'undefined'
    ? (window.location.protocol === 'https:' ? 'wss://' : 'ws://') + window.location.host
    : 'ws://localhost:3000');

// ---- Helpers ----
function getMyId() {
  const token = localStorage.getItem('access_token');
  if (!token) return null;
  try {
    const payload = decodeJwtPayload(token);
    return payload.uid;
  } catch {
    return null;
  }
}

function decodeJwtPayload(token) {
  const parts = token.split('.');
  if (parts.length !== 3) return null;
  const base64 = parts[1].replace(/-/g, '+').replace(/_/g, '/');
  const padded = base64 + '='.repeat((4 - (base64.length % 4 || 4)) % 4);
  return JSON.parse(atob(padded));
}

function hasValidAccessToken() {
  const token = localStorage.getItem('access_token');
  if (!token) return false;
  try {
    const payload = decodeJwtPayload(token);
    if (!payload || typeof payload.exp !== 'number') return false;
    const now = Math.floor(Date.now() / 1000);
    return payload.exp > now;
  } catch {
    return false;
  }
}

function normalizeArray(value) {
  return Array.isArray(value) ? value : [];
}

function clearAuthTokens() {
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
}

function formatTime(ts) {
  if (!ts) return '';
  const d = new Date(ts * 1000);
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

function avatarChar(name) {
  return (name || '?')[0].toUpperCase();
}

class AppErrorBoundary extends Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }

  componentDidCatch(error, info) {
    console.error('App crashed:', error, info);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="bootstrap-state error-state">
          <div className="bootstrap-card">
            <h2>页面加载失败</h2>
            <p>前端发生了运行时错误，已经被拦截，避免白屏。</p>
            <pre>{String(this.state.error?.message || this.state.error || 'Unknown error')}</pre>
            <button type="button" onClick={() => window.location.reload()}>刷新重试</button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

// ---- Auth Page ----
function AuthPage({ onLogin }) {
  const [mode, setMode] = useState('login'); // 'login' | 'register'
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      let data;
      if (mode === 'login') {
        data = await api.login(username, password);
      } else {
        data = await api.register(username, password);
      }
      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('refresh_token', data.refresh_token);
      onLogin();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h1>💬 IM System</h1>
        <p>{mode === 'login' ? '登录您的账号' : '注册新账号'}</p>
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="用户名"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="密码"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          {error && <div className="auth-error">{error}</div>}
          <button type="submit" disabled={loading}>
            {loading ? '请稍候…' : mode === 'login' ? '登录' : '注册'}
          </button>
        </form>
        <div className="toggle-link">
          {mode === 'login' ? (
            <>还没有账号？<span onClick={() => { setMode('register'); setError(''); }}>立即注册</span></>
          ) : (
            <>已有账号？<span onClick={() => { setMode('login'); setError(''); }}>立即登录</span></>
          )}
        </div>
      </div>
    </div>
  );
}

// ---- WebSocket hook ----
function useWebSocket(token, onMessage) {
  const wsRef = useRef(null);
  const connectRef = useRef(null);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    if (!token) return;

    function connect() {
      const ws = new WebSocket(`${WS_BASE}/ws/connect?token=${token}`);
      wsRef.current = ws;

      ws.onopen = () => setConnected(true);
      ws.onclose = () => {
        setConnected(false);
        // Reconnect after 3 s
        setTimeout(() => connectRef.current?.(), 3000);
      };
      ws.onerror = () => ws.close();
      ws.onmessage = (e) => {
        try {
          const msg = JSON.parse(e.data);
          onMessage(msg);
        } catch {
          /* ignore */
        }
      };
    }

    connectRef.current = connect;
    connect();

    return () => {
      connectRef.current = null;
      wsRef.current?.close();
    };
  }, [token, onMessage]);

  const send = useCallback((obj) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(obj));
    }
  }, []);

  return { connected, send };
}

// ---- Chat App ----
function ChatApp({ onLogout }) {
  const token = localStorage.getItem('access_token');
  const myId = getMyId();

  const [tab, setTab] = useState('friends'); // 'friends' | 'requests' | 'search'
  const [friends, setFriends] = useState([]);
  const [requests, setRequests] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [selectedFriend, setSelectedFriend] = useState(null); // { friend_id, username }
  const [messages, setMessages] = useState([]);
  const [inputText, setInputText] = useState('');
  const [unread, setUnread] = useState({}); // { peer_id: count }
  const [bootstrapped, setBootstrapped] = useState(false);
  const messagesEndRef = useRef(null);

  // Handle incoming WebSocket messages
  const handleWsMessage = useCallback((msg) => {
    if (msg.type === 'chat_message' && msg.payload) {
      const { from_user_id, to_user_id, content, timestamp, message_id } = msg.payload;
      const peerId = from_user_id === myId ? to_user_id : from_user_id;

      setSelectedFriend((cur) => {
        if (cur && cur.friend_id === peerId) {
          setMessages((prev) => [...prev, { from_user_id, to_user_id, content, timestamp, message_id }]);
        } else {
          setUnread((prev) => ({ ...prev, [peerId]: (prev[peerId] || 0) + 1 }));
        }
        return cur;
      });
    }
  }, [myId]);

  const { connected, send } = useWebSocket(token, handleWsMessage);

  // Load friends
  useEffect(() => {
    let cancelled = false;

    Promise.all([
      api.getFriends(),
      api.getUnreadCount(),
    ])
      .then(([friendsData, unreadData]) => {
        if (cancelled) return;
        setFriends(normalizeArray(friendsData));
        const map = {};
        normalizeArray(unreadData).forEach((c) => {
          if (c && typeof c.peer_id !== 'undefined') {
            map[c.peer_id] = c.unread_count || 0;
          }
        });
        setUnread(map);
      })
      .catch(() => {
        if (!cancelled) {
          setFriends([]);
          setUnread({});
        }
      })
      .finally(() => {
        if (!cancelled) setBootstrapped(true);
      });

    return () => {
      cancelled = true;
    };
  }, []);

  // Load pending requests when switching to requests tab
  useEffect(() => {
    if (tab === 'requests') {
      api.getFriendRequests().then((data) => setRequests(normalizeArray(data))).catch(() => {});
    }
  }, [tab]);

  // Search - debounced, timeout handles both empty and non-empty cases
  useEffect(() => {
    const t = setTimeout(() => {
      if (!searchQuery.trim()) {
        setSearchResults([]);
      } else {
          api.searchUser(searchQuery).then((data) => setSearchResults(normalizeArray(data))).catch(() => {});
      }
    }, 400);
    return () => clearTimeout(t);
  }, [searchQuery]);

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  if (!bootstrapped) {
    return (
      <div className="chat-layout">
        <div className="sidebar">
          <div className="sidebar-header">
            <h2>💬 IM System</h2>
            <button className="logout-btn" onClick={handleLogout}>退出</button>
          </div>
          <div className="friend-list">
            <div className="loading-state">正在加载联系人和未读消息…</div>
          </div>
        </div>

        <div className="chat-window">
          <div className="empty-state">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="#667eea" strokeWidth="1.5">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
            <p>正在连接聊天服务…</p>
          </div>
        </div>
      </div>
    );
  }

  // Select friend → load history
  function selectFriend(f) {
    setSelectedFriend(f);
    setUnread((prev) => ({ ...prev, [f.friend_id]: 0 }));
    api.getChatHistory(f.friend_id).then((msgs) => {
      setMessages(msgs || []);
    }).catch(() => setMessages([]));
    // batch ack
    send({ type: 'batch_ack', peer_id: f.friend_id });
  }

  // Send message
  function sendMessage() {
    if (!inputText.trim() || !selectedFriend) return;
    send({ type: 'chat', to: selectedFriend.friend_id, content: inputText.trim() });
    // Optimistic update
    setMessages((prev) => [
      ...prev,
      {
        from_user_id: myId,
        to_user_id: selectedFriend.friend_id,
        content: inputText.trim(),
        timestamp: Math.floor(Date.now() / 1000),
        message_id: `local-${Date.now()}`,
      },
    ]);
    setInputText('');
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  }

  async function handleRespond(requestId, action) {
    try {
      await api.respondFriendRequest(requestId, action);
      setRequests((prev) => prev.filter((r) => r.request_id !== requestId));
      if (action === 'accept') api.getFriends().then(setFriends).catch(() => {});
    } catch (err) {
      alert(err.message);
    }
  }

  async function handleAddFriend(userId) {
    try {
      await api.sendFriendRequest(userId, '你好，我想加你为好友');
      alert('好友申请已发送！');
    } catch (err) {
      alert(err.message);
    }
  }

  function handleLogout() {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    onLogout();
  }

  return (
    <div className="chat-layout">
      {/* Sidebar */}
      <div className="sidebar">
        <div className="sidebar-header">
          <h2>💬 IM System</h2>
          <button className="logout-btn" onClick={handleLogout}>退出</button>
        </div>

        <div className="sidebar-tabs">
          {['friends', 'requests', 'search'].map((t) => (
            <div
              key={t}
              className={`sidebar-tab ${tab === t ? 'active' : ''}`}
              onClick={() => setTab(t)}
            >
              {t === 'friends' ? '好友' : t === 'requests' ? '申请' : '搜索'}
            </div>
          ))}
        </div>

        {tab === 'search' && (
          <div className="search-bar">
            <input
              placeholder="搜索用户名…"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
        )}

        <div className="friend-list">
          {/* Friends tab */}
          {tab === 'friends' && friends.map((f) => (
            <div
              key={f.friend_id}
              className={`friend-item ${selectedFriend?.friend_id === f.friend_id ? 'active' : ''}`}
              onClick={() => selectFriend(f)}
            >
              <div className="avatar">{avatarChar(f.username)}</div>
              <div className="friend-info">
                <div className="friend-name">{f.username}</div>
                <div className="friend-id">ID: {f.friend_id}</div>
              </div>
              {unread[f.friend_id] > 0 && (
                <span className="unread-badge">{unread[f.friend_id]}</span>
              )}
            </div>
          ))}
          {tab === 'friends' && friends.length === 0 && (
            <div style={{ padding: '20px', color: '#aaa', textAlign: 'center', fontSize: 13 }}>
              暂无好友，去搜索添加吧
            </div>
          )}

          {/* Requests tab */}
          {tab === 'requests' && requests.filter(r => r.status === 'PENDING').map((req) => (
            <div key={req.request_id} className="request-item">
              <div className="request-info">来自：{req.from_username || `用户 ${req.from_user_id}`}</div>
              {req.remark && <div className="request-remark">{req.remark}</div>}
              <div className="request-actions">
                <button className="btn-accept" onClick={() => handleRespond(req.request_id, 'accept')}>接受</button>
                <button className="btn-reject" onClick={() => handleRespond(req.request_id, 'reject')}>拒绝</button>
              </div>
            </div>
          ))}
          {tab === 'requests' && requests.filter(r => r.status === 'PENDING').length === 0 && (
            <div style={{ padding: '20px', color: '#aaa', textAlign: 'center', fontSize: 13 }}>
              暂无好友申请
            </div>
          )}

          {/* Search tab */}
          {tab === 'search' && searchResults.map((u) => (
            <div key={u.id} className="search-result-item">
              <div className="avatar">{avatarChar(u.nickname || String(u.id))}</div>
              <div className="friend-info">
                <div className="friend-name">{u.nickname || `用户 ${u.id}`}</div>
                <div className="friend-id">ID: {u.id}</div>
              </div>
              <button className="add-btn" onClick={() => handleAddFriend(u.id)}>添加</button>
            </div>
          ))}
        </div>
      </div>

      {/* Chat window */}
      {selectedFriend ? (
        <div className="chat-window">
          <div className="chat-header">
            <div className="avatar">{avatarChar(selectedFriend.username)}</div>
            <h3>{selectedFriend.username}</h3>
            <span className={`ws-status ${connected ? 'connected' : 'disconnected'}`}>
              {connected ? '● 已连接' : '○ 断开'}
            </span>
          </div>

          <div className="messages-area">
            {messages.map((msg, i) => {
              const isSent = msg.from_user_id === myId;
              return (
                <div key={msg.message_id || i} className={`message-bubble ${isSent ? 'sent' : 'received'}`}>
                  <div className="message-text">{msg.content}</div>
                  <div className="message-time">{formatTime(msg.timestamp)}</div>
                </div>
              );
            })}
            <div ref={messagesEndRef} />
          </div>

          <div className="message-input-area">
            <textarea
              rows={1}
              placeholder="输入消息… (Enter 发送)"
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              onKeyDown={handleKeyDown}
            />
            <button className="send-btn" onClick={sendMessage} disabled={!inputText.trim() || !connected}>
              发送
            </button>
          </div>
        </div>
      ) : (
        <div className="chat-window">
          <div className="empty-state">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="#667eea" strokeWidth="1.5">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
            <p>从左侧选择好友开始聊天</p>
          </div>
        </div>
      )}
    </div>
  );
}

// ---- Root ----
export default function App() {
  const [loggedIn, setLoggedIn] = useState(() => hasValidAccessToken());

  useEffect(() => {
    if (!hasValidAccessToken()) {
      clearAuthTokens();
    }
  }, []);

  if (!loggedIn) {
    return <AuthPage onLogin={() => setLoggedIn(true)} />;
  }

  return (
    <AppErrorBoundary>
      <ChatApp onLogout={() => setLoggedIn(false)} />
    </AppErrorBoundary>
  );
}
