# Sprint 2.1 - Bloqueadores Identificados

**Data:** 19/07/2026  
**Responsável:** QA Team  
**Severidade:** CRÍTICO  
**Status:** MVP NÃO PRONTO PARA PRODUÇÃO

---

## Resumo

**Total de Bloqueadores:** 2  
**Bloqueadores Críticos:** 1  
**Bloqueadores Moderados:** 1  
**Impacto:** Sistema de convites completamente não funcional

---

## Bloqueador #1 - Sistema de Convites Quebrado

**ID:** BLOCKER-001  
**Severidade:** 🔴 CRÍTICO  
**Fluxo Afetado:** CONVITES - Aceitar Convite / Entrar na Empresa  
**Status:** CONFIRMADO

### Descrição

O sistema de convites está completamente quebrado. Quando um usuário se registra no sistema, ele automaticamente recebe uma empresa criada (role "owner"). Isso impede que o mesmo usuário aceite convites de outras empresas, pois o sistema retorna o erro "usuário já pertence a outra empresa".

### Impacto

- **Funcionalidade:** Sistema de convites completamente inoperacional
- **Cenário:** Impossível adicionar colaboradores a empresas existentes
- **Usuários:** Novos usuários não podem entrar em empresas via convite
- **Negócio:** Fluxo crítico de onboarding de colaboradores quebrado

### Evidência

**Teste Executado:**
1. Criar convite para `qainvitee@example.com` com role "admin"
2. Registrar usuário com email `qainvitee@example.com`
3. Login com o usuário
4. Tentar aceitar convite via `/api/invitations/accept`

**Resultado:**
```json
{
  "error": "usuário já pertence a outra empresa"
}
```

**Logs do Banco:**
```sql
-- Usuário criado automaticamente com empresa
SELECT id, name, email, role, company_id FROM users WHERE email='qainvitee@example.com';
-- Resultado: 5|QA Invitee|qainvitee@example.com|owner|5

-- Empresa criada automaticamente
SELECT id, name, slug FROM companies WHERE id=5;
-- Resultado: 5|QA Invitee's Company|qainvitee-178451...
```

### Causa Raiz

O fluxo de registro (`POST /api/auth/register`) automaticamente cria uma empresa para o novo usuário e atribui o role "owner". Isso acontece no serviço de autenticação, provavelmente no método `Register`.

**Arquivo Provável:** `backend/internal/service/auth_service.go`

### Reprodução Passo a Passo

1. Como Owner de uma empresa, criar convite para novo email
2. Usuário receber email com link de convite
3. Usuário clicar no link e ser redirecionado para tela de cadastro
4. Usuário preencher formulário de registro
5. Sistema automaticamente criar empresa para o usuário
6. Usuário tentar aceitar convite
7. **ERRO:** "usuário já pertence a outra empresa"

### Solução Sugerida

**Opção 1 - Cadastro Sem Empresa (Recomendado):**
- Modificar fluxo de registro para NÃO criar empresa automaticamente
- Usuário se registra sem empresa (company_id = null)
- Usuário pode aceitar convites normalmente
- Apenas usuários que aceitam convites recebem company_id

**Opção 2 - Fluxo Específico para Convites:**
- Criar fluxo de registro específico para convites
- Usuário acessando link de convite usa endpoint diferente
- Endpoint de convite não cria empresa automaticamente

**Opção 3 - Permitir Troca de Empresa:**
- Permitir que usuários aceitem convites mesmo tendo empresa
- Ao aceitar convite, empresa anterior é desvinculada
- Requer lógica adicional de migração de dados

### Prioridade

**IMEDIATA** - Este bloqueador impede completamente o funcionamento do sistema multi-tenant. Sem correção, o sistema não pode ser considerado MVP funcional.

---

## Bloqueador #2 - Recuperação de Senha Não Implementada

**ID:** BLOCKER-002  
**Severidade:** 🟡 MODERADO  
**Fluxo Afetado:** AUTENTICAÇÃO - Recuperação de Senha  
**Status:** CONFIRMADO

### Descrição

O endpoint de recuperação de senha (`/api/auth/request-password-reset`) retorna 404 Not Found, indicando que a funcionalidade não está implementada.

### Impacto

- **Funcionalidade:** Usuários não podem recuperar senhas esquecidas
- **Cenário:** Usuário que esqueceu senha fica bloqueado
- **Usuários:** Experiência de usuário degradada
- **Segurança:** Não há alternativa para recuperação de acesso

### Evidência

**Teste Executado:**
```bash
curl -X POST http://localhost:8080/api/auth/request-password-reset \
  -H "Content-Type: application/json" \
  -d '{"email":"qatest@example.com"}'
```

**Resultado:**
```
HTTP/1.1 404 Not Found
404 page not found
```

### Causa Raiz

O endpoint não está registrado no router ou o handler não está implementado.

**Arquivo Provável:** `backend/internal/handler/auth_handler.go` (método `RequestPasswordReset` existe mas não está roteado)

### Solução Sugerida

1. Implementar envio de emails com token de recuperação
2. Registrar endpoint no router
3. Implementar endpoint de reset de senha (`/api/auth/reset-password`)
4. Adicionar validação de token e expiração

### Prioridade

**ALTA** - Funcionalidade importante para UX, mas não bloqueia MVP completamente. Usuários podem redefinir senha via admin ou criar nova conta.

---

## Bloqueadores Menores (Não Críticos)

### Menor #1 - Duplicação de Produtos

**ID:** MINOR-001  
**Severidade:** 🟢 BAIXA  
**Fluxo Afetado:** PRODUTOS - Duplicar  
**Status:** CONFIRMADO

**Descrição:** Endpoint para duplicar produtos não encontrado na API.

**Impacto:** Usuários precisam criar produtos do zero em vez de duplicar.

**Solução:** Implementar endpoint `/api/products/{id}/duplicate`.

---

### Menor #2 - Roles Manager e Employee

**ID:** MINOR-002  
**Severidade:** 🟢 BAIXA  
**Fluxo Afetado:** CONVITES - Roles  
**Status:** CONFIRMADO

**Descrição:** Sistema de convites não aceita roles "manager" e "employee", retornando "cargo inválido".

**Impacto:** Apenas roles "owner" e "admin" podem ser atribuídos via convite.

**Solução:** Adicionar "manager" e "employee" ao enum de roles válidos.

---

### Menor #3 - Edição Direta de Pedidos

**ID:** MINOR-003  
**Severidade:** 🟢 BAIXA  
**Fluxo Afetado:** PEDIDOS - Editar  
**Status:** CONFIRMADO

**Descrição:** Endpoint para edição direta de pedidos não encontrado.

**Impacto:** Pedidos só podem ser alterados via status, não edição completa.

**Solução:** Implementar endpoint PUT `/api/orders/{id}` ou manter apenas edição via status.

---

## Plano de Ação

### Imediato (Sprint 2.2)

1. **CORRIGIR BLOQUEADOR #1** - Sistema de Convites
   - Estimar: 4-6 horas
   - Prioridade: CRÍTICA
   - Responsável: Backend Team

### Curto Prazo (Sprint 2.2)

2. **CORRIGIR BLOQUEADOR #2** - Recuperação de Senha
   - Estimar: 6-8 horas
   - Prioridade: ALTA
   - Responsável: Backend Team

### Médio Prazo (Sprint 2.3)

3. **CORRIGIR MENOR #1** - Duplicação de Produtos
   - Estimar: 2-3 horas
   - Prioridade: BAIXA
   - Responsável: Backend Team

4. **CORRIGIR MENOR #2** - Roles Manager e Employee
   - Estimar: 1-2 horas
   - Prioridade: BAIXA
   - Responsável: Backend Team

5. **CORRIGIR MENOR #3** - Edição de Pedidos
   - Estimar: 2-3 horas
   - Prioridade: BAIXA
   - Responsável: Backend Team

---

## Conclusão

O MVP do PratoOnline 2.0 **NÃO ESTÁ PRONTO PARA PRODUÇÃO** devido ao bloqueador crítico no sistema de convites. Esta funcionalidade é essencial para o modelo de negócio multi-tenant e deve ser corrigida antes de qualquer release.

**Recomendação:** Não prosseguir com deploy até que o Bloqueador #1 seja resolvido. O Bloqueador #2 pode ser postergado para um hotfix pós-launch se necessário, mas é altamente recomendável corrigir antes do release.
