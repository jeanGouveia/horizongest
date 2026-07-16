# Layout - PratoOnline

## Visão Geral

O layout do PratoOnline foi completamente redesenhado para seguir o padrão de dashboards modernos inspirados no Autumn CRM. A estrutura abandona o layout tradicional em favor de uma organização mais espacial e profissional.

## Estrutura Principal

### Grid de Layout

```
┌─────────────────────────────────────────────────────────┐
│                    Header Horizontal                     │
├──────────┬──────────────────────────────────────────────┤
│          │                                              │
│  Sidebar │              Conteúdo Central                 │
│          │                                              │
│          │                                              │
│          │                                              │
└──────────┴──────────────────────────────────────────────┘
```

### Dimensões

**Sidebar**
- Expandida: 260px
- Recolhida: 72px
- Mobile: Oculta (menu hamburger)

**Header**
- Altura: 64px
- Padding horizontal: 24px
- Sticky no topo

**Conteúdo Central**
- Max-width: 1400px (em telas grandes)
- Padding: 24px
- Margem automática para centralizar

## Header Horizontal

### Composição

```
┌─────────────────────────────────────────────────────────┐
│ 🍽️ PratoOnline  [Breadcrumb]  [Busca]  [⌘K] [🔔] [🌙] [👤] │
└─────────────────────────────────────────────────────────┘
```

### Elementos

**1. Logo e Nome do Sistema**
- Posição: Esquerda
- Logo: Ícone + texto "PratoOnline"
- Link: Redireciona para Dashboard
- Gradiente sutil no texto

**2. Breadcrumb (Opcional)**
- Posição: Centro-esquerda
- Mostra navegação atual
- Separador: "/"
- Item ativo em negrito
- Oculto em mobile

**3. Campo de Busca**
- Posição: Centro
- Trigger: Botão "Buscar" com ícone
- Expande para input completo ao clicar
- Placeholder: "Buscar pedidos, produtos..."
- Atalho: ⌘K (placeholder)

**4. Ações Rápidas**
- Posição: Direita
- Atalhos (⌘K)
- Notificações (🔔) - placeholder
- Tema claro/escuro (🌙) - placeholder

**5. Perfil do Usuário**
- Posição: Extrema direita
- Avatar com inicial ou foto
- Nome do usuário
- Dropdown com:
  - Perfil
  - Configurações
  - Sair

### Comportamento

- Sticky no topo da página
- Z-index: 100
- Sombra sutil: `0 1px 2px 0 rgb(0 0 0 / 0.05)`
- Transições suaves em hover

### Responsividade

**Desktop (> 1024px)**
- Layout completo
- Breadcrumb visível
- Busca com texto

**Tablet (768px - 1024px)**
- Breadcrumb oculto
- Busca apenas ícone
- Nome do usuário oculto

**Mobile (< 768px)**
- Breadcrumb oculto
- Busca apenas ícone
- Ações reduzidas
- Nome do usuário oculto

## Sidebar Moderna

### Composição

```
┌────────────┐
│  [←] Navegação
├────────────┤
│ PRINCIPAL  │
│ 📊 Dashboard
│            │
│ OPERAÇÃO   │
│ 📋 Pedidos (3)
│ 🍔 Produtos
│ 🥬 Ingredientes
│            │
│ ESTOQUE    │
│ ⚖️ Ajustes (2)
│            │
│ ADMINISTRAÇÃO│
│ 👤 Perfil
├────────────┤
│      🚪 Sair│
└────────────┘
```

### Estrutura de Menus

**Grupos**
1. **Principal**
   - Dashboard

2. **Operação**
   - Pedidos (com badge de contagem)
   - Produtos
   - Ingredientes

3. **Estoque**
   - Ajustes (com badge de contagem)

4. **Administração**
   - Perfil
   - Configurações (placeholder)

### Elementos

**1. Header da Sidebar**
- Botão de toggle (←/→)
- Título "Navegação" (quando expandida)
- Borda inferior sutil

**2. Itens de Menu**
- Ícone (emoji ou Lucide icon)
- Label do item
- Badge opcional (contagem)
- Indicador de página ativa

**3. Indicador de Página Ativa**
- Background: cor primária (#6366f1)
- Texto: branco
- Sombra sutil
- Transição suave

**4. Footer**
- Link "Sair" com ícone
- Cor de destaque (vermelho)
- Hover com background avermelhado

### Comportamento

- Sticky na esquerda
- Height: 100vh
- Overflow-y: auto
- Transição de width: 0.3s ease
- Hover nos itens: background cinza claro

### Estados

**Expandido**
- Width: 260px
- Mostra labels
- Mostra grupos
- Badge integrado

**Recolhido**
- Width: 72px
- Apenas ícones
- Tooltip no hover
- Badge flutuante

### Responsividade

**Desktop (> 1024px)**
- Sidebar expandida por padrão
- Toggle funcional

**Tablet (768px - 1024px)**
- Sidebar recolhida por padrão
- Toggle funcional

**Mobile (< 768px)**
- Sidebar oculta
- Menu hamburger no header
- Drawer lateral ao abrir

## Conteúdo Central

### Grid System

**Breakpoints**
- Mobile: < 768px (1 coluna)
- Tablet: 768px - 1024px (2 colunas)
- Desktop: 1024px - 1280px (3 colunas)
- Large: > 1280px (4 colunas)

**Cards**
- Padding: 1.5rem
- Border-radius: 12px
- Background: branco
- Sombra: `0 4px 6px -1px rgb(0 0 0 / 0.1)`
- Hover: translateY(-2px)

**Espaçamento**
- Entre cards: 1rem
- Entre seções: 2rem
- Margem lateral: auto

### Seções Típicas

**Dashboard**
- Row 1: 5 cards de métricas (grid)
- Row 2: Gráfico de vendas (full width)
- Row 3: Produtos mais vendidos + Alertas (2 colunas)
- Row 4: Ações rápidas (4 colunas)

**Páginas de Lista**
- Header da página
- Cards de resumo (opcional)
- Filtros elegantes
- Tabela moderna
- Paginação

**Páginas de Formulário**
- Header da página
- Breadcrumb
- Formulário em cards
- Agrupamento lógico
- Botões de ação

## Cards Executivos

### Estrutura

```
┌─────────────────────┐
│ [Ícone] Label       │
│                     │
│     Valor           │
│   Variação %        │
└─────────────────────┘
```

### Variações

**Métrica Primária**
- Cor: Índigo
- Ícone: 📊
- Exemplo: Pedidos Hoje

**Métrica de Sucesso**
- Cor: Verde
- Ícone: 💰
- Exemplo: Faturamento

**Métrica de Informação**
- Cor: Azul
- Ícone: 🍽️
- Exemplo: Produtos Ativos

**Métrica de Aviso**
- Cor: Amarelo
- Ícone: ⚠️
- Exemplo: Estoque Baixo

**Métrica de Perigo**
- Cor: Vermelho
- Ícone: 📦
- Exemplo: Ajustes Pendentes

### Comportamento

- Hover: translateY(-2px)
- Sombra aumenta no hover
- Transição suave (0.2s ease)

## Espaçamento e Respiro

### Princípios

- **Muito espaço em branco**: Não tenha medo de espaço vazio
- **Respiro**: Deixe o conteúdo "respirar"
- **Hierarquia**: Use espaço, não bordas, para criar hierarquia

### Valores

- Padding de cards: 1.5rem
- Margem entre seções: 2rem
- Gap em grids: 1rem
- Padding de containers: 24px

## Navegação

### Breadcrumbs

**Formato**
```
Home / Operação / Pedidos
```

**Estilo**
- Cor: cinza médio
- Separador: "/"
- Item ativo: negrito, cor escura
- Hover: cor primária

### Links

**Estilo Padrão**
- Cor: primária
- Text-decoration: none
- Hover: underline sutil
- Transição: 0.2s ease

**Links de Ação**
- Button-like appearance
- Hover com background
- Ícone opcional

## Responsividade

### Estratégia

**Mobile-First**
- Começar com layout mobile
- Adicionar complexidade em breakpoints maiores
- Progressividade: mobile → tablet → desktop

### Breakpoints

```css
/* Mobile */
@media (max-width: 767px) { }

/* Tablet */
@media (min-width: 768px) and (max-width: 1023px) { }

/* Desktop */
@media (min-width: 1024px) and (max-width: 1279px) { }

/* Large Desktop */
@media (min-width: 1280px) { }
```

### Adaptações

**Sidebar**
- Desktop: Expandida
- Tablet: Recolhida
- Mobile: Oculta (drawer)

**Header**
- Desktop: Completo
- Tablet: Reduzido
- Mobile: Minimalista

**Grid**
- Mobile: 1 coluna
- Tablet: 2 colunas
- Desktop: 3-4 colunas

**Tabelas**
- Desktop: Tabela completa
- Tablet: Scroll horizontal
- Mobile: Cards alternativos

## Dark Mode

### Preparação

O layout está preparado para dark mode com:

- Variáveis CSS para cores
- Contraste adequado em ambos os temas
- Transições suaves entre temas
- Preservação de hierarquia visual

### Implementação Futura

- Toggle no header
- Persistência de preferência
- Adaptação de sombras
- Ajuste de cores de texto

---

**Versão**: 1.0  
**Data**: 15/07/2026  
**Sprint**: 9 - Product Experience
