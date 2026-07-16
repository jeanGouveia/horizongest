# Components - PratoOnline

## Visão Geral

O sistema de componentes do PratoOnline foi redesenhado para criar uma biblioteca consistente, reutilizável e moderna. Todos os componentes seguem os princípios do design language e são construídos com performance e acessibilidade em mente.

## Princípios de Componentes

### 1. Atomicidade

Cada componente deve:
- Ter uma única responsabilidade
- Ser independente e reutilizável
- Ter props bem definidas
- Ser testável isoladamente

### 2. Consistência

Todos os componentes devem:
- Seguir o mesmo sistema de design tokens
- Usar as mesmas convenções de nomenclatura
- Ter comportamentos previsíveis
- Manter coerência visual

### 3. Componibilidade

Componentes devem:
- Ser combináveis para criar interfaces complexas
- Ter slots/flexibilidade para composição
- Não ter dependências acopladas
- Ser agnósticos ao contexto

### 4. Acessibilidade

Todo componente deve:
- Ter suporte a navegação por teclado
- Ter labels descritivos
- Ter estados de foco visíveis
- Ter contraste adequado
- Suportar screen readers

## Biblioteca de Componentes

### Layout Components

#### PageContainer

Container principal para páginas da aplicação.

**Props**
- `maxWidth?: string` - Largura máxima (default: "1400px")
- `padding?: string` - Padding (default: "24px")

**Uso**
```svelte
<PageContainer>
  <PageHeader title="Título" subtitle="Subtítulo" />
  <!-- Conteúdo -->
</PageContainer>
```

#### PageHeader

Header de página com título e subtítulo.

**Props**
- `title: string` - Título da página
- `subtitle?: string` - Subtítulo opcional
- `actions?: Snippet` - Ações no header (botões, etc)

**Uso**
```svelte
<PageHeader 
  title="Produtos" 
  subtitle="Gerencie seu cardápio"
>
  <Button href="/products/new">Novo Produto</Button>
</PageHeader>
```

#### Header

Header horizontal da aplicação.

**Props**
- `userName?: string` - Nome do usuário
- `userAvatar?: string` - URL do avatar
- `breadcrumb?: string[]` - Array de breadcrumb

**Features**
- Logo com gradiente
- Campo de busca expansível
- Menu de usuário com dropdown
- Responsivo

#### Sidebar

Sidebar de navegação moderna.

**Props**
- `currentPath?: string` - Path atual para highlight
- `collapsed?: boolean` - Estado recolhido
- `onToggle?: () => void` - Callback de toggle

**Features**
- Grupos de navegação
- Badges de contagem
- Indicador de página ativa
- Responsivo (recolhe em tablet)

### UI Components

#### Button

Botão com múltiplas variantes e tamanhos.

**Variantes**
- `primary` - Cor principal, destaque
- `secondary` - Cor secundária
- `ghost` - Transparente com hover
- `danger` - Vermelho para ações destrutivas

**Tamanhos**
- `sm` - Pequeno
- `md` - Médio (default)
- `lg` - Grande

**Props**
- `variant?: 'primary' | 'secondary' | 'ghost' | 'danger'`
- `size?: 'sm' | 'md' | 'lg'`
- `href?: string` - Se presente, renderiza como link
- `disabled?: boolean`
- `loading?: boolean`
- `icon?: string` - Ícone (emoji ou componente)

**Uso**
```svelte
<Button variant="primary" size="lg">Salvar</Button>
<Button variant="ghost" href="/cancel">Cancelar</Button>
<Button variant="danger" onclick={delete}>Excluir</Button>
```

#### Card

Container card com sombra e borda arredondada.

**Props**
- `padding?: string` - Padding interno
- `hover?: boolean` - Efeito de hover
- `border?: boolean` - Mostrar borda

**Uso**
```svelte
<Card padding="1.5rem" hover>
  <h3>Título</h3>
  <p>Conteúdo</p>
</Card>
```

#### Input

Campo de input com estados e validação.

**Variantes**
- `default` - Padrão
- `error` - Estado de erro
- `success` - Estado de sucesso

**Props**
- `type?: string` - Tipo do input
- `placeholder?: string`
- `value?: string`
- `error?: string` - Mensagem de erro
- `label?: string` - Label acima do input
- `helper?: string` - Texto de ajuda abaixo

**Uso**
```svelte
Input 
  label="Nome do Produto"
  placeholder="Digite o nome..."
  bind:value={name}
  error={nameError}
/>
```

#### Select

Select dropdown estilizado.

**Props**
- `options: {value: string, label: string}[]`
- `value?: string`
- `label?: string`
- `placeholder?: string`

**Uso**
```svelte
Select 
  label="Categoria"
  options={categories}
  bind:value={category}
/>
```

#### Textarea

Campo de texto multiline.

**Props**
- `value?: string`
- `placeholder?: string`
- `rows?: number` - Número de linhas (default: 4)
- `label?: string`

#### Checkbox

Checkbox estilizado com label.

**Props**
- `checked?: boolean`
- `label?: string`
- `disabled?: boolean`

**Uso**
```svelte
Checkbox bind:checked={active} label="Produto ativo" />
```

#### Badge

Badge para status e categorias.

**Variantes**
- `default` - Cinza
- `primary` - Índigo
- `success` - Verde
- `warning` - Amarelo
- `danger` - Vermelho
- `info` - Azul

**Props**
- `variant?: 'default' | 'primary' | 'success' | 'warning' | 'danger' | 'info'`
- `size?: 'sm' | 'md'`

**Uso**
```svelte
<Badge variant="success">Ativo</Badge>
<Badge variant="danger">Cancelado</Badge>
<Badge variant="warning">Estoque Baixo</Badge>
```

#### Alert

Componente de alerta com ícone e dismiss.

**Variantes**
- `info` - Azul
- `success` - Verde
- `warning` - Amarelo
- `error` - Vermelho

**Props**
- `variant?: 'info' | 'success' | 'warning' | 'error'`
- `dismissible?: boolean` - Botão de fechar
- `onDismiss?: () => void`

**Uso**
```svelte
<Alert variant="error" dismissible onDismiss={() => error = ''}>
  Erro ao salvar dados
</Alert>
```

#### Modal

Modal/dialog com overlay.

**Props**
- `open: boolean` - Estado aberto/fechado
- `title?: string` - Título do modal
- `onClose?: () => void` - Callback ao fechar
- `size?: 'sm' | 'md' | 'lg'` - Tamanho

**Uso**
```svelte
<Modal bind:open={showModal} title="Confirmar" onClose={() => showModal = false}>
  <p>Tem certeza?</p>
  <Button onclick={confirm}>Confirmar</Button>
</Modal>
```

#### Table

Tabela moderna com hover e ações.

**Props**
- `columns: {key: string, label: string}[]` - Colunas
- `data: any[]` - Dados
- `hover?: boolean` - Efeito hover nas linhas
- `actions?: (row: any) => Snippet` - Ações por linha

**Uso**
```svelte
Table 
  columns={columns}
  data={products}
  hover
  actions={(product) => (
    <Button href={`/products/${product.id}`}>Editar</Button>
  )}
/>
```

#### Divider

Divisor visual horizontal.

**Props**
- `orientation?: 'horizontal' | 'vertical'`
- `margin?: string` - Margem

#### Loading

Componente de loading com skeleton.

**Props**
- `type?: 'spinner' | 'skeleton' | 'dots'`
- `size?: 'sm' | 'md' | 'lg'`

**Uso**
```svelte
<Loading type="skeleton" />
```

#### EmptyState

Estado vazio elegante com call-to-action.

**Props**
- `icon: string` - Ícone (emoji)
- `title: string` - Título
- `description?: string` - Descrição
- `action?: Snippet` - Botão de ação

**Uso**
```svelte
<EmptyState 
  icon="📦"
  title="Nenhum produto"
  description="Comece adicionando seu primeiro produto"
>
  <Button href="/products/new">Adicionar Produto</Button>
</EmptyState>
```

#### ConfirmDialog

Modal de confirmação para ações destrutivas.

**Props**
- `open: boolean`
- `title?: string`
- `message?: string`
- `onConfirm?: () => void`
- `onCancel?: () => void`

### Data Display Components

#### MetricCard

Card especial para métricas executivas.

**Props**
- `icon: string` - Ícone
- `label: string` - Label da métrica
- `value: string | number` - Valor
- `change?: string` - Variação (ex: "+12%")
- `changeType?: 'positive' | 'negative' | 'neutral'`
- `variant?: 'primary' | 'success' | 'info' | 'warning' | 'danger'`

**Uso**
```svelte
<MetricCard 
  icon="📊"
  label="Pedidos Hoje"
  value={42}
  change="+12%"
  changeType="positive"
  variant="primary"
/>
```

#### StatusBadge

Badge específico para status de pedidos/produtos.

**Props**
- `status: string` - Status (pending, confirmed, paid, etc)
- `size?: 'sm' | 'md'`

**Mapeamento de Status**
- `pending` → Amarelo
- `confirmed` → Azul
- `preparing` → Roxo
- `ready` → Verde
- `delivered` → Verde escuro
- `cancelled` → Vermelho
- `paid` → Verde
- `active` → Verde
- `inactive` → Cinza

### Form Components

#### FormField

Wrapper para campos de formulário com label e erro.

**Props**
- `label?: string`
- `error?: string`
- `required?: boolean`
- `helper?: string`

**Uso**
```svelte
<FormField label="Nome" error={nameError} required>
  <Input bind:value={name} />
</FormField>
```

#### FormSection

Seção de formulário com título e borda.

**Props**
- `title: string`
- `description?: string`

**Uso**
```svelte
<FormSection title="Informações Básicas">
  <FormField label="Nome">
    <Input bind:value={name} />
  </FormField>
</FormSection>
```

## Sistema de Slots

Muitos componentes suportam slots para flexibilidade:

```svelte
<Card>
  <slot name="header" /> <!-- Header customizado -->
  <slot /> <!-- Conteúdo principal -->
  <slot name="footer" /> <!-- Footer customizado -->
</Card>
```

## Estados de Componentes

### Estados Comuns

**Default**
- Estado normal do componente

**Hover**
- Mouse sobre o componente
- Feedback visual sutil

**Focus**
- Componente focado (teclado)
- Ring ou borda destacada

**Active**
- Componente sendo clicado
- Feedback de pressão

**Disabled**
- Componente desabilitado
- Opacidade reduzida
- Cursor not-allowed

**Loading**
- Componente carregando
- Spinner ou skeleton
- Interatividade bloqueada

**Error**
- Estado de erro
- Cor vermelha
- Mensagem de erro

**Success**
- Estado de sucesso
- Cor verde
- Feedback positivo

## Animações e Transições

### Durações Padrão

- **Fast**: 150ms - Microinterações rápidas
- **Base**: 200ms - Transições padrão
- **Slow**: 300ms - Transições mais complexas

### Easing Functions

- **Ease-out**: `cubic-bezier(0, 0, 0.2, 1)` - Saída suave
- **Ease-in-out**: `cubic-bezier(0.4, 0, 0.2, 1)` - Entrada e saída
- **Spring**: Simulação de mola (para animações mais naturais)

### Tipos de Animação

**Fade**
- Opacity de 0 → 1
- Para entrada de elementos

**Slide**
- Transform translate
- Para drawers e modais

**Scale**
- Transform scale
- Para botões e cards

**Rotate**
- Transform rotate
- Para ícones e indicadores

## Convenções de Nomenclatura

### Arquivos

- PascalCase para componentes: `Button.svelte`
- Kebab-case para utilitários: `format-currency.ts`

### Props

- camelCase: `userName`, `isLoading`
- Booleanos com prefixo: `is`, `has`, `should`

### Classes CSS

- BEM modificado: `component--modifier`, `component__element`
- Prefixo de componente: `btn-`, `card-`, `input-`

## Performance

### Otimizações

1. **Lazy Loading** → Componentes pesados carregados sob demanda
2. **Code Splitting** → Divisão de código por rota
3. **Virtual Scrolling** → Para listas longas
4. **Memoização** → Evitar re-renders desnecessários
5. **CSS-in-JS** → Estilos críticos inline

### Best Practices

- Evitar props excessivas (max 5-7 props)
- Usar slots para composição complexa
- Componentes pequenos e focados
- Evitar lógica de negócio em componentes UI

## Acessibilidade

### ARIA Attributes

- `aria-label` - Labels descritivos
- `aria-describedby` - Relacionar com helper text
- `aria-invalid` - Estado de invalidade
- `aria-expanded` - Estado expandido/colapsado
- `role` - Role semântico

### Keyboard Navigation

- Tab order lógico
- Enter/Space para ativação
- Escape para cancelar/modais
- Arrow keys para navegação

### Focus Management

- Focus visível (outline ou ring)
- Focus trap em modais
- Focus restoration após fechamento
- Skip links para navegação

## Testabilidade

### Estratégia

- Componentes isolados
- Props controláveis
- Eventos observáveis
- Snapshot testing para UI

### Ferramentas

- Vitest para unit tests
- Playwright para E2E
- Testing Library para componentes

---

**Versão**: 1.0  
**Data**: 15/07/2026  
**Sprint**: 9 - Product Experience
