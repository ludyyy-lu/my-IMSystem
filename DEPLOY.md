# my-IMSystem 项目部署指南

## 目录

1. [系统架构概览](#1-系统架构概览)
2. [前置依赖](#2-前置依赖)
3. [配置文件说明与校验结果](#3-配置文件说明与校验结果)
4. [本地开发部署（推荐）](#4-本地开发部署推荐)
5. [端口规划](#5-端口规划)
6. [服务启动顺序](#6-服务启动顺序)
7. [常见问题排查](#7-常见问题排查)

---

## 1. 系统架构概览

```
客户端 (WebSocket)
      │
      ▼
ws-gateway  (:8888)        ← HTTP/WS 接入网关，负责认证、消息路由、推送
      │            │
      │ gRPC       │ gRPC
      ▼            ▼
auth-service     chat-service (message.rpc)
   (:8080)           (:8081)
                      │ gRPC
                      ▼
               friend-service
                   (:8082)

user-service (:8083)       ← 独立用户信息 RPC 服务

──── 基础设施 ────
MySQL      :3307  (宿主机) / :3306 (容器内)
Redis      :6379
etcd       :2379
Kafka      :9094  (宿主机外部监听) / :9092 (容器内部)
ZooKeeper  :2181
```

### 消息流转路径

#### 发送消息（客户端 → 服务端）

```
客户端 WS.send({"type":"chat","to":2,"content":"hello"})
  → ws-gateway readLoop 收帧
  → router.HandleMessage 解析 type="chat"
  → kafka.SendMessage("im-chat-topic", msg)      ← 发布到 Kafka
  → chat-service Kafka Consumer 消费
  → 持久化到 MySQL messages 表
  → (同时) ws-gateway Kafka Consumer 消费
  → push.Service.PushToUser(toUserId, ...)
  → session.Manager.SendTo(toUserId, data)
  → session.writeLoop → conn.WriteMessage        ← 推送到目标客户端
```

#### 好友事件流（申请/接受/拒绝）

```
客户端 HTTP → friend-service
  → 处理业务逻辑 + 写 MySQL
  → kafka.SendMessage("im-friend-topic", event)
  → ws-gateway FriendTopic Consumer 消费
  → push.Service.PushToUser(toUserId, "friend_event", event)
  → 目标客户端实时收到好友通知
```

---

## 2. 前置依赖

| 工具 | 版本要求 | 用途 |
|------|---------|------|
| Go | ≥ 1.23 | 编译运行服务 |
| Docker | ≥ 24 | 运行基础设施依赖 |
| Docker Compose | v2 | 编排基础设施 |

---

## 3. 配置文件说明与校验结果

### 3.1 配置文件汇总

| 服务 | 配置文件 | 监听端口 |
|------|---------|---------|
| auth-service | `auth-service/etc/auth.yaml` | gRPC :8080 |
| chat-service | `chat-service/etc/message.yaml` | gRPC :8081 |
| friend-service | `friend-service/etc/friend.yaml` | gRPC :8082 |
| user-service | `user-service/etc/user.yaml` | gRPC :8083 |
| ws-gateway | `ws-gateway/etc/ws.yaml` | HTTP :8888 |

### 3.2 已修复的配置问题

在本次审查中发现并修复了以下问题：

#### ❌ Bug 1 — 所有 RPC 服务端口冲突（已修复）

原始配置中 auth / chat / friend / user 四个 RPC 服务全部配置了
`ListenOn: 0.0.0.0:8080`，导致后启动的服务绑定端口失败直接崩溃。

**修复**：为每个服务分配独立端口：
- auth-service → `:8080`
- chat-service (message.rpc) → `:8081`
- friend-service → `:8082`
- user-service → `:8083`

#### ❌ Bug 2 — ws-gateway ChatRpcConf 格式错误（已修复）

原始 `ws.yaml` 中 `ChatRpcConf` 使用了错误的格式：
```yaml
# 错误：将 etcd 地址填入了 Endpoints（直连 gRPC 字段），且缺少 Key
ChatRpcConf:
  Endpoints:
    - 127.0.0.1:2379
  Timeout: 3000
```
`Endpoints` 字段是 go-zero 直连 gRPC 的字段（填服务实际监听地址），
而 `127.0.0.1:2379` 是 etcd 的端口，不是 gRPC 服务端口。
同时缺少 `Key` 导致 etcd 服务发现无法工作。

**修复**：使用正确的 etcd 服务发现格式：
```yaml
ChatRpcConf:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: message.rpc
  Timeout: 3000
```

#### ❌ Bug 3 — ws-gateway AuthRpcConf Key 层级错误（已修复）

原始配置：
```yaml
# 错误：Key 不是 RpcClientConf 的顶级字段，会被 go-zero 静默忽略
AuthRpcConf:
  Endpoints:
    - 127.0.0.1:2379
  Key: auth.rpc
```

**修复**：
```yaml
AuthRpcConf:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: auth.rpc
```

#### ❌ Bug 4 — friend-service 硬编码 Kafka 地址（已修复）

`friend-service/internal/svc/service_context.go` 中 Kafka 地址硬编码为
`[]string{"kafka:9092"}`，忽略了配置文件中的 `Kafka.Brokers` 字段。

**修复**：改为 `kafka.InitKafkaProducer(c.Kafka.Brokers)`。

#### ❌ Bug 5 — Kafka 从宿主机无法连接（已修复）

原始 docker-compose 只配置了容器内部监听 `PLAINTEXT://kafka:9092`，
宿主机服务连接 `localhost:9092`（或映射的 9093）后，Kafka 会告知客户端
重连到 `kafka:9092`，但宿主机无法解析 `kafka` 域名，导致连接失败。

**修复**：docker-compose 增加外部监听器：
```yaml
KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,EXTERNAL://:9094
KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,EXTERNAL://localhost:9094
KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,EXTERNAL:PLAINTEXT
KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT
```
宿主机服务使用 `localhost:9094` 连接 Kafka（已更新所有服务配置）。

---

## 4. 本地开发部署（推荐）

> **运行模式**：基础设施（MySQL / Redis / etcd / Kafka）跑在 Docker 中，
> Go 服务进程直接在宿主机运行。

### 4.1 启动基础设施

```bash
# 进入 deploy 目录
cd deploy

# 首次启动（下载镜像 + 创建容器）
docker compose up -d

# 检查所有容器状态
docker compose ps
```

预期输出：

```
NAME            STATUS
im-mysql        Up
im-redis        Up
im-etcd         Up
im-zookeeper    Up
im-kafka        Up
```

### 4.2 创建 Kafka Topic

Kafka 容器启动后需要预先创建两个 Topic（只需执行一次）：

```bash
# im-chat-topic：聊天消息
docker exec im-kafka kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --topic im-chat-topic \
  --partitions 3 --replication-factor 1

# im-friend-topic：好友事件
docker exec im-kafka kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --topic im-friend-topic \
  --partitions 3 --replication-factor 1

# 验证 Topic 已创建
docker exec im-kafka kafka-topics.sh \
  --bootstrap-server localhost:9092 --list
```

### 4.3 验证基础设施连通性

```bash
# MySQL
mysql -h 127.0.0.1 -P 3307 -u imuser -pim123456 im -e "SELECT 1"

# Redis
redis-cli -h 127.0.0.1 ping   # 期望: PONG

# etcd
etcdctl --endpoints=127.0.0.1:2379 endpoint health

# Kafka（从宿主机连接外部端口）
docker exec im-kafka kafka-topics.sh \
  --bootstrap-server localhost:9092 --list
```

### 4.4 编译并启动各 Go 服务

每个服务在独立终端窗口中启动，**务必按照以下顺序**（见第 6 节说明）。

```bash
# 终端 1 — auth-service（其他服务通过 etcd 发现它）
cd auth-service
go run auth.go

# 终端 2 — friend-service（chat-service 依赖它）
cd friend-service
go run friend.go

# 终端 3 — chat-service
cd chat-service
go run message.go

# 终端 4 — user-service（独立，可任意顺序）
cd user-service
go run user.go

# 终端 5 — ws-gateway（最后启动，依赖 auth + chat）
cd ws-gateway
go run ws.go
```

服务启动成功的日志示例：

```
# auth-service
Starting rpc server at 0.0.0.0:8080...

# chat-service
Starting rpc server at 0.0.0.0:8081...

# friend-service
Starting rpc server at 0.0.0.0:8082...

# user-service
Starting rpc server at 0.0.0.0:8083...

# ws-gateway
Starting server at 0.0.0.0:8888...
```

---

## 5. 端口规划

| 组件 | 宿主机端口 | 容器端口 | 说明 |
|------|-----------|---------|------|
| ws-gateway | 8888 | — | WebSocket / HTTP 接入 |
| auth-service | 8080 | — | gRPC，注册到 etcd key=`auth.rpc` |
| chat-service | 8081 | — | gRPC，注册到 etcd key=`message.rpc` |
| friend-service | 8082 | — | gRPC，注册到 etcd key=`friend.rpc` |
| user-service | 8083 | — | gRPC，注册到 etcd key=`user.rpc` |
| MySQL | 3307 | 3306 | 宿主机通过 3307 访问 |
| Redis | 6379 | 6379 | |
| etcd | 2379 | 2379 | |
| Kafka（内部） | — | 9092 | 仅供容器内（服务间）使用 |
| Kafka（外部） | 9094 | 9094 | 宿主机 Go 服务连接此端口 |
| ZooKeeper | 2181 | 2181 | |

---

## 6. 服务启动顺序

依赖关系决定了启动顺序：

```
基础设施（MySQL + Redis + etcd + Kafka）
         ↓
     auth-service          ← ws-gateway 启动时需要从 etcd 发现它
         ↓
    friend-service         ← chat-service 启动时需要从 etcd 发现它
         ↓
     chat-service          ← ws-gateway 启动时需要从 etcd 发现它
         ↓
     user-service          ← 独立服务，可与 chat-service 并行
         ↓
     ws-gateway            ← 最后启动
```

> **说明**：go-zero zrpc 客户端在启动时会从 etcd 解析依赖服务地址。
> 若依赖服务尚未注册到 etcd，`MustNewClient` 会 panic 并终止进程。
> 因此必须先确保被依赖的服务已完全启动并注册到 etcd。

---

## 7. 常见问题排查

### Q1: `bind: address already in use`

说明有另一个进程已占用该端口。检查并杀掉冲突进程：

```bash
# 查找占用 8081 端口的进程
lsof -i :8081
kill <PID>
```

### Q2: `failed to connect DB`

- 检查 MySQL 容器是否正常运行：`docker compose ps`
- 确认端口映射：宿主机 3307 → 容器 3306
- 手动测试连接：`mysql -h 127.0.0.1 -P 3307 -u imuser -pim123456`
- 确认数据库 `im` 已自动创建（由 docker-compose 环境变量 `MYSQL_DATABASE` 控制）

### Q3: Kafka 连接超时 / `dial tcp: no route to host`

- 确认 `im-kafka` 容器运行中：`docker compose ps`
- 宿主机服务配置使用的是 `localhost:9094`（外部监听端口），不是 `kafka:9092`
- 验证 9094 端口可达：`nc -zv localhost 9094`

### Q4: etcd 服务发现失败 / `context deadline exceeded`

- 检查 etcd 容器：`docker compose ps`
- 验证连接：`etcdctl --endpoints=127.0.0.1:2379 endpoint health`
- 检查对应服务（如 auth-service）是否已启动并成功注册到 etcd

### Q5: WebSocket 连接被拒绝 / `invalid token`

- 确认 auth-service 已启动
- 获取 token：先通过 auth-service 的 `Register` 或 `Login` 接口获取 `access_token`
- 连接时在 `Authorization` 头携带：`Authorization: Bearer <token>`
  或通过 query 参数：`ws://localhost:8888/ws/connect?token=<token>`

### Q6: 消息发送后对方未收到

按以下顺序排查：

1. **Kafka Topic 是否存在**：`docker exec im-kafka kafka-topics.sh --bootstrap-server localhost:9092 --list`
2. **ws-gateway 日志**：查看是否有 `Kafka chat topic is not configured` 错误（说明 `ws.yaml` 中 `Kafka.Topic` 为空）
3. **chat-service 日志**：查看 Kafka 消费者是否正常消费（`[Kafka] Chat consumer started.`）
4. **接收方是否在线**：消息路由到 `session.Manager.SendTo`，若接收方不在线则投递失败；离线消息需接收方重连后通过 `OfflineStore` 拉取

### Q7: `Kafka writer is not initialized`

common/kafka 的 `InitKafkaProducer` 未被调用。各服务的 `ServiceContext` 初始化时会调用此方法；若在测试中直接调用业务逻辑，需先手动初始化：

```go
kafka.InitKafkaProducer([]string{"localhost:9094"})
```
