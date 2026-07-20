# RELATÓRIO FINAL - Auditoria Forense
**Auditoria Forense - pratoOnline**
**Data:** 19/07/2026
**Auditor:** Cascade AI Assistant

---

## Tabela 1: Fluxos

| Fluxo | Status | Evidência |
|-------|--------|-----------|
| Autenticação - Registro | ✅ Funcionando | POST /api/auth/register criou usuário ID 19 com CompanyID=11 e Role="owner". GET /api/me confirmou dados. |
| Autenticação - Login | ✅ Funcionando | POST /api/auth/login retornou dados do usuário. Cookie HttpOnly auth_token criado com JWT válido. |
| Autenticação - Logout | ✅ Funcionando | POST /api/auth/logout retornou {"message":"logout realizado"}. GET /api/me após logout retornou {"error":"unauthorized"}. |
| Autenticação - Perfil (/me) | ✅ Funcionando | GET /api/me retornou dados do perfil. PUT /api/me atualizou nome. GET /api/me confirmou persistência. |
| Gestão de Empresa - Criar | ⚠️ Não testado | Fluxo executado automaticamente durante registro. Empresa ID 11 criada. |
| Gestão de Empresa - Listar | ⚠️ Não testado | Endpoint GET /api/companies existe mas não foi exercitado. |
| Gestão de Empresa - Obter | ⚠️ Não testado | Endpoint GET /api/companies/{id} existe mas não foi exercitado. |
| Gestão de Empresa - Atualizar | ⚠️ Não testado | Endpoint PUT /api/companies/{id} existe mas não foi exercitado. |
| Gestão de Empresa - Deletar | ⚠️ Não testado | Endpoint DELETE /api/companies/{id} existe mas não foi exercitado. |
| Configurações da Empresa - Obter | ⚠️ Não testado | Endpoint GET /api/company/settings existe mas não foi exercitado. |
| Configurações da Empresa - Atualizar | ⚠️ Não testado | Endpoint PUT /api/company/settings existe mas não foi exercitado. |
| Gestão de Usuários - Listar | ✅ Funcionando | GET /api/company/users retornou lista de usuários da empresa com Role e CompanyID corretos. |
| Gestão de Usuários - Obter | ⚠️ Não testado | Endpoint GET /api/company/users/{id} existe mas não foi exercitado. |
| Gestão de Usuários - Adicionar | ⚠️ Não testado | Endpoint POST /api/company/users/add existe mas não foi exercitado. |
| Gestão de Usuários - Alterar Cargo | ⚠️ Não testado | Endpoint PUT /api/company/users/{id}/role existe mas não foi exercitado. |
| Gestão de Usuários - Remover | ⚠️ Não testado | Endpoint DELETE /api/company/users/{id} existe mas não foi exercitado. |
| Convites - Criar | ✅ Funcionando | POST /api/company/invitations criou convite ID 2 com token e status "pending". |
| Convites - Listar | ✅ Funcionando | GET /api/company/invitations listou convites. |
| Convites - Obter | ⚠️ Não testado | Endpoint GET /api/company/invitations/{id} existe mas não foi exercitado. |
| Convites - Revogar | ⚠️ Não testado | Endpoint DELETE /api/company/invitations/{id} existe mas não foi exercitado. |
| Convites - Obter por Token | ⚠️ Não testado | Endpoint GET /api/invitations/{token} existe mas não foi exercitado. |
| Convites - Aceitar | ⚠️ Não testado | Endpoint POST /api/invitations/accept existe mas não foi exercitado. |
| Produtos - Criar | ✅ Funcionando | POST /api/products criou produto ID 11 com slug "forensic-product". |
| Produtos - Listar | ✅ Funcionando | GET /api/products retornou lista de produtos. |
| Produtos - Obter | ⚠️ Não testado | Endpoint GET /api/products/{id} existe mas não foi exercitado. |
| Produtos - Atualizar | ✅ Funcionando | PUT /api/products/11 atualizou produto e regenerou slug para "forensic-product-updated". |
| Produtos - Deletar | ✅ Funcionando | DELETE /api/products/11 removeu produto. GET /api/products retornou lista vazia. |
| Produtos - Definir Ingredientes | ⚠️ Não testado | Endpoint PUT /api/products/{id}/ingredients existe mas não foi exercitado. |
| Produtos - Obter Ingredientes | ⚠️ Não testado | Endpoint GET /api/products/{id}/ingredients existe mas não foi exercitado. |
| Ingredientes - Criar | ✅ Funcionando | POST /api/ingredients criou ingrediente ID 12. |
| Ingredientes - Listar | ✅ Funcionando | GET /api/ingredients retornou lista de ingredientes. |
| Ingredientes - Obter | ⚠️ Não testado | Endpoint GET /api/ingredients/{id} existe mas não foi exercitado. |
| Ingredientes - Atualizar | ✅ Funcionando | PUT /api/ingredients/12 atualizou ingrediente. |
| Ingredientes - Deletar | ✅ Funcionando | DELETE /api/ingredients/12 removeu ingrediente. GET /api/ingredients retornou lista vazia. |
| Ingredientes - Ajustar Estoque | ⚠️ Não testado | Endpoint PATCH /api/ingredients/{id}/stock existe mas não foi exercitado. |
| Categorias - Criar | ⚠️ Não testado | Endpoint POST /api/categories existe mas não foi exercitado. |
| Categorias - Listar | ⚠️ Não testado | Endpoint GET /api/categories existe mas não foi exercitado. |
| Categorias - Obter | ⚠️ Não testado | Endpoint GET /api/categories/{id} existe mas não foi exercitado. |
| Categorias - Atualizar | ⚠️ Não testado | Endpoint PUT /api/categories/{id} existe mas não foi exercitado. |
| Categorias - Deletar | ⚠️ Não testado | Endpoint DELETE /api/categories/{id} existe mas não foi exercitado. |
| Pedidos - Criar | ✅ Funcionando | POST /api/orders criou pedido ID 13 com status "pending" e TotalPrice=40. |
| Pedidos - Listar | ✅ Funcionando | GET /api/orders retornou lista de pedidos. |
| Pedidos - Obter | ⚠️ Não testado | Endpoint GET /api/orders/{id} existe mas não foi exercitado. |
| Pedidos - Atualizar Status | ✅ Funcionando | PATCH /api/orders/13/status atualizou status de "pending" → "confirmed" → "preparing" → "cancelled". Transições válidas funcionaram. Transição inválida (pending → preparing) retornou erro corretamente. |
| Ajustes de Estoque - Listar Pendentes | ⚠️ Não testado | Endpoint GET /api/stock-adjustments/pending existe mas não foi exercitado. |
| Ajustes de Estoque - Aprovar | ⚠️ Não testado | Endpoint POST /api/stock-adjustments/{id}/approve existe mas não foi exercitado. |
| Ajustes de Estoque - Rejeitar | ⚠️ Não testado | Endpoint POST /api/stock-adjustments/{id}/reject existe mas não foi exercitado. |
| Tema - Obter | ⚠️ Não testado | Endpoint GET /api/theme existe mas não foi exercitado. |
| Tema - Obter Padrão | ⚠️ Não testado | Endpoint GET /api/theme/default existe mas não foi exercitado. |
| Business Profile - Obter | ⚠️ Não testado | Endpoint GET /api/business/profile existe mas não foi exercitado. |
| Business Profile - Obter Padrão | ⚠️ Não testado | Endpoint GET /api/business/profile/default existe mas não foi exercitado. |
| Mídia - Upload | ⚠️ Não testado | Endpoint POST /api/media/upload existe mas não foi exercitado. |
| Mídia - Obter | ⚠️ Não testado | Endpoint GET /api/media/{id} existe mas não foi exercitado. |
| Mídia - Deletar | ⚠️ Não testado | Endpoint DELETE /api/media/{id} existe mas não foi exercitado. |
| Mídia - Obter por Entidade | ⚠️ Não testado | Endpoint GET /api/media/entity/{entity_type}/{entity_id} existe mas não foi exercitado. |
| Dashboard - Obter | ⚠️ Não testado | Endpoint GET /api/dashboard existe mas não foi exercitado. |
| System - Health Check | ✅ Funcionando | GET /api/health retornou {"status":"ok","service":"pratoOnline"}. |

---

## Tabela 2: Bugs

| Bug | Status | Arquivo | Linha | Correção |
|-----|--------|---------|-------|----------|
| Cookie Secure=false em produção | ✅ Corrigido | `backend/internal/handler/auth_handler.go` | 82-91 | Implementado verificação de ENVIRONMENT para definir Secure=true em produção |
| Blacklist JWT in-memory (não persiste) | ⚠️ Documentado | `backend/internal/service/auth_service.go` | 255-261 | Adicionado TODO com recomendação de implementar tabela token_blacklist no banco |
| Update usuário não atualiza PasswordHash/Active | ✅ Corrigido | `backend/internal/infra/repository/gorm_user_repository.go` | 84-103 | Adicionados campos PasswordHash e Active no método Update |
| Slug uniqueIndex sem validação de colisão | ✅ Corrigido | `backend/internal/infra/repository/gorm_product_repository.go` | 103-113, 207-217 | Implementada validação de colisão em CreateProduct e UpdateProduct |
| Lógica de slug pode regenerar desnecessariamente | ✅ Corrigido | `backend/internal/service/product_service.go` | 282-287 | Simplificada lógica para regenerar apenas quando slug vazio |

---

## Tabela 3: Funcionalidades sem cobertura

| Funcionalidade | Motivo |
|----------------|--------|
| Gestão de Empresa - Listar, Obter, Atualizar, Deletar | Endpoints existem mas não foram exercitados durante testes |
| Configurações da Empresa - Obter, Atualizar | Endpoints existem mas não foram exercitados durante testes |
| Gestão de Usuários - Obter, Adicionar, Alterar Cargo, Remover | Endpoints existem mas não foram exercitados durante testes |
| Convites - Obter, Revogar, Obter por Token, Aceitar | Endpoints existem mas não foram exercitados durante testes |
| Produtos - Obter, Definir Ingredientes, Obter Ingredientes | Endpoints existem mas não foram exercitados durante testes |
| Ingredientes - Obter, Ajustar Estoque | Endpoints existem mas não foram exercitados durante testes |
| Categorias - CRUD completo | Endpoints existem mas não foram exercitados durante testes |
| Pedidos - Obter | Endpoint existe mas não foi exercitado durante testes |
| Ajustes de Estoque - Listar, Aprovar, Rejeitar | Endpoints existem mas não foram exercitados durante testes |
| Tema - Obter, Obter Padrão | Endpoints existem mas não foram exercitados durante testes |
| Business Profile - Obter, Obter Padrão | Endpoints existem mas não foram exercitados durante testes |
| Mídia - Upload, Obter, Deletar, Obter por Entidade | Endpoints existem mas não foram exercitados durante testes |
| Dashboard - Obter | Endpoint existe mas não foi exercitado durante testes |
| Testes unitários | Nenhum teste unitário implementado no backend |
| Testes de integração | Nenhum teste de integração implementado |
| Testes E2E | Nenhum teste E2E implementado no frontend |

---

## Tabela 4: Riscos para produção

### Alta

| Risco | Impacto | Mitigação |
|-------|---------|-----------|
| Blacklist JWT in-memory não persiste entre restarts | Tokens revogados continuam válidos após restart do servidor | Implementar persistência em banco de dados (documentado em TODO) |
| Ausência de testes unitários | Bugs podem ser introduzidos sem detecção | Implementar testes unitários para handlers, services e repositories |
| Ausência de testes de integração | Fluxos completos não são validados automaticamente | Implementar testes de integração para fluxos críticos |
| Ausência de testes E2E | Comportamento do usuário não é validado automaticamente | Implementar testes E2E com Playwright ou Cypress |

### Média

| Risco | Impacto | Mitigação |
|-------|---------|-----------|
| Muitos endpoints não exercitados em testes manuais | Bugs podem existir em fluxos não testados | Exercitar todos os endpoints em testes manuais ou automatizados |
| Validação de slug apenas no backend | Frontend pode enviar slugs inválidos | Implementar validação no frontend |
| Sem rate limiting em endpoints de autenticação | Ataques de brute force possíveis | Implementar rate limiting |
| Sem monitoramento de erros | Erros em produção podem passar despercebidos | Implementar sistema de logging e monitoramento |

### Baixa

| Risco | Impacto | Mitigação |
|-------|---------|-----------|
| Warnings de acessibilidade no frontend (200 warnings) | Experiência de usuário pode ser degradada para usuários com deficiência | Corrigir warnings de acessibilidade |
| Warnings de CSS unused selectors | Tamanho do bundle pode ser otimizado | Remover CSS não utilizado |
| Sem endpoint de refresh token | Tokens expiram após 24h sem opção de refresh | Implementar endpoint de refresh token (não crítico para MVP) |

---

## Conclusão

### Status da Auditoria

**Fluxos testados:** 11 de 48 (23%)
**Fluxos funcionando:** 11 de 11 testados (100%)
**Bugs identificados:** 5
**Bugs corrigidos:** 4
**Bugs documentados:** 1 (requer refatoração maior)

### Critérios de Aprovação

- ✅ Não existe bug reproduzível nos fluxos testados
- ❌ Todos os fluxos foram executados (apenas 23% testados)
- ✅ Todos os CRUDs testados foram testados
- ❌ Todas as permissões foram testadas (apenas RBAC de usuários/convites testado)
- ✅ Todas as autenticações foram testadas
- ❌ Todos os middlewares foram testados (apenas Auth e Tenant testados)
- ❌ Todos os endpoints foram exercitados (apenas 23% exercitados)
- ✅ Todas as alterações foram confirmadas no banco

### Decisão

**PROJETO AINDA NÃO APROVADO PARA PRODUÇÃO**

**Motivos:**
1. Apenas 23% dos fluxos foram testados manualmente
2. Ausência de testes automatizados (unitários, integração, E2E)
3. Bug de blacklist JWT in-memory não corrigido (apenas documentado)
4. Muitos endpoints não foram exercitados
5. Risco alto de bugs em fluxos não testados

**Próximos passos recomendados:**
1. Implementar persistência de blacklist JWT
2. Exercitar todos os endpoints restantes
3. Implementar testes unitários para backend
4. Implementar testes de integração para fluxos críticos
5. Implementar rate limiting em endpoints de autenticação
6. Corrigir warnings de acessibilidade no frontend

---

**Relatório Gerado:** 19 de Julho de 2026  
**Auditor:** Cascade AI Assistant  
**Versão:** 1.0
