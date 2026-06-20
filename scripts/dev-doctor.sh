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

fail() {
  printf '[missing] %s\n' "$1"
}

has_command() {
  command -v "$1" >/dev/null 2>&1
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

printf 'Mugo development setup doctor\n'
printf 'Project: %s\n\n' "$ROOT_DIR"

if has_command tmux; then
  ok "tmux: $(tmux -V)"
else
  fail "tmux is required for the one-command dev dashboard"
fi

if has_command infisical; then
  ok "infisical: $(infisical --version 2>/dev/null | head -n 1)"
else
  fail "infisical CLI is required by project commands"
fi

if has_command pnpm; then
  ok "pnpm: $(pnpm --version 2>/dev/null)"
else
  fail "pnpm is required for the Expo app"
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
  fail "docker command is unavailable"
  if grep -qi microsoft /proc/version 2>/dev/null; then
    warn "WSL detected. Enable Docker Desktop WSL integration for this distro."
  fi
fi

if has_command adb; then
  ok "adb: $(command -v adb)"
  printf '\nConnected Android devices:\n'
  adb devices 2>/dev/null || true
else
  fail "adb is unavailable; install Android platform-tools and add them to PATH"
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
  fail "Android emulator binary is unavailable"
  warn "Set ANDROID_HOME or ANDROID_SDK_ROOT, or add the SDK emulator directory to PATH."
fi

printf '\nEnvironment hints:\n'
printf 'ANDROID_HOME=%s\n' "${ANDROID_HOME:-}"
printf 'ANDROID_SDK_ROOT=%s\n' "${ANDROID_SDK_ROOT:-}"
printf 'MUGO_AVD=%s\n' "$AVD_NAME"
