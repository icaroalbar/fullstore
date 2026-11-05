# PRD do FullStore

Cópia versionada do PRD aprovado em 05/11/2026 no documento "FullStore: plano de estudos, PRD, SDD e harness".
Mudou uma regra? Atualize este arquivo primeiro, depois a spec e o código.

**Problema.** Uma loja online pequena precisa vender produtos, reservar estoque e cobrar pagamentos sem vender o que não tem nem cobrar duas vezes. O FullStore resolve isso e cresce de um monolito até microsserviços em Kubernetes, para servir de laboratório do curso.

**Usuários.** Cliente (navega no catálogo, faz pedido, paga); Administrador (cadastra produtos e estoque); Gateway de pagamento externo (simulado).

**Fora do escopo.** Frontend, frete, cupons, notas fiscais, login social. Autenticação entra só na etapa 5, no gateway.

## Requisitos funcionais

| ID | Requisito | Etapa |
| --- | --- | --- |
| RF-01 | Cadastrar, editar, ativar e desativar produto | 1 |
| RF-02 | Listar produtos ativos com preço e disponibilidade | 1 |
| RF-03 | Repor e consultar estoque de um produto | 1 |
| RF-04 | Criar pedido com itens e calcular o total | 1 |
| RF-05 | Pagar pedido e confirmar ou recusar | 1 |
| RF-06 | Cancelar pedido e devolver estoque reservado | 1 |
| RF-07 | Consultar pedido com itens, status e pagamento | 1 (composto na 5) |

## Requisitos não funcionais

| ID | Requisito | Etapa |
| --- | --- | --- |
| RNF-01 | Cobertura de testes do domínio acima de 80% | 1 |
| RNF-02 | Listagem do catálogo em menos de 50 ms (p95) com 200 usuários no k6 | 2 |
| RNF-03 | Pedidos e Pagamentos funcionam se o outro serviço cair (eventos) | 3 |
| RNF-04 | Nenhum pagamento duplicado, mesmo com reenvio da mesma requisição | 4 |
| RNF-05 | Toda requisição rastreável ponta a ponta (trace id) | 5 |
| RNF-06 | Deploy automático a cada merge na main | 6 |

## Regras de negócio

Cada regra tem um ID que vai para a spec, para o nome do teste e para a mensagem de commit (ex.: `TestOrder_PE02_RejectsEmptyOrder`). Uma regra só está pronta quando seu critério de aceite vira teste passando.

### Produto (catálogo)

| ID | Regra | Critério de aceite |
| --- | --- | --- |
| PR-01 | Nome é obrigatório, de 3 a 100 caracteres | Nome "TV" retorna erro ErrInvalidName |
| PR-02 | Preço é maior que zero, em centavos (int64), moeda BRL | Preço 0 ou negativo retorna ErrInvalidPrice |
| PR-03 | Produto nasce desativado e só ativa com preço válido | Novo produto tem status disabled |
| PR-04 | Produto desativado não aparece no catálogo nem entra em pedido | Pedido com produto desativado retorna ErrProductUnavailable |
| PR-05 | Mudar o preço não altera pedidos já criados | Pedido guarda o preço unitário do momento da compra |

### Pedido

| ID | Regra | Critério de aceite |
| --- | --- | --- |
| PE-01 | Pedido tem de 1 a 20 itens distintos | Pedido vazio ou com 21 itens retorna erro |
| PE-02 | Quantidade por item é de 1 a 10 | Quantidade 0 ou 11 retorna ErrInvalidQuantity |
| PE-03 | Total = soma de preço unitário × quantidade, calculado pelo domínio | Total nunca vem do cliente; teste com 3 itens confere a soma |
| PE-04 | Status segue: pending → paid → shipped, ou pending → canceled | Pagar pedido cancelado retorna ErrInvalidTransition |
| PE-05 | Só pedido pending pode ser cancelado pelo cliente | Cancelar pedido paid retorna erro |
| PE-06 | Pedido pending sem pagamento expira em 30 minutos e é cancelado | Teste com relógio falso avançado 31 min |

### Estoque

| ID | Regra | Critério de aceite |
| --- | --- | --- |
| ES-01 | Criar pedido reserva o estoque de cada item | Disponível cai e reservado sobe na mesma quantidade |
| ES-02 | Sem estoque suficiente, o pedido inteiro é recusado (sem reserva parcial) | Pedido com 1 item faltando não reserva nenhum |
| ES-03 | Cancelar ou expirar pedido devolve a reserva | Disponível volta ao valor anterior |
| ES-04 | Pagamento aprovado baixa a reserva definitivamente | Reservado cai; disponível não muda |
| ES-05 | Estoque nunca fica negativo, mesmo com pedidos simultâneos | Teste com 50 goroutines disputando 10 unidades vende exatamente 10 |

### Pagamento

| ID | Regra | Critério de aceite |
| --- | --- | --- |
| PG-01 | Valor do pagamento = total do pedido | Valor diferente retorna ErrAmountMismatch |
| PG-02 | Um pedido tem no máximo um pagamento aprovado | Segundo pagamento aprovado é recusado |
| PG-03 | Requisição com a mesma chave de idempotência devolve o mesmo resultado | Duas chamadas iguais geram um só registro |
| PG-04 | Pagamento recusado mantém o pedido pending até expirar | Após recusa, status continua pending |
| PG-05 | Gateway sem resposta em 3 s conta como falha e pode ser tentado de novo | Teste com gateway falso lento |

## Decisões

- Moeda: só BRL, em centavos (int64). Multimoeda fora do escopo.
- Estorno (administrador cancelar pedido pago): só na etapa 4.
- PE-06, ES-05 e PG-05 mantidas como estão.
