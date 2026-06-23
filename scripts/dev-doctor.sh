#!/usr/bin/env bash
set -u

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
AVD_NAME="${MUGO_AVD:-Medium_Phone_API_36.1}"

ok() {
  printf '[ok] %s\n' "$1"
}

warn() {
  printf '[warn] %s\n' "$1"
}

missing() {
  printf '[missing] %s\n' "$1"
}

has_command() {
  command -v "$1" >/dev/null 2>&1
}

check_command() {
  local command_name="$1"
  local message="$2"

  if has_command "$command_name"; then
    ok "$command_name: $(command -v "$command_name")"
  else
    missing "$message"
  fi
}

check_port() {
  local port="$1"
  local label="$2"
  local occupied_status="${3:-warn}"

  if ! has_command ss; then
    warn "cannot inspect port ${port}; ss is unavailable"
    return 0
  fi

  local listener
  listener="$(ss -H -ltnp "sport = :${port}" 2>/dev/null || true)"
  if [ -z "$listener" ]; then
    ok "port ${port} (${label}) is free"
    return 0
  fi

  if [ "$occupied_status" = "ok" ]; then
    ok "port ${port} (${label}) is already in use"
  else
    warn "port ${port} (${label}) is already in use"
  fi
  printf '%s\n' "$listener"
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

check_android() {
  printf '\nAndroid tooling for optional `make emulator` flow:\n'

  if has_command adb; then
    ok "adb: $(command -v adb)"
    printf '\nConnected Android devices:\n'
    adb devices 2>/dev/null || true
  else
    missing "adb is unavailable; install Android platform-tools if you want CLI emulator control"
  fi

  if emulator_bin="$(resolve_emulator)"; then
    ok "emulator: ${emulator_bin}"
    printf '\nConfigured AVD: %s\n' "$AVD_NAME"
    if "$emulator_bin" -list-avds 2>/dev/null | grep -qx "$AVD_NAME"; then
      ok "configured AVD exists"
    else
      warn "configured AVD was not found; set MUGO_AVD to one of the available names"
      "$emulator_bin" -list-avds 2>/dev/null || true
    fi
  else
    missing "Android emulator binary is unavailable"
    warn "Set ANDROID_HOME or ANDROID_SDK_ROOT, or add the SDK emulator directory to PATH."
  fi

  printf '\nAndroid environment hints:\n'
  printf 'ANDROID_HOME=%s\n' "${ANDROID_HOME:-}"
  printf 'ANDROID_SDK_ROOT=%s\n' "${ANDROID_SDK_ROOT:-}"
  printf 'MUGO_AVD=%s\n' "$AVD_NAME"
}

case "${1:-}" in
  --android)
    printf 'Mugo Android setup doctor\n'
    printf 'Project: %s\n' "$ROOT_DIR"
    check_android
    exit 0
    ;;
  -h|--help|help)
    cat <<USAGE
Usage: scripts/dev-doctor.sh [--android]

Default checks whether the repo's `make dev` flow can start cleanly.
Use --android only for the optional CLI-managed emulator flow.
USAGE
    exit 0
    ;;
esac

printf 'Mugo development setup doctor\n'
printf 'Project: %s\n\n' "$ROOT_DIR"

check_command tmux "tmux is required because make dev splits the current tmux window"
if [ -n "${TMUX:-}" ]; then
  ok "running inside tmux"
else
  warn "not running inside tmux; make dev expects to be launched from your existing tmux session"
fi

if has_command infisical; then
  ok "infisical: $(infisical --version 2>/dev/null | head -n 1)"
else
  missing "infisical CLI is required by project commands"
fi

if has_command pnpm; then
  ok "pnpm: $(pnpm --version 2>/dev/null)"
else
  missing "pnpm is required for the Expo app"
fi

if [ -f "$ROOT_DIR/.env" ]; then
  ok ".env exists"
else
  warn ".env is missing; Makefile includes it and local commands may fail"
fi

if has_command docker; then
  if docker info >/dev/null 2>&1; then
    ok "docker daemon is reachable"
  else
    warn "docker command exists but the daemon is not reachable"
    if grep -qi microsoft /proc/version 2>/dev/null; then
      warn "WSL detected. Start Docker Desktop and enable WSL integration for this distro."
    fi
  fi
else
  missing "docker command is unavailable"
  if grep -qi microsoft /proc/version 2>/dev/null; then
    warn "WSL detected. Enable Docker Desktop WSL integration for this distro."
  fi
fi

printf '\nDev port status:\n'
check_port 8888 "Go API"
check_port 8081 "Expo/Metro"
check_port 5432 "Postgres" ok
check_port 9000 "Whisper" ok

printf '\nTip: use `make stop-dev` to stop Docker Compose services. Stale API/Expo processes may need to be killed separately.\n'
