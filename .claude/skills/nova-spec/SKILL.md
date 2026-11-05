---
name: nova-spec
description: Cria a spec de uma funcionalidade do FullStore a partir do template e das regras do PRD. Use quando o Icaro pedir para criar ou escrever uma spec, como "vamos criar a spec de pedido".
---

# Nova spec

1. Veja o maior número em `docs/specs/` e use o próximo (`NNN` com 3 dígitos) para criar `docs/specs/NNN-nome/`.
2. Copie `docs/specs/_templates/spec.md` para `docs/specs/NNN-nome/spec.md`, com `Status: rascunho`.
3. Leia `docs/prd.md` e preencha `Regras cobertas` e `Etapa` com os IDs que se aplicam (ex.: PR-01 a PR-05). Copie os critérios de aceite para os cenários Dado/Quando/Então e para a tabela de erros.
4. O primeiro rascunho do contexto e dos cenários é do Icaro: é assim que ele aprende. Peça o rascunho dele e proponha melhorias, em vez de escrever tudo sozinho.
5. Faça perguntas, poucas de cada vez, até `Dúvidas em aberto` ficar vazia. Registre cada resposta na spec.
6. Mantenha a spec sobre o quê e por quê. Tipos, pacotes e assinaturas ficam para o plano.
7. Uma regra precisa mudar? Proponha a mudança em `docs/prd.md` primeiro.
8. Pare e peça a aprovação do Icaro. Só ele muda o status para `aprovada`, ou você muda quando ele disser que aprovou.

Não crie `plan.md`, `tasks.md` nem código nesta skill.
