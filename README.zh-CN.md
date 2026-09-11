# gouno

[English](./README.md) | [文档](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/)

---

**轻量级 Go Web 项目启动器**。Gouno 提供可复用的启动、HTTP、响应、安全与项目工具机制；具体项目结构和开发约定由 Template 决定。

gouno **不是全栈框架**。它不规定数据库、缓存、消息队列、认证产品，也不规定业务项目必须采用 DDD、Clean Architecture 或某一种分层方式。

```text
Gouno 负责                          Template / 业务项目负责
├── 可复用运行时基础能力             ├── 项目架构
├── Codegen 协议与执行引擎            ├── Generator 清单与代码模板
├── 安全的模板渲染                    ├── 数据库/缓存/消息队列选型
├── 中间件与安全基础能力              ├── 认证实现
└── 通用项目工具机制                  └── 应用/领域策略
```

## 快速开始

```bash
go install github.com/rushairer/gouno-cli@latest

gouno-cli new my-service -m github.com/you/my-service
cd my-service && make dev
```

## Template 定义代码生成

代码生成属于 **Template Capability**，不再由 Gouno Core 内置某一套 DDD 代码结构。

如果当前 Template 提供 `.gouno/codegen.yaml`，它可以定义：

```bash
gouno gen suite user
gouno gen task send_email
gouno gen controller auth
```

另一个 Template 完全可以定义 `handler`、`usecase`、`module` 等不同命令；不需要代码生成的 Template 可以完全不提供 manifest，也就无需在项目 CLI 中出现 `gen`。

Gouno 只负责 manifest 发现与校验、动态构建 CLI、参数/Flag 解析、安全渲染、Generator 组合、覆盖策略以及 Go 代码格式化。具体 Generator 名称、参数、输出目录和源码模板全部属于 Template 的开发规范。

详见 [Gouno Template Codegen Specification v1](./docs/codegen-template-spec.md)。

## 设计理念

**Gouno 是启动器和项目工具宿主，不是业务框架。** Gouno 统一机制，Template 表达策略。

```text
gouno-cli           -> 从任意完整 Project Template 创建项目
gouno               -> 提供可复用机制 + Template 驱动的项目工具协议
Project Template    -> 项目初始代码 + 架构选择 + Codegen Policy
业务项目             -> 产品与领域逻辑
```

生成后的项目可以按需依赖 Gouno 提供的运行时基础能力，但 Gouno 不拥有业务项目的架构观点。`domain`、`repository`、`service`、`controller`、`task`、`suite` 等名称都属于 Template 的词汇，而不是 Gouno Core 的固定概念。

## 相关项目

| 仓库 | 说明 |
|------|------|
| [gouno](https://github.com/rushairer/gouno) | 核心库与项目工具协议（本仓库） |
| [gouno-cli](https://github.com/rushairer/gouno-cli) | 项目创建 CLI |
| [gouno-template](https://github.com/rushairer/gouno-template) | 默认项目模板与默认 Codegen Policy |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | 文档 |

## 许可证

MIT License。
