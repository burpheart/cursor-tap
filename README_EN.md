# Cursor-Tap

[中文](./README.md) | English

A tool for intercepting and analyzing Cursor IDE's gRPC traffic. Decrypts TLS, deserializes protobuf, and displays AI conversations in real-time.

## Why

Cursor talks to its backend entirely via gRPC (Connect Protocol). The body is binary protobuf. Burp and Fiddler just show unreadable bytes. Cursor doesn't publish proto definitions either.

This tool decrypts traffic into readable JSON and shows each streaming frame in real-time.

## How It Works

1. **MITM Proxy** - Sits between Cursor and api2.cursor.sh, decrypts TLS with self-signed CA
2. **Proto Extraction** - Extracts proto definitions from Cursor's JS bundle (`protobuf-es` compiled output)
3. **Real-time Parsing** - Parses Connect Protocol envelope framing, deserializes each protobuf frame
4. **WebUI** - Pushes frames via WebSocket, four-panel layout for service tree / calls / frames / details

## Quick Start

### 1. Start the Proxy

```bash
go run ./cmd/cursor-tap start --http-parse
```

`--http-parse` enables HTTP traffic parsing and WebSocket streaming, which the WebUI depends on.

Listens on `localhost:8080` (HTTP proxy), `localhost:1080` (SOCKS5 proxy), and `localhost:9090` (API + WebSocket).

### 2. Configure Cursor

Set environment variables to route Cursor through the proxy and trust the self-signed CA:

```bash
# Windows
set HTTP_PROXY=http://localhost:8080
set HTTPS_PROXY=http://localhost:8080
set http_proxy=http://localhost:8080
set https_proxy=http://localhost:8080
set NODE_EXTRA_CA_CERTS=C:\path\to\ca.crt

# macOS/Linux
export HTTP_PROXY=http://localhost:8080
export HTTPS_PROXY=http://localhost:8080
export http_proxy=http://localhost:8080
export https_proxy=http://localhost:8080
export NODE_EXTRA_CA_CERTS=~/.cursor-tap/ca/ca.crt
```

> ⚠️ **Both uppercase and lowercase versions must be set**: Node.js prioritizes lowercase `http_proxy`/`https_proxy`; setting only the uppercase versions may not work.

CA certificate is auto-generated at `~/.cursor-tap/ca/ca.crt` on first run.

#### macOS Extra Step: Install CA to System Keychain

Cursor is an Electron app — its network requests use the Chromium networking stack, which **ignores `NODE_EXTRA_CA_CERTS`** and only trusts the system certificate store. Install the CA to the system keychain:

```bash
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/.cursor-tap/ca/ca.crt
```

#### Launch Cursor from Terminal

Cursor must be launched from the terminal to inherit environment variables:

```bash
# macOS
/Applications/Cursor.app/Contents/MacOS/Cursor &

# Or run in background with nohup
nohup /Applications/Cursor.app/Contents/MacOS/Cursor > /dev/null 2>&1 &
```

> ⚠️ Using `open -a Cursor` will not inherit terminal environment variables.

### 3. Start WebUI

```bash
cd web
npm install
npm run dev
```

Open `http://localhost:3000`.

## Project Structure

```
├── cmd/cursor-tap/    # Proxy entry point
├── internal/
│   ├── api/           # API and WebSocket service
│   ├── ca/             # CA certificate management
│   ├── httpstream/     # gRPC parsing core
│   ├── mitm/          # MITM interceptor
│   └── proxy/          # HTTP/SOCKS5 proxy
├── cursor_proto/       # Extracted proto definitions
└── web/                # Next.js frontend
```

## What You Can See

- `AiService/RunSSE` - AI conversation channel (thinking, text, tool calls)
- `BidiService/BidiAppend` - User messages and tool results
- `AiService/StreamCpp` - Code completion
- `CppService/RecordCppFate` - Completion accept/reject feedback
- `AiService/Batch` - User behavior telemetry
- And dozens more...

## Disclaimer

For educational and research purposes only.

## Related

Detailed reverse engineering notes: [cursor-true-reverse-notes-1.md](./cursor-true-reverse-notes-1.md) (Chinese)
