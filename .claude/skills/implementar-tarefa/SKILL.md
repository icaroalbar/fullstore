---
name: implementar-tarefa
description: Conduz a próxima tarefa aberta de uma spec do FullStore em TDD. Escreve o teste vermelho, espera o Icaro implementar, roda make check e marca a tarefa. Use quando o Icaro disser "próxima tarefa".
---

# Implementar tarefa

1. Encontre a primeira tarefa `- [ ]` em `docs/specs/*/tasks.md` cuja spec esteja `aprovada`. Sem tarefa aberta ou sem spec aprovada, avise e pare.
2. Leia a tarefa e a regra correspondente na spec. Código e spec divergem? Pare e pergunte.
3. Escreva o teste, nomeado com o ID da regra (`TestOrder_PE02_...`), e rode só esse teste para mostrar que está vermelho. Em Go, não compilar porque o tipo ainda não existe conta como vermelho.
4. **Pare e espere o Icaro implementar.** Explique o que o teste exige e o idioma de Go envolvido, comparando com OO. Não escreva o código de produção, nem como exemplo pronto. Dê dicas se ele pedir.

   Enquanto o teste estiver vermelho, os hooks vão apontar `go vet` ou `make check` falhando. É o esperado: diga isso em uma frase e continue esperando, sem corrigir.
5. Quando ele avisar que terminou, rode `make check`. Se falhar, explique o erro e deixe a correção com ele.
6. Com tudo verde, revise o código dele: uma sugestão de cada vez, explicando o motivo.
7. Marque a tarefa como `- [x]` e proponha o commit no formato `feat(order): PE-02 limita quantidade por item`. Faça um commit por tarefa, e só quando ele pedir.

Só escreva a implementação de uma tarefa se o Icaro pedir explicitamente e ela não for código de domínio (por exemplo, infraestrutura repetitiva).
