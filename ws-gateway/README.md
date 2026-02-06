# ws-gateway

`ws-gateway` 是 IM 系统的 WebSocket 网关，负责长连接接入、消息路由与 Kafka 推送。

## 架构分层

```
ws-gateway/
├── ws.go                    # 启动入口（只做 wiring）
├── etc/ws.yaml
└── internal/
    ├── transport/           # 传输层（WebSocket）
    │   ├── server.go        # WS server 启动
    │   ├── handler.go       # onConnect / onMessage / onClose
    │   └── protocol.go      # WS 协议解析（JSON）
    ├── session/             # 会话层（核心）
    │   ├── manager.go       # ConnManager：uid -> session
    │   ├── session.go       # 单连接 Session（心跳、写通道）
    │   ├── state.go         # 在线 / 下线 / 重连
    │   └── offline_store.go # 离线消息存储
    ├── push/                # 推送层（只负责“投递”）
    │   ├── service.go       # PushService
    │   └── dispatcher.go    # uid / device / broadcast
    ├── consume/             # 消费层（Kafka）
    │   └── message.go       # 消费 im-message-topic
    ├── model/               # WS 内部消息结构
    │   ├── ws_message.go
    │   └── push_message.go
    ├── rpc/                 # 调用 IM 后端服务
    │   └── auth.go          # token 校验
    └── svc/
        └── service_context.go
```

## 连接与消息流程

1. 客户端连接 `/ws/connect`
2. transport 层解析 token（Authorization / query / Sec-WebSocket-Protocol）
3. rpc/auth 校验 token，创建 session
4. session 管理心跳、读写与离线消息加载
5. 收到消息后由 transport/router 分发至 chat/ack 逻辑
6. consume 从 Kafka 接收推送消息，交由 push 层投递

## 消息格式

客户端发送：

```json
{
  "type": "chat",
  "content": "hello",
  "to": 10002
}
```

服务端推送：

```json
{
  "type": "chat_message",
  "payload": { "from": 10001, "to": 10002, "content": "..." }
}
```

## 配置文件

`etc/ws.yaml` 示例：

```yaml
Name: ws-api
Host: 0.0.0.0
Port: 8888
Redis:
  Addr: localhost:6379
  Password: ""
  DB: 0
Kafka:
  Brokers:
    - kafka:9092
  Topic: im-chat-topic
ChatRpcConf:
  Endpoints:
    - 127.0.0.1:2379
  Timeout: 3000
AuthRpcConf:
  Endpoints:
    - 127.0.0.1:2379
  Key: auth.rpc
```

## 启动方式

```bash
go run ws.go -f etc/ws.yaml
```
