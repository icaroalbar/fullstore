#!/usr/bin/env bash
# PostToolUse: depois de editar um arquivo .go, roda gofmt e go vet.
# Em falha, sai com código 2 para o erro voltar ao Claude.
set -uo pipefail

# O Claude Code pode não herdar o PATH do .zshrc.
export PATH="$HOME/.local/go/bin:$HOME/go/bin:$PATH"

file="$(jq -r '.tool_input.file_path // empty')"
[[ "$file" == *.go ]] || exit 0

cd "${CLAUDE_PROJECT_DIR:-.}" || exit 1

status=0
unformatted="$(gofmt -l "$file" 2>&1)"
if [[ -n "$unformatted" ]]; then
  echo "gofmt: arquivo fora do formato ou com erro de sintaxe: $unformatted" >&2
  status=2
fi
if ! vet_output="$(go vet ./... 2>&1)"; then
  echo "go vet falhou:" >&2
  echo "$vet_output" >&2
  status=2
fi
exit "$status"
