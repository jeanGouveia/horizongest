# Spacing - PratoOnline

## Visão Geral

O sistema de espaçamento do PratoOnline foi criado para garantir consistência visual, hierarquia clara e uma interface respirável. Espaçamentos generosos são usados para transmitir leveza e profissionalismo.

## Filosofia

### Princípios

**Respiro**
- Muito espaço em branco
- Não ter medo de espaço vazio
- Deixar o conteúdo "respirar"

**Consistência**
- Escala de espaçamento consistente
- Mesmos valores em toda a aplicação
- Multiplos de 4px base

**Hierarquia**
- Espaço para criar relacionamentos
- Agrupamento lógico através de espaçamento
- Distinção visual através de distância

**Proporção**
- Relação harmoniosa entre elementos
- Espaçamento proporcional ao tamanho
- Golden ratio onde aplicável

## Escala de Espaçamento

### Base Tokens

```css
--spacing-0: 0px;
--spacing-1: 4px;
--spacing-2: 8px;
--spacing-3: 12px;
--spacing-4: 16px;
--spacing-5: 20px;
--spacing-6: 24px;
--spacing-7: 28px;
--spacing-8: 32px;
--spacing-9: 36px;
--spacing-10: 40px;
--spacing-11: 44px;
--spacing-12: 48px;
--spacing-14: 56px;
--spacing-16: 64px;
--spacing-20: 80px;
--spacing-24: 96px;
--spacing-28: 112px;
--spacing-32: 128px;
--spacing-36: 144px;
--spacing-40: 160px;
--spacing-44: 176px;
--spacing-48: 192px;
--spacing-52: 208px;
--spacing-56: 224px;
--spacing-60: 240px;
--spacing-64: 256px;
--spacing-72: 288px;
--spacing-80: 320px;
--spacing-96: 384px;
```

### Tokens Semânticos

```css
--space-none: 0px;
--space-xs: 4px;
--space-sm: 8px;
--space-base: 12px;
--space-md: 16px;
--space-lg: 24px;
--space-xl: 32px;
--space-2xl: 48px;
--space-3xl: 64px;
--space-4xl: 80px;
--space-5xl: 96px;
```

## Uso por Contexto

### Espaçamento Interno (Padding)

**Componentes Pequenos**
- Badge: 4px 8px
- Button small: 4px 12px
- Icon button: 8px

**Componentes Médios**
- Button: 8px 16px
- Input: 12px 16px
- Select: 12px 16px

**Componentes Grandes**
- Card: 24px
- Modal: 32px
- Page container: 24px

**Containers**
- Section: 48px 24px
- Hero: 64px 24px

### Espaçamento Externo (Margin)

**Entre Elementos Irmãos**
- Ícones em botão: 4px
- Badge e texto: 8px
- Items em lista: 12px
- Cards em grid: 16px
- Seções: 24px

**Entre Grupos**
- Seções principais: 32px
- Blocos de conteúdo: 48px
- Seções maiores: 64px

### Espaçamento em Grid

**Gap**
- Grid apertado: 8px
- Grid normal: 16px
- Grid espaçoso: 24px
- Grid muito espaçoso: 32px

## Padrões de Componentes

### Cards

**Card Padrão**
```css
padding: 24px;
gap: 16px;
```

**Card Compacto**
```css
padding: 16px;
gap: 12px;
```

**Card Espaçoso**
```css
padding: 32px;
gap: 24px;
```

### Formulários

**Campo de Formulário**
```css
/* Label */
margin-bottom: 8px;

/* Input */
padding: 12px 16px;

/* Helper text */
margin-top: 8px;
```

**Seção de Formulário**
```css
padding: 24px;
gap: 16px;
margin-bottom: 24px;
```

### Tabelas

**Célula**
```css
padding: 12px 16px;
```

**Header da Tabela**
```css
padding: 16px;
```

### Botões

**Button Small**
```css
padding: 6px 12px;
gap: 6px;
```

**Button Medium**
```css
padding: 10px 16px;
gap: 8px;
```

**Button Large**
```css
padding: 14px 20px;
gap: 10px;
```

### Modais

**Modal**
```css
padding: 32px;
gap: 24px;
```

**Modal Header**
```css
margin-bottom: 24px;
```

**Modal Footer**
```css
margin-top: 24px;
padding-top: 24px;
border-top: 1px solid var(--border-default);
```

## Layout

### Container

**Page Container**
```css
padding: 24px;
max-width: 1400px;
margin: 0 auto;
```

**Section**
```css
padding: 48px 24px;
```

### Grid

**Grid de Cards**
```css
display: grid;
gap: 16px;
grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
```

**Grid de Métricas**
```css
display: grid;
gap: 16px;
grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
```

### Flex

**Lista Horizontal**
```css
display: flex;
gap: 12px;
align-items: center;
```

**Lista Vertical**
```css
display: flex;
flex-direction: column;
gap: 12px;
```

## Espaçamento Responsivo

### Mobile (< 768px)

**Padding**
- Page: 16px
- Card: 16px
- Section: 32px 16px

**Gap**
- Grid: 12px
- List: 8px

### Tablet (768px - 1024px)

**Padding**
- Page: 20px
- Card: 20px
- Section: 40px 20px

**Gap**
- Grid: 16px
- List: 12px

### Desktop (> 1024px)

**Padding**
- Page: 24px
- Card: 24px
- Section: 48px 24px

**Gap**
- Grid: 16px
- List: 12px

## Padrões Específicos

### Dashboard

**Métricas**
```css
gap: 16px;
margin-bottom: 24px;
```

**Gráfico**
```css
margin-bottom: 24px;
padding: 24px;
```

**Cards Secundários**
```css
gap: 16px;
margin-bottom: 24px;
```

### Páginas de Lista

**Header**
```css
margin-bottom: 24px;
```

**Filtros**
```css
margin-bottom: 16px;
gap: 12px;
```

**Tabela**
```css
margin-bottom: 24px;
```

### Páginas de Detalhe

**Breadcrumb**
```css
margin-bottom: 16px;
```

**Header**
```css
margin-bottom: 32px;
```

**Seções**
```css
margin-bottom: 24px;
padding: 24px;
```

## Espaçamento Negativo

### Uso

Para sobrepor elementos ou criar layouts especiais:

```css
.negative-margin {
  margin-top: -8px;
  margin-left: -16px;
}
```

### Cuidados

- Usar com moderação
- Considerar z-index
- Testar em diferentes tamanhos de tela
- Documentar o propósito

## Espaçamento em Animações

### Transições

Ao animar espaçamento, usar transições suaves:

```css
.card {
  padding: 24px;
  transition: padding 0.2s ease;
}

.card:hover {
  padding: 32px;
}
```

### Performance

Evitar animar layout properties quando possível:
- Prefira transform e opacity
- Espaçamento animado pode causar reflow
- Use will-change com cuidado

## Tokens CSS

### Implementação

```css
:root {
  /* Base */
  --space-0: 0px;
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;
  --space-8: 32px;
  --space-10: 40px;
  --space-12: 48px;
  --space-16: 64px;
  --space-20: 80px;
  --space-24: 96px;
  
  /* Semantic */
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  --space-2xl: 48px;
  --space-3xl: 64px;
}
```

### Uso em Componentes

```css
.card {
  padding: var(--space-lg);
  gap: var(--space-md);
}

.button {
  padding: var(--space-sm) var(--space-md);
  gap: var(--space-sm);
}
```

## Utilities

### Classes de Espaçamento

```css
.p-0 { padding: var(--space-0); }
.p-1 { padding: var(--space-1); }
.p-2 { padding: var(--space-2); }
.p-3 { padding: var(--space-3); }
.p-4 { padding: var(--space-4); }
.p-6 { padding: var(--space-6); }
.p-8 { padding: var(--space-8); }

.m-0 { margin: var(--space-0); }
.m-1 { margin: var(--space-1); }
.m-2 { margin: var(--space-2); }
.m-3 { margin: var(--space-3); }
.m-4 { margin: var(--space-4); }
.m-6 { margin: var(--space-6); }
.m-8 { margin: var(--space-8); }

.gap-1 { gap: var(--space-1); }
.gap-2 { gap: var(--space-2); }
.gap-3 { gap: var(--space-3); }
.gap-4 { gap: var(--space-4); }
.gap-6 { gap: var(--space-6); }
.gap-8 { gap: var(--space-8); }
```

## Regras de Uso

### DO ✅

- Usar escala consistente
- Espaçamento generoso para respiro
- Agrupar elementos relacionados
- Considerar responsividade
- Usar tokens, não valores hardcoded

### DON'T ❌

- Valores aleatórios de espaçamento
- Espaçamento muito apertado
- Misturar escalas
- Ignorar mobile
- Hardcode valores

## Referências

### Inspirações

- **Tailwind CSS**: Escala de espaçamento de 4px
- **8pt Grid System**: Múltiplos de 8px
- **Material Design**: Sistema de espaçamento de 4px

### Ferramentas

- Grid calculators
- Spacing generators
- Design systems documentation

## Checklist

### Review de Espaçamento

- [ ] Escala consistente usada
- [ ] Hierarquia visual clara
- [ ] Respiro adequado
- [ ] Responsivo testado
- [ ] Tokens CSS usados
- [ ] Acessibilidade considerada

---

**Versão**: 1.0  
**Data**: 15/07/2026  
**Sprint**: 9 - Product Experience
