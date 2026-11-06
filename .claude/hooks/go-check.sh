#!/usr/bin/env bash
# PostToolUse: depois de editar um arquivo .go, formata com gofmt e roda go vet.
# Em falha, sai com código 2 para o erro voltar ao Claude.
set -uo pipefail

# O Claude Code pode não herdar o PATH do .zshrc.
export PATH="$HOME/.local/go/bin:$HOME/go/bin:$PATH"

# Verifica se jq está disponível
if ! command -v jq &> /dev/null; then
  echo "ERROR: jq not found in PATH. Install jq to use this hook." >&2
  exit 2
fi

file="$(jq -r '.tool_input.file_path // empty')"
[[ "$file" == *.go ]] || exit 0

cd "${CLAUDE_PROJECT_DIR:-.}" || exit 1

status=0
# Formata o arquivo com gofmt -w (write in place)
if ! gofmt -w "$file" > /dev/null 2>&1; then
  echo "gofmt: erro ao formatar $file" >&2
  status=2
fi
if ! vet_output="$(go vet ./... 2>&1)"; then
  echo "go vet falhou:" >&2
  echo "$vet_output" >&2
  status=2
fi
exit "$status"
