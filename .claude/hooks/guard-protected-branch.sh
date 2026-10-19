#!/usr/bin/env bash
# PreToolUse (Bash): bloquea git commit/push directo a ramas protegidas.
# Sin dependencias (bash puro). Falla "abierto" (exit 0) si no puede determinar la rama.
PROTECTED="main"

input=$(cat)

case "$input" in
  *"git commit"*|*"git push"*)
    branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null) || exit 0
    for p in $PROTECTED; do
      if [ "$branch" = "$p" ]; then
        echo "BLOQUEADO: no se commitea/pushea directo a '$branch'. Crea una rama (feat/, fix/, docs/...) y abri un PR a main." >&2
        exit 2
      fi
    done
    ;;
esac
exit 0
