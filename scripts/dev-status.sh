#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUN_DIR="$ROOT_DIR/.run"

SERVICES=(
  "auth-service"
  "user-service"
  "friend-service"
  "chat-service"
  "ws-gateway"
  "api-gateway"
  "frontend"
)

is_running() {
  local pid="$1"
  kill -0 "$pid" >/dev/null 2>&1
}

printf "%-15s %-10s %-8s %s\n" "SERVICE" "STATUS" "PID" "LOG"
printf "%-15s %-10s %-8s %s\n" "-------" "------" "---" "---"

for name in "${SERVICES[@]}"; do
  pid_file="$RUN_DIR/${name}.pid"
  log_file="$RUN_DIR/logs/${name}.log"

  if [[ -f "$pid_file" ]]; then
    pid="$(cat "$pid_file")"
    if [[ -n "$pid" ]] && is_running "$pid"; then
      printf "%-15s %-10s %-8s %s\n" "$name" "RUNNING" "$pid" "$log_file"
    else
      printf "%-15s %-10s %-8s %s\n" "$name" "EXITED" "${pid:-N/A}" "$log_file"
    fi
  else
    printf "%-15s %-10s %-8s %s\n" "$name" "STOPPED" "-" "$log_file"
  fi
done
