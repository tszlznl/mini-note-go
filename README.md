# Minimalist Web Notepad（极简网页记事本）

一个极简的基于网页的记事本，URL 就是笔记本身。无需登录，无需数据库，纯文件系统存储。

## 功能特性

- **URL 即是笔记** - 分享 URL 即可分享笔记
- **零依赖** - 单个静态二进制文件
- **即时保存** - 输入即保存
- **纯文件系统** - 笔记以纯文本文件存储
- **轻量级** - Docker 镜像约 6-12MB
- **跨平台** - 支持 amd64 和 arm64

## 快速开始

### Docker（推荐）

```bash
docker run -d \
  -p 20008:20008 \
  -v ./data:/tmp \
  minimalist-web-notepad:latest
```

### Docker Compose

**方式一：本地构建（开发调试）**

```bash
docker-compose up -d
```

**方式二：从 Docker Hub 拉取（生产环境）**

```bash
docker-compose -f docker-compose.hub.yml up -d
```

### 从源码构建

```bash
go build -o app .
./app
```

## 配置

所有配置通过环境变量进行：

| 变量名 | 默认值 | 描述 |
|--------|--------|------|
| `LISTEN_ADDR` | `:20008` | 监听地址 |
| `DATA_DIR` | `/tmp` | 笔记存储目录 |
| `MAX_NOTE_SIZE` | `1048576` | 最大笔记大小（字节），默认 1MB |
| `READ_ONLY` | `false` | 只读模式 |

## API 接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/` | 首页 |
| GET | `/note/{id}` | 根据 ID 读取笔记 |
| POST | `/note/{id}` | 根据 ID 保存笔记 |
| GET | `/list` | 列出所有笔记 |

## 开发

```bash
go run .
```

访问 http://localhost:20008

## 构建 Docker 镜像

```bash
docker build -t minimalist-web-notepad:latest .
```

## 许可证

MIT