# Yuuna Danmu

这是一个用 Go 语言编写的轻量级 B 站直播弹幕监听工具。基于 **Wails** 和 **Svelte** 前端。

参考文档基本出自 [哔哩哔哩 - API 收集整理](https://github.com/SocialSisterYi/bilibili-API-collect)

### 核心功能

* **双模运行**：支持命令行核心逻辑，同时也提供直观的 **Wails GUI** 桌面客户端。
* **实时弹幕读取**：秒级获取直播间观众发送的弹幕、牌子信息。
* **礼物记录**：自动识别送礼行为。
* **配置持久化**：自动保存房间号及 Cookie，支持跨平台路径。

### 技术栈

* **Backend**: Go 1.25
* **Frontend**: Svelte + HTML/CSS (Via Wails)
* **Communication**: WebSocket (Gorilla) & JSON-RPC

### 如何使用

#### 开发环境准备

1. 安装 [Wails CLI](https://wails.io/docs/gettingstarted/installation)。
2. 确保已安装 Node.js 和 Go 环境。

#### 运行与编译

1. **开发模式**（支持前端热更新）：
```bash
wails dev

```

2. **编译正式版**：
```bash
wails build

```

编译后的程序将出现在 `build/bin` 目录下。

### 界面预览

![img](./doc/preview.png)

### 开发计划

* [x] **Wails 界面集成**：CLI 到 GUI。
* [ ] **更多消息类型**：完善舰长开通、醒目留言 (SC) 等解析。
* [ ] **自动重连机制**：针对 WebSocket 网络波动的自动恢复。
* [ ] **主题自定义**：支持针对不同直播风格的皮肤切换。
