# Quality Gate - PratoOnline

## Definition of Done

Uma tarefa só pode ser considerada **CONCLUÍDA** quando todas as verificações abaixo passarem:

### Backend (Go)

- [ ] `go fmt ./...` - Formatação do código Go
- [ ] `go vet ./...` - Verificação de erros estáticos
- [ ] `go test ./...` - Execução de testes
- [ ] `go build ./...` - Compilação do projeto

### Frontend (Svelte/TypeScript)

- [ ] `npm install` - Instalação de dependências
- [ ] `npm run check` - Verificação de tipos Svelte/TypeScript
- [ ] `npm run build` - Build de produção
- [ ] `npm run lint` (se disponível) - Verificação de lint

### Verificação de Execução

- [ ] Backend inicia normalmente (sem panic/fatal)
- [ ] Frontend inicia normalmente
- [ ] Console do navegador sem erros
- [ ] Console do backend sem panic/fatal/stacktrace

### Smoke Test Manual

Executar manualmente os fluxos principais:

**Produto**
- [ ] Listar produtos
- [ ] Criar produto
- [ ] Editar produto
- [ ] Ativar produto
- [ ] Desativar produto

**Ingrediente**
- [ ] Listar ingredientes
- [ ] Criar ingrediente
- [ ] Editar ingrediente

**Pedido**
- [ ] Abrir tela de pedidos
- [ ] Adicionar item ao carrinho
- [ ] Remover item do carrinho
- [ ] Confirmar pedido

**Login**
- [ ] Autenticar usuário
- [ ] Logout

**Perfil**
- [ ] Abrir página de perfil
- [ ] Atualizar informações

### Regra Permanente

**NUNCA** responder "Concluído" sem executar TODAS as verificações acima.

**NUNCA** assumir que o código funciona. **SEMPRE** provar que funciona.

### Correções de Erros

Caso qualquer verificação falhe:

1. **NÃO** prosseguir
2. Corrigir o problema
3. Executar novamente as verificações
4. Somente continuar quando TODAS passarem

### GitHub Actions (quando implementado)

O pipeline deve executar automaticamente:

**Backend**
- go fmt
- go vet
- go test
- go build

**Frontend**
- npm install
- npm run check
- npm run build

O pipeline deve falhar caso qualquer etapa falhe.
