#!/usr/bin/env bash
# PostToolUse (Edit|Write): formatea el archivo editado segun su extension. No-op si no matchea
# o si el formateador no esta instalado. Sin dependencias (bash puro).
input=$(cat)
path=$(printf '%s' "$input" | sed -n 's/.*"file_path"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)
path=${path//\\\\//}   # rutas Windows escapadas en JSON (C:\\Users\\...) -> C:/Users/...
[ -n "$path" ] && [ -f "$path" ] || exit 0

case "$path" in
  *.go)
    command -v gofmt >/dev/null 2>&1 && gofmt -w "$path" >/dev/null 2>&1
    ;;
  */frontend/*.ts|*/frontend/*.tsx|*/frontend/*.js|*/frontend/*.jsx|*/frontend/*.css)
    root=${path%%/frontend/*}/frontend
    (cd "$root" && npx --no-install prettier --write "$path") >/dev/null 2>&1
    ;;
esac
exit 0
