# ETAPA 4: Bugs Identificados - Auditoria Forense
**Auditoria Forense - pratoOnline**
**Data:** 19/07/2026

---

# Bug 1: Cookie Secure=false em produção

## Local
**Arquivo:** `backend/internal/handler/auth_handler.go`
**Linha:** 87
**Função:** `Login`
**Gravidade:** Alta

## Fluxo afetado
Autenticação - Login

## Causa raiz
O cookie JWT é definido com `Secure: false` hardcoded, independente do ambiente. Em produção com HTTPS, o cookie não será enviado pelo navegador, causando falha de autenticação.

## Por que acontece
A linha 87 define `Secure: false` como valor fixo, sem verificar a variável de ambiente `ENVIRONMENT`.

## Qual fluxo quebra
Em produção com HTTPS, o usuário faz login com sucesso, o cookie é definido, mas o navegador não envia o cookie em requisições subsequentes (porque Secure=false em HTTPS), causando falha de autenticação em todas as requisições protegidas.

## Qual o impacto
- **Alta:** Usuários não conseguem autenticar em produção
- **Segurança:** Cookie pode ser interceptado em HTTP não criptografado
- **Disponibilidade:** Sistema inutilizável em produção com HTTPS

## Correção mínima
Verificar variável de ambiente `ENVIRONMENT` e definir `Secure=true` quando `ENVIRONMENT=="production"`.

## Correção ideal
Implementar verificação de ambiente com fallback seguro, adicionar testes unitários para diferentes ambientes, documentar configuração necessária.

## Evidência
**Estado antes (código original):**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    result.Token,
    Path:     "/",
    HttpOnly: true,
    Secure:   false, // true em produção (HTTPS)
    SameSite: http.SameSiteLaxMode,
    Expires:  time.Now().Add(24 * time.Hour),
})
```

**Estado após (correção aplicada):**
```go
secureCookie := os.Getenv("ENVIRONMENT") == "production"
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    result.Token,
    Path:     "/",
    HttpOnly: true,
    Secure:   secureCookie,
    SameSite: http.SameSiteLaxMode,
    Expires:  time.Now().Add(24 * time.Hour),
})
```

## Validação após correção
**Status:** ✅ Corrigido
**Teste:** Backend compilou com sucesso após correção. Código verifica ENVIRONMENT.

---

# Bug 2: Blacklist JWT in-memory (não persiste)

## Local
**Arquivo:** `backend/internal/service/auth_service.go`
**Linha:** 255-261
**Função:** `Logout`
**Gravidade:** Alta

## Fluxo afetado
Autenticação - Logout

## Causa raiz
A blacklist de tokens JWT é armazenada em um `map[string]time.Time` in-memory na estrutura `AuthService`. Quando o servidor é reiniciado, o map é perdido e todos os tokens revogados anteriormente tornam-se válidos novamente.

## Por que acontece
A estrutura `AuthService` define `blacklist map[string]time.Time` sem persistência em banco de dados. O método `Logout` apenas adiciona o token ao map in-memory.

## Qual fluxo quebra
1. Usuário faz logout → token adicionado ao blacklist in-memory
2. Servidor é reiniciado → blacklist é perdida
3. Token revogado anteriormente torna-se válido novamente
4. Usuário pode acessar o sistema com token que deveria estar revogado

## Qual o impacto
- **Alta:** Tokens revogados continuam válidos após restart do servidor
- **Segurança:** Violação de segurança - logout não é permanente
- **Compliance:** Não atende requisitos de segurança para revogação de sessões

## Correção mínima
Documentar TODO com recomendação de implementar persistência em banco de dados.

## Correção ideal
Criar tabela `token_blacklist` com campos (token, revoked_at, expires_at), implementar repository para persistência, modificar `Logout` para persistir no banco, modificar `ValidateToken` para consultar o banco, implementar limpeza automática de tokens expirados.

## Evidência
**Estado antes (código original):**
```go
func (s *AuthService) Logout(ctx context.Context, tokenStr string) error {
    s.blacklist[tokenStr] = time.Now()
    return nil
}
```

**Estado após (correção aplicada):**
```go
func (s *AuthService) Logout(ctx context.Context, tokenStr string) error {
    // TODO: Implementar persistência de blacklist em banco de dados
    // Atualmente usa map in-memory que não persiste entre restarts
    // Recomendação: criar tabela token_blacklist com campos (token, revoked_at, expires_at)
    s.blacklist[tokenStr] = time.Now()
    return nil
}
```

## Validação após correção
**Status:** ⚠️ Documentado (não corrigido - requer refatoração maior)
**Teste:** TODO adicionado ao código. Correção requer implementação de tabela e repository.

---

# Bug 3: Update usuário não atualiza PasswordHash e Active

## Local
**Arquivo:** `backend/internal/infra/repository/gorm_user_repository.go`
**Linha:** 84-103
**Função:** `Update`
**Gravidade:** Alta

## Fluxo afetado
Gestão de Usuários - Atualizar usuário

## Causa raiz
O método `Update` do repository não inclui os campos `PasswordHash` e `Active` no model GORM ao persistir. Mesmo que o service atualize esses campos no domain.User, eles não são persistidos no banco.

## Por que acontece
A linha 85-92 constrói o `GormUserModel` apenas com campos ID, Name, Email, CompanyID e Role, omitindo PasswordHash e Active.

## Qual fluxo quebra
1. Service atualiza PasswordHash (ex: change password) ou Active (ex: desativar usuário)
2. Repository.Update é chamado
3. PasswordHash e Active não são persistidos no banco
4. Próxima leitura do usuário retorna valores antigos
5. Senha não é alterada ou usuário não é desativado

## Qual o impacto
- **Alta:** Alteração de senha não funciona
- **Alta:** Desativação de usuário não funciona
- **Segurança:** Usuários não podem ser desativados
- **Disponibilidade:** Função crítica de gestão de usuários inoperante

## Correção mínima
Adicionar campos `PasswordHash` e `Active` ao model GormUserModel no método Update.

## Correção ideal
Adicionar campos, implementar testes unitários para Update, validar que todos os campos do domain.User são persistidos.

## Evidência
**Estado antes (código original):**
```go
func (r *GormUserRepository) Update(ctx context.Context, user *domain.User) error {
    model := GormUserModel{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CompanyID: user.CompanyID,
        Role:      nil,
    }
    if user.Role != nil {
        role := user.Role.String()
        model.Role = &role
    }
    if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
        return fmt.Errorf("UserRepository.Update: %w", err)
    }
    user.UpdatedAt = time.Unix(model.UpdatedAt, 0)
    return nil
}
```

**Estado após (correção aplicada):**
```go
func (r *GormUserRepository) Update(ctx context.Context, user *domain.User) error {
    model := GormUserModel{
        ID:           user.ID,
        Name:         user.Name,
        Email:        user.Email,
        PasswordHash: user.PasswordHash,
        Active:       user.Active,
        CompanyID:    user.CompanyID,
        Role:         nil,
    }
    if user.Role != nil {
        role := user.Role.String()
        model.Role = &role
    }
    if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
        return fmt.Errorf("UserRepository.Update: %w", err)
    }
    user.UpdatedAt = time.Unix(model.UpdatedAt, 0)
    return nil
}
```

## Validação após correção
**Status:** ✅ Corrigido
**Teste:** Backend compilou com sucesso após correção. Campos PasswordHash e Active agora são persistidos.

---

# Bug 4: Slug uniqueIndex sem validação de colisão

## Local
**Arquivo:** `backend/internal/infra/repository/gorm_product_repository.go`
**Linha:** 96-155 (CreateProduct), 196-244 (UpdateProduct)
**Função:** `CreateProduct`, `UpdateProduct`
**Gravidade:** Alta

## Fluxo afetado
Produtos - Criar Produto, Atualizar Produto

## Causa raiz
O campo `Slug` tem `uniqueIndex` no banco, mas não há validação antes de inserir/atualizar. Quando ocorre colisão, o banco retorna erro de constraint, mas a aplicação não trata especificamente, retornando erro genérico ao usuário.

## Por que acontece
Os métodos `CreateProduct` e `UpdateProduct` não verificam se o slug já existe antes de tentar persistir. Apenas tentam inserir/atualizar e deixam o banco retornar erro de constraint.

## Qual fluxo quebra
1. Usuário cria produto com slug "produto-a"
2. Outro usuário tenta criar produto com slug "produto-a"
3. Banco retorna erro de unique constraint
4. Aplicação retorna erro genérico "CreateProduct: UNIQUE constraint failed"
5. Usuário não recebe feedback claro sobre o problema

## Qual o impacto
- **Média:** Experiência do usuário degradada
- **Média:** Erros genéricos dificultam debug
- **Baixa:** Funcionalidade funciona, mas com feedback ruim

## Correção mínima
Adicionar verificação de colisão antes de inserir/atualizar, retornar erro específico quando slug já existe.

## Correção ideal
Adicionar verificação em CreateProduct e UpdateProduct, implementar validação no frontend, adicionar sugestão de slug alternativo, implementar testes unitários.

## Evidência
**Estado antes (código original - CreateProduct):**
```go
func (r *GormProductRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
    companyID, err := GetCompanyIDFromContext(ctx)
    if err != nil {
        return fmt.Errorf("CreateProduct: %w", err)
    }

    m := GormProduct{
        // ... campos ...
        Slug: p.Slug,
        // ... campos ...
    }
    // ... persiste sem verificar colisão ...
}
```

**Estado após (correção aplicada - CreateProduct):**
```go
func (r *GormProductRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
    companyID, err := GetCompanyIDFromContext(ctx)
    if err != nil {
        return fmt.Errorf("CreateProduct: %w", err)
    }

    // Check for slug collision
    if p.Slug != "" {
        var existing GormProduct
        err := r.db.WithContext(ctx).Where("slug = ?", p.Slug).First(&existing).Error
        if err == nil {
            return fmt.Errorf("CreateProduct: slug '%s' já está em uso", p.Slug)
        }
        if !errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("CreateProduct: verificar slug: %w", err)
        }
    }

    m := GormProduct{
        // ... campos ...
        Slug: p.Slug,
        // ... campos ...
    }
    // ... persiste ...
}
```

**Estado após (correção aplicada - UpdateProduct):**
```go
func (r *GormProductRepository) UpdateProduct(ctx context.Context, p *domain.Product) error {
    // ... verifica produto existe ...

    // Check for slug collision (excluding current product)
    if p.Slug != "" && p.Slug != existing.Slug {
        var slugConflict GormProduct
        err := r.db.WithContext(ctx).Where("slug = ? AND id != ?", p.Slug, p.ID).First(&slugConflict).Error
        if err == nil {
            return fmt.Errorf("UpdateProduct: slug '%s' já está em uso", p.Slug)
        }
        if !errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("UpdateProduct: verificar slug: %w", err)
        }
    }

    // ... atualiza ...
}
```

## Validação após correção
**Status:** ✅ Corrigido
**Teste:** Backend compilou com sucesso após correção. Validação de colisão implementada em CreateProduct e UpdateProduct.

---

# Bug 5: Lógica de slug pode regenerar desnecessariamente

## Local
**Arquivo:** `backend/internal/service/product_service.go`
**Linha:** 282-287
**Função:** `UpdateProduct`
**Gravidade:** Média

## Fluxo afetado
Produtos - Atualizar Produto

## Causa raiz
A condição `in.Slug == "" || in.Slug != p.Slug` causa regeneração do slug mesmo quando o usuário fornece um slug válido mas diferente do atual. Isso pode causar colisões inesperadas.

## Por que acontece
A linha 283 verifica se o slug está vazio OU se é diferente do atual. Se o usuário fornecer um slug diferente (mas válido), o código regenera o slug baseado no nome, ignorando o slug fornecido.

## Qual fluxo quebra
1. Usuário cria produto com slug "produto-a"
2. Usuário atualiza produto fornecendo slug "produto-b" (intencionalmente diferente)
3. Código ignora "produto-b" e regenera slug baseado no nome
4. Slug gerado pode colidir com outro produto
5. Usuário não tem controle sobre o slug

## Qual o impacto
- **Média:** Usuário não pode controlar o slug ao atualizar
- **Média:** Slug pode mudar inesperadamente
- **Baixa:** Funcionalidade funciona, mas com comportamento inesperado

## Correção mínima
Simplificar condição para regenerar apenas quando slug está vazio.

## Correção ideal
Regenerar apenas quando vazio, permitir controle total do usuário quando fornecido, adicionar validação de formato de slug, implementar testes unitários.

## Evidência
**Estado antes (código original):**
```go
// Gerar slug automaticamente se não fornecido ou se nome mudou
if in.Slug == "" || in.Slug != p.Slug {
    p.Slug = generateSlug(in.Name)
} else {
    p.Slug = in.Slug
}
```

**Estado após (correção aplicada):**
```go
// Gerar slug automaticamente se não fornecido
if in.Slug == "" {
    p.Slug = generateSlug(in.Name)
} else {
    p.Slug = in.Slug
}
```

## Validação após correção
**Status:** ✅ Corrigido
**Teste:** Backend compilou com sucesso após correção. Lógica simplificada para regenerar apenas quando vazio.

---

# Resumo dos Bugs

| ID | Bug | Gravidade | Status |
|----|-----|-----------|--------|
| 1 | Cookie Secure=false em produção | Alta | ✅ Corrigido |
| 2 | Blacklist JWT in-memory | Alta | ⚠️ Documentado |
| 3 | Update usuário não atualiza PasswordHash/Active | Alta | ✅ Corrigido |
| 4 | Slug uniqueIndex sem validação de colisão | Alta | ✅ Corrigido |
| 5 | Lógica de slug pode regenerar desnecessariamente | Média | ✅ Corrigido |

**Total:** 5 bugs identificados
**Corrigidos:** 4
**Documentados:** 1 (requer refatoração maior)
