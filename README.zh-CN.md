# gouno

[English](./README.md) | [文档](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/)

---

**轻量级 Go Web 项目启动器**。它为你搭建项目结构、启动流程、Web 层和响应格式——这些每个微服务都要写、但跟业务无关的代码——让你从第一行就专注业务逻辑。

gouno **不是框架**。它不绑定数据库驱动、缓存客户端或消息队列。技术选型由你决定，gouno 只负责其余部分。

```
gouno 负责的                         gouno 不负责的
├── 项目结构（DDD）                   ├── 数据库（pgx? gorm? ent?）
├── CLI + 配置（Cobra + Viper）       ├── 缓存（redis? memcached?）
├── Web 引擎（Gin + 中间件）          ├── 消息队列（kafka? rabbitmq?）
├── 响应格式（统一 JSON）             └── 认证（JWT? OAuth? session?）
└── 代码生成器 + 模板集
```

## 快速开始

```bash
# 安装
go install github.com/rushairer/gouno-cli@latest

# 创建、构建、运行
gouno-cli new my-service -m github.com/you/my-service
cd my-service && go mod tidy && make dev
# → http://localhost:8080
```

## 代码生成

```bash
gouno gen suite user   # → domain + repository + service
gouno gen task send_email
gouno gen controller auth
```

[完整指南 →](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/code-generation.md)

## 模板集

自定义 `gouno gen` 的输出。不同团队、不同代码风格——无需修改 gouno 源码。

```bash
gouno-cli template install gorm https://github.com/myorg/gouno-template-gorm
gouno-cli new order-service --template-set gorm -m github.com/myorg/order-service
```

[创建自定义模板集 →](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/template-sets.md)

## 文档

| 指南 | 说明 |
|------|------|
| [快速开始](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/getting-started.md) | 安装、创建项目、运行 |
| [代码生成](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/code-generation.md) | 生成 DDD 模块 |
| [模板集](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/template-sets.md) | 创建和分享自定义模板 |
| [配置管理](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/configuration.md) | 多环境 YAML 配置 |
| [中间件](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/middleware.md) | 内置中间件与自定义扩展 |

## 设计理念

**gouno 是启动器，不是框架。** 它给你一个标准化的起点，然后让开。真正的定制通过**模板集**完成——你的代码风格、技术选型、团队规范——全部封装为可复用的模板。

```
gouno（核心）         → 标准化：项目结构 + 启动流程 + Web 层 + 响应格式
模板集                → 定制化：代码风格 + 技术选型 + 团队规范
你的业务代码          → 实现：产品逻辑
```

## 相关项目

| 仓库 | 说明 |
|------|------|
| [gouno](https://github.com/rushairer/gouno) | 核心库（本仓库） |
| [gouno-cli](https://github.com/rushairer/gouno-cli) | CLI 工具 |
| [gouno-template](https://github.com/rushairer/gouno-template) | 默认模板集 |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | 文档 |

## 许可证

MIT License。
