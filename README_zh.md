# Nightingale MCP Server

[English](README.md) | 中文

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![MCP](https://img.shields.io/badge/MCP-Compatible-green.svg)](https://modelcontextprotocol.io/)

[Nightingale](https://github.com/ccfos/nightingale) 的 MCP Server。此 MCP Server 允许 AI 助手通过自然语言与夜莺 API 交互，实现告警管理、监控和可观测性任务。

## 兼容性

- **Nightingale**：v8.0.0+

## 主要用途

- **告警管理**：查询活跃告警和历史告警，查看告警规则和订阅
- **目标监控**：浏览和搜索被监控的主机，分析目标状态
- **事件响应**：创建和管理告警屏蔽规则、通知规则和事件流水线
- **团队协作**：查询用户、团队和业务组

## 快速开始

### 1.获取 API Token
1. 确保在 config.toml 中，启用了 HTTP.TokenAuth
  ```toml
    [HTTP.TokenAuth]
    Enable = true
  ```
2. 登录夜莺 Web 界面
3. 进入 **个人设置** > **个人信息** > **Token 管理**
4. 创建一个具有适当权限的新 Token

![image-20260205172354525](./doc/img/image-20260205172354525.png)

> **安全提示**：请妥善保管 API Token。切勿将 Token 提交到版本控制系统。请使用环境变量或安全的密钥管理系统。

### 2.与 MCP 客户端配合使用

#### Cursor（stdio 模式，默认）

在 `~/.cursor/mcp.json` 中添加：

```json
{
  "mcpServers": {
    "nightingale": {
      "command": "npx",
      "args": ["-y", "@n9e/n9e-mcp-server", "stdio"],
      "env": {
        "N9E_TOKEN": "your-api-token",
        "N9E_BASE_URL": "http://your-n9e-server:17000"
      }
    }
  }
}
```

#### HTTP 模式（可选）

以 HTTP 方式运行服务端（MCP streamable 传输，仅 JSON，无 SSE）时，使用 `http` 子命令启动。

**共享模式 vs 非共享（仅 HTTP）：**

- **`--shared=false`**（默认）：启动时可不提供 token 和 base URL。每个客户端在 mcp.json 里通过 `X-User-Token`、`X-N9e-Base-Url` 提供自己的夜莺身份或实例；若启动时设置了 `N9E_TOKEN` 和 `N9E_BASE_URL`，则作为默认值，客户端仍可通过 header 覆盖。
- **`--shared=true`**：启动时**必须**设置 `N9E_TOKEN` 和 `N9E_BASE_URL`。服务端仅使用该配置，**忽略**客户端请求头中的 `X-User-Token` 和 `X-N9e-Base-Url`。适用于组织统一提供的 MCP 服务、不允许用户覆盖凭证的场景。

```bash
# 非共享：由用户在 mcp.json 提供 token/URL（或启动时设默认值）
n9e-mcp-server http --listen :8080

# 共享：统一凭证，启动时必填，忽略客户端 header
N9E_TOKEN=xxx N9E_BASE_URL=https://n9e.example.com n9e-mcp-server http --listen :8080 --shared
```

**Cursor 连接 HTTP 服务端**

若服务端已以 HTTP 模式运行（例如在 `http://localhost:8080`），在 `~/.cursor/mcp.json` 中添加以 URL 方式配置的条目（无需 `command`/`args`，Cursor 会使用 streamable HTTP 传输）。

**Token 传递**：二选一即可，不必在 mcp.json 里传 token（只要服务端启动时配了 `N9E_TOKEN`）。

1. **仅服务端启动时**：启动时设置 `N9E_TOKEN`（如 `N9E_TOKEN=xxx ./n9e-mcp-server http`），所有连接该服务的客户端都会用这个 token，Cursor 里**无需**配置任何 header。
2. **客户端请求头（可选）**：可携带：
   - `X-User-Token`：用该 token 调夜莺 API，替代启动时的 `N9E_TOKEN`；
   - `X-N9e-Base-Url`：用该 URL 作为夜莺 API 地址（如 `https://n9e.other-env.com`），替代服务端启动时的 `N9E_BASE_URL`。
   这样每人可使用自己的 token 或指向不同夜莺环境（或同时覆盖两者）。

若服务端**已用** `N9E_TOKEN` 启动（Cursor 里不必写 header）：

```json
{
  "mcpServers": {
    "nightingale": {
      "url": "http://localhost:8080"
    }
  }
}
```

若由 Cursor 通过请求头传 token 和/或夜莺地址（例如服务未设 `N9E_TOKEN`、或连到不同夜莺环境时）：

```json
{
  "mcpServers": {
    "nightingale": {
      "url": "http://localhost:8080",
      "headers": {
        "X-User-Token": "你的夜莺-api-token",
        "X-N9e-Base-Url": "http://your-n9e-server:17000"
      }
    }
  }
}
```

可只写其中一个 header；未写的项会使用服务端启动时的 `N9E_TOKEN` / `N9E_BASE_URL`。**若服务端以 `--shared` 启动，则这些 header 会被忽略，请勿依赖。**

若 MCP 服务前还有网关等需要认证，可同时配置对应 headers（如 `Authorization: Bearer your-gateway-token`）。服务端仅用 `X-User-Token` 作为调用夜莺 API 的凭证。

### 3.重启 Cursor 等进程，即可使用

## 可用工具

| 工具集 | 工具 | 说明 |
|-------|------|------|
| alerts | `list_active_alerts` | 列出当前活跃告警，支持过滤条件 |
| alerts | `get_active_alert` | 根据事件 ID 获取活跃告警详情 |
| alerts | `list_history_alerts` | 列出历史告警，支持过滤条件 |
| alerts | `get_history_alert` | 获取历史告警详情 |
| alerts | `list_alert_rules` | 列出业务组的告警规则 |
| alerts | `get_alert_rule` | 获取告警规则详情 |
| targets | `list_targets` | 列出被监控主机/目标，支持过滤条件 |
| datasource | `list_datasources` | 列出所有可用数据源 |
| mutes | `list_mutes` | 列出业务组的告警屏蔽规则 |
| mutes | `get_mute` | 获取告警屏蔽规则详情 |
| mutes | `create_mute` | 创建告警屏蔽规则 |
| mutes | `update_mute` | 更新告警屏蔽规则 |
| notify_rules | `list_notify_rules` | 列出所有通知规则 |
| notify_rules | `get_notify_rule` | 获取通知规则详情 |
| alert_subscribes | `list_alert_subscribes` | 列出业务组的告警订阅 |
| alert_subscribes | `list_alert_subscribes_by_gids` | 列出多个业务组的订阅 |
| alert_subscribes | `get_alert_subscribe` | 获取订阅详情 |
| event_pipelines | `list_event_pipelines` | 列出所有事件流水线 |
| event_pipelines | `get_event_pipeline` | 获取事件流水线详情 |
| event_pipelines | `list_event_pipeline_executions` | 列出指定流水线的执行记录 |
| event_pipelines | `list_all_event_pipeline_executions` | 列出所有流水线的执行记录 |
| event_pipelines | `get_event_pipeline_execution` | 获取执行记录详情 |
| users | `list_users` | 列出用户，支持过滤条件 |
| users | `get_user` | 获取用户详情 |
| users | `list_user_groups` | 列出用户组/团队 |
| users | `get_user_group` | 获取用户组详情（包含成员） |
| busi_groups | `list_busi_groups` | 列出当前用户可访问的业务组 |

## 示例提示词

配置完成后，您可以使用自然语言与夜莺交互：

- "显示过去 24 小时内所有紧急告警"
- "当前有哪些告警正在触发？"
- "列出所有离线超过 5 分钟的监控目标"
- "业务组 1 配置了哪些告警规则？"
- "由于维护原因，为 service=api 的告警创建一个 2 小时的屏蔽规则"
- "查看事件流水线的执行历史"
- "运维团队有哪些成员？"

## 配置

### 运行模式

- **stdio**（默认）：通过 stdin/stdout 进行 MCP 通信。适用于 Cursor 等会拉起服务进程的客户端。
- **http**：通过 HTTP 使用 MCP streamable 传输（仅 JSON 请求/响应，无 SSE）。使用 `n9e-mcp-server http` 启动，客户端需支持 streamable HTTP（如 `StreamableClientTransport`）。

### 环境变量

| 变量 | 命令行参数 | 说明 | 默认值 |
|-----|-----------|------|-------|
| `N9E_TOKEN` | `--token` | 夜莺 API Token（必需） | - |
| `N9E_BASE_URL` | `--base-url` | 夜莺 API 地址 | `http://localhost:17000` |
| `N9E_READ_ONLY` | `--read-only` | 禁用写操作 | `false` |
| `N9E_TOOLSETS` | `--toolsets` | 启用的工具集（逗号分隔） | `all` |
| `N9E_LISTEN` | `--listen` | HTTP 模式：监听地址 | `:8080` |
| `N9E_SESSION_TIMEOUT` | `--session-timeout` | HTTP 模式：空闲会话超时（0 表示不超时） | `0` |
| `N9E_SHARED` | `--shared` | HTTP 模式：为 true 时启动必须提供 N9E_TOKEN 和 N9E_BASE_URL，并忽略客户端 header | `false` |

### 工具集选择

默认启用所有工具集。可以通过 `--toolsets` 参数或 `N9E_TOOLSETS` 环境变量只启用需要的工具集，减少暴露给 AI 助手的工具数量，节省上下文窗口的 token 消耗。

可用工具集：`alerts`、`targets`、`datasource`、`mutes`、`busi_groups`、`notify_rules`、`alert_subscribes`、`event_pipelines`、`users`

例如，只启用告警和监控目标相关工具：

```json
{
  "mcpServers": {
    "nightingale": {
      "command": "npx",
      "args": ["-y", "@n9e/n9e-mcp-server", "stdio"],
      "env": {
        "N9E_TOKEN": "your-api-token",
        "N9E_BASE_URL": "http://your-n9e-server:17000",
        "N9E_TOOLSETS": "alerts,targets"
      }
    }
  }
}
```

## 开源协议

Apache License 2.0

## 相关项目

- [Nightingale](https://github.com/ccfos/nightingale) - 企业级云原生监控系统
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - 官方 MCP Go SDK
- [MCP 规范](https://modelcontextprotocol.io/) - Model Context Protocol 规范
