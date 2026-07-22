# DOMÍNIO DO PRATOONLINE

> Este documento define os princípios permanentes do domínio do sistema.
>
> Nenhuma implementação deve contrariar estes princípios.
> Quando houver conflito entre implementação e este documento, este documento prevalece.

---

# 1. O domínio é a fonte de verdade

Todo desenvolvimento deve começar pelo domínio.

Nunca pela interface.

Nunca pelo banco.

Nunca pelo framework.

O domínio representa as regras do negócio.

Todo o restante existe apenas para atender o domínio.

---

# 2. Identidade é imutável

Toda entidade possui uma identidade permanente.

Exemplo:

Product
Ingredient
Order
Customer
Company

O ID nunca muda.

Todos os demais atributos podem evoluir.

Nunca reutilizar IDs.

Nunca alterar identidade.

---

# 3. Produto Vivo x Produto Vendido

Existe uma diferença fundamental.

## Produto Vivo

Representa o cadastro atual.

Pode sofrer alterações:

- nome
- preço
- descrição
- ficha técnica
- imagem
- categoria
- disponibilidade

Ele representa o presente.

---

## Produto Vendido

Após uma venda o produto deixa de depender do cadastro.

Ele passa a existir como parte do histórico.

O histórico nunca consulta novamente o cadastro do produto.

O histórico pertence ao pedido.

---

# 4. Histórico é imutável

Nenhum evento futuro pode alterar o passado.

Mudanças em:

- produtos
- ingredientes
- categorias
- preços
- fichas técnicas

Nunca alteram pedidos antigos.

O histórico representa exatamente o estado existente no momento da venda.

---

# 5. Snapshot

Toda venda gera um Snapshot.

O Snapshot pertence ao Pedido.

Nunca ao Produto.

O Snapshot representa:

- identificação comercial
- preço
- composição utilizada
- demais informações necessárias para reconstrução da venda

Após criado, nunca poderá ser alterado.

---

# 6. Estoque

O estoque representa o estado atual.

Nunca o histórico.

O histórico registra movimentações.

O estoque representa disponibilidade.

Baixas ocorrem somente durante a venda.

Reposições somente através dos fluxos previstos.

---

# 7. Active

active significa exclusivamente:

"Pode ser utilizado pelo negócio?"

Exemplos:

Produto disponível

Ingrediente disponível

Categoria disponível

Não representa exclusão.

---

# 8. Deleted At

deleted_at representa:

"O registro foi removido logicamente."

Um registro deletado:

- não aparece nas consultas normais
- continua existindo para auditoria
- continua preservando relacionamentos históricos

deleted_at nunca substitui active.

São responsabilidades distintas.

---

# 9. Separação entre Disponibilidade e Existência

Exemplos.

Produto fora do cardápio:

active = false

deleted_at = NULL

Produto removido:

active = false

deleted_at != NULL

Nunca existirão registros:

active = true

deleted_at != NULL

---

# 10. Banco acompanha o domínio

O banco nunca define regras.

Ele apenas materializa o domínio.

Toda alteração estrutural começa:

Domínio

↓

Models

↓

Repository

↓

Migration

↓

Banco

Nunca o inverso.

---

# 11. Evolução

Nenhuma nova funcionalidade poderá:

quebrar histórico

quebrar snapshots

quebrar APIs públicas

quebrar módulos

quebrar internacionalização futura

quebrar multiempresa futura

quebrar identidade das entidades

---

# 12. Generalização

A generalização será feita apenas quando prevista no Roadmap.

Até esse momento o sistema continuará especializado no domínio Restaurante.

Não serão antecipadas abstrações desnecessárias.

---

# 13. Princípio da Evolução Contínua

Toda evolução deve reduzir débito técnico.

Nunca aumentar.

Sempre que possível:

- simplificar
- padronizar
- documentar
- eliminar duplicidade

---

# 14. Responsabilidade Única

Cada conceito possui apenas uma responsabilidade.

active → disponibilidade

deleted_at → remoção lógica

snapshot → histórico

product → cadastro atual

order → transação comercial

stock → disponibilidade física

Nunca misturar responsabilidades.

---

# 15. Filosofia do Projeto

O sistema deve ser capaz de evoluir durante muitos anos sem necessidade de reescrita.

A arquitetura sempre terá prioridade sobre velocidade de implementação.

Implementações rápidas são aceitas apenas quando não violarem estes princípios.