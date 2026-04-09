#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUN_DIR="$ROOT_DIR/.run"
LOG_DIR="$RUN_DIR/logs"
COMPOSE_FILE="$ROOT_DIR/deploy/docker-compose.yaml"

mkdir -p "$RUN_DIR" "$LOG_DIR"

SERVICES=(
  "auth-service|$ROOT_DIR/auth-service|go run auth.go"
  "user-service|$ROOT_DIR/user-service|go run user.go"
  "friend-service|$ROOT_DIR/friend-service|go run friend.go"
  "chat-service|$ROOT_DIR/chat-service|go run message.go"
  "ws-gateway|$ROOT_DIR/ws-gateway|go run ws.go"
  "api-gateway|$ROOT_DIR/api-gateway|go run api.go"
  "frontend|$ROOT_DIR/frontend|npm run dev"
)

is_running() {
  local pid="$1"
  kill -0 "$pid" >/dev/null 2>&1
}

start_infra() {
  if ! command -v docker >/dev/null 2>&1; then
    echo "[ERROR] docker not found"
    exit 1
  fi

  echo "[INFO] Starting infrastructure containers..."
  docker compose -f "$COMPOSE_FILE" up -d mysql redis etcd zookeeper kafka
}

ensure_frontend_deps() {
  if [[ ! -d "$ROOT_DIR/frontend/node_modules" ]]; then
    echo "[INFO] Installing frontend dependencies..."
    (cd "$ROOT_DIR/frontend" && npm install)
  fi
}

kill_existing_frontend_vite() {
  local vite_pids pid
  vite_pids="$( { lsof -tiTCP:3000 -sTCP:LISTEN 2>/dev/null || true; lsof -tiTCP:3001 -sTCP:LISTEN 2>/dev/null || true; } | sort -u)"
  if [[ -n "$vite_pids" ]]; then
    echo "[INFO] Clearing existing frontend Vite processes on ports 3000/3001..."
    for pid in $vite_pids; do
      if ps -p "$pid" -o command= | grep -q '/frontend/node_modules/.bin/vite'; then
        kill "$pid" >/dev/null 2>&1 || true
      fi
    done
  fi
}

start_service() {
  local name="$1"
  local dir="$2"
  local cmd="$3"
  local pid_file="$RUN_DIR/${name}.pid"
  local log_file="$LOG_DIR/${name}.log"

  if [[ -f "$pid_file" ]]; then
    local old_pid
    old_pid="$(cat "$pid_file")"
    if [[ -n "$old_pid" ]] && is_running "$old_pid"; then
      echo "[SKIP] $name is already running (pid=$old_pid)"
      return
    fi
    rm -f "$pid_file"
  fi

  echo "[INFO] Starting $name ..."
  nohup bash -lc "cd '$dir' && $cmd" >"$log_file" 2>&1 &
  local pid=$!
  echo "$pid" >"$pid_file"
  echo "[OK]   $name started (pid=$pid, log=$log_file)"
}

print_summary() {
  echo
  echo "[DONE] All start commands have been issued."
  echo "[TIP]  Check status:  ./scripts/dev-status.sh"
  echo "[TIP]  Stop all:      ./scripts/dev-down.sh"
}

start_infra
ensure_frontend_deps
kill_existing_frontend_vite

for item in "${SERVICES[@]}"; do
  IFS="|" read -r name dir cmd <<<"$item"
  start_service "$name" "$dir" "$cmd"
done

print_summary
