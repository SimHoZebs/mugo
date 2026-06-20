#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SESSION_NAME="${MUGO_TMUX_SESSION:-mugo-dev}"
AVD_NAME="${MUGO_AVD:-Medium_Phone_API_36.1}"
BOOT_TIMEOUT_SECONDS="${MUGO_ANDROID_BOOT_TIMEOUT:-180}"

usage() {
  cat <<USAGE
Usage: scripts/dev.sh [all|server|mobile|android|stop]

Daily startup:
  make dev

Optional environment:
  MUGO_TMUX_SESSION      tmux session name (default: mugo-dev)
  MUGO_AVD               Android virtual device name (default: Medium_Phone_API_36.1)
  MUGO_SKIP_MIGRATIONS   set to 1 to skip make migrate-up during backend startup
USAGE
}

has_command() {
  command -v "$1" >/dev/null 2>&1
}

fail() {
  printf '\n%s\n' "$1" >&2
  printf 'Run `make doctor` for setup diagnostics.\n' >&2
  exit 1
}

ensure_command() {
  local command_name="$1"
  local help_text="$2"

  if ! has_command "$command_name"; then
    fail "Missing required command: ${command_name}. ${help_text}"
  fi
}

maybe_start_docker_desktop() {
  if ! grep -qi microsoft /proc/version 2>/dev/null; then
    return 0
  fi

  if ! has_command powershell.exe; then
    return 0
  fi

  local docker_desktop='C:\Program Files\Docker\Docker\Docker Desktop.exe'
  if [ -f "/mnt/c/Program Files/Docker/Docker/Docker Desktop.exe" ]; then
    printf 'Docker is not reachable. Trying to start Docker Desktop...\n'
    powershell.exe -NoProfile -Command "Start-Process '${docker_desktop}'" >/dev/null 2>&1 || true
  fi
}

wait_for_docker() {
  ensure_command docker "Enable Docker Desktop WSL integration for this distro or install Docker in WSL."

  if docker info >/dev/null 2>&1; then
    return 0
  fi

  maybe_start_docker_desktop

  printf 'Waiting for Docker to become reachable'
  for _ in $(seq 1 60); do
    if docker info >/dev/null 2>&1; then
      printf '\n'
      return 0
    fi
    printf '.'
    sleep 2
  done
  printf '\n'

  fail "Docker is still not reachable. If you use Docker Desktop, start it and enable WSL integration for this distro."
}

wait_for_container_health() {
  local container_name="$1"
  local label="$2"

  printf 'Waiting for %s' "$label"
  for _ in $(seq 1 60); do
    local status
    status="$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "$container_name" 2>/dev/null || true)"
    if [ "$status" = "healthy" ] || [ "$status" = "running" ]; then
      printf '\n'
      return 0
    fi
    printf '.'
    sleep 2
  done
  printf '\n'

  fail "${label} did not become healthy. Check the backend tmux pane or run docker compose logs."
}

resolve_emulator() {
  if has_command emulator; then
    command -v emulator
    return 0
  fi

  if [ -n "${ANDROID_SDK_ROOT:-}" ] && [ -x "${ANDROID_SDK_ROOT}/emulator/emulator" ]; then
    printf '%s\n' "${ANDROID_SDK_ROOT}/emulator/emulator"
    return 0
  fi

  if [ -n "${ANDROID_HOME:-}" ] && [ -x "${ANDROID_HOME}/emulator/emulator" ]; then
    printf '%s\n' "${ANDROID_HOME}/emulator/emulator"
    return 0
  fi

  return 1
}

android_device_connected() {
  has_command adb || return 1
  adb devices | grep -qE '[[:space:]]device$'
}

wait_for_android_boot() {
  ensure_command adb "Install Android platform-tools and add them to PATH."

  printf 'Waiting for Android emulator boot'
  local waited=0
  while [ "$waited" -lt "$BOOT_TIMEOUT_SECONDS" ]; do
    if adb shell getprop sys.boot_completed 2>/dev/null | tr -d '\r' | grep -q '^1$'; then
      printf '\nAndroid emulator is ready.\n'
      return 0
    fi
    printf '.'
    sleep 3
    waited=$((waited + 3))
  done
  printf '\n'

  return 1
}

start_android() {
  ensure_command adb "Install Android platform-tools and add them to PATH."

  if android_device_connected; then
    if wait_for_android_boot; then
      return 0
    fi

    fail "An Android device is connected, but it did not finish booting within ${BOOT_TIMEOUT_SECONDS}s."
  fi

  local emulator_bin
  if ! emulator_bin="$(resolve_emulator)"; then
    fail "Android emulator binary was not found. Set ANDROID_HOME or ANDROID_SDK_ROOT, or add the SDK emulator directory to PATH."
  fi

  if ! "$emulator_bin" -list-avds | grep -qx "$AVD_NAME"; then
    printf 'Configured AVD `%s` was not found. Available AVDs:\n' "$AVD_NAME" >&2
    "$emulator_bin" -list-avds >&2 || true
    fail "Set MUGO_AVD to an existing Android virtual device name."
  fi

  printf 'Starting Android emulator `%s`...\n' "$AVD_NAME"
  "$emulator_bin" -avd "$AVD_NAME" -netdelay none -netspeed full >/tmp/mugo-android-emulator.log 2>&1 &

  if ! wait_for_android_boot; then
    fail "Android emulator did not finish booting within ${BOOT_TIMEOUT_SECONDS}s. See /tmp/mugo-android-emulator.log."
  fi
}

start_server() {
  cd "$ROOT_DIR"
  ensure_command infisical "Install the Infisical CLI."
  wait_for_docker

  printf 'Starting Docker Compose services...\n'
  infisical run -- docker compose up -d
  wait_for_container_health db Postgres

  if [ "${MUGO_SKIP_MIGRATIONS:-0}" != "1" ]; then
    printf 'Running database migrations...\n'
    make migrate-up
  fi

  printf 'Starting Go API...\n'
  cd "$ROOT_DIR/server"
  if has_command air; then
    exec infisical run -- air
  fi

  printf 'air was not found; falling back to go run without hot reload.\n'
  exec infisical run -- go run ./cmd/api/main.go
}

start_mobile() {
  cd "$ROOT_DIR/mobile"
  ensure_command infisical "Install the Infisical CLI."
  ensure_command pnpm "Install pnpm."

  if has_command adb && wait_for_android_boot; then
    printf 'Starting Expo and opening Android...\n'
    exec infisical run -- pnpm expo start --android
  fi

  printf 'Android emulator is not ready; starting Expo only. Press `a` in this pane after Android boots.\n'
  exec infisical run -- pnpm expo start
}

start_tmux() {
  ensure_command tmux "Install tmux."
  ensure_command infisical "Install the Infisical CLI."
  wait_for_docker

  if tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
    printf 'Attaching to existing tmux session `%s`.\n' "$SESSION_NAME"
    exec tmux attach-session -t "$SESSION_NAME"
  fi

  tmux new-session -d -s "$SESSION_NAME" -n backend "cd '$ROOT_DIR' && scripts/dev.sh server; exec bash"
  tmux new-window -t "$SESSION_NAME" -n android "cd '$ROOT_DIR' && scripts/dev.sh android; exec bash"
  tmux new-window -t "$SESSION_NAME" -n mobile "cd '$ROOT_DIR' && scripts/dev.sh mobile; exec bash"
  tmux new-window -t "$SESSION_NAME" -n doctor "cd '$ROOT_DIR' && scripts/dev-doctor.sh; printf '\nPress Ctrl-b then w to switch panes/windows.\n'; exec bash"
  tmux select-window -t "$SESSION_NAME:backend"

  exec tmux attach-session -t "$SESSION_NAME"
}

stop_dev() {
  cd "$ROOT_DIR"
  if has_command tmux && tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
    tmux kill-session -t "$SESSION_NAME"
    printf 'Stopped tmux session `%s`.\n' "$SESSION_NAME"
  fi

  if has_command docker && docker info >/dev/null 2>&1; then
    docker compose stop
    printf 'Stopped Docker Compose services.\n'
  fi
}

case "${1:-all}" in
  all)
    start_tmux
    ;;
  server)
    start_server
    ;;
  mobile)
    start_mobile
    ;;
  android)
    start_android
    ;;
  stop)
    stop_dev
    ;;
  -h|--help|help)
    usage
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac
