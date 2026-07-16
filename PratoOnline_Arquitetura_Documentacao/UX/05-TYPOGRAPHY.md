# Typography - PratoOnline

## Visão Geral

O sistema de tipografia do PratoOnline foi desenhado para criar hierarquia visual clara, legibilidade excepcional e uma aparência profissional. A fonte Inter é usada como base por sua excelência em interfaces digitais.

## Fontes

### Font Family

**Primária: Inter**
```css
font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
```

**Secundária: JetBrains Mono** (para código)
```css
font-family: 'JetBrains Mono', "Fira Code", Consolas, Monaco, monospace;
```

### Por que Inter?

- Desenhada especificamente para interfaces digitais
- Excelente legibilidade em tamanhos pequenos
- Variações de peso bem definidas
- Open source e amplamente adotada
- Performance otimizada (variable font)

### Carregamento

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
```

## Escala de Tamanhos

### Font Size Tokens

```css
--font-size-xs: 0.75rem;      /* 12px */
--font-size-sm: 0.875rem;     /* 14px */
--font-size-base: 1rem;       /* 16px */
--font-size-lg: 1.125rem;     /* 18px */
--font-size-xl: 1.25rem;      /* 20px */
--font-size-2xl: 1.5rem;      /* 24px */
--font-size-3xl: 1.875rem;    /* 30px */
--font-size-4xl: 2.25rem;     /* 36px */
--font-size-5xl: 3rem;        /* 48px */
--font-size-6xl: 3.75rem;     /* 60px */
--font-size-7xl: 4.5rem;      /* 72px */
--font-size-8xl: 6rem;        /* 96px */
--font-size-9xl: 8rem;        /* 128px */
```

### Uso por Tamanho

**xs (12px)**
- Captions
- Labels pequenas
- Metadata
- Badges

**sm (14px)**
- Texto secundário
- Labels de formulário
- Helper text
- Botões pequenos

**base (16px)**
- Texto de corpo (padrão)
- Inputs
- Links
- Botões normais

**lg (18px)**
- Texto de corpo grande
- Subtítulos
- Descrições

**xl (20px)**
- Títulos de seção
- Cards grandes
- Headers de componente

**2xl (24px)**
- Títulos de página (h2)
- Cards de destaque
- Headers importantes

**3xl (30px)**
- Títulos principais (h1)
- Hero text
- Números grandes

**4xl+ (36px+)**
- Display text
- Landing pages
- Marketing

## Font Weight

### Escala de Peso

```css
--font-weight-thin: 100;
--font-weight-extralight: 200;
--font-weight-light: 300;
--font-weight-normal: 400;      /* Padrão para corpo */
--font-weight-medium: 500;      /* Labels, subtítulos */
--font-weight-semibold: 600;    /* Títulos, botões */
--font-weight-bold: 700;        /* Títulos importantes */
--font-weight-extrabold: 800;
--font-weight-black: 900;
```

### Uso por Peso

**400 (Normal)**
- Texto de corpo
- Descrições
- Conteúdo longo

**500 (Medium)**
- Labels de formulário
- Subtítulos
- Texto secundário importante

**600 (Semibold)**
- Títulos de seção
- Botões
- Links importantes
- Navegação

**700 (Bold)**
- Títulos principais
- Números de destaque
- Call-to-action

## Line Height

### Tokens

```css
--line-height-none: 1;
--line-height-tight: 1.2;
--line-height-snug: 1.35;
--line-height-normal: 1.5;      /* Padrão para corpo */
--line-height-relaxed: 1.625;
--line-height-loose: 2;
```

### Uso por Contexto

**Tight (1.2)**
- Títulos grandes
- Números
- Texto curto

**Snug (1.35)**
- Títulos médios
- Subtítulos
- Labels

**Normal (1.5)**
- Texto de corpo (padrão)
- Parágrafos
- Listas

**Relaxed (1.625)**
- Texto longo
- Descrições
- Conteúdo editorial

**Loose (2)**
- Citações
- Texto muito longo
- Acessibilidade

## Letter Spacing

### Tokens

```css
--letter-spacing-tighter: -0.05em;
--letter-spacing-tight: -0.025em;
--letter-spacing-normal: 0em;
--letter-spacing-wide: 0.025em;
--letter-spacing-wider: 0.05em;
--letter-spacing-widest: 0.1em;
```

### Uso por Contexto

**Tight (-0.025em)**
- Títulos grandes
- Texto em uppercase

**Normal (0em)**
- Texto de corpo
- Padrão para maioria

**Wide (0.05em)**
- Labels em uppercase
- Navegação
- Botões

## Tipografia Semântica

### Hierarquia de Títulos

**H1 - Título da Página**
```css
font-size: 1.875rem;      /* 30px */
font-weight: 700;
line-height: 1.2;
letter-spacing: -0.025em;
```

**H2 - Título de Seção**
```css
font-size: 1.5rem;        /* 24px */
font-weight: 600;
line-height: 1.2;
letter-spacing: 0;
```

**H3 - Subtítulo**
```css
font-size: 1.25rem;       /* 20px */
font-weight: 600;
line-height: 1.2;
letter-spacing: 0;
```

**H4 - Título de Componente**
```css
font-size: 1.125rem;     /* 18px */
font-weight: 500;
line-height: 1.2;
letter-spacing: 0;
```

**H5 - Título Pequeno**
```css
font-size: 1rem;          /* 16px */
font-weight: 500;
line-height: 1.2;
letter-spacing: 0;
```

**H6 - Label de Grupo**
```css
font-size: 0.875rem;      /* 14px */
font-weight: 600;
line-height: 1.2;
letter-spacing: 0;
```

### Texto de Corpo

**Body (Padrão)**
```css
font-size: 1rem;          /* 16px */
font-weight: 400;
line-height: 1.5;
letter-spacing: 0;
```

**Body Large**
```css
font-size: 1.125rem;      /* 18px */
font-weight: 400;
line-height: 1.5;
letter-spacing: 0;
```

**Body Small**
```css
font-size: 0.875rem;      /* 14px */
font-weight: 400;
line-height: 1.5;
letter-spacing: 0;
```

### Texto Auxiliar

**Caption**
```css
font-size: 0.75rem;       /* 12px */
font-weight: 400;
line-height: 1.5;
letter-spacing: 0;
```

**Overline (Uppercase)**
```css
font-size: 0.75rem;       /* 12px */
font-weight: 500;
line-height: 1.2;
letter-spacing: 0.05em;
text-transform: uppercase;
```

### Componentes

**Button**
```css
font-size: 0.875rem;      /* 14px */
font-weight: 500;
line-height: 1.2;
letter-spacing: 0;
```

**Button Large**
```css
font-size: 1rem;          /* 16px */
font-weight: 500;
line-height: 1.2;
letter-spacing: 0;
```

**Button Small**
```css
font-size: 0.75rem;       /* 12px */
font-weight: 500;
line-height: 1.2;
letter-spacing: 0;
```

**Label (Formulário)**
```css
font-size: 0.875rem;      /* 14px */
font-weight: 500;
line-height: 1.2;
letter-spacing: 0;
```

**Input**
```css
font-size: 1rem;          /* 16px */
font-weight: 400;
line-height: 1.5;
letter-spacing: 0;
```

**Helper Text**
```css
font-size: 0.75rem;       /* 12px */
font-weight: 400;
line-height: 1.5;
letter-spacing: 0;
color: var(--text-tertiary);
```

### Código

**Code (Bloco)**
```css
font-size: 0.875rem;      /* 14px */
font-weight: 400;
line-height: 1.5;
letter-spacing: 0;
font-family: 'JetBrains Mono', monospace;
```

**Code Inline**
```css
font-size: 0.875rem;      /* 14px */
font-weight: 400;
line-height: 1.5;
letter-spacing: 0;
font-family: 'JetBrains Mono', monospace;
background: var(--neutral-100);
padding: 0.125rem 0.375rem;
border-radius: 4px;
```

## Cores de Texto

### Hierarquia por Cor

**Primary**
```css
color: var(--text-primary);  /* #0f172a */
```
- Títulos
- Texto importante
- Links ativos

**Secondary**
```css
color: var(--text-secondary);  /* #475569 */
```
- Texto de corpo
- Descrições
- Conteúdo secundário

**Tertiary**
```css
color: var(--text-tertiary);  /* #64748b */
```
- Helper text
- Metadata
- Labels desabilitados

**Inverse**
```css
color: var(--text-inverse);  /* #ffffff */
```
- Texto sobre fundo escuro
- Texto sobre botões coloridos

**Disabled**
```css
color: var(--text-disabled);  /* #94a3b8 */
```
- Estados desabilitados
- Texto inativo

## Uso Prático

### Exemplos de Componentes

**Card Header**
```css
.card-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  line-height: var(--line-height-tight);
  color: var(--text-primary);
}
```

**Metric Value**
```css
.metric-value {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  line-height: var(--line-height-tight);
  color: var(--text-primary);
}
```

**Table Cell**
```css
.table-cell {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-normal);
  line-height: var(--line-height-normal);
  color: var(--text-secondary);
}
```

**Badge**
```css
.badge {
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  line-height: var(--line-height-tight);
  letter-spacing: var(--letter-spacing-wide);
  text-transform: uppercase;
}
```

## Responsividade

### Escala Mobile

**Mobile (< 768px)**
- H1: 1.5rem (24px)
- H2: 1.25rem (20px)
- Body: 0.875rem (14px)

**Tablet (768px - 1024px)**
- H1: 1.75rem (28px)
- H2: 1.375rem (22px)
- Body: 0.9375rem (15px)

**Desktop (> 1024px)**
- H1: 1.875rem (30px)
- H2: 1.5rem (24px)
- Body: 1rem (16px)

### Fluid Typography

Para títulos grandes, usar `clamp()`:

```css
h1 {
  font-size: clamp(1.5rem, 4vw, 2.5rem);
}
```

## Acessibilidade

### Line Height

- Mínimo de 1.5 para texto de corpo
- Mínimo de 1.2 para títulos
- Espaçamento entre parágrafos: 1.5x line-height

### Font Size

- Mínimo de 16px para texto de corpo
- Mínimo de 14px para labels
- Escalável até 200% sem quebra de layout

### Contraste

- Contraste mínimo 4.5:1 para texto normal
- Contraste mínimo 3:1 para texto grande
- Validar com ferramentas de acessibilidade

## Performance

### Otimizações

1. **Font Display**: `swap` para renderização imediata
2. **Subset**: Carregar apenas caracteres necessários
3. **Variable Font**: Usar quando possível
4. **Preload**: Carregar fontes críticas antecipadamente

### Implementação

```css
@font-face {
  font-family: 'Inter';
  font-style: normal;
  font-weight: 400 700;
  font-display: swap;
  src: url('/fonts/inter-variable.woff2') format('woff2-variations');
}
```

## Regras de Uso

### DO ✅

- Manter hierarquia visual clara
- Usar pesos consistentemente
- Respeitar line-height para legibilidade
- Usar tamanhos apropriados ao contexto
- Considerar legibilidade em mobile

### DON'T ❌

- Usar muitos tamanhos diferentes
- Misturar pesos sem propósito
- Line-height muito apertado
- Texto muito pequeno (< 12px)
- Ignorar escalabilidade

## Tokens CSS

### Implementação Completa

```css
:root {
  /* Font Family */
  --font-sans: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  --font-mono: 'JetBrains Mono', "Fira Code", Consolas, Monaco, monospace;
  
  /* Font Size */
  --text-xs: 0.75rem;
  --text-sm: 0.875rem;
  --text-base: 1rem;
  --text-lg: 1.125rem;
  --text-xl: 1.25rem;
  --text-2xl: 1.5rem;
  --text-3xl: 1.875rem;
  
  /* Font Weight */
  --font-normal: 400;
  --font-medium: 500;
  --font-semibold: 600;
  --font-bold: 700;
  
  /* Line Height */
  --leading-tight: 1.2;
  --leading-normal: 1.5;
  --leading-relaxed: 1.625;
  
  /* Letter Spacing */
  --tracking-tight: -0.025em;
  --tracking-normal: 0em;
  --tracking-wide: 0.05em;
}
```

---

**Versão**: 1.0  
**Data**: 15/07/2026  
**Sprint**: 9 - Product Experience
