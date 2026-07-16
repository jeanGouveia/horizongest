# Icons - PratoOnline

## Visão Geral

O sistema de ícones do PratoOnline foi desenhado para proporcionar consistência visual, clareza semântica e uma aparência profissional. A biblioteca Lucide Icons é usada como base por sua excelência em interfaces modernas.

## Biblioteca de Ícones

### Lucide Icons

**Por que Lucide?**
- Desenhado especificamente para interfaces modernas
- Estilo consistente e refinado
- Open source e amplamente adotado
- Performance otimizada (SVG)
- Fácil customização via CSS
- Grande variedade de ícones

**Instalação**
```bash
npm install lucide-svelte
```

**Importação**
```svelte
<script>
  import { Search, Menu, X, ChevronRight } from 'lucide-svelte';
</script>
```

## Categorias de Ícones

### Navegação

**Menu e Layout**
- `menu` - Menu hamburger
- `x` - Fechar
- `chevron-left` - Voltar
- `chevron-right` - Avançar
- `chevron-down` - Expandir
- `chevron-up` - Colapsar
- `layout-dashboard` - Dashboard
- `layout-sidebar` - Sidebar

### Ações

**Comuns**
- `search` - Buscar
- `filter` - Filtrar
- `plus` - Adicionar
- `minus` - Remover
- `edit` - Editar
- `trash-2` - Excluir
- `save` - Salvar
- `download` - Download
- `upload` - Upload
- `refresh-cw` - Atualizar
- `more-horizontal` - Mais opções
- `more-vertical` - Mais opções (vertical)

### Status

**Indicadores**
- `check-circle` - Sucesso
- `x-circle` - Erro
- `alert-circle` - Aviso
- `info` - Informação
- `help-circle` - Ajuda
- `clock` - Pendente
- `loader-2` - Carregando

### Usuário

**Perfil e Conta**
- `user` - Usuário
- `users` - Usuários
- `settings` - Configurações
- `log-out` - Sair
- `shield` - Admin
- `bell` - Notificações
- `mail` - Email

### Negócio

**Operação**
- `shopping-cart` - Carrinho/Pedidos
- `package` - Produtos
- `utensils` - Cardápio/Pratos
- `leaf` - Ingredientes
- `warehouse` - Estoque
- `scale` - Ajustes
- `trending-up` - Crescimento
- `trending-down` - Decréscimo
- `bar-chart` - Gráficos
- `pie-chart` - Gráficos

### Arquivos

**Documentos**
- `file` - Arquivo
- `folder` - Pasta
- `image` - Imagem
- `download-cloud` - Download
- `upload-cloud` - Upload
- `file-text` - Documento

### Interface

**UI Elements**
- `eye` - Visualizar
- `eye-off` - Ocultar
- `lock` - Bloqueado
- `unlock` - Desbloqueado
- `star` - Favorito
- `heart` - Curtir
- `bookmark` - Salvar
- `share` - Compartilhar
- `link` - Link
- `copy` - Copiar
- `clipboard` - Clipboard

### Setas

**Navegação**
- `arrow-left` - Esquerda
- `arrow-right` - Direita
- `arrow-up` - Cima
- `arrow-down` - Baixo
- `arrow-up-right` - Externo
- `external-link` - Link externo

## Tamanhos de Ícone

### Escala

```css
--icon-xs: 14px;
--icon-sm: 16px;
--icon-md: 20px;
--icon-lg: 24px;
--icon-xl: 32px;
```

### Uso por Contexto

**xs (14px)**
- Badges
- Buttons pequenos
- Inline com texto pequeno

**sm (16px)**
- Buttons normais
- Inputs (prefixo/sufixo)
- Lista items

**md (20px)**
- Buttons grandes
- Cards
- Tabela headers

**lg (24px)**
- Page headers
- Cards grandes
- Seções importantes

**xl (32px)**
- Hero sections
- Empty states
- Destaques especiais

## Cores de Ícone

### Semânticas

**Default**
```css
color: var(--text-secondary);  /* #475569 */
```

**Primary**
```css
color: var(--primary-500);  /* #6366f1 */
```

**Success**
```css
color: var(--success-500);  /* #10b981 */
```

**Warning**
```css
color: var(--warning-500);  /* #f59e0b */
```

**Error**
```css
color: var(--error-500);  /* #ef4444 */
```

**Inverse**
```css
color: var(--text-inverse);  /* #ffffff */
```

**Disabled**
```css
color: var(--text-disabled);  /* #94a3b8 */
```

## Uso em Componentes

### Buttons

**Com Ícone**
```svelte
<Button>
  <Search size={16} />
  Buscar
</Button>
```

**Ícone Apenas**
```svelte
<Button variant="ghost">
  <Settings size={20} />
</Button>
```

### Inputs

**Prefixo**
```svelte
<div class="input-wrapper">
  <Search size={16} class="input-icon" />
  <Input placeholder="Buscar..." />
</div>
```

**Sufixo**
```svelte
<div class="input-wrapper">
  <Input type="password" />
  <Eye size={16} class="input-icon" />
</div>
```

### Cards

**Header**
```svelte
<Card>
  <div class="card-header">
    <Package size={24} />
    <h3>Produtos</h3>
  </div>
</Card>
```

### Badges

**Com Ícone**
```svelte
<Badge variant="success">
  <CheckCircle size={12} />
  Ativo
</Badge>
```

### Menu Items

**Sidebar**
```svelte
<a href="/orders" class="nav-link">
  <ShoppingCart size={20} />
  <span>Pedidos</span>
</a>
```

## Ícones Específicos do PratoOnline

### Mapeamento de Funcionalidades

**Dashboard**
- `layout-dashboard` - Dashboard principal
- `bar-chart-3` - Gráficos
- `trending-up` - Métricas positivas
- `trending-down` - Métricas negativas

**Pedidos**
- `shopping-cart` - Pedidos
- `receipt` - Detalhes do pedido
- `clock` - Pendente
- `check-circle` - Confirmado
- `chef-hat` - Preparando
- `truck` - Em entrega
- `check-circle-2` - Entregue
- `x-circle` - Cancelado

**Produtos**
- `utensils` - Produtos/Pratos
- `star` - Destaque
- `image` - Foto do produto
- `tag` - Categoria

**Ingredientes**
- `leaf` - Ingredientes
- `beaker` - Estoque
- `alert-triangle` - Estoque baixo

**Estoque**
- `warehouse` - Estoque
- `scale` - Ajustes
- `package-plus` - Adicionar estoque
- `package-minus` - Remover estoque

**Administração**
- `user` - Perfil
- `users` - Usuários
- `settings` - Configurações
- `shield` - Permissões

## Estados de Ícone

### Default
```css
color: var(--text-secondary);
```

### Hover
```css
color: var(--text-primary);
transition: color 0.2s ease;
```

### Active
```css
color: var(--primary-500);
```

### Disabled
```css
color: var(--text-disabled);
pointer-events: none;
opacity: 0.5;
```

### Loading
```css
animation: spin 1s linear infinite;
```

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

## Animações

### Spin
```css
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.icon-spin {
  animation: spin 1s linear infinite;
}
```

### Pulse
```css
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.icon-pulse {
  animation: pulse 2s ease-in-out infinite;
}
```

### Bounce
```css
@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-25%); }
}

.icon-bounce {
  animation: bounce 1s ease infinite;
}
```

## Responsividade

### Tamanhos Adaptativos

**Mobile**
```css
.icon { width: 16px; height: 16px; }
```

**Tablet**
```css
@media (min-width: 768px) {
  .icon { width: 20px; height: 20px; }
}
```

**Desktop**
```css
@media (min-width: 1024px) {
  .icon { width: 24px; height: 24px; }
}
```

## Acessibilidade

### ARIA

**Decorativo**
```svelte
<Search aria-hidden="true" />
```

**Significativo**
```svelte
<button aria-label="Buscar">
  <Search aria-hidden="true" />
</button>
```

### Focus Visible

```css
.icon-button:focus-visible {
  outline: 2px solid var(--primary-500);
  outline-offset: 2px;
}
```

### Color Contrast

- Contraste mínimo 3:1 para ícones
- Considerar fundo ao escolher cor
- Diferenciação além de cor quando necessário

## Performance

### SVG Inline

**Vantagens**
- Sem request HTTP adicional
- Customização via CSS
- Carregamento imediato
- Cache do HTML

**Uso**
```svelte
<Search size={20} class="icon" />
```

### SVG Sprite

**Para ícones muito usados**
- Reduz tamanho do bundle
- Cache eficiente
- Reutilização

## Customização

### Via CSS

```css
.icon {
  color: var(--primary-500);
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}
```

### Via Props (Lucide)

```svelte
<Search 
  size={20}
  color="#6366f1"
  strokeWidth={2}
  absoluteStrokeWidth
/>
```

## Emojis vs Ícones

### Quando Usar Emojis

- Contexto casual
- Comunicação direta
- Estados simples
- Feedback rápido

### Quando Usar Ícones

- Interface profissional
- Consistência visual
- Acessibilidade
- Customização necessária

### Decisão do PratoOnline

**Ícones (Lucide) para:**
- Navegação
- Ações principais
- Interface consistente
- Componentes UI

**Emojis para:**
- Estados rápidos (loading, sucesso)
- Feedback visual
- Comunicação informal
- Métricas (opcional)

## Padrões de Uso

### Consistência

**Mesmo ícone para mesma ação**
- Editar sempre = `edit`
- Excluir sempre = `trash-2`
- Salvar sempre = `save`

**Ícones semelhantes para ações relacionadas**
- Download = `download`
- Upload = `upload`
- Refresh = `refresh-cw`

### Intuição

**Ícones reconhecíveis**
- Busca = `search`
- Configurações = `settings`
- Perfil = `user`

**Evitar ambiguidade**
- Não usar ícones muito similares
- Contexto claro
- Labels quando necessário

## Tokens CSS

### Implementação

```css
:root {
  /* Sizes */
  --icon-xs: 14px;
  --icon-sm: 16px;
  --icon-md: 20px;
  --icon-lg: 24px;
  --icon-xl: 32px;
  
  /* Colors */
  --icon-default: var(--text-secondary);
  --icon-primary: var(--primary-500);
  --icon-success: var(--success-500);
  --icon-warning: var(--warning-500);
  --icon-error: var(--error-500);
  --icon-inverse: var(--text-inverse);
  --icon-disabled: var(--text-disabled);
}
```

### Uso

```css
.button-icon {
  width: var(--icon-md);
  height: var(--icon-md);
  color: var(--icon-default);
}
```

## Regras de Uso

### DO ✅

- Usar ícones consistentemente
- Escolher ícones intuitivos
- Considerar acessibilidade
- Tamanho apropriado ao contexto
- Cor semântica

### DON'T ❌

- Misturar estilos de ícone
- Usar ícones ambíguos
- Ignorar acessibilidade
- Tamanho muito pequeno
- Cor sem propósito

## Referências

### Recursos

- [Lucide Icons](https://lucide.dev/)
- [Heroicons](https://heroicons.com/) - Alternativa
- [Feather Icons](https://feathericons.com/) - Inspiração

### Ferramentas

- Figma icons plugin
- Icon picker
- SVG optimizer

---

**Versão**: 1.0  
**Data**: 15/07/2026  
**Sprint**: 9 - Product Experience
