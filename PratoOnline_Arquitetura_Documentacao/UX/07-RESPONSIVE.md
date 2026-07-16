# Responsive - PratoOnline

## Visão Geral

O sistema responsivo do PratoOnline foi desenhado para proporcionar uma experiência consistente e otimizada em todos os dispositivos. A abordagem mobile-first garante que a interface funcione perfeitamente desde telas pequenas até monitores grandes.

## Estratégia

### Mobile-First

**Filosofia**
- Começar com layout mobile
- Adicionar complexidade em breakpoints maiores
- Progressividade: mobile → tablet → desktop
- Performance otimizada para mobile

**Benefícios**
- Performance melhor em mobile
- Código mais limpo
- Foco no conteúdo essencial
- Experiência consistente

## Breakpoints

### Escala de Breakpoints

```css
/* Mobile First */
--breakpoint-sm: 480px;    /* Mobile grande */
--breakpoint-md: 768px;    /* Tablet */
--breakpoint-lg: 1024px;   /* Desktop pequeno */
--breakpoint-xl: 1280px;   /* Desktop */
--breakpoint-2xl: 1536px;  /* Desktop grande */
```

### Media Queries

```css
/* Mobile (< 768px) - Base styles */
@media (max-width: 767px) { }

/* Tablet (768px - 1023px) */
@media (min-width: 768px) and (max-width: 1023px) { }

/* Desktop (1024px - 1279px) */
@media (min-width: 1024px) and (max-width: 1279px) { }

/* Desktop Large (>= 1280px) */
@media (min-width: 1280px) { }
```

### Mobile-First Approach

```css
/* Base: Mobile styles */
.container {
  padding: 16px;
  grid-template-columns: 1fr;
}

/* Tablet */
@media (min-width: 768px) {
  .container {
    padding: 20px;
    grid-template-columns: repeat(2, 1fr);
  }
}

/* Desktop */
@media (min-width: 1024px) {
  .container {
    padding: 24px;
    grid-template-columns: repeat(3, 1fr);
  }
}
```

## Layout Responsivo

### Sidebar

**Mobile (< 768px)**
- Oculta por padrão
- Menu hamburger no header
- Drawer lateral ao abrir
- Overlay escuro

**Tablet (768px - 1024px)**
- Recolhida por padrão (72px)
- Toggle funcional
- Apenas ícones
- Tooltip no hover

**Desktop (> 1024px)**
- Expandida por padrão (260px)
- Toggle funcional
- Ícones + labels
- Grupos visíveis

### Header

**Mobile (< 768px)**
- Layout minimalista
- Breadcrumb oculto
- Busca apenas ícone
- Ações reduzidas
- Nome do usuário oculto

**Tablet (768px - 1024px)**
- Breadcrumb oculto
- Busca apenas ícone
- Nome do usuário oculto
- Ações completas

**Desktop (> 1024px)**
- Layout completo
- Breadcrumb visível
- Busca com texto
- Todos os elementos

### Grid System

**Mobile (< 768px)**
- 1 coluna
- Cards full width
- Stack vertical

**Tablet (768px - 1024px)**
- 2 colunas
- Cards 50% width
- Grid adaptativo

**Desktop (1024px - 1279px)**
- 3 colunas
- Cards 33.33% width
- Grid balanceado

**Desktop Large (>= 1280px)**
- 4 colunas
- Cards 25% width
- Grid espaçoso

## Componentes Responsivos

### Cards

**Mobile**
```css
.card {
  padding: 16px;
  margin-bottom: 16px;
}
```

**Tablet**
```css
@media (min-width: 768px) {
  .card {
    padding: 20px;
  }
}
```

**Desktop**
```css
@media (min-width: 1024px) {
  .card {
    padding: 24px;
  }
}
```

### Tabelas

**Mobile**
- Scroll horizontal
- Cards alternativos (opcional)
- Informações essenciais apenas

**Tablet**
- Scroll horizontal
- Colunas adaptativas
- Informações importantes

**Desktop**
- Tabela completa
- Todas as colunas
- Hover effects

### Formulários

**Mobile**
- Stack vertical
- Inputs full width
- Botões full width

**Tablet**
- 2 colunas quando apropriado
- Inputs adaptativos
- Botões side-by-side

**Desktop**
- Grid multi-coluna
- Inputs otimizados
- Botões posicionados

### Modais

**Mobile**
- Full screen ou 90% width
- Bottom sheet quando apropriado
- Botões full width

**Tablet**
- 80% width
- Centralizado
- Botões side-by-side

**Desktop**
- 600px width fixo ou 50%
- Centralizado
- Botões side-by-side

## Tipografia Responsiva

### Escala Fluida

**Mobile (< 768px)**
```css
h1 { font-size: 1.5rem; }      /* 24px */
h2 { font-size: 1.25rem; }    /* 20px */
body { font-size: 0.875rem; } /* 14px */
```

**Tablet (768px - 1024px)**
```css
h1 { font-size: 1.75rem; }    /* 28px */
h2 { font-size: 1.375rem; }   /* 22px */
body { font-size: 0.9375rem; } /* 15px */
```

**Desktop (> 1024px)**
```css
h1 { font-size: 1.875rem; }   /* 30px */
h2 { font-size: 1.5rem; }     /* 24px */
body { font-size: 1rem; }     /* 16px */
```

### Fluid Typography com clamp()

```css
h1 {
  font-size: clamp(1.5rem, 4vw, 2.5rem);
}
```

## Espaçamento Responsivo

### Padding

**Mobile**
```css
.page-container { padding: 16px; }
.card { padding: 16px; }
.section { padding: 32px 16px; }
```

**Tablet**
```css
@media (min-width: 768px) {
  .page-container { padding: 20px; }
  .card { padding: 20px; }
  .section { padding: 40px 20px; }
}
```

**Desktop**
```css
@media (min-width: 1024px) {
  .page-container { padding: 24px; }
  .card { padding: 24px; }
  .section { padding: 48px 24px; }
}
```

### Gap

**Mobile**
```css
.grid { gap: 12px; }
.list { gap: 8px; }
```

**Tablet**
```css
@media (min-width: 768px) {
  .grid { gap: 16px; }
  .list { gap: 12px; }
}
```

**Desktop**
```css
@media (min-width: 1024px) {
  .grid { gap: 16px; }
  .list { gap: 12px; }
}
```

## Navegação Responsiva

### Breadcrumb

**Mobile**
- Oculto ou simplificado
- Máximo 2 níveis

**Tablet**
- Parcialmente visível
- Truncar com "..." se necessário

**Desktop**
- Completo
- Todos os níveis visíveis

### Tabs

**Mobile**
- Scroll horizontal
- Swipe para navegar
- Bottom navigation (opcional)

**Tablet**
- Scroll horizontal
- Todas as tabs visíveis se possível

**Desktop**
- Todas as tabs visíveis
- Hover effects

## Imagens Responsivas

### Imagens

```css
img {
  max-width: 100%;
  height: auto;
}
```

### srcset para diferentes resoluções

```html
<img 
  src="image-800.jpg"
  srcset="image-400.jpg 400w,
          image-800.jpg 800w,
          image-1200.jpg 1200w"
  sizes="(max-width: 768px) 100vw,
         (max-width: 1024px) 50vw,
         33vw"
  alt="Descrição"
/>
```

## Touch Targets

### Tamanhos Mínimos

**Mobile**
- Botões: mínimo 44px de altura
- Links: mínimo 44px de altura
- Touch targets: mínimo 44x44px

**Tablet**
- Botões: mínimo 40px de altura
- Links: mínimo 40px de altura

**Desktop**
- Botões: mínimo 36px de altura
- Links: mínimo 36px de altura

### Espaçamento entre Touch Targets

```css
.button {
  min-height: 44px;
  min-width: 44px;
  margin: 8px;
}
```

## Performance

### Otimizações Mobile

1. **Lazy Loading** → Carregar imagens sob demanda
2. **Code Splitting** → Dividir código por rota
3. **Critical CSS** → CSS crítico inline
4. **Font Display** → swap para renderização imediata
5. **Reduce Payload** → Minificar e comprimir

### Imagens

- WebP quando suportado
- Tamanhos apropriados
- Lazy loading
- Placeholder com blur

## Acessibilidade

### Zoom

- Suportar zoom até 200%
- Layout não quebrar
- Texto permanecer legível
- Touch targets usáveis

### Orientação

- Suportar portrait e landscape
- Layout adaptativo
- Conteúdo acessível em ambos

### Keyboard

- Navegação por teclado funcional
- Focus visível
- Tab order lógico
- Skip links

## Padrões de Página

### Dashboard

**Mobile**
- Métricas: 2 colunas
- Gráfico: full width
- Cards secundários: stack
- Ações rápidas: 2 colunas

**Tablet**
- Métricas: 3 colunas
- Gráfico: full width
- Cards secundários: 2 colunas
- Ações rápidas: 4 colunas

**Desktop**
- Métricas: 5 colunas
- Gráfico: full width
- Cards secundários: 2 colunas
- Ações rápidas: 4 colunas

### Lista (Produtos/Pedidos)

**Mobile**
- Filtros: drawer ou accordion
- Tabela: cards ou scroll
- Paginação: simplificada

**Tablet**
- Filtros: visíveis
- Tabela: scroll horizontal
- Paginação: completa

**Desktop**
- Filtros: visíveis
- Tabela: completa
- Paginação: completa

### Formulário

**Mobile**
- Stack vertical
- Seções como accordion
- Botões: full width

**Tablet**
- 2 colunas quando apropriado
- Seções visíveis
- Botões: side-by-side

**Desktop**
- Grid multi-coluna
- Seções visíveis
- Botões: posicionados

## Testes

### Dispositivos para Testar

**Mobile**
- iPhone SE (375px)
- iPhone 12/13 (390px)
- iPhone 14 Pro Max (430px)
- Android pequeno (360px)

**Tablet**
- iPad Mini (768px)
- iPad (1024px)
- Android tablet (800px)

**Desktop**
- Laptop (1366px)
- Desktop (1920px)
- 4K monitor (2560px)

### Ferramentas

- Chrome DevTools
- Firefox Responsive Design Mode
- BrowserStack
- Device lab (se disponível)

## Best Practices

### DO ✅

- Mobile-first approach
- Testar em múltiplos dispositivos
- Otimizar performance mobile
- Touch targets adequados
- Conteúdo prioritário visível

### DON'T ❌

- Desktop-first approach
- Ignorar mobile
- Touch targets pequenos
- Conteúdo escondido sem motivo
- Performance ruim em mobile

## CSS Utility Classes

### Responsive Utilities

```css
/* Display */
.hidden-mobile { display: none; }
@media (min-width: 768px) { .hidden-mobile { display: block; } }

.hidden-tablet { display: none; }
@media (min-width: 1024px) { .hidden-tablet { display: block; } }

/* Grid */
.grid-1 { grid-template-columns: 1fr; }
@media (min-width: 768px) { .grid-2 { grid-template-columns: repeat(2, 1fr); } }
@media (min-width: 1024px) { .grid-3 { grid-template-columns: repeat(3, 1fr); } }
@media (min-width: 1280px) { .grid-4 { grid-template-columns: repeat(4, 1fr); } }

/* Text */
.text-sm-mobile { font-size: 0.875rem; }
@media (min-width: 768px) { .text-base-tablet { font-size: 1rem; } }
@media (min-width: 1024px) { .text-lg-desktop { font-size: 1.125rem; } }
```

## Checklist

### Review Responsivo

- [ ] Mobile-first implementado
- [ ] Breakpoints definidos
- [ ] Layout testado em todos os breakpoints
- [ ] Touch targets adequados
- [ ] Performance otimizada
- [ ] Acessibilidade verificada
- [ ] Imagens responsivas
- [ ] Tipografia fluida
- [ ] Navegação adaptativa
- [ ] Formulários usáveis

---

**Versão**: 1.0  
**Data**: 15/07/2026  
**Sprint**: 9 - Product Experience
