# ws-gateway

`ws-gateway` 是 IM 系统的 WebSocket 接入网关，负责：

- 客户端长连接建立与鉴权
- 会话（Session）生命周期管理与心跳保活
- 离线消息在用户上线时自动投递
- 客户端上行消息路由（发送至 Kafka / 调用 RPC）
- Kafka 消费并将消息实时推送到已连接用户

---

## 架构分层

```
ws-gateway/
├── ws.go                        # 启动入口（仅做依赖注入/wiring，无业务逻辑）
├── etc/ws.yaml                  # 配置文件
└── internal/
    ├── config/
    │   └── config.go            # 配置结构体
    ├── model/
    │   ├── ws_message.go        # 客户端上行消息结构
    │   └── push_message.go      # 服务端下行推送结构及消息类型常量
    ├── session/                 # 连接会话层
    │   ├── manager.go           # SessionManager：uid → Session 并发安全映射
    │   ├── session.go           # 单连接 Session：读写循环、心跳、优雅关闭
    │   ├── state.go             # 会话状态枚举（在线/离线/重连）
    │   └── offline_store.go     # 离线消息 Redis 存储
    ├── router/
    │   └── router.go            # 上行消息路由：chat → Kafka，ack → ChatRPC
    ├── transport/               # 传输层
    │   ├── server.go            # 路由注册
    │   └── handler.go           # HTTP→WebSocket 升级 + Token 鉴权 + 离线消息投递
    ├── push/
    │   └── service.go           # PushService：向在线用户推送消息 / 广播
    ├── consume/
    │   └── message.go           # Kafka 消费：chat-topic / friend-topic
    ├── rpc/
    │   └── auth.go              # AuthService：封装 gRPC token 校验
    └── svc/
        └── service_context.go   # ServiceContext：统一持有所有依赖
```

### 层间依赖规则

```
main (ws.go)
  └── svc.ServiceContext          ← 持有全部依赖，向下注入
       ├── transport               ← 升级 + 鉴权，调用 router
       │     └── router            ← 消息路由，调用 kafka / rpc
       ├── push                    ← 消息投递，调用 session.Manager
       ├── consume                 ← Kafka 消费，调用 push
       └── session                 ← 纯 I/O，无业务逻辑
```

每一层**只向下依赖**，不向上回调，确保职责单一：

| 层 | 职责 | 不做什么 |
|---|---|---|
| `transport` | WebSocket 升级、鉴权、离线消息触发 | 不含任何业务逻辑 |
| `router` | 根据消息类型分发到 Kafka / RPC | 不做 I/O 管理 |
| `session` | 读写循环、心跳、优雅关闭 | 不知道 Manager、不做业务 |
| `push` | 序列化并写入 session.send 通道 | 不读 Kafka |
| `consume` | 从 Kafka 读取并调用 push | 不直接操作 session |

---

## 请求与消息流程

### 1. 建立连接

```
Client ──GET /ws/connect──► transport.ConnectHandler
           ├── extractToken(r)          # 从 Header / Query / Sec-WebSocket-Protocol 提取 token
           ├── AuthService.VerifyToken  # gRPC 调用 auth-service 校验
           ├── upgrader.Upgrade         # HTTP 升级为 WebSocket
           ├── session.NewSession       # 创建 Session（纯 I/O 结构）
           ├── SessionManager.Add       # 注册到在线用户表
           ├── sess.Start               # 启动 readLoop / writeLoop goroutine
           └── go deliverOfflineMessages # 异步推送离线消息
```

### 2. 客户端上行消息

```
Client ──WS frame──► session.readLoop
                       └── onMessage(userID, payload)
                             └── router.HandleMessage
                                   ├── "chat" → kafka.SendMessage(im-chat-topic)
                                   └── "ack"  → ChatRPC.AckMessage
```

### 3. 服务端下行推送（Kafka → Client）

```
chat-service / friend-service ──► Kafka topic
                                     └── consume.startConsumer
                                           └── push.Service.PushToUser
                                                 └── session.Manager.SendTo
                                                       └── session.send channel
                                                             └── session.writeLoop
                                                                   └── WS frame → Client
```

### 4. 心跳保活

Session 每 50 秒向客户端发送 `Ping` 帧；客户端须在 60 秒内回复 `Pong`，否则连接被关闭。

### 5. 断线处理

```
session.readLoop / writeLoop 遇到错误
  └── session.Close()（closeOnce 保证幂等）
        ├── cancel ctx（通知 writeLoop 退出）
        ├── onClose(userID) → SessionManager.Remove
        └── conn.Close()
```

---

## 消息格式

### 客户端上行

| 字段 | 类型 | 说明 |
|---|---|---|
| `type` | string | `"chat"` \| `"ack"` |
| `to` | int64 | 接收方 userId（chat 消息必填）|
| `content` | string | 消息内容 / 被 ACK 的 messageId |

```json
{ "type": "chat", "to": 10002, "content": "hello" }
{ "type": "ack",  "content": "msg-uuid-xxx" }
```

### 服务端下行

| 字段 | 类型 | 说明 |
|---|---|---|
| `type` | string | `"chat_message"` \| `"friend_event"` \| `"offline_message"` |
| `payload` | object | 对应类型的完整结构体 |

```json
{ "type": "chat_message",  "payload": { "message_id": "...", "from_user_id": 10001, "to_user_id": 10002, "content": "hello", "timestamp": 1700000000 } }
{ "type": "friend_event",  "payload": { "event_type": "FriendRequestReceived", "from_user": 10001, "to_user": 10002, "timestamp": 1700000000 } }
{ "type": "offline_message", "payload": { ... } }
```

---

## Token 提取优先级

连接时按以下顺序读取 token（取第一个非空值）：

1. `Authorization: Bearer <token>` 请求头
2. `?token=<token>` 查询参数
3. `Sec-WebSocket-Protocol: <token>` 请求头（兼容浏览器 WebSocket API）

---

## 配置项说明（etc/ws.yaml）

```yaml
Name: ws-gateway
Host: 0.0.0.0
Port: 8888            # HTTP/WebSocket 监听端口

Redis:
  Addr: localhost:6379
  Password: ""
  DB: 0               # 离线消息使用的 Redis DB

Kafka:
  Brokers:
    - kafka:9092
  Topic: im-chat-topic        # 上行聊天消息写入 / 下行聊天消息消费
  FriendTopic: im-friend-topic # 好友事件消费

ChatRpcConf:                  # chat-service gRPC 端点
  Endpoints:
    - 127.0.0.1:2379
  Timeout: 3000

AuthRpcConf:                  # auth-service gRPC 端点
  Endpoints:
    - 127.0.0.1:2379
  Key: auth.rpc
```

---

## 离线消息

用户不在线时，chat-service 将消息写入 Redis List（key：`offline:msg:<userID>`）。  
用户重新连接后，`transport.deliverOfflineMessages` 在 goroutine 中异步将离线消息全部读出并推送给客户端，然后清空 Redis 列表。

---

## 启动方式

```bash
go run ws.go -f etc/ws.yaml
```
