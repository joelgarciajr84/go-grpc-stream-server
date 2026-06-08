# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server
go run cmd/main.go

# Build
go build ./...

# Run tests
go test ./...

# Recompile protobuf definitions (requires Docker)
make compile-proto-go
```

## Architecture

This is a minimal gRPC server implementing **server-side streaming** in Go.

**Proto contract** (`proto/stream.proto`): defines `StreamService` with a single RPC `FetchResponse(Request) returns (stream Response)`. The generated Go code lives in `pkg/pb/` and should not be edited manually — regenerate with `make compile-proto-go`.

**Server** (`cmd/main.go`): listens on `:50005`. On each `FetchResponse` call it spawns 42 goroutines concurrently, each sleeping for `count` seconds then sending one `Response` back over the stream. The server uses `sync.WaitGroup` to block until all goroutines finish before closing the stream.

**Proto codegen**: uses a Dockerized `protoc` via `thethingsindustries/protoc`. The Makefile wipes and recreates `pkg/pb/` on each run, so any manual changes there are lost.
