# IM即时通讯系统 - API接口文档

## 一、接口概述

### 1.1 接口规范

**基础信息：**
- **协议**：HTTP/1.1 (REST API), WebSocket (实时通信)
- **编码**：UTF-8
- **数据格式**：JSON
- **认证方式**：JWT Bearer Token

**通用响应格式：**
```json
{
    "code": 200,
    "message": "success",
    "data": {}
}
```

**错误码规范：**
| Code | 说明 |
|------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证（Token 无效/过期） |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 二、认证服务 (Auth Service)

### 2.1 用户注册

**接口地址：** `POST /api/auth/register`

**请求参数：**
```json
{
    "username": "zhangsan",       // 必填，3-20字符，字母数字下划线
    "password": "123456",         // 必填，6-20字符
    "nickname": "张三",           // 可选，1-50字符
    "device_id": "iPhone-12345"   // 必填，设备唯一标识
}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "注册成功",
    "data": {
        "user_id": 10001,
        "username": "zhangsan",
        "nickname": "张三",
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_in": 900  // Access Token 有效期（秒）
    }
}
```

**错误示例：**
```json
{
    "code": 400,
    "message": "用户名已存在",
    "data": null
}
```

---

### 2.2 用户登录

**接口地址：** `POST /api/auth/login`

**请求参数：**
```json
{
    "username": "zhangsan",
    "password": "123456",
    "device_id": "iPhone-12345"
}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "user_id": 10001,
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_in": 900
    }
}
```

---

### 2.3 刷新 Token

**接口地址：** `POST /api/auth/refresh`

**请求头：**
```
Authorization: Bearer {refresh_token}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "刷新成功",
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_in": 900
    }
}
```

---

### 2.4 登出

**接口地址：** `POST /api/auth/logout`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求参数：**
```json
{
    "device_id": "iPhone-12345"  // 可选，不传则登出所有设备
}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "登出成功",
    "data": null
}
```

---

### 2.5 查询会话列表

**接口地址：** `GET /api/auth/sessions`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "sessions": [
            {
                "device_id": "iPhone-12345",
                "device_name": "iPhone 12",
                "login_time": "2026-02-03T10:00:00Z",
                "last_active": "2026-02-03T12:30:00Z",
                "is_current": true
            },
            {
                "device_id": "Chrome-67890",
                "device_name": "Chrome on Windows",
                "login_time": "2026-02-02T15:00:00Z",
                "last_active": "2026-02-03T09:00:00Z",
                "is_current": false
            }
        ]
    }
}
```

---

## 三、用户服务 (User Service)

### 3.1 获取用户信息

**接口地址：** `GET /api/users/{user_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "user_id": 10001,
        "username": "zhangsan",
        "nickname": "张三",
        "avatar": "https://cdn.example.com/avatar/10001.jpg",
        "bio": "这个人很懒，什么都没写",
        "gender": 1,  // 0=未知 1=男 2=女
        "is_online": true,
        "created_at": "2026-01-01T00:00:00Z"
    }
}
```

---

### 3.2 更新用户信息

**接口地址：** `PUT /api/users/{user_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求参数：**
```json
{
    "nickname": "张三三",
    "avatar": "https://cdn.example.com/avatar/new.jpg",
    "bio": "热爱生活",
    "gender": 1
}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "更新成功",
    "data": {
        "user_id": 10001,
        "nickname": "张三三",
        "avatar": "https://cdn.example.com/avatar/new.jpg",
        "bio": "热爱生活",
        "gender": 1
    }
}
```

---

### 3.3 搜索用户

**接口地址：** `GET /api/users/search`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求参数：**
```
keyword=张三&page=1&page_size=20
```

**响应示例：**
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "total": 2,
        "users": [
            {
                "user_id": 10001,
                "username": "zhangsan",
                "nickname": "张三",
                "avatar": "https://cdn.example.com/avatar/10001.jpg"
            },
            {
                "user_id": 10002,
                "username": "zhangsan2",
                "nickname": "张三2",
                "avatar": "https://cdn.example.com/avatar/10002.jpg"
            }
        ]
    }
}
```

---

## 四、好友服务 (Friend Service)

### 4.1 获取好友列表

**接口地址：** `GET /api/friends`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "friends": [
            {
                "user_id": 10002,
                "username": "lisi",
                "nickname": "李四",
                "avatar": "https://cdn.example.com/avatar/10002.jpg",
                "is_online": true,
                "remark": "大学同学"  // 备注
            },
            {
                "user_id": 10003,
                "username": "wangwu",
                "nickname": "王五",
                "avatar": "https://cdn.example.com/avatar/10003.jpg",
                "is_online": false,
                "remark": ""
            }
        ]
    }
}
```

---

### 4.2 添加好友

**接口地址：** `POST /api/friends/requests`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求参数：**
```json
{
    "receiver_id": 10002,
    "message": "你好，我是张三，可以加个好友吗？"
}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "好友请求已发送",
    "data": {
        "request_id": 1001,
        "status": 0  // 0=待处理 1=已同意 2=已拒绝
    }
}
```

---

### 4.3 获取好友请求列表

**接口地址：** `GET /api/friends/requests`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求参数：**
```
status=0&page=1&page_size=20
```

**响应示例：**
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "total": 3,
        "requests": [
            {
                "request_id": 1001,
                "sender_id": 10002,
                "sender_username": "lisi",
                "sender_nickname": "李四",
                "sender_avatar": "https://cdn.example.com/avatar/10002.jpg",
                "message": "你好，我是李四",
                "status": 0,
                "created_at": "2026-02-03T10:00:00Z"
            }
        ]
    }
}
```

---

### 4.4 响应好友请求

**接口地址：** `PUT /api/friends/requests/{request_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求参数：**
```json
{
    "action": "accept"  // accept=同意, reject=拒绝
}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "已同意好友请求",
    "data": null
}
```

---

### 4.5 删除好友

**接口地址：** `DELETE /api/friends/{friend_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "删除成功",
    "data": null
}
```

---

### 4.6 拉黑用户

**接口地址：** `POST /api/friends/block`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求参数：**
```json
{
    "blocked_id": 10002
}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "拉黑成功",
    "data": null
}
```

---

### 4.7 取消拉黑

**接口地址：** `DELETE /api/friends/block/{blocked_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "取消拉黑成功",
    "data": null
}
```

---

### 4.8 获取黑名单列表

**接口地址：** `GET /api/friends/blocked`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "blocked_users": [
            {
                "user_id": 10005,
                "username": "spammer",
                "nickname": "广告用户",
                "blocked_at": "2026-02-01T10:00:00Z"
            }
        ]
    }
}
```

---

## 五、聊天服务 (Chat Service)

### 5.1 发送消息（通过 WebSocket）

**WebSocket 连接：** `ws://localhost:8080/ws`

**连接时携带 Token：**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws?token=' + accessToken);
```

**发送消息格式：**
```json
{
    "type": "chat",
    "to_user_id": 10002,
    "content": "你好，在吗？",
    "msg_type": 1,  // 1=文本 2=图片 3=文件
    "client_msg_id": "uuid-1234-5678"  // 客户端生成的唯一ID
}
```

**服务端推送消息格式：**
```json
{
    "type": "chat",
    "message_id": "uuid-server-1234",
    "from_user_id": 10002,
    "from_username": "lisi",
    "from_nickname": "李四",
    "from_avatar": "https://cdn.example.com/avatar/10002.jpg",
    "to_user_id": 10001,
    "content": "我在，怎么了？",
    "msg_type": 1,
    "timestamp": 1738562400
}
```

---

### 5.2 获取聊天记录

**接口地址：** `GET /api/chat/history`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求参数：**
```
peer_id=10002&last_msg_time=1738562400&limit=20
```
- `peer_id`: 对方用户 ID
- `last_msg_time`: 上次查询的最后一条消息时间戳（游标分页）
- `limit`: 每页条数

**响应示例：**
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "has_more": true,
        "messages": [
            {
                "message_id": "uuid-1234",
                "from_user_id": 10001,
                "to_user_id": 10002,
                "content": "你好",
                "msg_type": 1,
                "status": 1,  // 0=未读 1=已读
                "timestamp": 1738562300
            },
            {
                "message_id": "uuid-5678",
                "from_user_id": 10002,
                "to_user_id": 10001,
                "content": "你好",
                "msg_type": 1,
                "status": 1,
                "timestamp": 1738562350
            }
        ]
    }
}
```

---

### 5.3 消息已读回执（通过 WebSocket）

**发送格式：**
```json
{
    "type": "ack",
    "message_id": "uuid-1234"
}
```

**服务端响应：**
```json
{
    "type": "ack_result",
    "message_id": "uuid-1234",
    "success": true
}
```

---

### 5.4 获取未读消息数

**接口地址：** `GET /api/chat/unread`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应示例：**
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "total_unread": 15,
        "conversations": [
            {
                "peer_id": 10002,
                "peer_nickname": "李四",
                "unread_count": 5,
                "last_message": {
                    "content": "在吗？",
                    "timestamp": 1738562400
                }
            },
            {
                "peer_id": 10003,
                "peer_nickname": "王五",
                "unread_count": 10,
                "last_message": {
                    "content": "收到了吗？",
                    "timestamp": 1738562350
                }
            }
        ]
    }
}
```

---

## 六、WebSocket 协议规范

### 6.1 连接建立

**连接地址：**
```
ws://localhost:8080/ws?token={access_token}
```

**握手成功响应：**
```json
{
    "type": "connected",
    "user_id": 10001,
    "server_time": 1738562400
}
```

---

### 6.2 心跳保活

**客户端发送（每 30 秒）：**
```json
{
    "type": "ping"
}
```

**服务端响应：**
```json
{
    "type": "pong",
    "server_time": 1738562400
}
```

---

### 6.3 消息类型汇总

| type | 说明 | 方向 |
|------|------|------|
| connected | 连接成功 | 服务端 → 客户端 |
| ping | 心跳请求 | 客户端 → 服务端 |
| pong | 心跳响应 | 服务端 → 客户端 |
| chat | 聊天消息 | 双向 |
| ack | 消息已读回执 | 客户端 → 服务端 |
| ack_result | 回执结果 | 服务端 → 客户端 |
| friend_request | 好友请求通知 | 服务端 → 客户端 |
| friend_accept | 好友请求被接受 | 服务端 → 客户端 |
| user_status | 好友上下线通知 | 服务端 → 客户端 |
| error | 错误信息 | 服务端 → 客户端 |

---

### 6.4 错误消息格式

```json
{
    "type": "error",
    "code": 401,
    "message": "Token 已过期，请重新登录"
}
```

---

## 七、客户端集成示例

### 7.1 JavaScript/TypeScript

```javascript
class IMClient {
    constructor(serverUrl, accessToken) {
        this.serverUrl = serverUrl;
        this.accessToken = accessToken;
        this.ws = null;
        this.retryCount = 0;
    }
    
    connect() {
        this.ws = new WebSocket(`${this.serverUrl}/ws?token=${this.accessToken}`);
        
        this.ws.onopen = () => {
            console.log('WebSocket 连接成功');
            this.retryCount = 0;
            this.startHeartbeat();
        };
        
        this.ws.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            this.handleMessage(msg);
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket 错误:', error);
        };
        
        this.ws.onclose = () => {
            console.log('WebSocket 连接关闭');
            this.reconnect();
        };
    }
    
    handleMessage(msg) {
        switch (msg.type) {
            case 'connected':
                console.log('连接成功，用户ID:', msg.user_id);
                break;
            case 'chat':
                this.onChatMessage(msg);
                break;
            case 'pong':
                console.log('心跳响应');
                break;
            case 'friend_request':
                this.onFriendRequest(msg);
                break;
            case 'error':
                console.error('错误:', msg.message);
                break;
        }
    }
    
    sendMessage(toUserId, content) {
        const msg = {
            type: 'chat',
            to_user_id: toUserId,
            content: content,
            msg_type: 1,
            client_msg_id: this.generateUUID()
        };
        this.ws.send(JSON.stringify(msg));
    }
    
    sendAck(messageId) {
        this.ws.send(JSON.stringify({
            type: 'ack',
            message_id: messageId
        }));
    }
    
    startHeartbeat() {
        this.heartbeatTimer = setInterval(() => {
            if (this.ws.readyState === WebSocket.OPEN) {
                this.ws.send(JSON.stringify({ type: 'ping' }));
            }
        }, 30000);
    }
    
    reconnect() {
        if (this.retryCount < 5) {
            const delay = Math.min(1000 * (2 ** this.retryCount), 30000);
            setTimeout(() => {
                this.retryCount++;
                this.connect();
            }, delay);
        }
    }
    
    generateUUID() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
            const r = Math.random() * 16 | 0;
            const v = c === 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }
    
    // 业务回调（由使用者实现）
    onChatMessage(msg) {
        console.log('收到消息:', msg);
    }
    
    onFriendRequest(msg) {
        console.log('收到好友请求:', msg);
    }
}

// 使用示例
const client = new IMClient('ws://localhost:8080', 'your_access_token');
client.connect();

// 发送消息
client.sendMessage(10002, '你好，在吗？');

// 收到消息后发送已读回执
client.onChatMessage = (msg) => {
    console.log('收到消息:', msg.content);
    client.sendAck(msg.message_id);
};
```

---

### 7.2 HTTP API 调用示例

```javascript
// 封装 API 客户端
class IMApiClient {
    constructor(baseUrl, accessToken) {
        this.baseUrl = baseUrl;
        this.accessToken = accessToken;
    }
    
    async request(method, path, data = null) {
        const options = {
            method: method,
            headers: {
                'Authorization': `Bearer ${this.accessToken}`,
                'Content-Type': 'application/json'
            }
        };
        
        if (data) {
            options.body = JSON.stringify(data);
        }
        
        const response = await fetch(`${this.baseUrl}${path}`, options);
        return await response.json();
    }
    
    // 用户相关
    async getUser(userId) {
        return this.request('GET', `/api/users/${userId}`);
    }
    
    async updateUser(userId, data) {
        return this.request('PUT', `/api/users/${userId}`, data);
    }
    
    // 好友相关
    async getFriends() {
        return this.request('GET', '/api/friends');
    }
    
    async sendFriendRequest(receiverId, message) {
        return this.request('POST', '/api/friends/requests', {
            receiver_id: receiverId,
            message: message
        });
    }
    
    async respondFriendRequest(requestId, action) {
        return this.request('PUT', `/api/friends/requests/${requestId}`, {
            action: action
        });
    }
    
    // 聊天相关
    async getChatHistory(peerId, lastMsgTime, limit = 20) {
        return this.request('GET', 
            `/api/chat/history?peer_id=${peerId}&last_msg_time=${lastMsgTime}&limit=${limit}`
        );
    }
    
    async getUnreadCount() {
        return this.request('GET', '/api/chat/unread');
    }
}

// 使用示例
const apiClient = new IMApiClient('http://localhost:8080', 'your_access_token');

// 获取好友列表
const friends = await apiClient.getFriends();
console.log('好友列表:', friends);

// 发送好友请求
await apiClient.sendFriendRequest(10002, '你好，加个好友吧');

// 查询聊天记录
const history = await apiClient.getChatHistory(10002, Date.now(), 20);
console.log('聊天记录:', history);
```

---

## 八、Postman 测试集合

可以导入以下 JSON 到 Postman 进行测试：

```json
{
    "info": {
        "name": "IM System API",
        "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
    },
    "variable": [
        {
            "key": "base_url",
            "value": "http://localhost:8080",
            "type": "string"
        },
        {
            "key": "access_token",
            "value": "",
            "type": "string"
        }
    ],
    "item": [
        {
            "name": "Auth",
            "item": [
                {
                    "name": "Register",
                    "request": {
                        "method": "POST",
                        "url": "{{base_url}}/api/auth/register",
                        "body": {
                            "mode": "raw",
                            "raw": "{\n  \"username\": \"testuser\",\n  \"password\": \"123456\",\n  \"nickname\": \"测试用户\",\n  \"device_id\": \"postman-test\"\n}"
                        }
                    }
                },
                {
                    "name": "Login",
                    "request": {
                        "method": "POST",
                        "url": "{{base_url}}/api/auth/login",
                        "body": {
                            "mode": "raw",
                            "raw": "{\n  \"username\": \"testuser\",\n  \"password\": \"123456\",\n  \"device_id\": \"postman-test\"\n}"
                        }
                    }
                }
            ]
        },
        {
            "name": "Users",
            "item": [
                {
                    "name": "Get User",
                    "request": {
                        "method": "GET",
                        "url": "{{base_url}}/api/users/10001",
                        "header": [
                            {
                                "key": "Authorization",
                                "value": "Bearer {{access_token}}"
                            }
                        ]
                    }
                }
            ]
        }
    ]
}
```

---

**文档版本：** v1.0  
**最后更新：** 2026-02-03  
**维护者：** Ludy
