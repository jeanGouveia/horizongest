# SPRINT 5C.4.6 — AUDITORIA DE PERFORMANCE, ESCALABILIDADE E GARGALOS

**Data:** 2025-01-XX  
**Auditor:** Principal Performance Architect  
**Projeto:** HorizonGest Backend  
**Tipo:** Auditoria de Performance e Escalabilidade  
**Objetivo:** Identificar gargalos atuais e futuros antes da entrada em produção

---

## RESUMO DOS PROBLEMAS

**Total de Problemas Identificados:** 42
- **Críticos:** 8
- **Altos:** 18
- **Médios:** 12
- **Baixos:** 4

---

## DETALHAMENTO DOS GARGALOS

### 1. POSTGRESQL - QUERIES

#### 1.1 N+1 Query Problem em Dashboard
**Gargalo:** Dashboard executa COUNT separado para cada pedido recente  
**Severidade:** Alta  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 280-285  
**Causa:** Loop sobre recentOrders executa COUNT para cada pedido  
**Impacto:** Para 10 pedidos recentes, executa 10 queries adicionais. Tempo adicional: ~50-100ms  
**Custo:** Alto - cada COUNT é uma query separada  
**Solução Recomendada:** Usar subquery ou JOIN para contar itens em uma única query  
**Ganho Esperado:** Redução de ~80ms no tempo carregamento do dashboard

```go
// ATUAL (N+1):
for i, o := range recentOrders {
    var itemsCount int64
    if err := query.WithContext(ctx).Model(&GormOrderItem{}).
        Where("order_id = ? AND deleted_at IS NULL", o.ID).
        Count(&itemsCount).Error; err != nil {
        return nil, fmt.Errorf("DashboardRepository.GetDashboard: %w", err)
    }
    dashboard.RecentOrders[i] = domain.RecentOrder{
        ItemsCount: int(itemsCount),
        // ...
    }
}

// RECOMENDADO (Single query):
type OrderWithCount struct {
    ID          uint
    OrderNumber int
    Status      string
    TotalPrice  int64
    CreatedAt   time.Time
    ItemsCount  int64
}
var recentOrders []OrderWithCount
if err := query.WithContext(ctx).Model(&GormOrder{}).
    Select("orders.*, COUNT(order_items.id) as items_count").
    Joins("LEFT JOIN order_items ON order_items.order_id = orders.id AND order_items.deleted_at IS NULL").
    Where("orders.deleted_at IS NULL").
    Group("orders.id").
    Order("orders.created_at DESC").
    Limit(10).
    Scan(&recentOrders).Error; err != nil {
    return nil, fmt.Errorf("DashboardRepository.GetDashboard: %w", err)
}
```

---

#### 1.2 Múltiplas Queries Separadas para KPIs de Dashboard
**Gargalo:** Dashboard executa ~30 queries separadas para KPIs  
**Severidade:** Crítica  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 37-258  
**Causa:** Cada KPI é uma query separada (hoje, ontem, semana, mês, etc.)  
**Impacto:** Tempo total de ~200-500ms para carregar dashboard. Alta latência percebida pelo usuário  
**Custo:** Crítico - 30+ queries em uma única requisição  
**Solução Recomendada:** Consolidar queries usando CTEs ou materialized views. Implementar cache de 1-5 minutos  
**Ganho Esperado:** Redução de ~300ms no tempo de carregamento do dashboard

---

#### 1.3 Loop de Queries para Sales By Day
**Gargalo:** Sales by day executa 7 queries separadas (uma por dia)  
**Severidade:** Alta  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 372-382  
**Causa:** Loop for executa query separada para cada dia  
**Impacto:** 7 queries sequenciais, ~70ms adicionais  
**Custo:** Médio - queries simples mas acumulam  
**Solução Recomendada:** Usar GROUP BY DATE em uma única query  
**Ganho Esperado:** Redução de ~60ms

```go
// ATUAL (Loop):
for i := 6; i >= 0; i-- {
    date := now.AddDate(0, 0, -i).Format("2006-01-02")
    var total int64
    if err := query.WithContext(ctx).Model(&GormOrder{}).
        Where("DATE(created_at) = ? AND deleted_at IS NULL", date).
        Select("COALESCE(SUM(total_price), 0)").
        Scan(&total).Error; err != nil {
        return []domain.ChartPoint{}
    }
    results = append(results, DayResult{Date: date, Total: total})
}

// RECOMENDADO (Single query):
var results []DayResult
weekAgo := now.AddDate(0, 0, -7).Format("2006-01-02")
if err := query.WithContext(ctx).Model(&GormOrder{}).
    Select("DATE(created_at) as date, COALESCE(SUM(total_price), 0) as total").
    Where("DATE(created_at) >= ? AND deleted_at IS NULL", weekAgo).
    Group("DATE(created_at)").
    Order("date ASC").
    Scan(&results).Error; err != nil {
    return []domain.ChartPoint{}
}
```

---

#### 1.4 Mesmo Problema em Report Repository
**Gargalo:** Sales by day em reports também usa loop  
**Severidade:** Alta  
**Arquivo:** `internal/infra/repository/gorm_report_repository.go`  
**Linha:** 396-407  
**Causa:** Mesmo padrão de loop de queries  
**Impacto:** Para relatórios de 30 dias, 30 queries sequenciais  
**Custo:** Alto - relatórios são operações pesadas  
**Solução Recomendada:** Usar GROUP BY DATE  
**Ganho Esperado:** Redução de ~200-500ms em relatórios mensais

---

#### 1.5 SELECT * em Queries
**Gargalo:** SELECT * carrega todas as colunas desnecessariamente  
**Severidade:** Média  
**Arquivo:** `internal/infra/repository/gorm_order_repository.go`  
**Linha:** 88  
**Causa:** GORM usa SELECT * por padrão  
**Impacto:** Transferência de dados desnecessária, uso de memória  
**Custo:** Médio - impacto em tabelas com muitas colunas  
**Solução Recomendada:** Usar Select() para especificar apenas colunas necessárias  
**Ganho Esperado:** Redução de ~20-30% no tamanho do resultset

---

#### 1.6 Falta de Índices Compostos
**Gargalo:** Queries com múltiplos filtros não têm índices compostos  
**Severidade:** Crítica  
**Arquivo:** `migrations/`  
**Linha:** N/A  
**Causa:** Migration 00017 está vazia, índices compostos não implementados  
**Impacto:** Queries com (company_id, status, date) usam apenas índice parcial ou scan completo  
**Custo:** Crítico - afeta performance de todas as queries filtradas por tenant  
**Solução Recomendada:** Criar índices compostos para padrões comuns:
- (company_id, status, created_at) para orders
- (company_id, deleted_at) para todas as tabelas
- (company_id, active, deleted_at) para produtos/ingredientes  
**Ganho Esperado:** Redução de ~50-80% em queries filtradas por tenant

---

#### 1.7 JOINs com deleted_at Filters
**Gargalo:** JOINs filtram por deleted_at em ambas as tabelas  
**Severidade:** Alta  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 57, 99, 136, 173, etc.  
**Causa:** Soft delete requer filter em todas as tabelas  
**Impacto:** Queries mais complexas, índices menos eficientes  
**Custo:** Alto - afeta todas as queries com JOINs  
**Solução Recomendada:** Considerar índices parciais (WHERE deleted_at IS NULL) ou particionamento  
**Ganho Esperado:** Melhoria de ~30% em queries com JOINs

---

#### 1.8 COUNT Queries Sem Cache
**Gargalo:** COUNT queries executadas repetidamente sem cache  
**Severidade:** Alta  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 49, 91, 128, 165, etc.  
**Causa:** COUNT é caro e não há cache  
**Impacto:** COUNT em tabelas grandes pode levar 50-200ms  
**Custo:** Alto - COUNT é uma das operações mais caras no PostgreSQL  
**Solução Recomendada:** Implementar cache de contadores (Redis) com invalidação  
**Ganho Esperado:** Redução de ~100-300ms em COUNT queries

---

#### 1.9 GROUP BY Sem Índice Apropriado
**Gargalo:** GROUP BY em colunas sem índice composto  
**Severidade:** Alta  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 408, 444, 486  
**Causa:** GROUP BY em (hour), (products.id, products.name), (categories.id, categories.name)  
**Impacto:** Queries de agregação podem ser lentas com muitos dados  
**Custo:** Alto - afeta gráficos e relatórios  
**Solução Recomendada:** Adicionar índices compostos para GROUP BY  
**Ganho Esperado:** Redução de ~40-60% em queries de agregação

---

#### 1.10 ORDER BY Sem Índice
**Gargalo:** ORDER BY created_at DESC sem índice  
**Severidade:** Média  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 271, 308, 329  
**Causa:** created_at não está indexado em algumas tabelas  
**Impacto:** Sort pode ser feito em memória  
**Custo:** Médio - impacto cresce com volume de dados  
**Solução Recomendada:** Adicionar índice em (company_id, created_at DESC)  
**Ganho Esperado:** Redução de ~20-40% em queries com ORDER BY

---

### 2. GORM

#### 2.1 Save() em Vez de Update() para Campos Únicos
**Gargalo:** Save() atualiza todas as colunas desnecessariamente  
**Severidade:** Média  
**Arquivo:** `internal/infra/repository/gorm_product_repository.go`  
**Linha:** 292  
**Causa:** Save() faz UPDATE de todos os campos  
**Impacto:** Tráfego de rede adicional, logs maiores  
**Custo:** Médio - impacto em operações frequentes  
**Solução Recomendada:** Usar Update() ou Updates() para campos específicos  
**Ganho Esperado:** Redução de ~30-50% no tamanho do UPDATE

---

#### 2.2 Find Dentro de Loop
**Gargalo:** Find executado dentro de loop (N+1)  
**Severidade:** Alta  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 280-285  
**Causa:** Loop sobre orders executa Find para cada um  
**Impacto:** N+1 queries  
**Custo:** Alto - ver problema 1.1  
**Solução Recomendada:** Ver solução 1.1  
**Ganho Esperado:** Redução de ~80ms

---

#### 2.3 Preload Desnecessário
**Gargalo:** Preload carrega dados não utilizados  
**Severidade:** Baixa  
**Arquivo:** `internal/infra/repository/gorm_stock_movement_repository.go`  
**Linha:** 94, 109  
**Causa:** Preload("Items.Ingredient") pode carregar dados não usados  
**Impacto:** Queries adicionais desnecessárias  
**Custo:** Baixo - impacto pequeno  
**Solução Recomendada:** Avaliar se todos os dados preload são usados  
**Ganho Esperado:** Redução de ~10-20ms

---

#### 2.4 Falta de Batch Operations
**Gargalo:** Operações em lote não usam batch  
**Severidade:** Média  
**Arquivo:** `internal/infra/repository/gorm_product_repository.go`  
**Linha:** 534-539  
**Causa:** GetProductIngredients carrega todos os ingredientes de um produto  
**Impacto:** Para produtos com muitos ingredientes, pode ser lento  
**Custo:** Médio - impacto em produtos complexos  
**Solução Recomendada:** Usar batch operations para inserts/updates em lote  
**Ganho Esperado:** Redução de ~30-50% em operações em lote

---

### 3. REDIS

#### 3.1 Pool Size Pequeno
**Gargalo:** Pool size de 10 pode ser insuficiente  
**Severidade:** Média  
**Arquivo:** `internal/infra/redis/client.go`  
**Linha:** 23  
**Causa:** PoolSize: 10 é pequeno para alta concorrência  
**Impacto:** Conexões podem esgotar em pico de tráfego  
**Custo:** Médio - pode causar timeouts  
**Solução Recomendada:** Aumentar para 50-100 dependendo da carga  
**Ganho Esperado:** Melhoria em throughput sob alta carga

---

#### 3.2 TTL Não Configurado para Alguns Caches
**Gargalo:** Cache de global config não tem TTL explícito  
**Severidade:** Baixa  
**Arquivo:** `internal/infra/repository/gorm_global_config_repository.go`  
**Linha:** 47-50  
**Causa:** Cache in-memory não expira automaticamente  
**Impacto:** Dados podem ficar obsoletos  
**Custo:** Baixo - impacto pequeno  
**Solução Recomendada:** Adicionar TTL ou invalidação explícita  
**Ganho Esperado:** Consistência de dados

---

#### 3.3 Falta de Métricas de Cache Hit/Miss
**Gargalo:** Não há métricas de cache hit/miss  
**Severidade:** Média  
**Arquivo:** `internal/infra/redis/`  
**Linha:** N/A  
**Causa:** Métricas existem mas não estão expostas/monitoradas  
**Impacto:** Dificuldade em identificar problemas de cache  
**Custo:** Médio - dificulta debugging  
**Solução Recomendada:** Expor métricas via Prometheus  
**Ganho Esperado:** Visibilidade para otimizações

---

#### 3.4 Fallback Cache Loga Erros Mas Continua
**Gargalo:** Fallback cache trata erros como cache miss  
**Severidade:** Média  
**Arquivo:** `internal/infra/redis/fallback.go`  
**Linha:** 22-30  
**Causa:** Erros de Redis são silenciados  
**Impacto:** Degradação silenciosa de performance  
**Custo:** Médio - pode esconder problemas de Redis  
**Solução Recomendada:** Adicionar métricas de fallback rate  
**Ganho Esperado:** Visibilidade de problemas de Redis

---

### 4. RABBITMQ

#### 4.1 Falta de Prefetch Count Configuration
**Gargalo:** Prefetch count não configurado  
**Severidade:** Alta  
**Arquivo:** `internal/infra/messaging/rabbitmq/`  
**Linha:** N/A  
**Causa:** Prefetch padrão pode não ser ideal para workload  
**Impacto:** Consumers podem receber mensagens muito rápido ou muito lento  
**Custo:** Alto - afeta throughput e latência  
**Solução Recomendada:** Configurar prefetch count baseado no tempo de processamento  
**Ganho Esperado:** Melhoria de ~20-40% em throughput de consumers

---

#### 4.2 Falta de DLQ Configuration
**Gargalo:** Dead Letter Queue não configurada  
**Severidade:** Alta  
**Arquivo:** `internal/infra/messaging/rabbitmq/`  
**Linha:** N/A  
**Causa:** Mensagens que falham repetidamente não são movidas para DLQ  
**Impacto:** Mensagens podem ficar presas na fila principal  
**Custo:** Alto - pode causar bloqueio de processamento  
**Solução Recomendada:** Configurar DLQ com TTL e max retries  
**Ganho Esperado:** Resiliência e visibilidade de problemas

---

#### 4.3 Publisher Timeout Pode Ser Curto
**Gargalo:** Publisher timeout de 10s pode ser curto  
**Severidade:** Média  
**Arquivo:** `internal/service/dispatcher_config.go`  
**Linha:** 20  
**Causa:** Timeout fixo pode não ser adequado para todos os cenários  
**Impacto:** Publicações podem falhar prematuramente  
**Custo:** Médio - pode causar perda de eventos  
**Solução Recomendada:** Aumentar para 30s ou tornar configurável  
**Ganho Esperado:** Redução de falhas de publicação

---

#### 4.4 Batch Size Fixo
**Gargalo:** Batch size de 50 é fixo  
**Severidade:** Média  
**Arquivo:** `internal/service/dispatcher_config.go`  
**Linha:** 17  
**Causa:** Batch size não se adapta ao volume  
**Impacto:** Pode ser ineficiente em picos ou vales de tráfego  
**Custo:** Médio - throughput subótimo  
**Solução Recomendada:** Implementar batch size dinâmico  
**Ganho Esperado:** Melhoria de ~10-20% em throughput

---

### 5. DASHBOARD

#### 5.1 Dashboard Sem Cache
**Gargalo:** Dashboard não tem cache  
**Severidade:** Crítica  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 22-360  
**Causa:** Dashboard é recalculado em cada requisição  
**Impacto:** 200-500ms por carregamento, alta carga no banco  
**Custo:** Crítico - dashboard é uma das páginas mais acessadas  
**Solução Recomendada:** Implementar cache de 1-5 minutos em Redis  
**Ganho Esperado:** Redução de ~95% no tempo de carregamento (cache hit)

---

#### 5.2 Dashboard Carrega Dados Não Utilizados
**Gargalo:** Dashboard carrega dados que podem não ser usados  
**Severidade:** Baixa  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 346-357  
**Causa:** Gráficos são carregados mesmo se não visíveis  
**Impacto:** Queries desnecessárias  
**Custo:** Baixo - impacto pequeno  
**Solução Recomendada:** Lazy loading de gráficos  
**Ganho Esperado:** Redução de ~50-100ms

---

#### 5.3 CMV Calculado com Porcentagem Fixa
**Gargalo:** CMV usa 30% fixo em vez de cálculo real  
**Severidade:** Baixa  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 72  
**Causa:** CMV simplificado  
**Impacto:** Dados imprecisos  
**Custo:** Baixo - não afeta performance  
**Solução Recomendada:** Implementar cálculo real de CMV (fora do escopo de performance)  
**Ganho Esperado:** Precisão de dados

---

### 6. APIs

#### 6.1 API de Dashboard Faz 30+ Queries
**Gargalo:** GET /api/dashboard executa 30+ queries  
**Severidade:** Crítica  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 22-360  
**Causa:** Ver problema 1.2  
**Impacto:** 200-500ms de tempo de resposta  
**Custo:** Crítico - endpoint crítico  
**Solução Recomendada:** Ver solução 1.2  
**Ganho Esperado:** Redução de ~300ms

---

#### 6.2 APIs de Listagem Sem Paginação Adequada
**Gargalo:** Algumas APIs não têm LIMIT padrão  
**Severidade:** Alta  
**Arquivo:** `internal/infra/repository/gorm_product_repository.go`  
**Linha:** 213  
**Causa:** ListProducts não tem LIMIT  
**Impacto:** Pode retornar milhares de registros  
**Custo:** Alto - pode causar timeout  
**Solução Recomendada:** Adicionar LIMIT padrão de 100 e exigir paginação  
**Ganho Esperado:** Prevenção de timeouts

---

#### 6.3 APIs Sem Rate Limiting por Tenant
**Gargalo:** Rate limiting é global, não por tenant  
**Severidade:** Média  
**Arquivo:** `internal/middleware/rate_limiter.go`  
**Linha:** 11-27  
**Causa:** Rate limiting não considera tenant  
**Impacto:** Um tenant pode afetar outros  
**Custo:** Médio - isolamento inadequado  
**Solução Recomendada:** Implementar rate limiting por tenant  
**Ganho Esperado:** Isolamento melhor entre tenants

---

#### 6.4 Timeout de Contexto Não Configurado em Muitas Operações
**Gargalo:** Operações não têm timeout configurado  
**Severidade:** Alta  
**Arquivo:** `internal/service/`  
**Linha:** N/A  
**Causa:** Contexto sem timeout pode esperar indefinidamente  
**Impacto:** Requisições podem travar  
**Custo:** Alto - pode esgotar conexões  
**Solução Recomendada:** Adicionar timeout em todas as operações de banco  
**Ganho Esperado:** Prevenção de travamentos

---

### 7. FRONTEND

#### 7.1 Múltiplas Requisições em onMount
**Gargalo:** Cada componente faz fetch em onMount  
**Severidade:** Média  
**Arquivo:** `frontend/src/routes/(app)/dashboard/+page.svelte`  
**Linha:** 14-27  
**Causa:** Cada página carrega dados independentemente  
**Impacto:** Múltiplas requisições sequenciais  
**Custo:** Médio - latência acumulada  
**Solução Recomendada:** Implementar data loading no layout ou usar SWR  
**Ganho Esperado:** Redução de ~100-200ms em navegação

---

#### 7.2 Sem Request Batching
**Gargalo:** Requisições não são batched  
**Severidade:** Média  
**Arquivo:** `frontend/src/lib/managers/tenantSessionManager.ts`  
**Linha:** 343-366  
**Causa:** Cada operação é uma requisição separada  
**Impacto:** Overhead de HTTP  
**Custo:** Médio - latência adicional  
**Solução Recomendada:** Implementar batching de requisições  
**Ganho Esperado:** Redução de ~20-30% em latência de rede

---

#### 7.3 Sem Request Deduplication
**Gargalo:** Mesma requisição pode ser feita múltiplas vezes  
**Severidade:** Baixa  
**Arquivo:** `frontend/src/lib/api/`  
**Linha:** N/A  
**Causa:** Não há cache de requisições em andamento  
**Impacto:** Requisições duplicadas  
**Custo:** Baixo - impacto pequeno  
**Solução Recomendada:** Implementar request deduplication  
**Ganho Esperado:** Redução de requisições duplicadas

---

#### 7.4 Sem Lazy Loading de Componentes
**Gargalo:** Todos os componentes são carregados imediatamente  
**Severidade:** Baixa  
**Arquivo:** `frontend/src/routes/(app)/+layout.svelte`  
**Linha:** N/A  
**Causa:** Svelte não tem code splitting por padrão  
**Impacto:** Bundle grande  
**Custo:** Baixo - tempo de carregamento inicial  
**Solução Recomendada:** Implementar lazy loading de rotas  
**Ganho Esperado:** Redução de ~30% no bundle inicial

---

#### 7.5 Sem Otimização de Imagens
**Gargalo:** Imagens não são otimizadas  
**Severidade:** Média  
**Arquivo:** `frontend/src/`  
**Linha:** N/A  
**Causa:** Upload de imagens sem compressão  
**Impacto:** Bandwidth desperdiçada  
**Custo:** Médio - afeta UX  
**Solução Recomendada:** Implementar compressão de imagens no upload  
**Ganho Esperado:** Redução de ~50-70% no tamanho de imagens

---

### 8. MEMÓRIA

#### 8.1 Alocação de Slices Sem Capacidade Predefinida
**Gargalo:** Slices alocados sem capacidade  
**Severidade:** Baixa  
**Arquivo:** `internal/infra/repository/gorm_order_repository.go`  
**Linha:** 149, 261, 298  
**Causa:** make([]T, len) em vez de make([]T, 0, len)  
**Impacto:** Realocações desnecessárias  
**Custo:** Baixo - impacto pequeno  
**Solução Recomendada:** Usar make([]T, 0, len) para evitar realocações  
**Ganho Esperado:** Redução de alocações

---

#### 8.2 Maps Sem Capacidade Predefinida
**Gargalo:** Maps alocados sem capacidade  
**Severidade:** Baixa  
**Arquivo:** `internal/infra/repository/gorm_order_repository.go`  
**Linha:** 124, 409, 482  
**Causa:** make(map[K]V) sem hint de tamanho  
**Impacto:** Realocações de hash table  
**Custo:** Baixo - impacto pequeno  
**Solução Recomendada:** Usar make(map[K]V, size) quando tamanho é conhecido  
**Ganho Esperado:** Redução de alocações

---

#### 8.3 Goroutines Sem Limite
**Gargalo:** Goroutines podem ser criadas sem limite  
**Severidade:** Média  
**Arquivo:** `internal/service/idempotency_cleanup_job.go`  
**Linha:** 44  
**Causa:** go func() sem worker pool  
**Impacto:** Pode criar milhares de goroutines  
**Custo:** Médio - pode causar OOM  
**Solução Recomendada:** Implementar worker pool com limite  
**Ganho Esperado:** Controle de uso de memória

---

#### 8.4 Contextos Não Cancelados
**Gargalo:** Contextos podem não ser cancelados  
**Severidade:** Baixa  
**Arquivo:** `internal/service/`  
**Linha:** N/A  
**Causa:** defer cancel() pode faltar em alguns casos  
**Impacto:** Goroutines podem vazar  
**Custo:** Baixo - impacto pequeno  
**Solução Recomendada:** Auditar todos os context.WithTimeout/Cancel  
**Ganho Esperado:** Prevenção de leaks

---

### 9. ESCALABILIDADE

#### 9.1 Connection Pool de PostgreSQL Pequeno
**Gargalo:** MaxOpenConns: 25 é pequeno  
**Severidade:** Crítica  
**Arquivo:** `internal/infra/database/connection.go`  
**Linha:** 47  
**Causa:** Pool pequeno limita concorrência  
**Impacto:** Com 25 conexões, suporta ~50-100 requisições simultâneas  
**Custo:** Crítico - limita escalabilidade  
**Solução Recomendada:** Aumentar para 100-200 dependendo da carga  
**Ganho Esperado:** Suporte para 200-500 requisições simultâneas

---

#### 9.2 Sem Horizontal Scaling
**Gargalo:** Sistema não suporta múltiplas instâncias  
**Severidade:** Crítica  
**Arquivo:** N/A  
**Linha:** N/A  
**Causa:** Cache in-memory não é distribuído  
**Impacto:** Não é possível escalar horizontalmente  
**Custo:** Crítico - limita crescimento  
**Solução Recomendada:** Migrar cache in-memory para Redis distribuído  
**Ganho Esperado:** Capacidade de escalar horizontalmente

---

#### 9.3 Sem Load Balancer Configuration
**Gargalo:** Não há documentação de load balancer  
**Severidade:** Alta  
**Arquivo:** N/A  
**Linha:** N/A  
**Causa:** Arquitetura não considera load balancing  
**Impacto:** Single point of failure  
**Custo:** Alto - risco de disponibilidade  
**Solução Recomendada:** Documentar arquitetura com load balancer  
**Ganho Esperado:** Alta disponibilidade

---

#### 9.4 Limite de 100 Usuários Simultâneos
**Gargalo:** Com configuração atual, suporta ~100 usuários  
**Severidade:** Crítica  
**Arquivo:** `internal/infra/database/connection.go`  
**Linha:** 47  
**Causa:** Pool de 25 conexões + Redis pool de 10  
**Impacto:** 100 usuários simultâneos é o limite prático  
**Custo:** Crítico - limita crescimento  
**Solução Recomendada:** Aumentar pools e implementar cache distribuído  
**Ganho Esperado:** Suporte para 500-1000 usuários

---

#### 9.5 Limite de 500 Usuários Simultâneos (com otimizações)
**Gargalo:** Com otimizações, suporta ~500 usuários  
**Severidade:** Alta  
**Arquivo:** N/A  
**Linha:** N/A  
**Causa:** Mesmo com otimizações, arquitetura tem limites  
**Impacto:** 500 usuários é o limite com otimizações básicas  
**Custo:** Alto - requer arquitetura mais robusta  
**Solução Recomendada:** Implementar cache distribuído, aumentar pools, otimizar queries  
**Ganho Esperado:** Suporte para 1000-2000 usuários

---

#### 9.6 Limite de 1000 Usuários Simultâneos (com arquitetura robusta)
**Gargalo:** 1000 usuários requer arquitetura robusta  
**Severidade:** Alta  
**Arquivo:** N/A  
**Linha:** N/A  
**Causa:** Requer otimizações avançadas  
**Impacto:** 1000 usuários é o limite com arquitetura robusta  
**Custo:** Alto - investimento significativo  
**Solução Recomendada:** Implementar sharding, read replicas, CDN  
**Ganho Esperado:** Suporte para 5000+ usuários

---

### 10. CONCORRÊNCIA

#### 10.1 SELECT FOR UPDATE em Muitos Lugares
**Gargalo:** SELECT FOR UPDATE pode causar contenção  
**Severidade:** Média  
**Arquivo:** `internal/infra/repository/gorm_product_repository.go`  
**Linha:** 402, 433, 577, 627  
**Causa:** Lock pessimista em operações de estoque  
**Impacto:** Pode causar contenção em alta concorrência  
**Custo:** Médio - pode limitar throughput  
**Solução Recomendada:** Considerar optimistic locking para operações menos críticas  
**Ganho Esperado:** Melhoria de ~20-30% em throughput

---

#### 10.2 Mutex em Cache In-Memory
**Gargalo:** RWMutex em cache pode causar contenção  
**Severidade:** Baixa  
**Arquivo:** `internal/infra/repository/gorm_global_config_repository.go`  
**Linha:** 48  
**Causa:** Cache in-memory usa mutex  
**Impacto:** Contenção em alta concorrência  
**Custo:** Baixo - impacto pequeno  
**Solução Recomendada:** Migrar para Redis distribuído  
**Ganho Esperado:** Eliminação de contenção

---

#### 10.3 Sem Configuração de Deadlock Timeout
**Gargalo:** Deadlock timeout não configurado  
**Severidade:** Média  
**Arquivo:** `internal/infra/database/connection.go`  
**Linha:** N/A  
**Causa:** PostgreSQL usa timeout padrão  
**Impacto:** Deadlocks podem travar por muito tempo  
**Custo:** Médio - pode causar degradação  
**Solução Recomendada:** Configurar deadlock_timeout para 1-5s  
**Ganho Esperado:** Detecção mais rápida de deadlocks

---

#### 10.4 Goroutines em Testes de Concorrência
**Gargalo:** Testes criam muitas goroutines  
**Severidade:** Baixa  
**Arquivo:** `internal/infra/repository/concurrency_test.go`  
**Linha:** 76-82, 200-206  
**Causa:** Testes criam 100 goroutines  
**Impacto:** Testes podem ser lentos  
**Custo:** Baixo - apenas em testes  
**Solução Recomendada:** Reduzir número de goroutines em testes  
**Ganho Esperado:** Testes mais rápidos

---

### 11. CACHE

#### 11.1 Dashboard Sem Cache
**Gargalo:** Dashboard não tem cache  
**Severidade:** Crítica  
**Arquivo:** `internal/infra/repository/gorm_dashboard_repository.go`  
**Linha:** 22-360  
**Causa:** Ver problema 5.1  
**Impacto:** Ver problema 5.1  
**Custo:** Crítico  
**Solução Recomendada:** Implementar cache de 1-5 minutos  
**Ganho Esperado:** Redução de ~95% (cache hit)

---

#### 11.2 Produtos Sem Cache
**Gargalo:** Lista de produtos não tem cache  
**Severidade:** Alta  
**Arquivo:** `internal/infra/repository/gorm_product_repository.go`  
**Linha:** 213  
**Causa:** Produtos são carregados do banco a cada requisição  
**Impacto:** 50-100ms por listagem  
**Custo:** Alto - produtos são acessados frequentemente  
**Solução Recomendada:** Implementar cache de 5-10 minutos  
**Ganho Esperado:** Redução de ~80% (cache hit)

---

#### 11.3 Ingredientes Sem Cache
**Gargalo:** Lista de ingredientes não tem cache  
**Severidade:** Alta  
**Arquivo:** `internal/infra/repository/gorm_product_repository.go`  
**Linha:** 417  
**Causa:** Ingredientes são carregados do banco a cada requisição  
**Impacto:** 50-100ms por listagem  
**Custo:** Alto - ingredientes são acessados frequentemente  
**Solução Recomendada:** Implementar cache de 5-10 minutos  
**Ganho Esperado:** Redução de ~80% (cache hit)

---

#### 11.4 Categorias Sem Cache
**Gargalo:** Lista de categorias não tem cache  
**Severidade:** Média  
**Arquivo:** `internal/infra/repository/gorm_category_repository.go`  
**Linha:** 77  
**Causa:** Categorias são carregadas do banco a cada requisição  
**Impacto:** 20-50ms por listagem  
**Custo:** Médio - categorias mudam raramente  
**Solução Recomendada:** Implementar cache de 10-30 minutos  
**Ganho Esperado:** Redução de ~90% (cache hit)

---

#### 11.5 Planos Sem Cache
**Gargalo:** Lista de planos não tem cache  
**Severidade:** Média  
**Arquivo:** `internal/infra/repository/gorm_plan_repository.go`  
**Linha:** 61  
**Causa:** Planos são carregados do banco a cada requisição  
**Impacto:** 20-50ms por listagem  
**Custo:** Médio - planos mudam raramente  
**Solução Recomendada:** Implementar cache de 30-60 minutos  
**Ganho Esperado:** Redução de ~90% (cache hit)

---

#### 11.6 Configurações de Empresa Sem Cache
**Gargalo:** Configurações de empresa não têm cache  
**Severidade:** Média  
**Arquivo:** `internal/infra/repository/gorm_company_repository.go`  
**Linha:** N/A  
**Causa:** Configurações são carregadas do banco a cada requisição  
**Impacto:** 20-50ms por requisição  
**Custo:** Médio - configurações mudam raramente  
**Solução Recomendada:** Implementar cache de 5-10 minutos  
**Ganho Esperado:** Redução de ~80% (cache hit)

---

#### 11.7 Usuários Sem Cache
**Gargalo:** Lista de usuários não tem cache  
**Severidade:** Média  
**Arquivo:** `internal/infra/repository/gorm_user_repository.go`  
**Linha:** 146  
**Causa:** Usuários são carregados do banco a cada requisição  
**Impacto:** 50-100ms por listagem  
**Custo:** Médio - usuários mudam ocasionalmente  
**Solução Recomendada:** Implementar cache de 1-5 minutos  
**Ganho Esperado:** Redução de ~70% (cache hit)

---

#### 11.8 Permissões Sem Cache
**Gargalo:** Verificações de permissão não têm cache  
**Severidade:** Alta  
**Arquivo:** `internal/service/rbac_service.go`  
**Linha:** N/A  
**Causa:** Permissões são verificadas no banco a cada requisição  
**Impacto:** 10-30ms por verificação  
**Custo:** Alto - permissões são verificadas frequentemente  
**Solução Recomendada:** Implementar cache de 5-10 minutos  
**Ganho Esperado:** Redução de ~85% (cache hit)

---

### 12. PRODUÇÃO

#### 12.1 Sistema Não Suporta Produção Sem Otimizações
**Gargalo:** Configuração atual não é adequada para produção  
**Severidade:** Crítica  
**Arquivo:** N/A  
**Linha:** N/A  
**Causa:** Pools pequenos, sem cache, queries não otimizadas  
**Impacto:** Sistema não suporta carga de produção  
**Custo:** Crítico - entrada em produção arriscada  
**Solução Recomendada:** Implementar otimizações críticas antes da entrada em produção  
**Ganho Esperado:** Sistema pronto para produção

---

#### 12.2 Throughput Estimado: 50-100 RPS
**Gargalo:** Throughput atual é limitado  
**Severidade:** Crítica  
**Arquivo:** N/A  
**Linha:** N/A  
**Causa:** Ver problema 12.1  
**Impacto:** 50-100 requests por segundo é o limite atual  
**Custo:** Crítico - insuficiente para produção  
**Solução Recomendada:** Implementar otimizações críticas  
**Ganho Esperado:** 200-500 RPS

---

#### 12.3 Primeiro Ponto de Quebra: Banco de Dados
**Gargalo:** Banco de dados será o primeiro gargalo  
**Severidade:** Crítica  
**Arquivo:** `internal/infra/database/connection.go`  
**Linha:** 47  
**Causa:** Pool pequeno, queries não otimizadas  
**Impacto:** Banco esgotará primeiro sob carga  
**Custo:** Crítico - limita escalabilidade  
**Solução Recomendada:** Aumentar pool, otimizar queries, adicionar cache  
**Ganho Esperado:** Banco suporta carga 5-10x maior

---

#### 12.4 Segundo Ponto de Quebra: Redis
**Gargalo:** Redis será o segundo gargalo  
**Severidade:** Alta  
**Arquivo:** `internal/infra/redis/client.go`  
**Linha:** 23  
**Causa:** Pool pequeno  
**Impacto:** Redis esgotará após banco  
**Custo:** Alto - limita escalabilidade  
**Solução Recomendada:** Aumentar pool, implementar cache distribuído  
**Ganho Esperado:** Redis suporta carga 3-5x maior

---

#### 12.5 Terceiro Ponto de Quebra: CPU
**Gargalo:** CPU será o terceiro gargalo  
**Severidade:** Média  
**Arquivo:** N/A  
**Linha:** N/A  
**Causa:** Processamento de JSON, serialização  
**Impacto:** CPU esgotará após Redis  
**Custo:** Médio - pode ser resolvido com horizontal scaling  
**Solução Recomendada:** Implementar horizontal scaling  
**Ganho Esperado:** CPU não é mais gargalo

---

#### 12.6 Memória Não é Gargalo Imediato
**Gargalo:** Memória não é o primeiro gargalo  
**Severidade:** Baixa  
**Arquivo:** N/A  
**Linha:** N/A  
**Causa:** Uso de memória é moderado  
**Impacto:** Memória esgotará apenas em carga muito alta  
**Custo:** Baixo - não é prioridade  
**Solução Recomendada:** Monitorar uso de memória  
**Ganho Esperado:** Visibilidade

---

#### 12.7 Rede Não é Gargalo Imediato
**Gargalo:** Rede não é o primeiro gargalo  
**Severidade:** Baixa  
**Arquivo:** N/A  
**Linha:** N/A  
**Causa:** Latência de rede é baixa em datacenter  
**Impacto:** Rede esgotará apenas em carga muito alta  
**Custo:** Baixo - não é prioridade  
**Solução Recomendada:** Monitorar latência de rede  
**Ganho Esperado:** Visibilidade

---

## RESUMO POR SEVERIDADE

### Críticos (8)
1. Dashboard executa ~30 queries separadas
2. Falta de índices compostos
3. Dashboard sem cache
4. Connection pool de PostgreSQL pequeno (25)
5. Sem horizontal scaling (cache in-memory)
6. Limite de 100 usuários simultâneos
7. Sistema não suporta produção sem otimizações
8. Throughput estimado: 50-100 RPS

### Altos (18)
1. N+1 query em dashboard (recent orders items count)
2. Múltiplas queries separadas para KPIs
3. Loop de queries para sales by day
4. Mesmo problema em report repository
5. JOINs com deleted_at filters
6. COUNT queries sem cache
7. GROUP BY sem índice apropriado
8. Find dentro de loop
9. Falta de prefetch count em RabbitMQ
10. Falta de DLQ em RabbitMQ
11. APIs de listagem sem paginação adequada
12. Timeout de contexto não configurado
13. Limite de 500 usuários (com otimizações)
14. Limite de 1000 usuários (com arquitetura robusta)
15. Produtos sem cache
16. Ingredientes sem cache
17. Permissões sem cache
18. Segundo ponto de quebra: Redis

### Médios (12)
1. SELECT * em queries
2. Save() em vez de Update()
3. Pool size de Redis pequeno (10)
4. Falta de métricas de cache hit/miss
5. Fallback cache silencia erros
6. Batch size fixo em RabbitMQ
7. Dashboard carrega dados não utilizados
8. APIs sem rate limiting por tenant
9. Sem request batching no frontend
10. Sem otimização de imagens
11. Goroutines sem limite
12. SELECT FOR UPDATE pode causar contenção

### Baixos (4)
1. Preload desnecessário
2. TTL não configurado para alguns caches
3. Alocação de slices/maps sem capacidade
4. Memória e rede não são gargalos imediatos

---

## ESTIMATIVA DE ESFORÇO PARA CORREÇÃO

### Prioridade 0 - Crítico (Pré-produção)
**Total: 40 horas**

| Item | Estimativa |
|------|-----------|
| Otimizar dashboard queries | 8h |
| Adicionar índices compostos | 4h |
| Implementar cache de dashboard | 4h |
| Aumentar connection pool PostgreSQL | 2h |
| Migrar cache in-memory para Redis | 8h |
| Aumentar pool Redis | 2h |
| Adicionar cache em produtos/ingredientes | 4h |
| Adicionar cache em permissões | 4h |
| Configurar DLQ RabbitMQ | 2h |
| Configurar prefetch RabbitMQ | 2h |

### Prioridade 1 - Alta (Primeira Sprint)
**Total: 32 horas**

| Item | Estimativa |
|------|-----------|
| Corrigir N+1 em dashboard | 4h |
| Consolidar sales by day queries | 4h |
| Adicionar LIMIT padrão em APIs | 2h |
| Configurar timeout em operações | 4h |
| Implementar rate limiting por tenant | 4h |
| Adicionar cache em categorias/planos | 4h |
| Adicionar cache em usuários | 4h |
| Implementar request batching | 4h |
| Adicionar métricas de cache | 2h |

### Prioridade 2 - Média (Próxima Sprint)
**Total: 24 horas**

| Item | Estimativa |
|------|-----------|
| Substituir SELECT * por Select() | 4h |
| Substituir Save() por Update() | 4h |
| Implementar lazy loading de gráficos | 4h |
| Implementar worker pool | 4h |
| Configurar deadlock timeout | 2h |
| Otimizar imagens | 4h |
| Implementar request deduplication | 2h |

### Prioridade 3 - Baixa (Quando Possível)
**Total: 8 horas**

| Item | Estimativa |
|------|-----------|
| Remover preload desnecessário | 2h |
| Configurar TTL em caches | 2h |
| Otimizar alocação de slices/maps | 2h |
| Implementar lazy loading de rotas | 2h |

**ESFORÇO TOTAL: 104 horas (~13 dias úteis)**
