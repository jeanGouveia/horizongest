# Sprint 3 - Plano de Testes

**Data:** 19/07/2026  
**Versão:** 3.0  
**Status:** Planejamento

---

## Resumo Executivo

Este documento define o plano completo de testes para validar a refatoração arquitetural multi-tenant do PratoOnline. Os testes cobrem novos fluxos de plataforma, validação de regras de negócio, e testes de regressão para garantir que funcionalidades existentes continuem funcionando.

**Total de Casos de Teste:** 87  
**Testes de Plataforma:** 35  
**Testes de Empresa:** 28  
**Testes de Regressão:** 24

---

## Estrutura de Testes

### Nível 1: Plataforma (Platform Admin/Support)
- Autenticação da plataforma
- Gestão de empresas
- Criação de owners
- Gestão de platform users

### Nível 2: Empresa (Owner/Admin/Manager/Employee)
- Autenticação da empresa
- Gestão de usuários
- Convites (opcional)
- Funcionalidades existentes

### Nível 3: Regressão
- Produtos
- Pedidos
- Ingredientes
- Tema
- RBAC
- JWT
- Soft Delete

---

## Testes de Plataforma

### TP-001: Login Plataforma - Sucesso
**Objetivo:** Validar login de PlatformAdmin com credenciais corretas  
**Pré-condições:** PlatformAdmin criado via migration  
**Passos:**
1. Acessar `/platform/auth`
2. Preencher email: `admin@pratoonline.com`
3. Preencher senha: `admin123`
4. Clicar em "Login"
5. Verificar redirecionamento para `/platform`
**Resultado Esperado:** Login bem-sucedido, redirecionado para dashboard plataforma  
**Prioridade:** CRÍTICA

### TP-002: Login Plataforma - Credenciais Inválidas
**Objetivo:** Validar rejeição de login com credenciais incorretas  
**Pré-condições:** PlatformAdmin criado via migration  
**Passos:**
1. Acessar `/platform/auth`
2. Preencher email: `admin@pratoonline.com`
3. Preencher senha: `senhaerrada`
4. Clicar em "Login"
**Resultado Esperado:** Erro "Credenciais inválidas", permanecer na tela de login  
**Prioridade:** ALTA

### TP-003: Login Plataforma - Usuário Inativo
**Objetivo:** Validar bloqueio de login para platform user inativo  
**Pré-condições:** PlatformAdmin criado e desativado  
**Passos:**
1. Acessar `/platform/auth`
2. Preencher email de platform user inativo
3. Preencher senha correta
4. Clicar em "Login"
**Resultado Esperado:** Erro "Usuário desativado"  
**Prioridade:** ALTA

### TP-004: Login Plataforma - Usuário de Empresa
**Objetivo:** Validar que company users não podem acessar plataforma  
**Pré-condições:** Owner de empresa criado  
**Passos:**
1. Tentar acessar `/platform/auth` com credenciais de Owner
**Resultado Esperado:** Erro "Credenciais inválidas" ou redirecionamento para login empresa  
**Prioridade:** CRÍTICA

### TP-005: Logout Plataforma
**Objetivo:** Validar logout e invalidação de token  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Clicar em "Logout"
2. Tentar acessar `/platform`
**Resultado Esperado:** Redirecionado para login, token invalidado  
**Prioridade:** ALTA

### TP-006: Criar Empresa - Sucesso
**Objetivo:** Validar criação de empresa pela plataforma  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Acessar `/platform/companies`
2. Clicar em "Criar Empresa"
3. Preencher: Nome, CNPJ, Email, Telefone, Plano, Status
4. Clicar em "Salvar"
**Resultado Esperado:** Empresa criada, redirecionado para detalhes da empresa  
**Prioridade:** CRÍTICA

### TP-007: Criar Empresa - Slug Duplicado
**Objetivo:** Validar prevenção de slug duplicado  
**Pré-condições:** Empresa com slug "empresa-teste" existe  
**Passos:**
1. Tentar criar empresa com mesmo nome
**Resultado Esperado:** Erro "Slug já existe"  
**Prioridade:** ALTA

### TP-008: Criar Empresa - Campos Obrigatórios
**Objetivo:** Validar obrigatoriedade de campos  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Tentar criar empresa sem nome
**Resultado Esperado:** Erro de validação "Nome é obrigatório"  
**Prioridade:** ALTA

### TP-009: Editar Empresa - Sucesso
**Objetivo:** Validar edição de empresa pela plataforma  
**Pré-condições:** Empresa criada  
**Passos:**
1. Acessar `/platform/companies/{id}`
2. Clicar em "Editar"
3. Modificar nome
4. Clicar em "Salvar"
**Resultado Esperado:** Empresa atualizada com sucesso  
**Prioridade:** ALTA

### TP-010: Suspender Empresa - Sucesso
**Objetivo:** Validar suspensão de empresa  
**Pré-condições:** Empresa ativa criada  
**Passos:**
1. Acessar `/platform/companies/{id}`
2. Clicar em "Suspender"
3. Confirmar ação
**Resultado Esperado:** Empresa suspensa, status atualizado  
**Prioridade:** ALTA

### TP-011: Suspender Empresa - Bloqueio de Acesso
**Objetivo:** Validar que usuários de empresa suspensa não podem acessar  
**Pré-condições:** Empresa suspensa, Owner existe  
**Passos:**
1. Tentar login com credenciais de Owner de empresa suspensa
**Resultado Esperado:** Erro "Empresa suspensa" ou "Usuário desativado"  
**Prioridade:** CRÍTICA

### TP-012: Excluir Empresa - Sucesso
**Objetivo:** Validar exclusão de empresa (soft delete)  
**Pré-condições:** Empresa criada sem dados importantes  
**Passos:**
1. Acessar `/platform/companies/{id}`
2. Clicar em "Excluir"
3. Confirmar ação
**Resultado Esperado:** Empresa excluída (soft delete), não aparece na listagem  
**Prioridade:** ALTA

### TP-013: Excluir Empresa - Com Usuários
**Objetivo:** Validar bloqueio de exclusão com usuários ativos  
**Pré-condições:** Empresa com Owner criado  
**Passos:**
1. Tentar excluir empresa
**Resultado Esperado:** Erro "Empresa possui usuários ativos"  
**Prioridade:** ALTA

### TP-014: Listar Empresas - Filtro Ativas
**Objetivo:** Validar filtro de empresas ativas  
**Pré-condições:** Múltiplas empresas criadas (ativas e suspensas)  
**Passos:**
1. Acessar `/platform/companies`
2. Aplicar filtro "Ativas"
**Resultado Esperado:** Apenas empresas ativas listadas  
**Prioridade:** MÉDIA

### TP-015: Listar Empresas - Filtro Suspensas
**Objetivo:** Validar filtro de empresas suspensas  
**Pré-condições:** Múltiplas empresas criadas (ativas e suspensas)  
**Passos:**
1. Acessar `/platform/companies`
2. Aplicar filtro "Suspensas"
**Resultado Esperado:** Apenas empresas suspensas listadas  
**Prioridade:** MÉDIA

### TP-016: Criar Owner - Sucesso
**Objetivo:** Validar criação de owner para empresa  
**Pré-condições:** Empresa criada sem owner  
**Passos:**
1. Acessar `/platform/companies/{id}/owner`
2. Preencher: Nome, Email, Senha
3. Clicar em "Criar Owner"
**Resultado Esperado:** Owner criado com CompanyID correto, RoleOwner  
**Prioridade:** CRÍTICA

### TP-017: Criar Owner - Email Duplicado
**Objetivo:** Validar prevenção de email duplicado  
**Pré-condições:** Owner criado com email `owner@test.com`  
**Passos:**
1. Tentar criar outro owner com mesmo email
**Resultado Esperado:** Erro "Email já cadastrado"  
**Prioridade:** ALTA

### TP-018: Criar Owner - Empresa Já Tem Owner
**Objetivo:** Validar prevenção de owner duplicado  
**Pré-condições:** Empresa já possui owner  
**Passos:**
1. Tentar criar outro owner para mesma empresa
**Resultado Esperado:** Erro "Empresa já possui owner"  
**Prioridade:** ALTA

### TP-019: Criar Owner - Senha Fraca
**Objetivo:** Validar validação de senha  
**Pré-condições:** Empresa criada  
**Passos:**
1. Tentar criar owner com senha "123"
**Resultado Esperado:** Erro "Senha deve ter no mínimo 6 caracteres"  
**Prioridade:** ALTA

### TP-020: Criar Platform User - Sucesso
**Objetivo:** Validar criação de platform user  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Acessar `/platform/users`
2. Clicar em "Criar Usuário"
3. Preencher: Nome, Email, Senha, Role (PlatformSupport)
4. Clicar em "Salvar"
**Resultado Esperado:** Platform user criado  
**Prioridade:** ALTA

### TP-021: Criar Platform User - Role Inválido
**Objetivo:** Validar restrição de roles de plataforma  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Tentar criar platform user com RoleOwner
**Resultado Esperado:** Erro "Role inválido para platform user"  
**Prioridade:** ALTA

### TP-022: Listar Platform Users
**Objetivo:** Validar listagem de platform users  
**Pré-condições:** Múltiplos platform users criados  
**Passos:**
1. Acessar `/platform/users`
**Resultado Esperado:** Todos os platform users listados  
**Prioridade:** MÉDIA

### TP-023: Editar Platform User - Sucesso
**Objetivo:** Validar edição de platform user  
**Pré-condições:** Platform user criado  
**Passos:**
1. Acessar `/platform/users/{id}`
2. Modificar nome
3. Clicar em "Salvar"
**Resultado Esperado:** Platform user atualizado  
**Prioridade:** MÉDIA

### TP-024: Desativar Platform User - Sucesso
**Objetivo:** Validar desativação de platform user  
**Pré-condições:** Platform user criado  
**Passos:**
1. Acessar `/platform/users/{id}`
2. Clicar em "Desativar"
3. Confirmar
**Resultado Esperado:** Platform user desativado  
**Prioridade:** ALTA

### TP-025: Desativar Platform User - Próprio Usuário
**Objetivo:** Validar prevenção de auto-desativação  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Tentar desativar a si mesmo
**Resultado Esperado:** Erro "Não é possível desativar o próprio usuário"  
**Prioridade:** CRÍTICA

### TP-026: Excluir Platform User - Sucesso
**Objetivo:** Validar exclusão de platform user  
**Pré-condições:** Platform user criado  
**Passos:**
1. Acessar `/platform/users/{id}`
2. Clicar em "Excluir"
3. Confirmar
**Resultado Esperado:** Platform user excluído (soft delete)  
**Prioridade:** MÉDIA

### TP-027: Excluir Platform User - Único Admin
**Objetivo:** Validar prevenção de exclusão do único admin  
**Pré-condições:** Apenas um PlatformAdmin existe  
**Passos:**
1. Tentar excluir o único PlatformAdmin
**Resultado Esperado:** Erro "Não é possível excluir o único administrador da plataforma"  
**Prioridade:** CRÍTICA

### TP-028: Dashboard Plataforma - Métricas
**Objetivo:** Validar exibição de métricas no dashboard  
**Pré-condições:** Múltiplas empresas criadas  
**Passos:**
1. Acessar `/platform`
**Resultado Esperado:** Métricas exibidas (total empresas, ativas, suspensas, usuários)  
**Prioridade:** MÉDIA

### TP-029: Acesso Não Autorizado - Sem Autenticação
**Objetivo:** Validar bloqueio de acesso sem autenticação  
**Pré-condições:** Nenhuma  
**Passos:**
1. Tentar acessar `/platform/companies` sem login
**Resultado Esperado:** Redirecionado para `/platform/auth`  
**Prioridade:** CRÍTICA

### TP-030: Acesso Não Autorizado - PlatformSupport
**Objetivo:** Validar que PlatformSupport não pode criar empresas  
**Pré-condições:** PlatformSupport logado  
**Passos:**
1. Tentar acessar `/platform/companies/create`
**Resultado Esperado:** Erro "Permissão negada"  
**Prioridade:** ALTA

### TP-031: Acesso Não Autorizado - Company User
**Objetivo:** Validar que company users não acessam /platform  
**Pré-condições:** Owner logado  
**Passos:**
1. Tentar acessar `/platform/companies`
**Resultado Esperado:** Erro "Permissão negada" ou redirecionamento  
**Prioridade:** CRÍTICA

### TP-032: Impersonation - Futuro
**Objetivo:** Validar funcionalidade de impersonation (futuro)  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Clicar em "Entrar como empresa" em detalhes da empresa
**Resultado Esperado:** Funcionalidade não implementada ou placeholder  
**Prioridade:** BAIXA

### TP-033: Buscar Empresa por ID
**Objetivo:** Validar busca de empresa específica  
**Pré-condições:** Empresa criada  
**Passos:**
1. Acessar `/platform/companies/{id}`
**Resultado Esperado:** Detalhes da empresa exibidos  
**Prioridade:** MÉDIA

### TP-034: Buscar Empresa Inexistente
**Objetivo:** Validar tratamento de empresa inexistente  
**Pré-condições:** Nenhuma  
**Passos:**
1. Acessar `/platform/companies/99999`
**Resultado Esperado:** Erro "Empresa não encontrada"  
**Prioridade:** MÉDIA

### TP-035: Validação de CNPJ
**Objetivo:** Validar formato de CNPJ  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Tentar criar empresa com CNPJ inválido
**Resultado Esperado:** Erro "CNPJ inválido"  
**Prioridade:** MÉDIA

---

## Testes de Empresa

### TE-001: Login Empresa - Sucesso
**Objetivo:** Validar login de Owner com credenciais corretas  
**Pré-condições:** Owner criado pela plataforma  
**Passos:**
1. Acessar `/api/auth/login`
2. Preencher email do owner
3. Preencher senha do owner
4. Clicar em "Login"
**Resultado Esperado:** Login bem-sucedido, redirecionado para `/dashboard`  
**Prioridade:** CRÍTICA

### TE-002: Login Empresa - Usuário Inativo
**Objetivo:** Validar bloqueio de login para usuário inativo  
**Pré-condições:** Owner desativado  
**Passos:**
1. Tentar login com credenciais de owner inativo
**Resultado Esperado:** Erro "Usuário desativado"  
**Prioridade:** ALTA

### TE-003: Login Empresa - Empresa Suspensa
**Objetivo:** Validar bloqueio de login para empresa suspensa  
**Pré-condições:** Empresa suspensa pela plataforma  
**Passos:**
1. Tentar login com credenciais de owner de empresa suspensa
**Resultado Esperado:** Erro "Empresa suspensa"  
**Prioridade:** CRÍTICA

### TE-004: Login Empresa - Platform User
**Objetivo:** Validar que platform users não acessam empresa  
**Pré-condições:** PlatformAdmin criado  
**Passos:**
1. Tentar login com credenciais de PlatformAdmin em `/api/auth/login`
**Resultado Esperado:** Erro "Credenciais inválidas"  
**Prioridade:** CRÍTICA

### TE-005: Logout Empresa
**Objetivo:** Validar logout e invalidação de token  
**Pré-condições:** Owner logado  
**Passos:**
1. Clicar em "Logout"
2. Tentar acessar `/dashboard`
**Resultado Esperado:** Redirecionado para login, token invalidado  
**Prioridade:** ALTA

### TE-006: Criar Usuário Diretamente - Sucesso
**Objetivo:** Validar criação de usuário por Owner  
**Pré-condições:** Owner logado  
**Passos:**
1. Acessar `/api/company/users`
2. Clicar em "Criar Usuário"
3. Preencher: Nome, Email, Senha, Role (Admin)
4. Clicar em "Salvar"
**Resultado Esperado:** Usuário criado com CompanyID do Owner  
**Prioridade:** CRÍTICA

### TE-007: Criar Usuário - Email Duplicado na Empresa
**Objetivo:** Validar prevenção de email duplicado  
**Pré-condições:** Usuário com email `user@test.com` existe na empresa  
**Passos:**
1. Tentar criar outro usuário com mesmo email
**Resultado Esperado:** Erro "Email já cadastrado na empresa"  
**Prioridade:** ALTA

### TE-008: Criar Usuário - Email em Outra Empresa
**Objetivo:** Validar que email em outra empresa é permitido  
**Pré-condições:** Usuário com email `user@test.com` em outra empresa  
**Passos:**
1. Criar usuário com mesmo email em empresa diferente
**Resultado Esperado:** Usuário criado com sucesso (emails são globais, mas CompanyID diferente)  
**Prioridade:** ALTA

### TE-009: Criar Usuário - Role Inválido
**Objetivo:** Validar restrição de roles  
**Pré-condições:** Owner logado  
**Passos:**
1. Tentar criar usuário com RolePlatformAdmin
**Resultado Esperado:** Erro "Role inválido"  
**Prioridade:** ALTA

### TE-010: Alterar Role - Sucesso
**Objetivo:** Validar alteração de role por Owner  
**Pré-condições:** Admin criado  
**Passos:**
1. Acessar `/api/company/users/{id}`
2. Alterar role para Manager
3. Clicar em "Salvar"
**Resultado Esperado:** Role alterada com sucesso  
**Prioridade:** ALTA

### TE-011: Alterar Role - Owner por Admin
**Objetivo:** Validar que Admin não pode alterar Owner  
**Pré-condições:** Admin logado, Owner existe  
**Passos:**
1. Admin tentar alterar role de Owner
**Resultado Esperado:** Erro "Apenas Owner pode alterar papel de Owner"  
**Prioridade:** CRÍTICA

### TE-012: Alterar Role - Admin por Admin
**Objetivo:** Validar que Admin não pode alterar Admin  
**Pré-condições:** Admin logado, outro Admin existe  
**Passos:**
1. Admin tentar alterar role de outro Admin
**Resultado Esperado:** Erro "Apenas Owner pode alterar papel de Admin"  
**Prioridade:** CRÍTICA

### TE-013: Desativar Usuário - Sucesso
**Objetivo:** Validar desativação de usuário  
**Pré-condições:** Manager criado  
**Passos:**
1. Acessar `/api/company/users/{id}`
2. Clicar em "Desativar"
3. Confirmar
**Resultado Esperado:** Usuário desativado  
**Prioridade:** ALTA

### TE-014: Desativar Usuário - Owner
**Objetivo:** Validar que Owner não pode ser desativado  
**Pré-condições:** Owner existe  
**Passos:**
1. Tentar desativar Owner
**Resultado Esperado:** Erro "Não é possível desativar Owner da empresa"  
**Prioridade:** CRÍTICA

### TE-015: Desativar Usuário - Próprio Usuário
**Objetivo:** Validar prevenção de auto-desativação  
**Pré-condições:** Admin logado  
**Passos:**
1. Admin tentar desativar a si mesmo
**Resultado Esperado:** Erro "Não é possível desativar o próprio usuário"  
**Prioridade:** CRÍTICA

### TE-016: Remover Usuário - Sucesso
**Objetivo:** Validar remoção de usuário da empresa  
**Pré-condições:** Manager criado  
**Passos:**
1. Acessar `/api/company/users/{id}`
2. Clicar em "Remover"
3. Confirmar
**Resultado Esperado:** Usuário removido (CompanyID NULL não permitido, usuário deve ser excluído ou reatribuído)  
**Prioridade:** ALTA

### TE-017: Remover Usuário - Owner
**Objetivo:** Validar que Owner não pode ser removido  
**Pré-condições:** Owner existe  
**Passos:**
1. Tentar remover Owner
**Resultado Esperado:** Erro "Não é possível remover Owner da empresa"  
**Prioridade:** CRÍTICA

### TE-018: Listar Usuários - Filtragem por Empresa
**Objetivo:** Validar que listagem apenas mostra usuários da empresa  
**Pré-condições:** Múltiplas empresas com usuários  
**Passos:**
1. Owner da Empresa A acessar `/api/company/users`
**Resultado Esperado:** Apenas usuários da Empresa A listados  
**Prioridade:** CRÍTICA

### TE-019: Criar Convite - Sucesso
**Objetivo:** Validar criação de convite  
**Pré-condições:** Owner logado  
**Passos:**
1. Acessar `/api/company/invitations`
2. Clicar em "Enviar Convite"
3. Preencher: Email, Role
4. Clicar em "Enviar"
**Resultado Esperado:** Convite criado com token  
**Prioridade:** ALTA

### TE-020: Aceitar Convite - Usuário Não Cadastrado
**Objetivo:** Validar aceitação de convite por usuário não cadastrado  
**Pré-condições:** Convite criado  
**Passos:**
1. Usuário acessar link do convite
2. Definir senha
3. Clicar em "Aceitar"
**Resultado Esperado:** Usuário criado com CompanyID da empresa  
**Prioridade:** CRÍTICA

### TE-021: Aceitar Convite - Usuário Cadastrado Sem Empresa
**Objetivo:** Validar aceitação por usuário cadastrado sem empresa (não deve existir mais)  
**Pré-condições:** Convite criado, usuário cadastrado (CompanyID NULL não permitido)  
**Passos:**
1. Este teste não é mais aplicável (usuários sem empresa não existem)  
**Resultado Esperado:** N/A  
**Prioridade:** N/A

### TE-022: Revogar Convite - Sucesso
**Objetivo:** Validar revogação de convite  
**Pré-condições:** Convite criado  
**Passos:**
1. Acessar `/api/company/invitations/{id}`
2. Clicar em "Revogar"
3. Confirmar
**Resultado Esperado:** Convite revogado  
**Prioridade:** ALTA

### TE-023: Aceitar Convite - Já Aceito
**Objetivo:** Validar prevenção de reuso de convite  
**Pré-condições:** Convite já aceito  
**Passos:**
1. Tentar aceitar mesmo convite novamente
**Resultado Esperado:** Erro "Convite já utilizado"  
**Prioridade:** ALTA

### TE-024: Aceitar Convite - Expirado
**Objetivo:** Validar prevenção de aceitação de convite expirado  
**Pré-condições:** Convite expirado  
**Passos:**
1. Tentar aceitar convite expirado
**Resultado Esperado:** Erro "Convite expirado"  
**Prioridade:** ALTA

### TE-025: Aceitar Convite - Email Incorreto
**Objetivo:** Validar que apenas email do convite pode aceitar  
**Pré-condições:** Convite para `invitee@test.com`  
**Passos:**
1. Usuário com email diferente tentar aceitar
**Resultado Esperado:** Erro "O convite não pertence a este usuário"  
**Prioridade:** ALTA

### TE-026: Configurações da Empresa - Sucesso
**Objetivo:** Validar atualização de configurações  
**Pré-condições:** Owner logado  
**Passos:**
1. Acessar `/api/company/settings`
2. Modificar nome, cores
3. Clicar em "Salvar"
**Resultado Esperado:** Configurações atualizadas  
**Prioridade:** ALTA

### TE-027: Tema - Persistência
**Objetivo:** Validar persistência de cores do tema  
**Pré-condições:** Owner logado  
**Passos:**
1. Alterar cor primária para `#ff0000`
2. Fazer logout
3. Login novamente
4. Verificar `/api/theme`
**Resultado Esperado:** Cor persistida  
**Prioridade:** ALTA

### TE-028: Validação de CompanyID
**Objetivo:** Validar que CompanyID nunca é NULL  
**Pré-condições:** Nenhuma  
**Passos:**
1. Verificar banco: `SELECT * FROM users WHERE company_id IS NULL`
**Resultado Esperado:** Nenhum resultado  
**Prioridade:** CRÍTICA

---

## Testes de Regressão

### TR-001: Produto - Criar
**Objetivo:** Validar criação de produto  
**Pré-condições:** Owner logado  
**Passos:**
1. Acessar `/api/products`
2. Criar produto com nome e preço
**Resultado Esperado:** Produto criado com slug gerado  
**Prioridade:** CRÍTICA

### TR-002: Produto - Editar
**Objetivo:** Validar edição de produto  
**Pré-condições:** Produto criado  
**Passos:**
1. Editar nome do produto
**Resultado Esperado:** Produto atualizado, slug atualizado  
**Prioridade:** CRÍTICA

### TR-003: Produto - Excluir (Soft Delete)
**Objetivo:** Validar soft delete de produto  
**Pré-condições:** Produto criado  
**Passos:**
1. Excluir produto
2. Verificar banco
**Resultado Esperado:** Produto com deleted_at preenchido, não aparece na listagem  
**Prioridade:** CRÍTICA

### TR-004: Produto - Slug Único
**Objetivo:** Validar geração de slug único  
**Pré-condições:** Produto com slug "produto-teste" existe  
**Passos:**
1. Criar outro produto com mesmo nome
**Resultado Esperado:** Slug diferente (ex: "produto-teste-2")  
**Prioridade:** ALTA

### TR-005: Ingrediente - Criar
**Objetivo:** Validar criação de ingrediente  
**Pré-condições:** Owner logado  
**Passos:**
1. Acessar `/api/ingredients`
2. Criar ingrediente com nome e estoque
**Resultado Esperado:** Ingrediente criado  
**Prioridade:** CRÍTICA

### TR-006: Ingrediente - Ajustar Estoque
**Objetivo:** Validar ajuste de estoque  
**Pré-condições:** Ingrediente criado  
**Passos:**
1. Ajustar estoque para 20
**Resultado Esperado:** Estoque atualizado  
**Prioridade:** CRÍTICA

### TR-007: Ingrediente - Excluir (Soft Delete)
**Objetivo:** Validar soft delete de ingrediente  
**Pré-condições:** Ingrediente criado  
**Passos:**
1. Excluir ingrediente
2. Verificar banco
**Resultado Esperado:** Ingrediente com deleted_at preenchido  
**Prioridade:** CRÍTICA

### TR-008: Pedido - Criar
**Objetivo:** Validar criação de pedido  
**Pré-condições:** Produto criado, Owner logado  
**Passos:**
1. Acessar `/api/orders`
2. Criar pedido com itens
**Resultado Esperado:** Pedido criado com itens  
**Prioridade:** CRÍTICA

### TR-009: Pedido - Alterar Status
**Objetivo:** Validar alteração de status  
**Pré-condições:** Pedido criado  
**Passos:**
1. Alterar status para "confirmed"
**Resultado Esperado:** Status atualizado  
**Prioridade:** CRÍTICA

### TR-010: Pedido - Validação de Transição
**Objetivo:** Validar transição de status inválida  
**Pré-condições:** Pedido com status "cancelled"  
**Passos:**
1. Tentar alterar para "confirmed"
**Resultado Esperado:** Erro "Transição de status inválida"  
**Prioridade:** ALTA

### TR-011: Pedido - Validação de Estoque
**Objetivo:** Validar validação de estoque ao criar pedido  
**Pré-condições:** Produto com estoque 5  
**Passos:**
1. Tentar criar pedido com quantidade 10
**Resultado Esperado:** Erro "Estoque insuficiente"  
**Prioridade:** CRÍTICA

### TR-012: RBAC - Owner Permissões
**Objetivo:** Validar permissões de Owner  
**Pré-condições:** Owner logado  
**Passos:**
1. Tentar acessar `/api/company/users`
2. Tentar criar usuário
3. Tentar alterar role
**Resultado Esperado:** Todas as ações permitidas  
**Prioridade:** CRÍTICA

### TR-013: RBAC - Admin Permissões
**Objetivo:** Validar permissões de Admin  
**Pré-condições:** Admin logado  
**Passos:**
1. Tentar acessar `/api/company/users`
2. Tentar criar usuário
3. Tentar alterar role de Manager
4. Tentar alterar role de Owner
**Resultado Esperado:** 1, 2, 3 permitidos; 4 negado  
**Prioridade:** CRÍTICA

### TR-014: RBAC - Manager Permissões
**Objetivo:** Validar permissões de Manager  
**Pré-condições:** Manager logado  
**Passos:**
1. Tentar acessar `/api/company/users`
2. Tentar criar usuário
**Resultado Esperado:** Ambos negados  
**Prioridade:** CRÍTICA

### TR-015: RBAC - Employee Permissões
**Objetivo:** Validar permissões de Employee  
**Pré-condições:** Employee logado  
**Passos:**
1. Tentar acessar `/api/company/users`
2. Tentar criar produto
**Resultado Esperado:** Ambos negados (Employee tem acesso limitado)  
**Prioridade:** CRÍTICA

### TR-016: JWT - Geração
**Objetivo:** Validar geração de JWT  
**Pré-condições:** Owner logado  
**Passos:**
1. Fazer login
2. Verificar cookie auth_token
**Resultado Esperado:** Token JWT válido gerado  
**Prioridade:** CRÍTICA

### TR-017: JWT - Validação
**Objetivo:** Validar validação de JWT  
**Pré-condições:** Token JWT válido  
**Passos:**
1. Fazer requisição com token
**Resultado Esperado:** Requisição autorizada  
**Prioridade:** CRÍTICA

### TR-018: JWT - Expiração
**Objetivo:** Validar expiração de JWT  
**Pré-condições:** Token JWT expirado  
**Passos:**
1. Tentar fazer requisição com token expirado
**Resultado Esperado:** Erro "Token expirado"  
**Prioridade:** ALTA

### TR-019: JWT - Blacklist
**Objetivo:** Validar blacklist de JWT  
**Pré-condições:** Usuário logado  
**Passos:**
1. Fazer logout
2. Tentar usar token anterior
**Resultado Esperado:** Erro "Token revogado"  
**Prioridade:** CRÍTICA

### TR-020: Soft Delete - Produtos
**Objetivo:** Validar soft delete em produtos  
**Pré-condições:** Produto criado  
**Passos:**
1. Excluir produto
2. Verificar banco: `SELECT * FROM products WHERE deleted_at IS NOT NULL`
**Resultado Esperado:** Produto com deleted_at preenchido  
**Prioridade:** CRÍTICA

### TR-021: Soft Delete - Ingredientes
**Objetivo:** Validar soft delete em ingredientes  
**Pré-condições:** Ingrediente criado  
**Passos:**
1. Excluir ingrediente
2. Verificar banco
**Resultado Esperado:** Ingrediente com deleted_at preenchido  
**Prioridade:** CRÍTICA

### TR-022: Soft Delete - Usuários
**Objetivo:** Validar soft delete em usuários  
**Pré-condições:** Usuário criado  
**Passos:**
1. Excluir usuário
2. Verificar banco
**Resultado Esperado:** Usuário com deleted_at preenchido  
**Prioridade:** CRÍTICA

### TR-023: Isolamento de Dados - CompanyID
**Objetivo:** Validar isolamento por CompanyID  
**Pré-condições:** Empresa A e Empresa B criadas  
**Passos:**
1. Owner da Empresa A criar produto
2. Owner da Empresa B listar produtos
**Resultado Esperado:** Produto da Empresa A não aparece na listagem da Empresa B  
**Prioridade:** CRÍTICA

### TR-024: Tema - Persistência Após Reload
**Objetivo:** Validar persistência de tema após reload da página  
**Pré-condições:** Tema configurado  
**Passos:**
1. Alterar cor primária
2. Recarregar página
3. Verificar `/api/theme`
**Resultado Esperado:** Cor persistida  
**Prioridade:** ALTA

---

## Testes de Integração

### TI-001: Fluxo Completo - Plataforma Cria Empresa e Owner
**Objetivo:** Validar fluxo completo de onboarding  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Criar empresa
2. Criar owner para empresa
3. Owner fazer login
4. Owner criar usuário
5. Usuário fazer login
**Resultado Esperado:** Fluxo completo sem erros  
**Prioridade:** CRÍTICA

### TI-002: Fluxo Completo - Convite
**Objetivo:** Validar fluxo completo de convite  
**Pré-condições:** Owner logado  
**Passos:**
1. Owner criar convite
2. Usuário receber email
3. Usuário aceitar convite
4. Usuário fazer login
**Resultado Esperado:** Fluxo completo sem erros  
**Prioridade:** CRÍTICA

### TI-003: Fluxo Completo - Suspensão de Empresa
**Objetivo:** Validar impacto de suspensão em usuários  
**Pré-condições:** Empresa com múltiplos usuários  
**Passos:**
1. PlatformAdmin suspender empresa
2. Usuários tentarem fazer login
**Resultado Esperado:** Todos os usuários bloqueados  
**Prioridade:** CRÍTICA

---

## Testes de Performance

### TPERF-001: Listagem de Empresas - 1000 Empresas
**Objetivo:** Validar performance de listagem com muitos dados  
**Pré-condições:** 1000 empresas criadas  
**Passos:**
1. Acessar `/platform/companies`
**Resultado Esperado:** Listagem carrega em < 2 segundos  
**Prioridade:** MÉDIA

### TPERF-002: Listagem de Usuários - 500 Usuários
**Objetivo:** Validar performance de listagem de usuários  
**Pré-condições:** 500 usuários na empresa  
**Passos:**
1. Acessar `/api/company/users`
**Resultado Esperado:** Listagem carrega em < 2 segundos  
**Prioridade:** MÉDIA

---

## Testes de Segurança

### TSEC-001: SQL Injection - Email
**Objetivo:** Validar proteção contra SQL injection  
**Pré-condições:** Nenhuma  
**Passos:**
1. Tentar login com email: `' OR '1'='1`
**Resultado Esperado:** Erro de credenciais, não erro de SQL  
**Prioridade:** CRÍTICA

### TSEC-002: XSS - Nome da Empresa
**Objetivo:** Validar proteção contra XSS  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Criar empresa com nome: `<script>alert('xss')</script>`
2. Listar empresas
**Resultado Esperado:** Script não executado, texto escapado  
**Prioridade:** CRÍTICA

### TSEC-003: CSRF - Criar Empresa
**Objetivo:** Validar proteção contra CSRF  
**Pré-condições:** PlatformAdmin logado  
**Passos:**
1. Enviar POST para `/platform/companies` sem token CSRF
**Resultado Esperado:** Erro "CSRF token inválido"  
**Prioridade:** ALTA

---

## Critérios de Aceite

### Critérios Gerais
- Todos os testes CRÍTICOS devem passar
- Pelo menos 90% dos testes ALTA prioridade devem passar
- Pelo menos 80% dos testes MÉDIA prioridade devem passar
- Zero vulnerabilidades de segurança CRÍTICAS

### Critérios Específicos
- **TP-004:** Company users não acessam plataforma
- **TP-011:** Usuários de empresa suspensa não acessam
- **TP-018:** Owner criado com CompanyID correto
- **TP-025:** Auto-desativação bloqueada
- **TP-027:** Único admin não pode ser excluído
- **TE-003:** Empresa suspensa bloqueia login
- **TE-004:** Platform users não acessam empresa
- **TE-011:** Admin não altera Owner
- **TE-014:** Owner não pode ser desativado
- **TE-028:** CompanyID nunca NULL
- **TR-023:** Isolamento por CompanyID
- **TI-001:** Fluxo completo de onboarding

---

## Plano de Execução

### Fase 1: Preparação (Dia 1)
- Configurar ambiente de testes
- Criar dados de teste (platform admin, empresas, usuários)
- Configurar ferramentas de automação (se aplicável)

### Fase 2: Testes de Plataforma (Dia 2)
- Executar TP-001 a TP-035
- Documentar resultados
- Corrigir bugs encontrados

### Fase 3: Testes de Empresa (Dia 3)
- Executar TE-001 a TE-028
- Documentar resultados
- Corrigir bugs encontrados

### Fase 4: Testes de Regressão (Dia 4)
- Executar TR-001 a TR-024
- Documentar resultados
- Corrigir bugs encontrados

### Fase 5: Testes de Integração (Dia 5)
- Executar TI-001 a TI-003
- Documentar resultados
- Corrigir bugs encontrados

### Fase 6: Testes de Performance e Segurança (Dia 5)
- Executar TPERF-001 a TPERF-002
- Executar TSEC-001 a TSEC-003
- Documentar resultados
- Corrigir bugs encontrados

### Fase 7: Validação Final (Dia 6)
- Re-executar testes CRÍTICOS
- Validar critérios de aceite
- Gerar relatório final

---

## Ferramentas de Teste

### Backend
- **Postman/Insomnia:** Testes manuais de API
- **curl:** Testes automatizados via script
- **sqlite3:** Validação de banco de dados

### Frontend
- **Playwright:** Testes E2E automatizados (futuro)
- **Testes manuais:** Validação visual e funcional

### Performance
- **Apache Benchmark (ab):** Testes de carga
- **Go pprof:** Profiling de performance

### Segurança
- **OWASP ZAP:** Varredura de vulnerabilidades
- **SQLMap:** Testes de SQL injection

---

## Relatório de Testes

Após execução, gerar relatório com:
- Status de cada teste (PASS/FAIL)
- Screenshots de evidências (quando aplicável)
- Logs de erros
- Tempo de execução
- Taxa de sucesso por categoria
- Lista de bugs encontrados
- Recomendações

---

## Conclusão

Este plano de testes abrange todos os aspectos da refatoração arquitetural, garantindo que:
- Novas funcionalidades de plataforma funcionam corretamente
- Funcionalidades existentes continuam operando
- Regras de negócio são respeitadas
- Segurança e performance são mantidas
- Isolamento entre níveis é garantido

A execução sistemática destes testes assegura uma migração segura e confiável para o modelo SaaS empresarial.
