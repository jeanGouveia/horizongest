# Dívida Técnica

Este documento registra decisões arquiteturais conscientes que representam dívida técnica a ser resolvida em evoluções futuras do sistema.

---

## Snapshot Builder

### Contexto Atual

Atualmente o Repository cria os Snapshots durante CreateOrder.

Esta foi uma decisão consciente para manter o MVP simples, evitando a criação de serviços adicionais e mantendo a responsabilidade de montagem de snapshots dentro da camada de persistência.

### Fluxo Atual

```
OrderService
        ↓
OrderRepository (cria snapshot)
        ↓
Banco de Dados
```

### Evolução Planejada

Na evolução da arquitetura (após o MVP), a responsabilidade de montar Snapshots deverá migrar para um Domain Service (SnapshotBuilder ou StockService), deixando o Repository responsável apenas pela persistência.

### Fluxo Futuro Desejado

```
OrderService
        ↓
SnapshotBuilder (monta snapshot)
        ↓
OrderRepository (persiste apenas)
        ↓
Banco de Dados
```

### Justificativa

- **Separação de Responsabilidades:** O Repository deve ser responsável apenas pela persistência, não pela lógica de negócio de montagem de snapshots.
- **Reutilização:** Um SnapshotBuilder centralizado pode ser reutilizado por diferentes contextos (pedidos, ajustes de estoque, relatórios).
- **Testabilidade:** A lógica de snapshot pode ser testada independentemente da persistência.
- **Manutenibilidade:** Regras de snapshot complexas podem evoluir sem impactar o Repository.

### Quando Implementar

Esta mudança deve ser considerada quando:
- O sistema evoluir para além do MVP
- A lógica de snapshot se tornar mais complexa
- Houver necessidade de reutilizar a lógica de snapshot em múltiplos contextos

### NÃO Implementar Agora

Esta mudança NÃO deve ser implementada agora. Ela fica registrada apenas como evolução arquitetural planejada.

**Restrições atuais:**
- Não mover código
- Não criar SnapshotBuilder
- Não criar novos Services
- Não criar interfaces
- Não alterar Repository
- Não alterar injeção de dependências
- Não alterar CreateOrder
- Não alterar arquitetura do MVP
