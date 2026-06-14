# Storage Service

A lightweight personal file storage service written in Go using only the standard library. Upload files via HTTP, retrieve them by ID, and serve them back — all stored on the local filesystem with metadata tracked in JSON files alongside each upload.

## Tech Stack

| Layer       | Technology                                                  |
| ----------- | ----------------------------------------------------------- |
| Language    | [Go 1.25](https://go.dev)                                   |
| HTTP server | `net/http` — stdlib router with Go 1.22+ pattern matching   |
| Storage     | Local filesystem via `os`, `io`, `path/filepath`            |
| Metadata    | JSON files via `encoding/json`                              |
| MIME detect | `http.DetectContentType` (reads first 512 bytes)            |
| UUID        | `crypto/rand` — no external package                         |
| Logging     | `log/slog` — structured JSON logs                           |
| Config      | `.env` file parsed with `bufio` — no external package       |

> Zero external dependencies. Everything is Go standard library.

## Requirements

- Go 1.25+
- [Task](https://taskfile.dev) (optional, for `task` commands)

## Quick Start

```bash
cp .env.example .env   # edit PORT if needed
task run               # or: go run ./cmd/server/
```

The server starts on the port set in `.env` (default `9000`).

## Configuration

Copy `.env.example` to `.env` and set the values:

| Variable | Default | Description                     |
| -------- | ------- | ------------------------------- |
| `PORT`   | `8080`  | Port the HTTP server listens on |

## Available Tasks

| Command      | Description                          |
| ------------ | ------------------------------------ |
| `task run`   | Run the server                       |
| `task build` | Build binary to `bin/server`         |
| `task check` | Format, vet, and build               |
| `task fmt`   | Format all Go source files           |
| `task vet`   | Run `go vet` across all packages     |
| `task tidy`  | Tidy `go.mod`                        |
| `task clean` | Remove build artifacts               |

## Project Structure

```text
storage-service/
├── cmd/server/          # Entry point — wires routes and starts the server
├── internal/
│   ├── api/             # HTTP handlers
│   ├── models/          # Media struct
│   └── storage/         # Filesystem and metadata logic
├── pkg/                 # Shared utilities (UUID generation, .env loader)
├── helper/              # Shared response helpers (JSON/error writers)
├── storage/             # Uploaded files (gitignored, created at startup)
│   ├── images/
│   ├── videos/
│   ├── audio/
│   ├── documents/
│   └── others/
├── .env.example
├── API.md               # Full API reference
└── Taskfile.yml
```

## Storage Layout

Each uploaded file gets its own directory named by UUID:

```text
storage/
└── images/
    └── a1b2c3d4-e5f6-4789-a012-b3c4d5e6f789/
        ├── photo.jpg
        └── meta.json
```

If an upload fails midway, the partial directory is automatically cleaned up.

## API

See [API.md](API.md) for the full endpoint reference.

### Quick reference

| Method   | Path                    | Description            |
| -------- | ----------------------- | ---------------------- |
| `GET`    | `/health`               | Health check           |
| `POST`   | `/api/media/upload`     | Upload a file          |
| `GET`    | `/api/media`            | List all files         |
| `GET`    | `/api/media/{id}`       | Get file metadata      |
| `GET`    | `/api/media/{id}/file`  | Download the file      |
| `DELETE` | `/api/media/{id}`       | Delete a file          |
