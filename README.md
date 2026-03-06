# ThreadSync Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/threadsync/threadsync-go.svg)](https://pkg.go.dev/github.com/threadsync/threadsync-go)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

The official Go SDK for the [ThreadSync](https://threadsync.io) API. Sync data between any source and destination with a few lines of code.

## Installation

> **Preview**: SDKs are in preview and installed from GitHub. Registry packages (npm, PyPI, etc.) will be available at GA.

```bash
go get github.com/threadsync-infrastructure/threadsync-go
```

Requires Go 1.21 or later.

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    threadsync "github.com/threadsync/threadsync-go"
)

func main() {
    // Token is read from THREADSYNC_API_TOKEN env var if empty string is passed
    client := threadsync.New("your-api-token")

    // Create a connection
    conn, err := client.Connections.Create("salesforce", map[string]interface{}{
        "name": "Production Salesforce",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Connection created: %s (status: %s)\n", conn.ID, conn.Status)

    // Create a sync
    result, err := client.Sync.Create(&threadsync.SyncConfig{
        Source: threadsync.Endpoint{
            Connection: conn.ID,
            Object:     "Contact",
        },
        Destination: threadsync.Endpoint{
            Connection: "dest-conn-id",
            Table:      "contacts",
        },
        Schedule: "0 * * * *", // hourly
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Sync %s started — status: %s\n", result.ID, result.Status)
}
```

## Authentication

Set your API token via environment variable (recommended):

```bash
export THREADSYNC_API_TOKEN=ts_live_xxxxxxxxxxxx
```

Or pass it directly to `New()`:

```go
client := threadsync.New("ts_live_xxxxxxxxxxxx")
```

## API Reference

### Client

```go
client := threadsync.New(token string) *Client
```

Creates a new ThreadSync client. If `token` is an empty string, the value of the `THREADSYNC_API_TOKEN` environment variable is used.

---

### Connections

#### Create a connection

```go
conn, err := client.Connections.Create(provider string, options map[string]interface{}) (*Connection, error)
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `provider` | `string` | Provider name (e.g. `"salesforce"`, `"postgres"`, `"hubspot"`) |
| `options` | `map[string]interface{}` | Additional provider-specific options (e.g. `"name"`) |

Returns a `*Connection`:

```go
type Connection struct {
    ID       string // Unique connection ID
    Provider string // Provider name
    Name     string // Display name
    Status   string // "active", "pending", "error"
}
```

#### Get a connection

```go
conn, err := client.Connections.Get(id string) (*Connection, error)
```

---

### Sync

#### Create a sync

```go
result, err := client.Sync.Create(config *SyncConfig) (*SyncResult, error)
```

```go
type SyncConfig struct {
    Source      Endpoint // Source connection and object/table
    Destination Endpoint // Destination connection and object/table
    Schedule    string   // Cron expression (e.g. "0 * * * *")
}

type Endpoint struct {
    Connection string // Connection ID
    Object     string // Source object name (e.g. "Contact")
    Table      string // Destination table name
}
```

#### Get a sync

```go
result, err := client.Sync.Get(id string) (*SyncResult, error)
```

Returns a `*SyncResult`:

```go
type SyncResult struct {
    ID            string // Sync job ID
    Status        string // "pending", "running", "completed", "failed"
    RecordsSynced int    // Number of records synced (when completed)
}
```

---

## Constants

| Constant | Value |
|----------|-------|
| `DefaultBaseURL` | `https://api.threadsync.io/v1` |
| `Version` | `0.1.0` |

## License

MIT — see [LICENSE](LICENSE).
