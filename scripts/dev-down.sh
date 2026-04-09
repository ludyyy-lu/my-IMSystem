#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUN_DIR="$ROOT_DIR/.run"
COMPOSE_FILE="$ROOT_DIR/deploy/docker-compose.yaml"
WITH_INFRA=false

if [[ "${1:-}" == "--with-infra" ]]; then
  WITH_INFRA=true
fi

SERVICES=(
  "frontend"
  "api-gateway"
  "ws-gateway"
  "chat-service"
  "friend-service"
  "user-service"
  "auth-service"
)

is_running() {
  local pid="$1"
  kill -0 "$pid" >/dev/null 2>&1
}

stop_service() {
  local name="$1"
  local pid_file="$RUN_DIR/${name}.pid"

  if [[ ! -f "$pid_file" ]]; then
    echo "[SKIP] $name is not tracked"
    return
  fi

  local pid
  pid="$(cat "$pid_file")"

  if [[ -z "$pid" ]]; then
    rm -f "$pid_file"
    echo "[SKIP] $name has empty pid file"
    return
  fi

  if is_running "$pid"; then
    echo "[INFO] Stopping $name (pid=$pid) ..."
    kill "$pid" >/dev/null 2>&1 || true
  else
    echo "[SKIP] $name process already exited (pid=$pid)"
  fi

  rm -f "$pid_file"
  echo "[OK]   $name stopped"
}

for name in "${SERVICES[@]}"; do
  stop_service "$name"
done

if $WITH_INFRA; then
  echo "[INFO] Stopping infrastructure containers..."
  docker compose -f "$COMPOSE_FILE" down
  echo "[OK]   infrastructure stopped"
else
  echo "[TIP]  Infra is still running. Use './scripts/dev-down.sh --with-infra' to stop it."
fi
