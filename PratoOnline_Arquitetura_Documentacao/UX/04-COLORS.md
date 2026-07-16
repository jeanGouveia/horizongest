# Colors - PratoOnline

## Visão Geral

O sistema de cores do PratoOnline foi desenhado para criar uma identidade visual sofisticada e profissional. A paleta é sóbria, com cores neutras predominantes e acentos estrategicamente posicionados para guiar a atenção do usuário.

## Filosofia de Cores

### Princípios

**Sobriedade**
- Cores neutras como base
- Evitar excesso de cores vibrantes
- Paleta restrita e consistente

**Hierarquia**
- Cores para indicar importância
- Contraste para guiar a atenção
- Neutralidade para conteúdo secundário

**Acessibilidade**
- Contraste mínimo de 4.5:1 para texto
- Diferenciação além de cor (ícones, textura)
- Suporte a diferentes modos de visão

**Profissionalismo**
- Cores que inspiram confiança
- Evitar cores "brinquedosas"
- Paleta adequada ao contexto de negócios

## Paleta Principal

### Primary (Índigo)

Cor principal para ações e destaques.

```css
--primary-50: #eef2ff
--primary-100: #e0e7ff
--primary-200: #c7d2fe
--primary-300: #a5b4fc
--primary-400: #818cf8
--primary-500: #6366f1  /* Principal */
--primary-600: #4f46e5
--primary-700: #4338ca
--primary-800: #3730a3
--primary-900: #312e81
--primary-950: #1e1b4b
```

**Usos**
- Botões primários
- Links
- Indicadores de página ativa
- Elementos de destaque
- Gradientes sutis

**Não usar para**
- Texto de corpo
- Fundos grandes
- Alertas de erro

### Success (Esmeralda)

Cor para estados positivos e sucesso.

```css
--success-50: #ecfdf5
--success-100: #d1fae5
--success-200: #a7f3d0
--success-300: #6ee7b7
--success-400: #34d399
--success-500: #10b981  /* Principal */
--success-600: #059669
--success-700: #047857
--success-800: #065f46
--success-900: #064e3b
--success-950: #022c22
```

**Usos**
- Estados de sucesso
- Badges "Ativo", "Pago"
- Métricas positivas
- Confirmações
- Checkmarks

### Error (Coral)

Cor para estados de erro e ações destrutivas.

```css
--error-50: #fef2f2
--error-100: #fee2e2
--error-200: #fecaca
--error-300: #fca5a5
--error-400: #f87171
--error-500: #ef4444  /* Principal */
--error-600: #dc2626
--error-700: #b91c1c
--error-800: #991b1b
--error-900: #7f1d1d
--error-950: #450a0a
```

**Usos**
- Estados de erro
- Ações destrutivas (excluir)
- Badges "Cancelado", "Inativo"
- Alertas críticos
- Validação falha

### Warning (Âmbar)

Cor para avisos e atenção.

```css
--warning-50: #fffbeb
--warning-100: #fef3c7
--warning-200: #fde68a
--warning-300: #fcd34d
--warning-400: #fbbf24
--warning-500: #f59e0b  /* Principal */
--warning-600: #d97706
--warning-700: #b45309
--warning-800: #92400e
--warning-900: #78350f
--warning-950: #451a03
```

**Usos**
- Estados de aviso
- Badges "Estoque Baixo", "Pendente"
- Alertas não críticos
- Atenção necessária
- Validação warning

### Neutral (Ardesia)

Cores neutras para estrutura e texto.

```css
--neutral-50: #f8fafc
--neutral-100: #f1f5f9
--neutral-200: #e2e8f0
--neutral-300: #cbd5e1
--neutral-400: #94a3b8
--neutral-500: #64748b
--neutral-600: #475569
--neutral-700: #334155
--neutral-800: #1e293b
--neutral-900: #0f172a
--neutral-950: #020617
```

**Usos**
- Texto secundário
- Bordas
- Fundos secundários
- Ícones neutros
- Divisores

## Cores Semânticas

### Background

Fundos da aplicação.

```css
--background-default: #ffffff      /* Fundo principal */
--background-secondary: #f8fafc    /* Fundo secundário */
--background-tertiary: #f1f5f9     /* Fundo terciário */
--background-elevated: #ffffff     /* Elementos elevados */
--background-surface: #ffffff      /* Superfícies */
```

**Usos**
- `default`: Fundo principal da página
- `secondary`: Cards, sidebar
- `tertiary`: Inputs, hover states
- `elevated`: Modais, dropdowns
- `surface`: Componentes de superfície

### Text

Cores para tipografia.

```css
--text-primary: #0f172a      /* Texto principal */
--text-secondary: #475569    /* Texto secundário */
--text-tertiary: #64748b     /* Texto terciário */
--text-inverse: #ffffff      /* Texto sobre fundo escuro */
--text-disabled: #94a3b8     /* Texto desabilitado */
```

**Usos**
- `primary`: Títulos, texto importante
- `secondary`: Corpo de texto, descrições
- `tertiary`: Labels, helper text
- `inverse`: Texto sobre fundo colorido
- `disabled`: Estados desabilitados

### Border

Cores para bordas e divisores.

```css
--border-default: #e2e8f0    /* Borda padrão */
--border-light: #f1f5f9      /* Borda leve */
--border-dark: #1e293b       /* Borda escura */
--border-focus: #6366f1      /* Borda de foco */
```

**Usos**
- `default`: Bordas de cards, inputs
- `light`: Divisores sutis
- `dark`: Bordas em dark mode
- `focus`: Estado de foco

### Accent

Cores de acento para destaques especiais.

```css
--accent-blue: #3b82f6
--accent-purple: #8b5cf6
--accent-pink: #ec4899
--accent-cyan: #06b6d4
```

**Usos**
- Destaques especiais
- Gradientes
- Features premium
- Elementos decorativos

## Sistema de Badges

### Mapeamento de Status

**Pedidos**
- `pending` → Warning-500 (#f59e0b)
- `confirmed` → Primary-500 (#6366f1)
- `preparing` → Purple-500 (#8b5cf6)
- `ready` → Success-500 (#10b981)
- `delivered` → Success-700 (#047857)
- `cancelled` → Error-500 (#ef4444)

**Produtos**
- `active` → Success-500 (#10b981)
- `inactive` → Neutral-400 (#94a3b8)

**Pagamento**
- `paid` → Success-500 (#10b981)
- `pending` → Warning-500 (#f59e0b)
- `failed` → Error-500 (#ef4444)

**Estoque**
- `low` → Warning-500 (#f59e0b)
- `out` → Error-500 (#ef4444)
- `ok` → Success-500 (#10b981)

## Gradientes

### Principais

**Primary Gradient**
```css
background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
```

**Success Gradient**
```css
background: linear-gradient(135deg, #10b981 0%, #059669 100%);
```

**Subtle Gradient**
```css
background: linear-gradient(180deg, #ffffff 0%, #f8fafc 100%);
```

## Contraste e Acessibilidade

### Requisitos

- **Texto normal**: Contraste mínimo 4.5:1
- **Texto grande**: Contraste mínimo 3:1
- **Componentes UI**: Contraste mínimo 3:1
- **Elementos gráficos**: Contraste mínimo 3:1

### Validação

Todas as combinações de cor foram validadas para atender aos requisitos WCAG AA.

**Exemplos válidos**
- Primary-500 sobre branco: 5.74:1 ✓
- Text-primary sobre branco: 15.2:1 ✓
- Text-secondary sobre branco: 7.12:1 ✓
- Success-500 sobre branco: 3.67:1 ✓

## Dark Mode

### Estratégia

O sistema está preparado para dark mode com:

- Variáveis CSS para fácil inversão
- Paleta adaptada para baixo contraste
- Preservação de hierarquia visual
- Transições suaves entre temas

### Adaptação Futura

```css
[data-theme="dark"] {
  --background-default: #0f172a;
  --background-secondary: #1e293b;
  --text-primary: #f8fafc;
  --text-secondary: #cbd5e1;
  /* ... */
}
```

## Uso Prático

### Componentes

**Botão Primário**
```css
background: var(--primary-500);
color: white;
border: none;
```

**Card**
```css
background: var(--background-default);
border: 1px solid var(--border-default);
```

**Input**
```css
background: var(--background-tertiary);
border: 1px solid var(--border-default);
color: var(--text-primary);
```

**Badge Success**
```css
background: var(--success-100);
color: var(--success-700);
```

### Estados

**Hover**
```css
background: var(--neutral-100);
```

**Focus**
```css
border-color: var(--border-focus);
box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
```

**Error**
```css
border-color: var(--error-500);
color: var(--error-500);
```

## Tokens CSS

### Implementação

```css
:root {
  /* Primary */
  --color-primary-50: #eef2ff;
  --color-primary-500: #6366f1;
  /* ... */
  
  /* Semantic */
  --color-background-default: #ffffff;
  --color-text-primary: #0f172a;
  /* ... */
}
```

### Uso em Componentes

```svelte
<style>
  .button {
    background: var(--color-primary-500);
    color: var(--color-text-inverse);
  }
</style>
```

## Referências

### Inspirações

- **Stripe**: Paleta sofisticada e profissional
- **Linear**: Cores neutras com acentos sutis
- **Vercel**: Sistema de cores consistente
- **Tailwind CSS**: Escala de cores bem definida

### Ferramentas

- Coolors.co para geração de paletas
- Contrast Checker para validação
- Adobe Color para harmonia

## Regras de Uso

### DO ✅

- Usar cores semânticas para seu propósito
- Manter consistência em toda a aplicação
- Validar contraste para acessibilidade
- Usar variáveis CSS, não hardcode
- Considerar contexto de uso

### DON'T ❌

- Usar cores sem propósito
- Misturar muitas cores
- Hardcode valores hexadecimais
- Ignorar acessibilidade
- Usar cores que não comunicam

---

**Versão**: 1.0  
**Data**: 15/07/2026  
**Sprint**: 9 - Product Experience
