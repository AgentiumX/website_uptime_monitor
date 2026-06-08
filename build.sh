#!/bin/bash
set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT_DIR"

echo "=== Building Web ==="
cd web
npm install
npm run build
cd "$ROOT_DIR"

echo "=== Copying Web build artifacts ==="
rm -rf server/cmd/web/dist 2>/dev/null || true
mkdir -p server/cmd/web/dist
cp -r web/dist/* server/cmd/web/dist/

echo "=== Building Server ==="
cd server
go build -o uptime-server cmd/main.go
cd "$ROOT_DIR"

echo "=== Building Agent ==="
cd agent
go build -o uptime-agent .
cd "$ROOT_DIR"

echo ""
echo "=== Build Complete ==="
echo "Server: server/uptime-server"
echo "Agent:  agent/uptime-agent"
echo ""
echo "Run Server:  ./server/uptime-server -c server/config.yaml"
echo "Run Agent:   ./agent/uptime-agent -c agent/agent.yaml"
