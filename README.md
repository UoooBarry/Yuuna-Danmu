# Yuuna Danmu

轻量级 B 站直播弹幕监听工具，使用 Go 后端，基于 Wails 框架结合 Svelte 前端实现。

API 参考主要来自 [哔哩哔哩-API-收集整理](https://github.com/SocialSisterYi/bilibili-API-collect)。

### 核心功能

- 实时接收直播间弹幕及相关信息（包括牌子、礼物等）
- 自动记录礼物投喂
- 配置自动保存（房间号、Cookie），跨平台路径处理
- 支持 Windows、Linux、macOS

### 技术栈

- 后端：Go 1.25
- 前端：Svelte（通过 Wails 集成）
- 通信：WebSocket（gorilla/websocket）与 JSON-RPC
- 扩展接口：gRPC（proto 文件位于 `api/grpc/pb`）

设计目标：保持简洁、无多余配置，支持跨平台，提供 gRPC 接口便于第三方扩展或集成。

### 使用方法

#### 开发环境

1. 安装 [Wails CLI](https://wails.io/docs/gettingstarted/installation)
2. 准备好 Node.js 和 Go 环境

#### 运行与打包

- 开发模式（支持前端热重载）：
  ```bash
  wails dev
  ```

- 打包正式版：
  ```bash
  wails build
  ```
  可执行文件会生成在 `build/bin` 目录。

#### 配置

配置文件路径（程序首次运行时自动创建）：

- Windows：`%APPDATA%\yuuna-danmu\config.json`
- macOS：`~/Library/Application Support/yuuna-danmu/config.json`
- Linux：`~/.config/yuuna-danmu/config.json`

示例内容：
```json
{
  "room_id": 1,
  "cookie": "",
  "debug": false
}
```
开发时可将 `debug` 设为 `true` 以开启调试模式。

### 界面预览

![preview](./doc/preview.png)

### 已完成

- [X] Wails GUI 界面
- [X] 礼物显示（含连击、图标）
- [X] WebSocket 自动重连
- [X] gRPC 服务支持
- [X] Cookie 自动刷新
- [X] 更多消息类型解析（舰长、醒目留言等）

### 计划中

- [] 主题/皮肤切换
- [] 插件系统
