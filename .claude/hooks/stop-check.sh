#!/usr/bin/env bash
# Stop: roda make check antes do Claude encerrar a resposta.
# Em falha, devolve o erro ao Claude uma única vez (stop_hook_active): um teste
# vermelho esperado, como o da skill implementar-tarefa, não pode virar um laço.
set -uo pipefail

export PATH="$HOME/.local/go/bin:$HOME/go/bin:$PATH"

[[ "$(jq -r '.stop_hook_active // false')" == "true" ]] && exit 0

cd "${CLAUDE_PROJECT_DIR:-.}" || exit 1

if ! command -v make >/dev/null; then
  echo "make não está instalado: make check não rodou (sudo apt install make)." >&2
  exit 1
fi

if ! output="$(make check 2>&1)"; then
  echo "make check falhou:" >&2
  echo "$output" >&2
  exit 2
fi
