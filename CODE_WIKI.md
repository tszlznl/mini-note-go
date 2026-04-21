# Minimalist Web Notepad - Code Wiki

本文档旨在提供 **Minimalist Web Notepad（极简网页记事本）** 项目的完整代码与架构参考。

## 1. 项目整体架构

本项目是一个极简的、零依赖的在线记事本应用。它采用 **Client-Server（客户端-服务器）** 架构，去除了所有不必要的中间件和数据库。

*   **设计理念**：轻量级、无需登录、纯文件系统存储。每一个 URL 路径都唯一映射到服务器磁盘上的一篇笔记。
*   **存储机制**：笔记以 `.txt` 纯文本文件形式存储于本地文件系统（默认存放在应用所在目录下的 `_tmp` 或容器内的 `/tmp`）。URL 的 ID 部分直接作为文件名。
*   **交互模式**：
    *   前端采用 **单页应用 (SPA)** 模式，不刷新页面即可加载和切换笔记。
    *   采用**自动保存机制**（500ms 防抖），在用户输入时自动异步调用后端的 POST 接口同步内容。
*   **UI 设计**：遵循 `DESIGN.md` 中定义的高级暗黑模式（Linear 风格），仅由单页面内的 HTML/CSS/JS 原生实现，无外部框架依赖。

---

## 2. 主要模块职责

项目代码非常紧凑，后端代码全部分布在三个主要的 Go 源文件中，主要模块职责如下：

### 2.1 入口与服务启动模块 (`main.go`)
*   **配置管理**：通过读取系统环境变量（如 `LISTEN_ADDR`, `DATA_DIR`）进行动态配置。
*   **目录初始化**：在服务启动前，检查并创建存放笔记数据的文件夹。
*   **服务启动**：调用 HTTP 处理器注册逻辑，启动并监听底层的 HTTP 服务。

### 2.2 路由与业务逻辑处理模块 (`handlers.go`)
*   **路由分发**：利用 Go 标准库的 `http.ServeMux` 注册路由规则（`/`, `/list`, `/note/`）。
*   **读写核心逻辑**：执行笔记的增删改查（CRUD）操作中的“读”和“写”，负责本地文件系统的 I/O 操作。
*   **安全与校验**：包含针对非法 URL 字符的正则过滤、内容大小（`MAX_NOTE_SIZE` 默认 1MB）的限制，以及全局的只读模式（`READ_ONLY`）保护。

### 2.3 前端与静态资源模块 (`embed.go`)
*   **资源内联**：利用 Go 的原生能力，将前端完整的 HTML、CSS 和 JavaScript 逻辑作为一个字符串字面量 (`indexHTML`) 硬编码在后端服务中，以实现单文件二进制分发。

---

## 3. 关键类与函数说明

### 后端核心函数 (Go)

*   **`main()`** _([main.go](file:///workspace/main.go#L9-L30))_
    *   程序的入口点。负责解析 `LISTEN_ADDR`（默认 `:20008`）和 `DATA_DIR`（默认 `./_tmp`），并启动 `http.ListenAndServe`。
*   **`init()`** _([handlers.go](file:///workspace/handlers.go#L21-L30))_
    *   在程序启动时自动运行，读取环境变量配置 `MAX_NOTE_SIZE` 与 `READ_ONLY` 并应用到全局变量。
*   **`registerHandlers(dataDir string)`** _([handlers.go](file:///workspace/handlers.go#L32-L44))_
    *   注册 HTTP 路由。将请求挂载到对应的处理函数上。
*   **`handleNote(...)`** _([handlers.go](file:///workspace/handlers.go#L46-L61))_
    *   `/note/` 路由的统一入口处理器。使用 `validIDRegex` 正则验证 ID 是否合法，并根据 HTTP 方法（GET / POST）将请求分发给 `readNote` 或 `writeNote`。
*   **`readNote(...)`** _([handlers.go](file:///workspace/handlers.go#L69-L85))_
    *   **GET** 请求处理：读取 `{id}.txt`。如果文件不存在则直接返回 200 OK 及空内容（创建新笔记）。
*   **`writeNote(...)`** _([handlers.go](file:///workspace/handlers.go#L87-L112))_
    *   **POST** 请求处理：对请求体大小进行上限拦截，检查只读模式。然后将请求内容覆盖写入到对应的 `{id}.txt` 文件中。
*   **`listNotes(...)`** _([handlers.go](file:///workspace/handlers.go#L114-L132))_
    *   **GET** `/list` 处理：遍历 `DATA_DIR` 目录，以 JSON 数组形式返回所有 `.txt` 文件的 ID 列表。

### 前端核心函数 (JavaScript in `embed.go`)

*   **`generateId()`**
    *   随机生成 12 位（包含大小写字母、数字和符号）的字符串，用作新建笔记的唯一 ID。
*   **`saveNote()`** / **`scheduleSave()`**
    *   `saveNote`：发起 Fetch POST 请求保存当前编辑器内容。
    *   `scheduleSave`：利用 `setTimeout` 实现 500 毫秒的防抖（Debounce）保存逻辑，确保用户连续输入时不会频繁请求接口。
*   **`loadNote(id)`**
    *   发起 Fetch GET 请求，获取指定 ID 的笔记内容，并更新到文本框内。
*   **`newNote()`**
    *   生成全新 ID，利用 HTML5 History API (`window.history.pushState`) 更新浏览器 URL（不触发刷新），并清空编辑器。

---

## 4. 依赖关系

本项目最大的特点是**极致轻量和零依赖**。

*   **后端依赖**：
    *   基于 **Go 1.23** 编写。
    *   完全依赖 Go 标准库 (`net/http`, `os`, `io`, `regexp`, `encoding/json` 等)。
    *   无任何第三方包依赖（`go.mod` 中无 `require` 块）。
*   **前端依赖**：
    *   纯原生 HTML5, CSS3, ES6 JavaScript。无 React/Vue 等框架。
    *   唯一的外部网络依赖：Google Fonts 提供的 `Inter` 字体包加载。
*   **构建与部署依赖**：
    *   依赖 **Docker** (多阶段构建 `scratch` 极小镜像)。
    *   依赖 **Docker Compose** 进行编排。

---

## 5. 项目运行方式

所有配置项均通过环境变量传入。支持以下几种运行方式：

### 5.1 源码本地运行（开发环境）
需要本地具备 Go 1.23+ 环境。
```bash
# 直接运行
go run .

# 或编译为二进制文件运行
go build -o app .
./app
```
运行后访问 `http://localhost:20008`。默认数据会存储在当前目录下的 `_tmp` 文件夹中。

### 5.2 Docker 运行
适用于快速部署测试，镜像极小 (基于 scratch)。
```bash
# 构建镜像
docker build -t minimalist-web-notepad:latest .

# 运行容器，将本地的 ./data 目录挂载为容器的 /tmp 存储目录
docker run -d \
  -p 20008:20008 \
  -v ./data:/tmp \
  minimalist-web-notepad:latest
```

### 5.3 Docker Compose 运行（推荐）
适合生产环境或需要固化参数的本地开发环境。

**方式一：基于本地源码构建并运行**
```bash
docker-compose up -d
```

**方式二：直接拉取远端 Docker Hub 镜像运行**
```bash
docker-compose -f docker-compose.hub.yml up -d
```

---
**核心环境变量配置参考**：
*   `LISTEN_ADDR`: 监听地址，默认 `:20008`
*   `DATA_DIR`: 数据存放路径，默认 `./_tmp` 或 `/tmp`
*   `MAX_NOTE_SIZE`: 单个笔记最大容量，默认 `1048576` (1MB)
*   `READ_ONLY`: 是否开启只读模式，设置 `"true"` 时禁止保存修改
