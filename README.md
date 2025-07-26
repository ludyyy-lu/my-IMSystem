# 💬my-IMSystem - 基于 Go-Zero 分布式 IM 即时通讯系统

一个基于 **Go 语言 + 微服务架构** 构建的高性能即时通讯系统，支持注册登录、好友管理、实时聊天、群聊功能、消息持久化与 Kafka 异步分发。系统模块清晰，服务解耦，部署灵活，具备企业级架构雏形。

---

## 📌 项目目标

本项目旨在构建一个结构合理、功能完善、具有工程复杂度的 IM 即时通讯后端系统，以实践 Go 微服务开发中的多种核心能力，包括：

- 微服务架构设计与拆分
- 实时通信与消息分发机制
- Kafka 高吞吐异步消息系统
- Redis 缓存设计与用户状态管理
- WebSocket 长连接通信管理
- Docker 容器化部署

---

## 📦 核心模块设计

| 服务名                | 模块职责                              | 技术栈                              |
| ------------------ | --------------------------------- | -------------------------------- |
| **user-service**   | 用户注册、登录、资料管理、用户状态查询               | MySQL + gRPC + GORM              |
| **auth-service**   | JWT 身份验证、Token 签发与刷新、多设备会话管理      | Redis + JWT + gRPC               |
| **friend-service** | 好友添加、删除、搜索、拉黑、验证、通知推送             | MySQL + Redis + gRPC             |
| **chat-service**   | 聊天消息处理：存储、分发、离线缓存、历史记录查询          | Kafka + Redis + MySQL            |
| **group-service**  | 群聊创建、成员管理、群权限控制、禁言、公告等            | MySQL + Redis + gRPC             |
| **ws-gateway**     | WebSocket 长连接通信、消息推送、心跳检测、用户上下线通知 | WebSocket + Redis Pub/Sub + gRPC |
| **api-gateway**    | 所有 RESTful 接口的统一入口、请求转发、限流、鉴权     | go-zero API 网关组件                 |


---

## 🛠️ 技术选型

| 类别 | 技术 | 用途 |
|------|------|------|
| 核心语言 | Go | 高性能后端语言 |
| 微服务框架 | go-zero | 快速构建 RPC / API 微服务 |
| 通信协议 | gRPC + Protobuf | 微服务内部通信 |
| 实时通信 | WebSocket | 客户端连接 & 消息推送 |
| 消息系统 | Kafka | 消息异步分发、削峰、解耦 |
| 缓存组件 | Redis | 用户状态、Session 缓存、未读消息等 |
| 数据存储 | MySQL | 用户/聊天/关系数据持久化 |
| ORM 框架 | GORM | 简化数据库访问 |
| 服务注册 | Etcd | 服务发现、健康检查 |
| 配置管理 | YAML | 每个服务独立配置文件 |
| 容器部署 | Docker + Compose | 本地一键部署 |

---

## 🗂️ 数据结构设计

### ✅ MySQL 数据表（核心）

- `users`：用户信息表
- `friends`：好友关系
- `groups`：群组基本信息
- `group_members`：群成员表
- `messages`：聊天记录表（单聊/群聊）

### ✅ Redis 缓存结构

| Key | Value | 用途 |
|-----|-------|------|
| `user:online:{uid}` | bool | 是否在线 |
| `auth:session:{uid}` | session 信息 | 登录设备 & token 信息 |
| `unread:{uid}` | 消息列表 | 未读消息缓存 |

### ✅ Kafka Topic

| Topic | 用途 |
|-------|------|
| `chat_message` | 聊天消息（异步发给接收方） |
| `user_status` | 用户上线/下线广播（可选） |
| `chat_ack` | 消息状态回执（可选） |

---

## 🧱 系统架构图（简图）
```
[ Web / App 客户端 ]
        ↓  (HTTP / WS)
    [API Gateway + WS Gateway]
        ↓
 ┌────────────┬─────────────┬─────────────┐
 │ user-svc   │ friend-svc  │ group-svc   │
 └────────────┴─────────────┴─────────────┘
        ↓
    [auth-svc + chat-svc]
        ↓
 ┌───────────────┬──────────────┐
 │ Kafka         │ Redis        │
 └───────────────┴──────────────┘
        ↓
      [MySQL]
```

---

## 🪜 开发计划与阶段划分

| 阶段                     | 模块                             | 目标                       | 预计时间   |
| ---------------------- | ------------------------------ | ------------------------ | ------ |
| **阶段一：系统初始化**          | 项目结构搭建、Docker 环境、go-zero 脚手架配置 | 构建可运行的 go-zero 微服务骨架     | 1 天    |
| **阶段二：用户体系**           | `user` + `auth` 模块             | 实现注册/登录、JWT鉴权、多设备登录支持    | 2\~3 天 |
| **阶段三：WebSocket 通信网关** | `ws-gateway` 模块                | 建立长连接、身份校验、上下线状态管理       | 2 天    |
| **阶段四：聊天功能构建**         | `chat` 模块 + Kafka              | 消息发送、离线投递、Kafka 流转、存储、回执 | 3\~4 天 |
| **阶段五：好友/群组模块**        | `friend` + `group` 模块          | 添加好友、验证消息、群聊功能等          | 2\~3 天 |
| **阶段六：部署上线**           | Docker Compose 全链路联调           | 支持一键启动整个 IM 系统           | 1 天    |
| **阶段七：前端或API文档扩展（可选）** | Vue / React 或 Swagger 文档       | 给出配套接口使用说明               | 1 天    |

---

## 🚀 可选扩展功能

| 功能          | 技术点                                 |
| ----------- | ----------------------------------- |
| 离线消息通知      | Kafka + Redis 延迟队列 / 缓存未读消息         |
| 消息回执（已读/未读） | ACK 机制，可使用 Kafka + Redis + 状态位      |
| 群组权限管理      | 群主/管理员权限控制、踢人、禁言                    |
| 多端同步        | WebSocket 连接管理 + session 机制         |
| 监控与追踪       | Prometheus + Grafana / Jaeger（链路追踪） |


---

## 🐳 快速部署说明（待完成）

项目支持通过 Docker Compose 一键部署：

```bash
docker-compose up --build
```

默认启动以下服务：
go-zero API 网关
WebSocket 网关服务
user / auth / friend / chat / group 各服务
Kafka + Zookeeper
Redis
MySQL
Etcd

## 📖 使用说明（后续补充）
- API 接口使用文档（Swagger）

- 消息格式说明（Proto 消息结构）

- WebSocket 协议交互规范

## 🤝 开发与维护
- 项目开发者：Ludy（superkiwi）
- 协助与架构顾问：ChatGPT（aka Riki）

- 本项目作为后端架构展示型项目，欢迎提出 issue 或 PR 交流。