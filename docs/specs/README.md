# Specs (SDD: spec-driven development)

Nenhum código de produção é escrito sem uma spec aprovada pelo Icaro. Cada funcionalidade tem uma pasta numerada, e cada arquivo tem um portão antes do próximo.

```
docs/specs/NNN-nome/
├── spec.md    # o quê e por quê
├── plan.md    # como
└── tasks.md   # passos
```

| Passo | Arquivo | Conteúdo | Portão |
| --- | --- | --- | --- |
| 1. Spec | `spec.md` | O quê e por quê, sem falar de código. Regras do PRD cobertas e critérios em Dado/Quando/Então | O Icaro aprova (`Status: aprovada`) |
| 2. Plano | `plan.md` | Pacotes, tipos, interfaces, decisões e alternativas descartadas | Respeita as camadas do `CLAUDE.md` |
| 3. Tarefas | `tasks.md` | Passos de até 1 hora, cada um com o teste que prova que terminou | Toda regra da spec tem pelo menos uma tarefa |
| 4. Implementação | código | Uma tarefa por vez em TDD (vermelho, código, verde, refatora), um commit por tarefa | `make check` verde e tarefa marcada |

Quando algo muda no meio do caminho, atualiza-se a spec primeiro e o código depois. A spec é a fonte da verdade: se o código e a spec divergem, pare e pergunte.

As regras vêm de [`docs/prd.md`](../prd.md). Os modelos estão em [`_templates/`](_templates/), e as skills `nova-spec`, `planejar-spec` e `implementar-tarefa` conduzem cada passo.
