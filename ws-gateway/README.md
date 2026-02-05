# ws-gateway

`ws-gateway` 是 IM 系统的 WebSocket 网关，负责长连接管理、身份认证、消息路由以及 Kafka 推送。

## 功能概览

- WebSocket 连接升级与鉴权（Auth RPC）
- 连接池管理（ConnManager）
- 心跳保活（Ping/Pong）
- 消息路由：chat/ack
- Kafka 消费推送（好友事件 & 聊天消息）
- Redis 离线消息加载

## 目录结构

```
ws-gateway/
├── ws.go                     # 服务入口
├── etc/ws-api.yaml           # 配置文件
├── api/ws.api                # API 定义
└── internal/
    ├── handler/              # WebSocket 接入与路由
    ├── logic/                # 消息解析与业务分发
    ├── conn/                 # 连接池 + 离线消息存储
    ├── connx/                # 单连接生命周期（心跳、离线加载）
    ├── kafka/                # Kafka 消费者
    ├── model/                # WebSocket 消息结构
    ├── svc/                  # ServiceContext
    ├── ws1/                  # 推送封装（Kafka -> WebSocket）
    └── types/                # go-zero 生成类型
```

## 连接与鉴权流程

1. 客户端发起连接：`GET /ws/connect`
2. 从以下任意位置获取 token：
   - `Authorization: Bearer <token>`
   - `?token=<token>` 查询参数
   - `Sec-Websocket-Protocol` 子协议
3. 调用 `auth-service` 校验 token，获取用户 ID
4. 升级为 WebSocket，并注册到 `ConnManager`
5. 启动心跳循环与离线消息加载
6. 连接关闭时移除连接

> 如果客户端使用 `Sec-Websocket-Protocol` 传递 token，服务端会在握手响应中回传该子协议。

## 消息格式

客户端发送的 JSON：

```json
{
  "type": "chat",
  "content": "hello",
  "to": 10002
}
```

目前支持的消息类型：

- `chat`：聊天消息（会写入 Kafka）
- `ack`：消息回执（通过 chat-service RPC 处理）

服务端推送消息采用统一包裹：

```json
{
  "type": "chat_message",
  "payload": { "from": 10001, "to": 10002, "content": "..." }
}
```

`type` 目前包括：

- `chat_message`
- `friend_event`

## Kafka 数据流

`ws.go` 启动两个消费者：

- `im-chat-topic` -> `StartChatConsumer`
- `im-friend-topic` -> `StartFriendConsumer`

消费者读取消息后统一调用 `ws1.PushToUser`，通过连接池推送给在线用户。

## 配置说明（etc/ws-api.yaml）

关键配置项：

- `Host` / `Port`：服务监听地址
- `Redis`：离线消息读取
- `Kafka.Brokers`：Kafka 集群地址
- `ChatRpcConf`：chat-service RPC 端点
- `AuthRpcConf`：auth-service RPC 端点

> `Kafka.Topic` 目前用于发送消息（logic/connection.go）。如需开启发送功能，请在配置中补充该字段。

## 启动方式

```bash
go run ws.go -f etc/ws-api.yaml
```

依赖服务：Kafka、Redis、auth-service、chat-service。

## 常见问题

- **invalid token**：确认 token 正确且 auth-service 可用。
- **User not connected**：客户端没有保持 WebSocket 连接或连接已关闭。
- **failed to upgrade**：浏览器跨域或协议头异常，请检查请求头。
