---
name: planejar-spec
description: Gera plan.md e tasks.md de uma spec aprovada do FullStore, checando que toda regra tem tarefa e que as camadas são respeitadas. Use quando uma spec estiver com status aprovada e o Icaro pedir para planejar.
---

# Planejar spec

1. Leia `docs/specs/NNN-nome/spec.md`. Se o status não for `aprovada`, pare e avise o Icaro.
2. Copie os templates de `docs/specs/_templates/` para `plan.md` e `tasks.md` na pasta da spec.
3. No `plan.md`, defina onde mora o código, os tipos e as interfaces (só assinaturas em Go, sem implementação), as decisões com a alternativa descartada e os riscos. Compare com OO quando a decisão for um idioma de Go.
4. No `tasks.md`, escreva passos de até 1 hora. Cada tarefa cita o ID da regra e o teste que prova que terminou (ex.: `TestOrder_PE02_...`).
5. Antes de entregar, confira e mostre o resultado:
   - **Cobertura:** cada ID em `Regras cobertas` tem pelo menos uma tarefa.
   - **Camadas:** `internal/domain` não importa `usecase` nem `infra`; `usecase` depende de portas que ele mesmo declara; `infra` implementa as portas.
6. Encontrou algo que a spec não cobre ou que a contradiz? Pare e pergunte. A spec muda primeiro.
7. Peça ao Icaro para revisar o plano.

Não escreva código nem testes nesta skill.
