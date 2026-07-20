<script lang="ts">
  import { MoreHorizontal, Star, Tag, Clock, Package } from '@lucide/svelte';
  import type { Product } from '$lib/types/product';
  import { Badge } from '$lib/components/ui';

  interface Props {
    product: Product;
    onClick?: () => void;
    onEdit?: () => void;
    onDuplicate?: () => void;
    onArchive?: () => void;
    onToggleActive?: () => void;
    onToggleFeatured?: () => void;
    loadingArchive?: boolean;
    loadingToggleActive?: boolean;
    loadingToggleFeatured?: boolean;
    loadingDuplicate?: boolean;
  }

  let { product, onClick, onEdit, onDuplicate, onArchive, onToggleActive, onToggleFeatured, loadingArchive, loadingToggleActive, loadingToggleFeatured, loadingDuplicate }: Props = $props();

  let showMenu = $state(false);

  function formatPrice(value: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(value);
  }

  function formatPromotionPrice(value: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(value);
  }

  function getAvailabilityStatus() {
    if (!product.Active) return 'archived';
    if (product.PromotionPrice && product.PromotionStart && product.PromotionEnd) {
      const now = new Date();
      const start = new Date(product.PromotionStart);
      const end = new Date(product.PromotionEnd);
      if (now >= start && now <= end) return 'promotion';
    }
    return 'available';
  }

  function getAvailabilityColor() {
    const status = getAvailabilityStatus();
    switch (status) {
      case 'available': return '#10b981'; // green
      case 'promotion': return '#f97316'; // orange
      case 'archived': return '#94a3b8'; // gray
      default: return '#94a3b8';
    }
  }

  function getAvailabilityLabel() {
    const status = getAvailabilityStatus();
    switch (status) {
      case 'available': return 'Disponível';
      case 'promotion': return 'Promoção';
      case 'archived': return 'Arquivado';
      default: return '';
    }
  }
</script>

<div class="product-card" onclick={onClick} role="button" tabindex="0">
  <div class="card-header">
    <div class="photo-container">
      {#if product.PhotoURL}
        <img src={product.PhotoURL} alt={product.Name} class="product-photo" loading="lazy" />
      {:else}
        <div class="product-photo-placeholder">
          <Package size={32} />
        </div>
      {/if}
      <div class="availability-indicator" style="background-color: {getAvailabilityColor()}"></div>
    </div>
    
    <div class="menu-container">
      <button 
        class="menu-button" 
        onclick={(e) => { e.stopPropagation(); showMenu = !showMenu; }}
        title="Mais opções"
      >
        <MoreHorizontal size={18} />
      </button>
      
      {#if showMenu}
        <div class="menu-dropdown">
          <button onclick={(e) => { e.stopPropagation(); showMenu = false; onEdit?.(); }}>
            Editar
          </button>
          <button onclick={(e) => { e.stopPropagation(); showMenu = false; onDuplicate?.(); }} disabled={loadingDuplicate}>
            {loadingDuplicate ? 'Duplicando...' : 'Duplicar'}
          </button>
          <button onclick={(e) => { e.stopPropagation(); showMenu = false; onToggleActive?.(); }} disabled={loadingToggleActive}>
            {loadingToggleActive ? 'Processando...' : (product.Active ? 'Desativar' : 'Ativar')}
          </button>
          <button onclick={(e) => { e.stopPropagation(); showMenu = false; onToggleFeatured?.(); }} disabled={loadingToggleFeatured}>
            {loadingToggleFeatured ? 'Processando...' : (product.Featured ? 'Remover destaque' : 'Destacar')}
          </button>
          <button onclick={(e) => { e.stopPropagation(); showMenu = false; onArchive?.(); }} disabled={loadingArchive}>
            {loadingArchive ? 'Arquivando...' : 'Arquivar'}
          </button>
        </div>
      {/if}
    </div>
  </div>

  <div class="card-body">
    <div class="badges">
      {#if product.IsNew}
        <Badge variant="success" size="sm">NOVO</Badge>
      {/if}
      {#if product.Featured}
        <Badge variant="warning" size="sm"><Star size={10} /> Destaque</Badge>
      {/if}
      {#if product.IsComposto}
        <Badge variant="info" size="sm">Composto</Badge>
      {/if}
    </div>

    <h3 class="product-name">{product.Name}</h3>
    
    <div class="product-meta">
      <span class="meta-item">
        <Tag size={12} />
        {product.CategoryID ? 'Categoria' : 'Sem categoria'}
      </span>
      {#if product.PreparationTimeMinutes > 0}
        <span class="meta-item">
          <Clock size={12} />
          {product.PreparationTimeMinutes} min
        </span>
      {/if}
    </div>

    <div class="price-section">
      {#if product.PromotionPrice}
        <span class="original-price">{formatPrice(product.Price)}</span>
        <span class="promotion-price">{formatPromotionPrice(product.PromotionPrice)}</span>
      {:else}
        <span class="regular-price">{formatPrice(product.Price)}</span>
      {/if}
    </div>
  </div>
</div>

<style>
  .product-card {
    background: white;
    border: 1px solid #f1f5f9;
    border-radius: 12px;
    overflow: hidden;
    transition: all 0.15s ease;
    cursor: pointer;
    position: relative;
  }

  .product-card:hover {
    border-color: #e2e8f0;
    box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
    transform: translateY(-2px);
  }

  .product-card:focus {
    outline: none;
    box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.5);
  }

  .card-header {
    position: relative;
    height: 160px;
  }

  .photo-container {
    width: 100%;
    height: 100%;
    position: relative;
  }

  .product-photo {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .product-photo-placeholder {
    width: 100%;
    height: 100%;
    background: #f8fafc;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #94a3b8;
  }

  .availability-indicator {
    position: absolute;
    top: 0.5rem;
    left: 0.5rem;
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }

  .menu-container {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
  }

  .menu-button {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.9);
    border: none;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: #64748b;
    transition: all 0.15s ease;
  }

  .menu-button:hover {
    background: white;
    color: #0f172a;
  }

  .menu-dropdown {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 0.25rem;
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
    min-width: 140px;
    z-index: 10;
  }

  .menu-dropdown button {
    width: 100%;
    padding: 0.5rem 1rem;
    text-align: left;
    border: none;
    background: none;
    font-size: 0.875rem;
    color: #0f172a;
    cursor: pointer;
    transition: background 0.15s ease;
  }

  .menu-dropdown button:hover {
    background: #f8fafc;
  }

  .card-body {
    padding: 1rem;
  }

  .badges {
    display: flex;
    gap: 0.25rem;
    flex-wrap: wrap;
    margin-bottom: 0.5rem;
  }

  .product-name {
    font-size: 0.875rem;
    font-weight: 600;
    color: #0f172a;
    margin: 0 0 0.5rem 0;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .product-meta {
    display: flex;
    gap: 0.75rem;
    margin-bottom: 0.75rem;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 0.75rem;
    color: #64748b;
  }

  .price-section {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .regular-price {
    font-size: 1rem;
    font-weight: 700;
    color: #0f172a;
  }

  .original-price {
    font-size: 0.75rem;
    color: #94a3b8;
    text-decoration: line-through;
  }

  .promotion-price {
    font-size: 1rem;
    font-weight: 700;
    color: #f97316;
  }
</style>
