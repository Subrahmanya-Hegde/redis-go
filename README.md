# Redis from Scratch in Go

A Redis server implementation built from scratch in Go, following the [CodeCrafters](https://codecrafters.io) "Build Your Own Redis" challenge.

## What this is

A learning project to understand:
- TCP server fundamentals
- The RESP (REdis Serialization Protocol) wire protocol
- Concurrent connection handling with goroutines
- Clean package design in Go

## Project Structure
```
app/
main.go — TCP listener, connection handling, command dispatch
resp/
    resp.go — RESP types (Value struct), constants
    reader.go — Deserialize: bytes → Value (recursive parser)
    writer.go — Serialize: Value → bytes (buffered writer)
```

## Run
```bash
go run app/main.go
```

## Testing
```bash
# PING
printf '*1\r\n$4\r\nPING\r\n' | nc localhost 6379
# +PONG

# ECHO
printf '*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n' | nc localhost 6379
# $5
# hello

# SET
- printf '*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$6\r\nSubrah\r\n' | nc localhost 6379
- printf '*5\r\n$3\r\nSET\r\n$4\r\nname\r\n$6\r\nSubrah\r\n$2\r\nEX\r\n$1\r\n5\r\n' | nc localhost 6379

# GET
printf '*2\r\n$3\r\nGET\r\n$4\r\nname\r\n' | nc localhost 6379
```
