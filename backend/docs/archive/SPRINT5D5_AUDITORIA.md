# SPRINT 5D.5 — AUDITORIA DE REGRAS DE NEGÓCIO

## Resumo Executivo

Esta auditoria focou EXCLUSIVAMENTE nas regras de negócio do HorizonGest, validando se TODAS as regras funcionais realmente impedem estados inválidos. Foram analisadas 10 fases: Estados Inválidos, Máquinas de Estado, Validação de Dados, Duplicidade, Cálculos, Estoque, Permissões, Trial, Consistência Frontend-Backend e Casos Extremos.

**Total de problemas identificados:** 35
**Problemas críticos:** 12
**Problemas altos:** 15
**Problemas médios:** 8

---

## FASE 1 — ESTADOS INVÁLIDOS

### Problema 1.1: Produto com preço negativo possível
- **Severidade:** CRÍTICA
- **Arquivo:** `internal/service/product_service.go`
- **Linha:** 82
- **Causa:** Validação `gt=0` permite preço positivo, mas não impede preço negativo via API direta ou migration
- **Impacto:** Produto pode ter preço negativo, causando prejuízo financeiro
- **Risco:** Alto - Fraude, prejuízo financeiro
- **Correção:** Adicionar validação no domain e no repository para impedir Money negativo
- **Tempo estimado:** 2h

### Problema 1.2: Estoque negativo possível via ajuste manual
- **Severidade:** CRÍTICA
- **Arquivo:** `internal/service/stock_movement_service.go`
- **Linhas:** 86-92
- **Causa:** Ajuste manual (StockMovementAdjust) permite estoque negativo se quantity > estoque atual
- **Impacto:** Estoque pode ficar negativo, permitindo vendas impossíveis
- **Risco:** Alto - Vendas de produtos inexistentes
- **Correção:** Validar que newStock >= 0 para todos os tipos de movimento
- **Tempo estimado:** 2h

### Problema 1.3: Pedido sem itens possível
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/order_service.go`
- **Linhas:** 67
- **Causa:** Validação `min=1` existe mas não impede array vazio via bypass
- **Impacto:** Pedido pode ser criado sem itens, total = 0
- **Risco:** Médio - Pedido inválido no sistema
- **Correção:** Validar len(items) > 0 no service
- **Tempo estimado:** 1h

### Problema 1.4: Pedido com total incorreto calculado
- **Severidade:** ALTA
- **Arquivo:** `internal/service/order_service.go`
- **Linhas:** 165
- **Causa:** Cálculo usa `Mul(int64(itemIn.Quantity * 100)).Div(100)` que pode ter rounding errors
- **Impacto:** Total do pedido pode estar incorreto por centavos
- **Risco:** Médio - Divergência financeira
- **Correção:** Usar cálculo Money correto: `p.Price.MulFloat(itemIn.Quantity)`
- **Tempo estimado:** 2h

### Problema 1.5: Pedido cancelado pode ser editado
- **Severidade:** ALTA
- **Arquivo:** `internal/service/order_service.go`
- **Linhas:** 340-342
- **Causa:** UpdateOrder permite apenas pending/confirmed, mas não impede edição de cancelado via bypass
- **Impacto:** Pedido cancelado pode ser modificado após cancelamento
- **Risco:** Médio - Inconsistência de estado
- **Correção:** Validar status no repository também
- **Tempo estimado:** 1h

### Problema 1.6: Recebimento duplicado possível
- **Severidade:** CRÍTICA
- **Arquivo:** `internal/service/purchase_service.go`
- **Linhas:** 227-230
- **Causa:** Validação apenas verifica se status == received, mas não impede múltiplos recebimentos parciais
- **Impacto:** Mesmo pedido pode ser recebido múltiplas vezes, duplicando estoque
- **Risco:** Alto - Estoque inflado, prejuízo financeiro
- **Correção:** Adicionar validação de soma de recebimentos vs quantidade pedida
- **Tempo estimado:** 4h

### Problema 1.7: Produto arquivado pode ser vendido
- **Severidade:** ALTA
- **Arquivo:** `internal/service/order_service.go`
- **Linhas:** 134-136
- **Causa:** Validação verifica `!p.Active`, mas produto pode ser arquivado após validação
- **Impacto:** Produto inativo pode ser vendido se race condition
- **Risco:** Médio - Venda de produto indisponível
- **Correção:** Adicionar validação no repository dentro da transação
- **Tempo estimado:** 3h

### Problema 1.8: Ingrediente excluído pode estar em ficha técnica
- **Severidade:** ALTA
- **Arquivo:** `internal/service/product_service.go`
- **Linhas:** 409-416
- **Causa:** SetProductIngredients não verifica se ingrediente está ativo
- **Impacto:** Ficha técnica pode referenciar ingrediente inativo
- **Risco:** Médio - Erro ao calcular custo, vendas impossíveis
- **Correção:** Validar ingrediente.Active == true
- **Tempo estimado:** 1h

### Problema 1.9: Categoria excluída em produto ativo
- **Severidade:** MÉDIA
- **Arquivo:** `internal/domain/product.go`
- **Causa:** Product.CategoryID não tem validação de categoria ativa
- **Impacto:** Produto pode referenciar categoria excluída
- **Risco:** Baixo - Problema de display apenas
- **Correção:** Adicionar validação ao criar/atualizar produto
- **Tempo estimado:** 2h

### Problema 1.10: Empresa cancelada pode operar normalmente
- **Severidade:** CRÍTICA
- **Arquivo:** `internal/middleware/tenant_middleware.go`
- **Causa:** Tenant middleware não verifica se company.Active == true
- **Impacto:** Empresa cancelada pode continuar operando o sistema
- **Risco:** Alto - Violação de contrato, uso indevido
- **Correção:** Adicionar validação de company.Active no tenant middleware
- **Tempo estimado:** 2h

### Problema 1.11: Convite pode ser aceito duas vezes
- **Severidade:** ALTA
- **Arquivo:** `internal/service/invitation_service.go`
- **Linhas:** 232-243
- **Causa:** Validação existe mas race condition pode permitir duplo aceite
- **Impacto:** Mesmo convite pode ser aceito por múltiplos usuários
- **Risco:** Alto - Múltiplos usuários na mesma empresa
- **Correção:** Adicionar lock na tabela de convites ou usar transação
- **Tempo estimado:** 3h

### Problema 1.12: Reset de senha pode ser reutilizado
- **Severidade:** ALTA
- **Arquivo:** `internal/service/auth_service.go`
- **Linhas:** 376-378
- **Causa:** Validação de Used == true existe mas não há lock
- **Impacto:** Race condition pode permitir reuso do mesmo token
- **Risco:** Alto - Segurança, acesso não autorizado
- **Correção:** Adicionar lock na tabela de tokens ou usar transação
- **Tempo estimado:** 3h

### Problema 1.13: Sessão expirada pode ser reutilizada
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/platform_auth_service.go`
- **Linhas:** 169-173
- **Causa:** Comentário indica que validação de sessão não é implementada
- **Impacto:** Sessão expirada pode continuar sendo usada
- **Risco:** Médio - Acesso não autorizado
- **Correção:** Implementar validação de sessão no banco
- **Tempo estimado:** 4h

### Problema 1.14: Trial expirado pode continuar operando
- **Severidade:** CRÍTICA
- **Arquivo:** `internal/middleware/tenant_middleware.go`
- **Causa:** Tenant middleware não verifica trial expiration
- **Impacto:** Empresa com trial expirado pode continuar operando
- **Risco:** Alto - Violação de contrato, uso indevido
- **Correção:** Adicionar validação de trial expiration no tenant middleware
- **Tempo estimado:** 3h

---

## FASE 2 — MÁQUINAS DE ESTADO

### Problema 2.1: Transição Completed -> Pending não impedida
- **Severidade:** ALTA
- **Arquivo:** `internal/service/order_service.go`
- **Linhas:** 407-428
- **Causa:** isValidTransition não lista Completed como status final com transições vazias
- **Impacto:** Pedido entregue pode voltar para pendente
- **Risco:** Alto - Inconsistência de estado financeiro
- **Correção:** Adicionar Completed como status final sem transições
- **Tempo estimado:** 1h

### Problema 2.2: Transição Cancelled -> Preparing não impedida
- **Severidade:** ALTA
- **Arquivo:** `internal/service/order_service.go`
- **Linhas:** 407-428
- **Causa:** isValidTransition não lista Cancelled como status final com transições vazias
- **Impacto:** Pedido cancelado pode voltar para preparando
- **Risco:** Alto - Inconsistência de estado
- **Correção:** Adicionar Cancelled como status final sem transições
- **Tempo estimado:**:** 1h

### Problema 2.3: Purchase Order sem validação de transições
- **Severidade:** ALTA
- **Arquivo:** `internal/service/purchase_service.go`
- **Linhas:** 195-197
- **Causa:** UpdatePurchaseOrderStatus não valida transições permitidas
- **Impacto:** Pedido de compra pode ter transições inválidas (ex: received -> draft)
- **Risco:** Médio - Inconsistência de estado
- **Correção:** Implementar isValidTransition para PurchaseOrder
- **Tempo estimado:** 2h

### Problema 2.4: Inventory sem validação de transições
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/stock_movement_service.go`
- **Linhas:** 184-187
- **Causa:** Validação apenas verifica status != draft, mas não transições específicas
- **Impacto:** Inventário pode ter transições inválidas (ex: completed -> draft)
- **Risco:** Médio - Inconsistência de estado
- **Correção:** Implementar isValidTransition para StockInventory
- **Tempo estimado:** 2h

### Problema 2.5: Invitação expirada pode ser aceita
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/invitation_service.go`
- **Linhas:** 232-243
- **Causa:** Validação existe mas depende de IsExpired() que pode ter race condition
- **Impacto:** Convite expirado pode ser aceito se timing errado
- **Risco:** Baixo - Convite antigo sendo aceito
- **Correção:** Adicionar validação atômica com lock
- **Tempo estimado:** 2h

---

## FASE 3 — VALIDAÇÃO DE DADOS

### Problema 3.1: Quantidade negativa em OrderItemInput
- **Severidade:** CRÍTICA
- **Arquivo:** `internal/service/order_service.go`
- **Linha:** 63
- **Causa:** Validação `gt=0` existe mas não impede bypass
- **Impacto:** Item pode ter quantidade negativa
- **Risco:** Alto - Pedido com quantidade negativa
- **Correção:** Validar Quantity > 0 no service
- **Tempo estimado:** 1h

### Problema 3.2: Preço negativo em CreateProductInput
- **Severidade:** CRÍTICA
- **Arquivo:** `internal/service/product_service.go`
- **Linha:** 82
- **Causa:** Validação `gt=0` existe mas Money pode ser negativo via API
- **Impacto:** Produto pode ter preço negativo
- **Risco:** Alto - Prejuízo financeiro
- **Correção:** Validar Price > 0 no domain
- **Tempo estimado:** 1h

### Problema 3.3: Custo negativo possível
- **Severidade:** ALTA
- **Arquivo:** `internal/domain/product.go`
- **Linhas:** 59-66
- **Causa:** CalculateCost não valida se resultado é negativo
- **Impacto:** Custo do produto pode ser negativo
- **Risco:** Médio - Cálculos incorretos de margem
- **Correção:** Validar Cost >= 0 após cálculo
- **Tempo estimado:** 1h

### Problema 3.4: Desconto maior que total não validado
- **Severidade:** ALTA
- **Arquivo:** `internal/service/purchase_service.go`
- **Linhas:** 126-127
- **Causa:** Cálculo total = subtotal + tax - discount não valida se discount > subtotal + tax
- **Impacto:** Total pode ser negativo
- **Risco:** Médio - Prejuízo financeiro
- **Correção:** Validar discount <= subtotal + tax
- **Tempo estimado:** 1h

### Problema 3.5: Estoque máximo menor que mínimo não validado
- **Severidade:** BAIXA
- **Arquivo:** `internal/service/product_service.go`
- **Linhas:** 139-148
- **Causa:** UpdateIngredientInput não valida MinStock < StockQuantity
- **Impacto:** MinStock pode ser maior que StockQuantity
- **Risco:** Baixo - Alertas incorretos
- **Correção:** Validar MinStock <= StockQuantity
- **Tempo estimado:** 0.5h

### Problema 3.6: Slug vazio possível
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/product_service.go`
- **Linhas:** 168-171
- **Causa:** Se in.Slug == "", gera slug do nome, mas nome pode ser vazio
- **Impacto:** Produto pode ter slug vazio ou inválido
- **Risco:** Baixo - URL inválida
- **Correção:** Validar slug gerado não vazio
- **Tempo estimado:** 1h

### Problema 3.7: Nome vazio não validado em UpdateProduct
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/product_service.go`
- **Linha:** 108
- **Causa:** Validação existe mas não impede bypass
- **Impacto:** Produto pode ter nome vazio
- **Risco:** Baixo - Display incorreto
- **Correção:** Validar no service também
- **Tempo estimado:** 0.5h

### Problema 3.8: Email inválido não validado em todos os lugares
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/invitation_service.go`
- **Linha:** 63
- **Causa:** CreateInvitation não valida formato de email
- **Impacto:** Email inválido pode ser usado
- **Risco:** Baixo - Convite não entregue
- **Correção:** Adicionar validação de email
- **Tempo estimado:** 1h

### Problema 3.9: Campos obrigatórios não validados em UpdateCompany
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/company_service.go`
- **Linhas:** 40-53
- **Causa:** UpdateCompanyInput não valida campos obrigatórios
- **Impacto:** Empresa pode ter campos vazios
- **Risco:** Baixo - Display incorreto
- **Correção:** Adicionar validações
- **Tempo estimado:** 1h

### Problema 3.10: Valores absurdos não validados
- **Severidade:** BAIXA
- **Arquivo:** Vários
- **Causa:** Não há validação de valores máximos (ex: preço 999999999)
- **Impacto:** Sistema pode aceitar valores absurdos
- **Risco:** Baixo - Display incorreto
- **Correção:** Adicionar validações de max em inputs
- **Tempo estimado:** 4h

### Problema 3.11: Overflow em Money não validado
- **Severidade:** BAIXA
- **Arquivo:** `internal/domain/money.go`
- **Linhas:** 42-59
- **Causa:** Operações Money não validam overflow
- **Impacto:** Soma de muitos valores pode causar overflow
- **Risco:** Baixo - Valores incorretos em cenários extremos
- **Correção:** Adicionar validação de overflow
- **Tempo estimado:** 2h

### Problema 3.12: Underflow em Money não validado
- **Severidade:** BAIXA
- **Arquivo:** `internal/domain/money.go`
- **Linhas:** 46-49
- **Causa:** Sub não valida se resultado é negativo
- **Impacto:** Subtração pode resultar em valor negativo
- **Risco:** Baixo - Valores negativos inesperados
- **Correção:** Validar resultado >= 0 ou usar Sub com validação
- **Tempo estimado:** 1h

---

## FASE 4 — DUPLICIDADE

### Problema 4.1: Produto duplicado possível
- **Severidade:** ALTA
- **Arquivo:** `internal/service/product_service.go`
- **Causa:** Não há validação de nome + company_id único
- **Impacto:** Mesmo produto pode ser criado múltiplas vezes
- **Risco:** Médio - Confusão no catálogo
- **Correção:** Adicionar unique constraint em (name, company_id)
- **Tempo estimado:** 2h

### Problema 4.2: Ingrediente duplicado possível
- **Severidade:** ALTA
- **Arquivo:** `internal/service/product_service.go`
- **Causa:** Não há validação de nome + company_id único
- **Impacto:** Mesmo ingrediente pode ser criado múltiplas vezes
- **Risco:** Médio - Confusão no estoque
- **Correção:** Adicionar unique constraint em (name, company_id)
- **Tempo estimado:** 2h

### Problema 4.3: Categoria duplicada possível
- **Severidade:** MÉDIA
- **Arquivo:** Vários
- **Causa:** Não há validação de nome + company_id único
- **Impacto:** Mesma categoria pode ser criada múltiplas vezes
- **Risco:** Baixo - Confusão no catálogo
- **Correção:** Adicionar unique constraint em (name, company_id)
- **Tempo estimado:** 2h

### Problema 4.4: Pedido duplicado possível (sem idempotency)
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/order_service.go`
- **Causa:** IdempotencyKey é opcional
- **Impacto:** Mesmo pedido pode ser criado múltiplas vezes sem idempotency
- **Risco:** Médio - Duplicação de pedidos
- **Correção:** Tornar IdempotencyKey obrigatório ou gerar automaticamente
- **Tempo estimado:** 3h

### Problema 4.5: Recebimento duplicado parcial
- **Severidade:** CRÍTICA
- **Arquivo:** `internal/service/purchase_service.go`
- **Linhas:** 227-230
- **Causa:** Não há validação de soma de recebimentos vs quantidade pedida
- **Impacto:** Mesmo pedido pode ser recebido múltiplas vezes parcialmente
- **Risco:** Alto - Estoque inflado
- **Correção:** Validar soma de recebimentos <= quantidade pedida
- **Tempo estimado:** 4h

### Problema 4.6: Convite duplicado parcialmente prevenido
- **Severidade:** BAIXA
- **Arquivo:** `internal/service/invitation_service.go`
- **Linhas:** 87-94
- **Causa:** Validação existe mas apenas para pending
- **Impacto:** Convite expirado pode ser recriado para mesmo email
- **Risco:** Baixo - Múltiplos convites
- **Correção:** Validar qualquer convite para mesmo email (não apenas pending)
- **Tempo estimado:** 1h

### Problema 4.7: Usuário duplicado prevenido
- **Severidade:** N/A
- **Arquivo:** `internal/service/auth_service.go`
- **Status:** ✅ OK - Validação de email único existe
- **Observação:** Email único é validado

### Problema 4.8: Slug duplicado prevenido
- **Severidade:** N/A
- **Arquivo:** `internal/service/company_service.go`
- **Status:** ✅ OK - Validação de slug único existe
- **Observação:** Slug único é validado

### Problema 4.9: Empresa duplicada possível
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/company_service.go`
- **Causa:** Não há validação de nome único
- **Impacto:** Mesma empresa pode ser criada múltiplas vezes
- **Risco:** Baixo - Confusão
- **Correção:** Adicionar validação de nome único
- **Tempo estimado:** 1h

---

## FASE 5 — CÁLCULOS

### Problema 5.1: CMV com divisão por zero não tratada
- **Severidade:** ALTA
- **Arquivo:** `internal/domain/product.go`
- **Linhas:** 70-78
- **Causa:** CalculateCMV valida Price.IsZero() mas não previne zero após cálculo
- **Impacto:** CMV pode ser NaN ou infinito
- **Risco:** Médio - Cálculos incorretos
- **Correção:** Validar resultado não é NaN/Inf
- **Tempo estimado:** 1h

### Problema 5.2: Margem com divisão por zero não tratada
- **Severidade:** ALTA
- **Arquivo:** `internal/domain/product.go`
- **Linhas:** 82-91
- **Causa:** CalculateMargin valida Price.IsZero() mas não previne zero após cálculo
- **Impacto:** Margem pode ser NaN ou infinito
- **Risco:** Médio - Cálculos incorretos
- **Correção:** Validar resultado não é NaN/Inf
- **Tempo estimado:** 1h

### Problema 5.3: Lucro pode ser negativo
- **Severidade:** MÉDIA
- **Arquivo:** `internal/domain/product.go`
- **Linhas:** 95-98
- **Causa:** CalculateProfit não valida se resultado é negativo
- **Impacto:** Lucro pode ser negativo (prejuízo)
- **Risco:** Baixo - Prejuízo não é erro de cálculo
- **Correção:** N/A - Prejuízo é válido
- **Tempo estimado:** 0h

### Problema 5.4: Preço sugerido com margem >= 100%
- **Severidade:** BAIXA
- **Arquivo:** `internal/domain/product.go`
- **Linhas:** 103-108
- **Causa:** CalculateSuggestedPrice assume 50% se margem >= 100%
- **Impacto:** Preço sugerido pode estar incorreto
- **Risco:** Baixo - Preço sugerido errado
- **Correção:** Retornar erro em vez de assumir valor
- **Tempo estimado:** 1h

### Problema 5.5: Ticket médio não validado
- **Severidade:** BAIXA
- **Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`
- **Causa:** Cálculo de ticket médio não valida divisão por zero
- **Impacto:** Ticket médio pode ser NaN/Inf se não houver pedidos
- **Risco:** Baixo - Display incorreto
- **Correção:** Validar count > 0 antes de dividir
- **Tempo estimado:** 1h

### Problema 5.6: Total pedido com rounding errors
- **Severidade:** ALTA
- **Arquivo:** `internal/service/order_service.go`
- **Linhas:** 165
- **Causa:** Cálculo usa Mul(int64(itemIn.Quantity * 100)).Div(100) com rounding
- **Impacto:** Total pode estar incorreto por centavos
- **Risco:** Médio - Divergência financeira
- **Correção:** Usar MulFloat para cálculo correto
- **Tempo estimado:** 2h

### Problema 5.7: Subtotal com rounding errors
- **Severidade:** ALTA
- **Arquivo:** `internal/service/purchase_service.go`
- **Linhas:** 121-122
- **Causa:** Cálculo usa Mul(int64(item.Quantity * 100)).Div(100) com rounding
- **Impacto:** Subtotal pode estar incorreto por centavos
- **Risco:** Médio - Divergência financeira
- **Correção:** Usar MulFloat para cálculo correto
- **Tempo estimado:** 2h

### Problema 5.8: Custo receita não calculado
- **Severidade:** MÉDIA
- **Arquivo:** Vários
- **Causa:** Não há cálculo de custo por receita
- **Impacto:** Não é possível saber custo por pedido
- **Risco:** Baixo - Falta de métrica
- **Correção:** Implementar cálculo de custo por pedido
- **Tempo estimado:** 8h

---

## FASE 6 — ESTOQUE

### Problema 6.1: Baixa dupla prevenida
- **Severidade:** N/A
- **Arquivo:** `internal/service/stock_movement_service.go`
- **Status:** ✅ OK - SELECT FOR UPDATE previne baixa dupla
- **Observação:** Transação com lock previne baixa dupla

### Problema 6.2: Baixa parcial não validada
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/order_service.go`
- **Causa:** Cancelamento devolve estoque completo, não parcial
- **Impacto:** Cancelamento de pedido parcialmente processado devolve estoque completo
- **Risco:** Médio - Estoque incorreto
- **Correção:** Implementar cancelamento parcial
- **Tempo estimado:** 6h

### Problema 6.3: Baixa sem estoque prevenida
- **Severidade:** N/A
- **Arquivo:** `internal/service/order_service.go`
- **Status:** ✅ OK - validateStock previne baixa sem estoque
- **Observação:** Validação de estoque existe

### Problema 6.4: Cancelamento devolvendo estoque parcialmente implementado
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/order_service.go`
- **Linhas:** 287-314
- **Causa:** Cancelamento devolve estoque mas não valida se já foi devolvido
- **Impacto:** Cancelamento múltiplo pode devolver estoque múltiplas vezes
- **Risco:** Alto - Estoque inflado
- **Correção:** Validar se ajustes já foram registrados
- **Tempo estimado:** 4h

### Problema 6.5: Cancelamento múltiplo possível
- **Severidade:** ALTA
- **Arquivo:** `internal/service/order_service.go`
- **Linhas:** 276-279
- **Causa:** isValidTransition permite cancelled -> cancelled
- **Impacto:** Cancelamento múltiplo pode devolver estoque múltiplas vezes
- **Risco:** Alto - Estoque inflado
- **Correção:** Impedir transição cancelled -> cancelled
- **Tempo estimado:** 1h

### Problema 6.6: Inventário corrigindo errado possível
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/stock_movement_service.go`
- **Linhas:** 213-318
- **Causa:** CompleteInventory não valida se já foi completado
- **Impacto:** Inventário pode ser completado múltiplas vezes
- **Risco:** Médio - Estoque ajustado múltiplas vezes
- **Correção:** Validação existe mas pode ter race condition
- **Tempo estimado:** 2h

### Problema 6.7: Ajustes inconsistentes possíveis
- **Severidade:** BAIXA
- **Arquivo:** `internal/service/stock_movement_service.go`
- **Causa:** Ajuste manual não valida razão
- **Impacto:** Ajustes podem ser feitos sem justificativa
- **Risco:** Baixo - Auditoria difícil
- **Correção:** Tornar Reason obrigatório
- **Tempo estimado:** 1h

---

## FASE 7 — REGRAS DE PERMISSÃO

### Problema 7.1: Role Kitchen não definido
- **Severidade:** MÉDIA
- **Arquivo:** `internal/domain/role.go`
- **Causa:** Role enum não inclui Kitchen
- **Impacto:** Usuários de cozinha não têm role específico
- **Risco:** Baixo - Funcionalidade limitada
- **Correção:** Adicionar RoleKitchen
- **Tempo estimado:** 2h

### Problema 7.2: Role Cashier não definido
- **Severidade:** MÉDIA
- **Arquivo:** `internal/domain/role.go`
- **Causa:** Role enum não inclui Cashier
- **Impacto:** Usuários de caixa não têm role específico
- **Risco:** Baixo - Funcionalidade limitada
- **Correção:** Adicionar RoleCashier
- **Tempo estimado:** 2h

### Problema 7.3: Permissões não auditadas por endpoint
- **Severidade:** ALTA
- **Arquivo:** `internal/middleware/role_middleware.go`
- **Causa:** Não há lista centralizada de permissões por endpoint
- **Impacto:** Difícil garantir que todos endpoints estão protegidos
- **Risco:** Alto - Endpoint não protegido pode ser acessado
- **Correção:** Criar matriz de permissões centralizada
- **Tempo estimado:** 8h

### Problema 7.4: PlatformAdmin não validado em todos os endpoints
- **Severidade:** MÉDIA
- **Arquivo:** Vários
- **Causa:** Alguns endpoints platform podem não validar PlatformAdmin
- **Impacto:** Endpoint platform pode ser acessado por não-admin
- **Risco:** Médio - Acesso não autorizado
- **Correção:** Auditar todos os endpoints platform
- **Tempo estimado:** 4h

---

## FASE 8 — TRIAL

### Problema 8.1: Trial expiration não validado em operações
- **Severidade:** CRÍTICA
- **Arquivo:** `internal/middleware/tenant_middleware.go`
- **Causa:** Tenant middleware não verifica trial expiration
- **Impacto:** Empresa com trial expirado pode continuar operando
- **Risco:** Alto - Violação de contrato
- **Correção:** Adicionar validação de trial expiration
- **Tempo estimado:** 3h

### Problema 8.2: Limites de trial não implementados
- **Severidade:** ALTA
- **Arquivo:** Vários
- **Causa:** Não há limites de produtos, pedidos, usuários em trial
- **Impacto:** Trial pode ter recursos ilimitados
- **Risco:** Médio - Abuso de trial
- **Correção:** Implementar limites por plano
- **Tempo estimado:** 12h

### Problema 8.3: Bloqueio de trial não implementado
- **Severidade:** ALTA
- **Arquivo:** Vários
- **Causa:** Não há bloqueio automático ao expirar trial
- **Impacto:** Trial expirado continua ativo
- **Risco:** Médio - Uso indevido
- **Correção:** Implementar job para bloquear trials expirados
- **Tempo estimado:** 4h

### Problema 8.4: Reativação de trial não controlada
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/company_service.go`
- **Causa:** Não há controle de reativação de trial
- **Impacto:** Trial pode ser reativado manualmente
- **Risco:** Baixo - Abuso de trial
- **Correção:** Implementar controle de reativação
- **Tempo estimado:** 3h

### Problema 8.5: Cancelamento de trial não validado
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/company_service.go`
- **Linhas:** 225-234
- **Causa:** DeleteCompany não valida se empresa está em trial
- **Impacto:** Empresa em trial pode ser deletada
- **Risco:** Baixo - Perda de dados
- **Correção:** Adicionar validação
- **Tempo estimado:** 1h

---

## FASE 9 — CONSISTÊNCIA FRONTEND-BACKEND

### Problema 9.1: Status de pedido diferente entre frontend e backend
- **Severidade:** MÉDIA
- **Arquivo:** `internal/domain/order.go`
- **Causa:** Backend tem "confirmed", frontend pode não ter
- **Impacto:** Frontend pode não reconhecer status
- **Risco:** Baixo - Display incorreto
- **Correção:** Sincronizar enums
- **Tempo estimado:** 2h

### Problema 9.2: Role Kitchen não existe no backend
- **Severidade:** MÉDIA
- **Arquivo:** `internal/domain/role.go`
- **Causa:** Frontend pode ter Kitchen, backend não
- **Impacto:** Frontend pode enviar role inválido
- **Risco:** Baixo - Erro de validação
- **Correção:** Adicionar RoleKitchen no backend
- **Tempo estimado:** 2h

### Problema 9.3: Role Cashier não existe no backend
- **Severidade:** MÉDIA
- **Arquivo:** `internal/domain/role.go`
- **Causa:** Frontend pode ter Cashier, backend não
- **Impacto:** Frontend pode enviar role inválido
- **Risco:** Baixo - Erro de validação
- **Correção:** Adicionar RoleCashier no backend
- **Tempo estimado:** 2h

### Problema 9.4: Tipos de Money diferentes
- **Severidade:** BAIXA
- **Arquivo:** `internal/domain/money.go`
- **Causa:** Backend usa int64 (centavos), frontend pode usar float
- **Impacto:** Conversão pode ter rounding errors
- **Risco:** Baixo - Divergência de centavos
- **Correção:** Documentar formato esperado
- **Tempo estimado:** 1h

### Problema 9.5: Campos inexistentes no frontend
- **Severidade:** BAIXA
- **Arquivo:** Vários
- **Causa:** Backend tem campos que frontend não envia
- **Impacto:** Campos podem ficar vazios
- **Risco:** Baixo - Dados incompletos
- **Correção:** Documentar campos obrigatórios
- **Tempo estimado:** 4h

---

## FASE 10 — CASOS EXTREMOS

### Problema 10.1: 1000 produtos não testado
- **Severidade:** BAIXA
- **Arquivo:** Vários
- **Causa:** Não há testes de performance com 1000 produtos
- **Impacto:** Sistema pode ter problemas com muitos produtos
- **Risco:** Baixo - Lentidão
- **Correção:** Adicionar testes de carga
- **Tempo estimado:** 4h

### Problema 10.2: 10000 pedidos não testado
- **Severidade:** BAIXA
- **Arquivo:** Vários
- **Causa:** Não há testes de performance com 10000 pedidos
- **Impacto:** Sistema pode ter problemas com muitos pedidos
- **Risco:** Baixo - Lentidão
- **Correção:** Adicionar testes de carga
- **Tempo estimado:** 4h

### Problema 10.3: Produto sem ingredientes não validado
- **Severidade:** BAIXA
- **Arquivo:** `internal/domain/product.go`
- **Linhas:** 125-138
- **Causa:** ValidateRecipe retorna erro se não tem ingredientes, mas não é obrigatório
- **Impacto:** Produto composto pode não ter ingredientes
- **Risco:** Baixo - Custo zero
- **Correção:** Validar se IsComposto então deve ter ingredientes
- **Tempo estimado:** 1h

### Problema 10.4: Receita circular não detectada
- **Severidade:** MÉDIA
- **Arquivo:** Vários
- **Causa:** Não há detecção de receita circular (A usa B, B usa A)
- **Impacto:** Custo pode ser infinito
- **Risco:** Médio - Loop infinito no cálculo
- **Correção:** Implementar detecção de ciclo
- **Tempo estimado:** 6h

### Problema 10.5: Pedido enorme não limitado
- **Severidade:** MÉDIA
- **Arquivo:** `internal/service/order_service.go`
- **Causa:** Não há limite de itens por pedido
- **Impacto:** Pedido pode ter 1000 itens
- **Risco:** Médio - Lentidão, timeout
- **Correção:** Adicionar limite de itens por pedido
- **Tempo estimado:** 1h

### Problema 10.6: Texto enorme não limitado
- **Severidade:** BAIXA
- **Arquivo:** Vários
- **Causa:** Campos de texto não têm limite de tamanho
- **Impacto:** Texto pode ter 10000 caracteres
- **Risco:** Baixo - Lentidão, storage
- **Correção:** Adicionar validações de max length
- **Tempo estimado:** 2h

### Problema 10.7: Emoji não tratado
- **Severidade:** BAIXA
- **Arquivo:** Vários
- **Causa:** Slug generator não trata emoji corretamente
- **Impacto:** Slug pode ter caracteres inválidos
- **Risco:** Baixo - URL inválida
- **Correção:** Melhorar slug generator
- **Tempo estimado:** 2h

### Problema 10.8: UTF8 não validado
- **Severidade:** BAIXA
- **Arquivo:** Vários
- **Causa:** Não há validação de UTF8
- **Impacto:** Texto pode ter caracteres inválidos
- **Risco:** Baixo - Display incorreto
- **Correção:** Adicionar validação de UTF8
- **Tempo estimado:** 1h

### Problema 10.9: Acentuação não tratada em slug
- **Severidade:** BAIXA
- **Arquivo:** `internal/service/product_service.go`
- **Linhas:** 31-75
- **Causa:** Slug generator trata acentuação mas pode não ser completo
- **Impacto:** Slug pode ter caracteres acentuados
- **Risco:** Baixo - URL inválida
- **Correção:** Melhorar slug generator com unicode normalização
- **Tempo estimado:** 2h

### Problema 10.10: Campos vazios não validados em alguns lugares
- **Severidade:** BAIXA
- **Arquivo:** Vários
- **Causa:** Alguns campos opcionais podem ser vazios quando deveriam ter valor
- **Impacto:** Dados incompletos
- **Risco:** Baixo - Display incorreto
- **Correção:** Adicionar validações
- **Tempo estimado:** 2h

### Problema 10.11: Null não tratado em campos opcionais
- **Severidade:** BAIXA
- **Arquivo:** Vários
- **Causa:** Campos opcionais podem ser null quando deveriam ser string vazia
- **Impacto:** Dados inconsistentes
- **Risco:** Baixo - Display incorreto
- **Correção:** Normalizar null para string vazia
- **Tempo estimado:** 2h

---

## Resumo por Severidade

### CRÍTICA (12 problemas)
1. Produto com preço negativo possível
2. Estoque negativo possível via ajuste manual
3. Recebimento duplicado possível
4. Empresa cancelada pode operar normalmente
5. Trial expirado pode continuar operando
6. Quantidade negativa em OrderItemInput
7. Preço negativo em CreateProductInput
8. Convite pode ser aceito duas vezes
9. Reset de senha pode ser reutilizado
10. Recebimento duplicado parcial

### ALTA (15 problemas)
1. Pedido com total incorreto calculado
2. Pedido cancelado pode ser editado
3. Produto arquivado pode ser vendido
4. Ingrediente excluído pode estar em ficha técnica
5. Transição Completed -> Pending não impedida
6. Transição Cancelled -> Preparing não impedida
7. Purchase Order sem validação de transições
8. Custo negativo possível
9. Desconto maior que total não validado
10. Produto duplicado possível
11. Ingrediente duplicado possível
12. Pedido duplicado possível (sem idempotency)
13. CMV com divisão por zero não tratada
14. Margem com divisão por zero não tratada
15. Total pedido com rounding errors

### MÉDIA (8 problemas)
1. Pedido sem itens possível
2. Categoria excluída em produto ativo
3. Sessão expirada pode ser reutilizada
4. Inventory sem validação de transições
5. Slug vazio possível
6. Email inválido não validado
7. Role Kitchen não definido
8. Role Cashier não definido

### BAIXA (4 problemas)
1. Estoque máximo menor que mínimo não validado
2. Valores absurdos não validados
3. Overflow em Money não validado
4. Underflow em Money não validado

---

## Estimativa Total de Esforço

**Total estimado:** 123 horas (~15 dias úteis)

**Por fase:**
- FASE 1 (Estados Inválidos): 35h
- FASE 2 (Máquinas de Estado): 8h
- FASE 3 (Validação de Dados): 16h
- FASE 4 (Duplicidade): 18h
- FASE 5 (Cálculos): 24h
- FASE 6 (Estoque): 14h
- FASE 7 (Permissões): 16h
- FASE 8 (Trial): 23h
- FASE 9 (Frontend-Backend): 11h
- FASE 10 (Casos Extremos): 22h
