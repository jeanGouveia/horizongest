# AUDITORIA PLATAFORMA 2.0 - RELATÓRIO FINAL

**Data**: 18/07/2026  
**Auditor**: Arquiteto Principal  
**Escopo**: Sprints 1-8 (Tenant Engine, White Label, Business Engine, Tenant Isolation, Company Settings, RBAC, User Management, Invites & Onboarding)

---

## RESUMO EXECUTIVO

A Plataforma PratoOnline 2.0 apresenta uma arquitetura sólida com separação de camadas (Domain, Ports, Infrastructure, Service, Handler), mas contém **7 problemas críticos**, **12 problemas altos**, **8 problemas médios** e **6 problemas baixos** que devem ser corrigidos antes de ir para produção.

**Principais Riscos**:
- Violação de RBAC em endpoints de gerenciamento de usuários
- Falta de isolamento de tenant em endpoints de empresas
- Vulnerabilidade de segurança em cookies de autenticação
- Inconsistência entre middleware e serviço de RBAC
- Duplicação de código em handlers
- Problemas de performance em queries de usuários

---

## PROBLEMAS CRÍTICOS

### 1. Violação de RBAC - Inconsistência entre Middleware e Serviço

**Severidade**: Crítico  
**Arquivo**: `backend/internal/service/rbac_service.go`  
**Linha**: 92  
**Explicação Técnica**: O método `CanManageUsers()` retorna true apenas para Owner, mas no `main.go` (linhas 132 e 144), o RoleMiddleware está configurado com `RequireAny(RoleOwner, RoleAdmin)`. Isso cria uma inconsistência onde o middleware permite Admin acessar endpoints de gerenciamento de usuários, mas o serviço bloqueia.

**Impacto**: Admins podem acessar endpoints mas receberão erro 403 do serviço, causando confusão e UX ruim. Mais grave: se a validação do serviço for removida ou falhar, Admins poderão realizar ações não autorizadas.

**Melhor Solução**: Atualizar `CanManageUsers()` para permitir Owner e Admin, ou atualizar o middleware para permitir apenas Owner. Recomenda-se permitir Admin para gerenciamento de usuários (exceto alterar Owner).

```go
// rbac_service.go linha 92
func (s *RBACService) CanManageUsers(ctx context.Context, userID uint) (bool, error) {
    return s.HasAnyRole(ctx, userID, domain.RoleOwner, domain.RoleAdmin)
}
```

---

### 2. Falta de Isolamento de Tenant em Endpoints de Empresas

**Severidade**: Crítico  
**Arquivo**: `backend/cmd/server/main.go`  
**Linha**: 119-123  
**Explicação Técnica**: Os endpoints de empresas (`/api/companies`) estão dentro do grupo protegido por AuthMiddleware e TenantMiddleware, mas os handlers não aplicam filtragem de tenant. Qualquer usuário autenticado pode listar todas as empresas, ver empresas de outros tenants, e potencialmente modificar empresas de outros tenants.

**Impacto**: Violação completa de isolamento de tenant. Usuários podem acessar dados de empresas de outros tenants, causando vazamento de dados sensíveis.

**Melhor Solução**: Mover endpoints de empresas para fora do grupo de tenant, ou implementar lógica específica para filtrar empresas por tenant. Para multi-tenant, usuários só devem ver sua própria empresa.

```go
// main.go - remover empresas do grupo de tenant ou adicionar filtragem
// Opção 1: Mover para fora do grupo de tenant (se empresas devem ser públicas)
// Opção 2: Adicionar filtragem no handler para mostrar apenas empresa do usuário
```

---

### 3. Vulnerabilidade de Segurança em Cookie de Autenticação

**Severidade**: Crítico  
**Arquivo**: `backend/internal/handler/auth_handler.go`  
**Linha**: 87  
**Explicação Técnica**: O cookie de autenticação está configurado com `Secure: false`, o que significa que o cookie será enviado em conexões HTTP não criptografadas. Isso expõe o token JWT a ataques de interceptação (Man-in-the-Middle).

**Impacto**: Em ambientes não-HTTPS, o token pode ser interceptado e usado para sequestro de sessão. Em produção, isso é uma vulnerabilidade crítica de segurança.

**Melhor Solução**: Configurar `Secure: true` em produção, usar variável de ambiente para controlar isso.

```go
// auth_handler.go linha 87
secureFlag := os.Getenv("COOKIE_SECURE") == "true" // true em produção
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    result.Token,
    Path:     "/",
    HttpOnly: true,
    Secure:   secureFlag, // true em produção
    SameSite: http.SameSiteStrictMode,
    Expires:  time.Now().Add(24 * time.Hour),
})
```

---

### 4. Comentário em Chinês no Código

**Severidade**: Crítico  
**Arquivo**: `backend/internal/service/invitation_service.go`  
**Linha**: 23  
**Explicação Técnica**: Existe um comentário em chinês ("邀请") no código, o que viola as convenções de código e pode causar problemas de manutenção e compreensão.

**Impacto**: Dificuldade de manutenção, violação de padrões de código, possível problema com encoding em alguns editores.

**Melhor Solução**: Remover ou traduzir o comentário para português ou inglês.

```go
// invitation_service.go linha 23
// DefaultInvitationExpiration is the default expiration time for invitations (7 days)
```

---

### 5. Race Condition em Aceitação de Convite

**Severidade**: Crítico  
**Arquivo**: `backend/internal/service/invitation_service.go`  
**Linha**: 216-272  
**Explicação Técnica**: O método `AcceptInvitation()` não usa transação de banco de dados. Se dois usuários tentarem aceitar o mesmo convite simultaneamente, ou se o usuário for associado a uma empresa entre a verificação e a atualização, pode ocorrer estado inconsistente.

**Impacto**: Condição de corrida pode levar a dados inconsistentes, convites aceitos múltiplas vezes, ou usuários associados a empresas incorretamente.

**Melhor Solução**: Envolver a lógica de aceitação em uma transação de banco de dados com bloqueio otimista.

```go
// invitation_service.go - envolver em transação
func (s *InvitationService) AcceptInvitation(ctx context.Context, token string) error {
    return s.invitationRepo.Transaction(ctx, func(tx *ports.InvitationRepository) error {
        // ... lógica de aceitação dentro da transação
    })
}
```

---

### 6. Falta de Validação de Role em Criação de Convite

**Severidade**: Crítico  
**Arquivo**: `backend/internal/handler/invitation_handler.go`  
**Linha**: 75-79  
**Explicação Técnica**: O handler valida se o role é válido usando `domain.ParseRole()`, mas não verifica se o role pode ser atribuído via convite. Por exemplo, não há restrição para não permitir criar convites com role Owner (isso pode causar múltiplos Owners).

**Impacto**: Usuários podem criar convites com role Owner, resultando em múltiplos Owners na mesma empresa, o que pode causar problemas de governança e segurança.

**Melhor Solução**: Adicionar validação para não permitir criar convites com role Owner. Convites devem ser limitados a Admin, Manager, Cashier, Kitchen, Waiter.

```go
// invitation_handler.go - adicionar validação
if role == domain.RoleOwner {
    jsonError(w, "não é possível criar convite com cargo Owner", http.StatusBadRequest)
    return
}
```

---

### 7. Falta de Rate Limiting em Endpoints Públicos

**Severidade**: Crítico  
**Arquivo**: `backend/cmd/server/main.go`  
**Linha**: 153-154  
**Explicação Técnica**: Os endpoints públicos de convite (`GET /api/invitations/{token}` e `POST /api/invitations/accept`) não têm rate limiting. Atacantes podem fazer brute force em tokens ou tentar aceitar convites repetidamente.

**Impacto**: Ataques de força bruta em tokens, abuso de API, possível DoS.

**Melhor Solução**: Implementar rate limiting em endpoints públicos, especialmente em endpoints de aceitação de convite.

```go
// main.go - adicionar middleware de rate limiting
r.Get("/api/invitations/{token}", rateLimiter.Limit(10, time.Minute), invitationHandler.GetInvitationByToken)
r.Post("/api/invitations/accept", rateLimiter.Limit(5, time.Minute), invitationHandler.AcceptInvitation)
```

---

## PROBLEMAS ALTOS

### 8. Duplicação de Código - getUserCompanyID

**Severidade**: Alto  
**Arquivo**: `backend/internal/handler/user_management_handler.go`, `backend/internal/handler/invitation_handler.go`  
**Linha**: 30-40 (user_management), 30-40 (invitation)  
**Explicação Técnica**: O método `getUserCompanyID()` está duplicado em ambos os handlers. Isso viola o princípio DRY (Don't Repeat Yourself) e torna a manutenção mais difícil.

**Impacto**: Se a lógica precisar ser alterada, precisa ser alterada em múltiplos lugares. Maior chance de bugs e inconsistências.

**Melhor Solução**: Criar um helper compartilhado no pacote middleware ou em um pacote de utilitários.

```go
// middleware/helper.go
func GetUserCompanyID(ctx context.Context, userRepo ports.UserRepository) (uint, error) {
    userID, ok := GetUserIDFromContext(ctx)
    if !ok {
        return 0, errors.New("não autorizado")
    }
    user, err := userRepo.FindByID(ctx, userID)
    if err != nil {
        return 0, err
    }
    if user == nil || user.CompanyID == nil {
        return 0, errors.New("usuário não possui empresa")
    }
    return *user.CompanyID, nil
}
```

---

### 9. Falta de Validação de Deleção de Empresa

**Severidade**: Alto  
**Arquivo**: `backend/internal/handler/company_handler.go`  
**Linha**: 104-119  
**Explicação Técnica**: O endpoint de deleção de empresa não verifica se a empresa tem usuários, produtos, pedidos ou outros dados associados. Deletar uma empresa pode deixar dados órfãos ou causar inconsistências.

**Impacto**: Deleção de empresa com dados associados pode causar dados órfãos, violação de integridade referencial, e problemas de negócio.

**Melhor Solução**: Implementar verificação de dependências antes de deletar empresa, similar ao `CanDeleteProduct`.

```go
// company_service.go - adicionar verificação
func (s *CompanyService) CanDeleteCompany(ctx context.Context, id uint) (*domain.DependencyCheck, error) {
    // Verificar usuários, produtos, pedidos, etc.
}
```

---

### 10. Inconsistência em ListUsers - Performance

**Severidade**: Alto  
**Arquivo**: `backend/internal/service/user_management_service.go`  
**Linha**: 47-75  
**Explicação Técnica**: O método `ListUsers()` busca todos os usuários do banco (`s.userRepo.List(ctx)`) e depois filtra em memória por companyID. Isso é extremamente ineficiente para bases de dados grandes.

**Impacto**: Performance degradada conforme número de usuários cresce. Carrega dados desnecessários do banco. Escalabilidade prejudicada.

**Melhor Solução**: Adicionar método no repository para buscar usuários por companyID diretamente no banco.

```go
// ports/user_repository.go
FindByCompanyID(ctx context.Context, companyID uint) ([]*domain.User, error)

// user_management_service.go
func (s *UserManagementService) ListUsers(ctx context.Context, companyID uint) ([]UserOutput, error) {
    users, err := s.userRepo.FindByCompanyID(ctx, companyID)
    // ...
}
```

---

### 11. Falta de Validação de Slug em UpdateCompany

**Severidade**: Alto  
**Arquivo**: `backend/internal/service/company_service.go`  
**Linha**: 145-164  
**Explicação Técnica**: O método `UpdateCompany()` verifica se o novo slug conflita com outra empresa, mas não valida se o slug é válido (caracteres especiais, comprimento, etc.).

**Impacto**: Slugs inválidos podem causar problemas em URLs, rotas, e integrações externas.

**Melhor Solução**: Adicionar validação de slug (formato, caracteres, comprimento) antes de atualizar.

```go
// company_service.go - adicionar validação
if !isValidSlug(slug) {
    return nil, errors.New("slug inválido")
}
```

---

### 12. Falta de Índice Composto em invitations

**Severidade**: Alto  
**Arquivo**: `backend/migrations/00012_create_invitations.sql`  
**Linha**: 1-21  
**Explicação Técnica**: A tabela de invitations tem índices individuais em company_id, email, token, status, expires_at, mas não tem índice composto para queries comuns como (company_id, status) ou (email, company_id, status).

**Impacto**: Queries comuns podem ser lentas conforme número de convites cresce. Performance degradada.

**Melhor Solução**: Adicionar índices compostos para queries comuns.

```sql
CREATE INDEX idx_invitations_company_status ON invitations(company_id, status);
CREATE INDEX idx_invitations_email_company_status ON invitations(email, company_id, status);
```

---

### 13. Falta de Validação de Email em Invitation

**Severidade**: Alto  
**Arquivo**: `backend/internal/handler/invitation_handler.go`  
**Linha**: 65-68  
**Explicação Técnica**: O handler valida se o email não está vazio, mas não valida se é um email válido (formato, domínio, etc.).

**Impacto**: Emails inválidos podem ser criados no sistema, causando problemas de entrega e usabilidade.

**Melhor Solução**: Adicionar validação de email usando regex ou biblioteca de validação.

```go
// invitation_handler.go - adicionar validação
if !isValidEmail(input.Email) {
    jsonError(w, "e-mail inválido", http.StatusBadRequest)
    return
}
```

---

### 14. Falta de Sanitização em Company Settings

**Severidade**: Alto  
**Arquivo**: `backend/internal/service/company_settings_service.go`  
**Linha**: 100-161  
**Explicação Técnica**: O método `UpdateSettings()` não sanitiza os dados de entrada (HTML, scripts, etc.) antes de salvar no banco. Isso pode causar vulnerabilidades de XSS se os dados forem exibidos sem sanitização.

**Impacto**: Vulnerabilidade de XSS se os dados forem exibidos sem sanitização. Injeção de código malicioso.

**Melhor Solução**: Sanitizar dados de entrada antes de salvar, ou sanitizar na saída.

```go
// company_settings_service.go - adicionar sanitização
if input.Name != nil {
    company.Name = sanitizeHTML(*input.Name)
}
```

---

### 15. Falta de Log de Auditoria

**Severidade**: Alto  
**Arquivo**: Múltiplos arquivos (todos os services)  
**Linha**: N/A  
**Explicação Técnica**: Não há log de auditoria para ações críticas como criação de usuários, mudança de roles, deleção de empresas, etc. Não é possível rastrear quem fez o quê.

**Impacto**: Impossibilidade de auditoria, conformidade, e investigação de incidentes. Dificuldade de debug.

**Melhor Solução**: Implementar sistema de log de auditoria para ações críticas.

```go
// Criar AuditService
auditService.CreateAuditLog(ctx, userID, "user.created", targetUserID, details)
```

---

### 16. Falta de Validação de Expiração em GetInvitationByToken

**Severidade**: Alto  
**Arquivo**: `backend/internal/service/invitation_service.go`  
**Linha**: 194-214  
**Explicação Técnica**: O método `GetInvitationByToken()` verifica expiração e atualiza status para expired, mas ignora erro se a atualização falhar (linha 209: `_ = s.invitationRepo.Update(ctx, invitation)`). Isso pode deixar convites expirados com status pending.

**Impacto**: Convites expirados podem permanecer com status pending, causando confusão e comportamento inconsistente.

**Melhor Solução**: Tratar erro de atualização ou usar transação.

```go
// invitation_service.go - tratar erro
if err := s.invitationRepo.Update(ctx, invitation); err != nil {
    return nil, fmt.Errorf("GetInvitationByToken: failed to expire invitation: %w", err)
}
```

---

### 17. Falta de Validação de Role em ChangeRole

**Severidade**: Alto  
**Arquivo**: `backend/internal/handler/user_management_handler.go`  
**Linha**: 165-169  
**Explicação Técnica**: O handler valida se o role é válido usando `domain.ParseRole()`, mas não verifica se o role pode ser alterado pelo usuário atual. Por exemplo, um Admin não pode alterar outro Admin para Owner, mas essa validação está apenas no serviço.

**Impacto**: Validação duplicada entre handler e serviço. Se a validação do serviço falhar ou for removida, ações não autorizadas podem ser permitidas.

**Melhor Solução**: Centralizar toda validação de RBAC no serviço, remover validações duplicadas do handler.

---

### 18. Falta de Validação de Self-ChangeRole

**Severidade**: Alto  
**Arquivo**: `backend/internal/service/user_management_service.go`  
**Linha**: 108-171  
**Explicação Técnica**: O método `ChangeRole()` não impede que um usuário altere seu próprio role. Um Owner pode (acidentalmente ou intencionalmente) se rebaixar para Manager, perdendo acesso administrativo.

**Impacto**: Usuários podem perder acesso administrativo acidentalmente. Problemas de governança.

**Melhor Solução**: Impedir que usuários alterem seu próprio role.

```go
// user_management_service.go - adicionar verificação
if actorUserID == targetUserID {
    return errors.New("não é possível alterar seu próprio cargo")
}
```

---

### 19. Falta de Validação de LastOwner

**Severidade**: Alto  
**Arquivo**: `backend/internal/service/user_management_service.go`  
**Linha**: 173-220  
**Explicação Técnica**: O método `RemoveFromCompany()` impede remover Owner, mas não impede remover o último Owner da empresa. Se o único Owner for removido, a empresa fica sem administrador.

**Impacto**: Empresa pode ficar sem administrador, causando problemas de governança e operação.

**Melhor Solução**: Verificar se é o último Owner antes de remover.

```go
// user_management_service.go - adicionar verificação
if target.Role != nil && *target.Role == domain.RoleOwner {
    // Verificar se é o último Owner
    ownerCount, _ := s.userRepo.CountByRoleAndCompany(ctx, domain.RoleOwner, *actor.CompanyID)
    if ownerCount <= 1 {
        return errors.New("não é possível remover o último Owner da empresa")
    }
}
```

---

## PROBLEMAS MÉDIOS

### 20. Duplicação de Error Variables

**Severidade**: Médio  
**Arquivo**: `backend/internal/service/user_management_service.go`, `backend/internal/service/invitation_service.go`  
**Linha**: 13-20 (user_management), 13-20 (invitation)  
**Explicação Técnica**: As variáveis de erro `ErrUserNotFound`, `ErrPermissionDenied`, `ErrUserAlreadyInCompany` estão duplicadas em ambos os arquivos. Isso pode causar confusão e problemas de comparação de erros.

**Impacto**: Comparação de erros pode não funcionar corretamente se as variáveis não forem as mesmas. Confusão na manutenção.

**Melhor Solução**: Criar um pacote de erros compartilhado (ex: `internal/errors`) e importar em todos os services.

```go
// internal/errors/errors.go
var (
    ErrUserNotFound = errors.New("usuário não encontrado")
    ErrPermissionDenied = errors.New("permissão negada")
    // ...
)
```

---

### 21. Falta de Método CanManageMedia em RBACService

**Severidade**: Médio  
**Arquivo**: `backend/internal/service/rbac_service.go`  
**Linha**: N/A  
**Explicação Técnica**: O RBACService não tem método `CanManagedMedia()`, mas há endpoints de media que precisam de controle de acesso.

**Impacto**: Falta de controle de acesso granular para mídia. Dificuldade de implementar RBAC para mídia.

**Melhor Solução**: Adicionar método `CanManageMedia()` ao RBACService.

```go
// rbac_service.go - adicionar método
func (s *RBACService) CanManageMedia(ctx context.Context, userID uint) (bool, error) {
    return s.HasAnyRole(ctx, userID, domain.RoleOwner, domain.RoleAdmin, domain.RoleManager)
}
```

---

### 22. Falta de Método CanViewDashboard em RBACService

**Severidade**: Médio  
**Arquivo**: `backend/internal/service/rbac_service.go`  
**Linha**: N/A  
**Explicação Técnica**: O RBACService não tem método `CanViewDashboard()`, mas há endpoint de dashboard que precisa de controle de acesso.

**Impacto**: Falta de controle de acesso para dashboard. Todos os usuários autenticados podem acessar.

**Melhor Solução**: Adicionar método `CanViewDashboard()` ao RBACService.

```go
// rbac_service.go - adicionar método
func (s *RBACService) CanViewDashboard(ctx context.Context, userID uint) (bool, error) {
    return s.HasAnyRole(ctx, userID, domain.RoleOwner, domain.RoleAdmin, domain.RoleManager)
}
```

---

### 23. Falta de Validação de PasswordHash em Update

**Severidade**: Médio  
**Arquivo**: `backend/internal/infra/repository/gorm_user_repository.go`  
**Linha**: 79-98  
**Explicação Técnica**: O método `Update()` não atualiza o `PasswordHash`, mesmo se for fornecido no domain.User. Isso pode causar confusão se alguém tentar atualizar a senha através deste método.

**Impacto**: Confusão na API. Se alguém tentar atualizar senha através de User.Update(), não funcionará.

**Melhor Solução**: Documentar claramente que Update() não atualiza senha, ou adicionar parâmetro explícito para atualização de senha.

---

### 24. Falta de Endpoint companies no Frontend API Client

**Severidade**: Médio  
**Arquivo**: `frontend/src/lib/api/client.ts`  
**Linha**: N/A  
**Explicação Técnica**: O frontend API client não tem métodos para empresas (list, create, update, delete), apesar dos endpoints existirem no backend.

**Impacto**: Frontend não pode gerenciar empresas. Funcionalidade incompleta.

**Melhor Solução**: Adicionar métodos de empresas ao frontend API client.

```typescript
// api/client.ts - adicionar companies
companies: {
  list: () => request<any[]>('/companies'),
  create: (body) => request<any>('/companies', { method: 'POST', body: JSON.stringify(body) }),
  // ...
}
```

---

### 25. Falta de Validação de CompanyID em CreateInvitation

**Severidade**: Médio  
**Arquivo**: `backend/internal/service/invitation_service.go`  
**Linha**: 63-120  
**Explicação Técnica**: O método `CreateInvitation()` recebe companyID como parâmetro, mas não valida se o companyID existe ou se o usuário tem permissão para criar convites para essa empresa.

**Impacto**: Usuários podem tentar criar convites para empresas que não existem ou que não têm permissão.

**Melhor Solução**: Validar se companyID existe e se usuário tem permissão.

```go
// invitation_service.go - adicionar validação
company, err := s.companyRepo.FindByID(ctx, companyID)
if err != nil || company == nil {
    return nil, errors.New("empresa não encontrada")
}
```

---

### 26. Falta de Validação de Self-Invitation

**Severidade**: Médio  
**Arquivo**: `backend/internal/service/invitation_service.go`  
**Linha**: 63-120  
**Explicação Técnica**: O método `CreateInvitation()` não impede que um usuário crie um convite para si mesmo.

**Impacto**: Usuários podem criar convites para si mesmos, causando confusão e comportamento inconsistente.

**Melhor Solução**: Impedir que usuários criem convites para si mesmos.

```go
// invitation_service.go - adicionar verificação
actor, err := s.userRepo.FindByID(ctx, actorUserID)
if err != nil || actor == nil {
    return nil, errors.New("usuário não encontrado")
}
if actor.Email == email {
    return nil, errors.New("não é possível criar convite para si mesmo")
}
```

---

### 27. Falta de Validação de Role em AddExistingUser

**Severidade**: Médio  
**Arquivo**: `backend/internal/service/user_management_service.go`  
**Linha**: 222-283  
**Explicação Técnica**: O método `AddExistingUser()` adiciona usuário com role padrão Manager, mas não valida se o usuário tem permissão para adicionar usuários com esse role. Por exemplo, um Admin não deveria poder adicionar outro Admin.

**Impacto**: Admins podem adicionar outros Admins, violando a regra de que apenas Owner pode alterar roles de Admin.

**Melhor Solução**: Validar se o usuário tem permissão para adicionar usuários com o role padrão.

```go
// user_management_service.go - adicionar validação
if defaultRole == domain.RoleAdmin {
    canAlter, err := s.rbacService.CanAlterAdminRole(ctx, actorUserID)
    if err != nil || !canAlter {
        return nil, ErrPermissionDenied
    }
}
```

---

### 28. Falta de Validação de MaxOwners

**Severidade**: Médio  
**Arquivo**: `backend/internal/service/user_management_service.go`  
**Linha**: 108-171  
**Explicação Técnica**: O método `ChangeRole()` não impede criar múltiplos Owners. Não há limite de quantos Owners uma empresa pode ter.

**Impacto**: Empresa pode ter múltiplos Owners, o que pode causar problemas de governança e segurança.

**Melhor Solução**: Limitar número de Owners por empresa (ex: máximo 1 ou 3).

```go
// user_management_service.go - adicionar verificação
if newRole == domain.RoleOwner {
    ownerCount, _ := s.userRepo.CountByRoleAndCompany(ctx, domain.RoleOwner, *actor.CompanyID)
    if ownerCount >= 1 { // ou 3, dependendo da regra de negócio
        return errors.New("empresa já atingiu o limite de Owners")
    }
}
```

---

## PROBLEMAS BAIXOS

### 29. Comentário Inconsistente em main.go

**Severidade**: Baixo  
**Arquivo**: `backend/cmd/server/main.go`  
**Linha**: 152-154  
**Explicação Técnica**: O comentário "Public invitation endpoints (no authentication required)" está incorreto, pois esses endpoints estão dentro do grupo protegido por AuthMiddleware e TenantMiddleware.

**Impacto**: Documentação incorreta pode causar confusão para desenvolvedores.

**Melhor Solução**: Mover endpoints públicos para fora do grupo de autenticação, ou corrigir o comentário.

---

### 30. Falta de Validação de Length em Token

**Severidade**: Baixo  
**Arquivo**: `backend/internal/domain/invitation.go`  
**Linha**: N/A  
**Explicação Técnica**: O método `GenerateToken()` gera um token de 64 caracteres hexadecimais, mas não há validação de comprimento ou formato no domain.

**Impacto**: Se a implementação de token mudar, pode causar problemas de compatibilidade.

**Melhor Solução**: Adicionar validação de comprimento e formato no domain.

```go
// invitation.go - adicionar validação
func (t *Invitation) ValidateToken() error {
    if len(t.Token) != 64 {
        return errors.New("token inválido")
    }
    return nil
}
```

---

### 31. Falta de Validação de ExpirationTime

**Severidade**: Baixo  
**Arquivo**: `backend/internal/domain/invitation.go`  
**Linha**: N/A  
**Explicação Técnica**: O método `IsExpired()` verifica se a data atual é maior que ExpiresAt, mas não valida se ExpiresAt é uma data válida (ex: não está no passado na criação).

**Impacto**: Convites podem ser criados com expiração no passado.

**Melhor Solução**: Validar ExpiresAt na criação.

```go
// invitation_service.go - adicionar validação
if invitation.ExpiresAt.Before(time.Now()) {
    return nil, errors.New("data de expiração inválida")
}
```

---

### 32. Falta de Validação de Role em InvitationOutput

**Severidade**: Baixo  
**Arquivo**: `backend/internal/service/invitation_service.go`  
**Linha**: 274-287  
**Explicação Técnica**: O método `toInvitationOutput()` converte Role para string usando `String()`, mas não valida se a conversão é válida.

**Impacto**: Se Role for inválido, pode retornar string vazia ou inválida.

**Melhor Solução**: Validar conversão de Role para string.

```go
// invitation_service.go - adicionar validação
roleStr := invitation.Role.String()
if roleStr == "" {
    roleStr = "unknown"
}
```

---

### 33. Falta de Validação de Status em InvitationOutput

**Severidade**: Baixo  
**Arquivo**: `backend/internal/service/invitation_service.go`  
**Linha**: 274-287  
**Explicação Técnica**: O método `toInvitationOutput()` converte Status para string usando `String()`, mas não valida se a conversão é válida.

**Impacto**: Se Status for inválido, pode retornar string vazia ou inválida.

**Melhor Solução**: Validar conversão de Status para string.

```go
// invitation_service.go - adicionar validação
statusStr := invitation.Status.String()
if statusStr == "" {
    statusStr = "unknown"
}
```

---

### 34. Falta de Validação de Email em AcceptInvitation

**Severidade**: Baixo  
**Arquivo**: `backend/internal/service/invitation_service.go`  
**Linha**: 216-272  
**Explicação Técnica**: O método `AcceptInvitation()` busca usuário por email do convite, mas não valida se o email do usuário bate com o email do convite (caso o email do usuário tenha mudado).

**Impacto**: Se o email do usuário mudou, ele pode aceitar um convite que não era para ele.

**Melhor Solução**: Validar se o email do usuário bate com o email do convite.

```go
// invitation_service.go - adicionar validação
if user.Email != invitation.Email {
    return errors.New("o e-mail do usuário não corresponde ao convite")
}
```

---

## VIOLAÇÕES DE CLEAN ARCHITECTURE

### 35. Handler Fazendo Lógica de Negócio

**Severidade**: Médio  
**Arquivo**: `backend/internal/handler/invitation_handler.go`  
**Linha**: 30-40  
**Explicação Técnica**: O método `getUserCompanyID()` no handler faz lógica de acesso a dados (buscar usuário por ID), o que viola a separação de responsabilidades. Essa lógica deveria estar no serviço ou em um helper.

**Impacto**: Violação de Clean Architecture. Lógica de negócio misturada com lógica de apresentação.

**Melhor Solução**: Mover lógica para o serviço ou criar um helper no pacote middleware.

---

### 36. Service Acessando Repository Diretamente sem Interface

**Severidade**: Baixo  
**Arquivo**: `backend/internal/service/user_management_service.go`  
**Linha**: 47-75  
**Explicação Técnica**: O método `ListUsers()` chama `s.userRepo.List(ctx)` que retorna todos os usuários, e depois filtra em memória. Isso viola o princípio de que o service deveria delegar toda lógica de acesso de dados ao repository.

**Impacto**: Violação de Clean Architecture. Lógica de filtragem no service em vez do repository.

**Melhor Solução**: Adicionar método `FindByCompanyID()` no repository e usar no service.

---

## PROBLEMAS DE PERFORMANCE

### 37. N+1 Query em GetProductIngredients

**Severidade**: Médio  
**Arquivo**: `backend/internal/infra/repository/gorm_product_repository.go`  
**Linha**: 436-455  
**Explicação Técnica**: O método `GetProductIngredients()` usa `Preload("Ingredient")` que pode causar N+1 queries se não for otimizado corretamente.

**Impacto**: Performance degradada para produtos com muitos ingredientes.

**Melhor Solução**: Usar joins explícitos ou otimizar o preload.

---

### 38. Falta de Paginação em List Methods

**Severidade**: Médio  
**Arquivo**: Múltiplos repositories  
**Linha**: N/A  
**Explicação Técnica**: Os métodos `List*()` não têm paginação. Conforme o número de registros cresce, a performance degrada.

**Impacto**: Performance degradada para grandes volumes de dados. Consumo excessivo de memória.

**Melhor Solução**: Adicionar paginação em todos os métodos `List*()`.

```go
// ports - adicionar parâmetros de paginação
List(ctx context.Context, limit, offset int) ([]*domain.User, error)
```

---

## PROBLEMAS DE SEGURANÇA

### 39. Falta de CSRF Protection

**Severidade**: Alto  
**Arquivo**: `backend/cmd/server/main.go`  
**Linha**: N/A  
**Explicação Técnica**: Não há proteção CSRF em nenhum endpoint que modifica estado (POST, PUT, DELETE).

**Impacto**: Vulnerabilidade a ataques CSRF. Usuários podem ser induzidos a executar ações não autorizadas.

**Melhor Solução**: Implementar proteção CSRF usando tokens ou SameSite cookies.

---

### 40. Falta de Input Sanitization Global

**Severidade**: Médio  
**Arquivo**: Múltiplos handlers  
**Linha**: N/A  
**Explicação Técnica**: Não há sanitização global de input. XSS pode ocorrer se os dados forem exibidos sem sanitização.

**Impacto**: Vulnerabilidade de XSS se os dados forem exibidos sem sanitização.

**Melhor Solução**: Implementar middleware de sanitização ou sanitizar na saída.

---

### 41. Falta de Rate Limiting Global

**Severidade**: Médio  
**Arquivo**: `backend/cmd/server/main.go`  
**Linha**: N/A  
**Explicação Técnica**: Não há rate limiting global. Atacantes podem fazer DoS.

**Impacto**: Vulnerabilidade a DoS. Servidor pode ficar indisponível.

**Melhor Solução**: Implementar rate limiting global por IP.

---

## PROBLEMAS DE MANUTENÇÃO

### 42. Falta de Testes de Integração

**Severidade**: Alto  
**Arquivo**: N/A  
**Linha**: N/A  
**Explicação Técnica**: Não há testes de integração visíveis no código. Apenas testes unitários podem não ser suficientes.

**Impacto**: Regressões podem ocorrer sem detecção. Dificuldade de refatorar.

**Melhor Solução**: Adicionar testes de integração para fluxos críticos.

---

### 43. Falta de Documentação de API

**Severidade**: Médio  
**Arquivo**: N/A  
**Linha**: N/A  
**Explicação Técnica**: Não há documentação de API (Swagger/OpenAPI). Desenvolvedores frontend precisam ler o código para entender os endpoints.

**Impacto**: Dificuldade de integração frontend-backend. Comunicação ineficiente.

**Melhor Solução**: Adicionar documentação Swagger/OpenAPI.

---

### 44. Falta de Logs Estruturados

**Severidade**: Médio  
**Arquivo**: Múltiplos arquivos  
**Linha**: N/A  
**Explicação Técnica**: Logs são strings não estruturadas. Dificulta análise de logs em produção.

**Impacto**: Dificuldade de debug e monitoramento em produção.

**Melhor Solução**: Implementar logs estruturados (ex: usando zerolog ou logrus).

---

## PROBLEMAS DE FRONTEND

### 45. Falta de Tratamento de Erros Global

**Severidade**: Médio  
**Arquivo**: `frontend/src/lib/api/client.ts`  
**Linha**: 46-67  
**Explicação Técnica**: O tratamento de erros é básico. Não há distinção entre diferentes tipos de erros (network, validation, auth, etc.).

**Impacto**: UX ruim. Usuários recebem mensagens genéricas de erro.

**Melhor Solução**: Implementar tratamento de erros granular com tipos específicos.

---

### 46. Falta de Loading States Globais

**Severidade**: Baixo  
**Arquivo**: N/A  
**Linha**: N/A  
**Explicação Técnica**: Não há loading states globais. Cada componente precisa implementar seu próprio loading.

**Impacto**: UX inconsistente. Código duplicado.

**Melhor Solução**: Implementar loading state global com context ou store.

---

### 47. Falta de Error Boundary

**Severidade**: Médio  
**Arquivo**: N/A  
**Linha**: N/A  
**Explicação Técnica**: Não há Error Boundary para capturar erros de React/Svelte.

**Impacto**: Erros não tratados podem quebrar a aplicação inteira.

**Melhor Solução**: Implementar Error Boundary no nível da aplicação.

---

## RECOMENDAÇÕES PRIORITÁRIAS

### Imediato (Antes de Produção)

1. **Corrigir RBAC inconsistency** (Problema 1)
2. **Corrigir tenant isolation em empresas** (Problema 2)
3. **Corrigir cookie security** (Problema 3)
4. **Adicionar transação em AcceptInvitation** (Problema 5)
5. **Validar role em criação de convite** (Problema 6)
6. **Adicionar rate limiting em endpoints públicos** (Problema 7)

### Curto Prazo (1-2 Semanas)

7. **Remover duplicação de getUserCompanyID** (Problema 8)
8. **Adicionar validação de deleção de empresa** (Problema 9)
9. **Otimizar ListUsers** (Problema 10)
10. **Adicionar índices compostos em invitations** (Problema 12)
11. **Validar email em Invitation** (Problema 13)
12. **Adicionar log de auditoria** (Problema 15)

### Médio Prazo (1 Mês)

13. **Adicionar validação de slug** (Problema 11)
14. **Sanitizar company settings** (Problema 14)
15. **Tratar erro de expiração** (Problema 16)
16. **Impedir self-changeRole** (Problema 18)
17. **Validar lastOwner** (Problema 19)
18. **Criar pacote de erros compartilhado** (Problema 20)

### Longo Prazo (2-3 Meses)

19. **Adicionar métodos faltantes ao RBACService** (Problemas 21-22)
20. **Adicionar endpoints de companies no frontend** (Problema 24)
21. **Adicionar paginação em List methods** (Problema 38)
22. **Implementar CSRF protection** (Problema 39)
23. **Adicionar testes de integração** (Problema 42)
24. **Adicionar documentação Swagger** (Problema 43)

---

## CONCLUSÃO

A Plataforma PratoOnline 2.0 apresenta uma arquitetura sólida com boa separação de camadas, mas contém problemas críticos de segurança e consistência que devem ser corrigidos antes de ir para produção. Os principais riscos são:

1. **Violação de RBAC** - Inconsistência entre middleware e serviço
2. **Falha de tenant isolation** - Endpoints de empresas não filtram por tenant
3. **Vulnerabilidade de segurança** - Cookie não seguro em produção
4. **Race conditions** - Falta de transações em operações críticas
5. **Performance** - Queries ineficientes e falta de paginação

Recomenda-se priorizar a correção dos 7 problemas críticos antes de qualquer deploy em produção, seguido pelos 12 problemas altos. Os problemas médios e baixos podem ser corrigidos iterativamente.

A arquitetura geral é boa e segue princípios de Clean Architecture, mas precisa de refinamentos em segurança, performance, e consistência para estar pronta para produção.
