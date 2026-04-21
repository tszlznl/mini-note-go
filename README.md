# Minimalist Web Notepad

A minimalist web-based notepad where the URL is the note. No login required, no database, pure file system storage.

## Features

- **URL is the note** - Share the URL to share your note
- **Zero dependencies** - Single static binary
- **Instant save** - Type and your changes are saved automatically
- **Pure file system** - Notes stored as plain text files
- **Lightweight** - Docker image ~6-12MB
- **Cross-platform** - Supports amd64 and arm64

## Quick Start

### Docker (Recommended)

```bash
docker run -d \
  -p 20008:20008 \
  -v ./data:/tmp \
  minimalist-web-notepad:latest
```

### Docker Compose

```bash
docker-compose up -d
```

### From Source

```bash
go build -o app .
./app
```

## Configuration

All configuration via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `LISTEN_ADDR` | `:20008` | Listen address |
| `DATA_DIR` | `/tmp` | Notes storage directory |
| `MAX_NOTE_SIZE` | `1048576` | Max note size in bytes (1MB) |
| `READ_ONLY` | `false` | Read-only mode |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Home page |
| GET | `/{id}` | Read note by ID |
| POST | `/{id}` | Save note by ID |
| GET | `/list` | List all notes |

## Development

```bash
go run .
```

Visit http://localhost:20008

## Building Docker Image

```bash
docker build -t minimalist-web-notepad:latest .
```

## License

MIT
