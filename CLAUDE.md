# FullStore

Projeto de estudo (curso Full Cycle). O dono, Icaro, aprende Go e pensa em OO.
O que construir está em `docs/prd.md`; como trabalhar está em `docs/specs/README.md`.

## Como trabalhar
- Siga o fluxo SDD: spec -> plano -> tarefas -> código. Sem spec aprovada, não escreva código de produção.
- Por enquanto, o Claude implementa todo o código, inclusive o domínio. O Icaro aprova specs e revisa diffs.
- Cada decisão explicada comparando com OO: struct+métodos = classes, interfaces implícitas = duck typing, composição = herança.
- Código diverge da spec? Pare e pergunte.

## Arquitetura
- `internal/domain` não importa nada de `usecase` nem `infra` (o depguard barra).
- `internal/usecase` depende de interfaces (portas) declaradas por ele.
- `internal/infra` implementa as portas (HTTP, Postgres, Redis).

## Convenções
- Dinheiro em centavos (`int64`). Erros de domínio como variáveis `Err...` no pacote domain.
- Testes nomeados com o ID da regra: `TestOrder_PE02_...`
- Commits: `feat(order): PE-02 limita quantidade por item`.

## Comandos
- `make check` (rode antes de dizer que terminou)
- `make test`, `make lint`, `make run`

## Skills
- `nova-spec`: criar a spec de uma funcionalidade.
- `planejar-spec`: gerar plano e tarefas de uma spec aprovada.
- `implementar-tarefa`: escrever teste vermelho, implementar até verde, rodar `make check`, mostrar diff com explicação de cada decisão (comparando com OO).
