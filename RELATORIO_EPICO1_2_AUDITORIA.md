# RELATÓRIO DE AUDITORIA — ÉPICO 1.2: PRODUTO COMERCIAL COMPLETO

**Data:** 16 de Julho de 2026  
**Objetivo:** Auditoria completa da arquitetura existente antes de implementar o Épico 1.2

---

## 1. ARQUITETURA BACKEND (Go + GORM)

### 1.1 Estrutura de Camadas (Clean Architecture)

```
backend/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── domain/                     # Entidades de domínio
│   │   ├── product.go              # Product entity
│   │   ├── category.go             # Category entity
│   │   ├── ingredient.go           # Ingredient entity
│   │   ├── product_ingredient.go   # ProductIngredient (ficha técnica)
│   │   ├── order.go                # Order entity
│   │   ├── order_item.go           # OrderItem entity (com snapshots)
│   │   ├── user.go                 # User entity
│   │   └── stock_adjustment_pending.go
│   ├── ports/                      # Interfaces de repositórios
│   │   ├── product_repository.go    # ProductRepository interface
│   │   ├── category_repository.go   # CategoryRepository interface
│   │   ├── order_repository.go      # OrderRepository interface
│   │   ├── user_repository.go       # UserRepository interface
│   │   └── stock_adjustment_repository.go
│   ├── infra/
│   │   ├── database/
│   │   │   └── db.go               # Conexão GORM
│   │   └── repository/
│   │       ├── gorm_product_repository.go
│   │       ├── gorm_category_repository.go
│   │       ├── gorm_order_repository.go
│   │       └── gorm_user_repository.go
│   ├── service/
│   │   ├── product_service.go      # Lógica de negócio
│   │   ├── category_service.go
│   │   ├── order_service.go
│   │   └── user_service.go
│   ├── handler/
│   │   ├── product_handler.go      # Endpoints HTTP
│   │   ├── category_handler.go
│   │   ├── order_handler.go
│   │   └── user_handler.go
│   └── middleware/
│       └── auth.go
```

### 1.2 Domain: Product

**Arquivo:** `backend/internal/domain/product.go`

```go
type Product struct {
    ID                     uint
    Name                   string
    Description            string
    Price                  float64
    IsComposto             bool
    Active                 bool // "Pode ser utilizado pelo negócio?"
    PhotoURL               string
    CategoryID             *uint
    DisplayOrder           int
    PreparationTimeMinutes int
    Featured               bool
    IsNew                  bool
    PromotionPrice         *float64
    PromotionStart         *time.Time
    PromotionEnd           *time.Time
    AvailableFrom          string
    AvailableUntil         string
    SKU                    string
    InternalNotes          string
    DeletedAt              *time.Time // "O registro foi removido logicamente"
    CreatedAt              time.Time
    UpdatedAt              time.Time

    // Preenchido sob demanda (não vem do banco direto)
    Ingredients []ProductIngredient
}
```

**Campos Atuais: 18**

**Campos Faltando para SEO:**
- Slug (string)
- MetaTitle (string)
- MetaDescription (string)
- AltImage (string)
- Canonical (string)

**Campos Faltando para iFood:**
- ExternalID (string)
- MarketplaceID (string)
- SyncStatus (string)
- LastSync (*time.Time)

**Campos Faltando para Arquivamento:**
- Archived (bool) - ou usar Active como arquivamento?

**Observações:**
- Soft delete implementado com DeletedAt
- Sem campo específico para arquivamento (pode usar Active ou criar Archived)
- Sem campos de SEO para Cardápio Digital
- Sem campos de integração com marketplaces

### 1.3 Domain: OrderItem (Snapshots)

**Arquivo:** `backend/internal/domain/order_item.go`

```go
type OrderItem struct {
    ID                    uint
    OrderID               uint
    ProductID             uint
    Quantity              float64
    UnitPrice             float64    // snapshot do preço
    ProductName           string     // snapshot do nome
    ProductDescription    string     // snapshot da descrição
    ProductIsComposto     bool       // snapshot da flag
    ProductPhotoURL       string     // snapshot da foto
    ProductCategoryID     *uint      // snapshot da categoria
    ProductPromotionPrice *float64   // snapshot do preço promocional
    ProductFeatured       bool       // snapshot do destaque
    ProductIsNew          bool       // snapshot do selo novo
    DeletedAt             *time.Time
    Product               *Product   // para navegação
}
```

**Observações:**
- Snapshots já implementados para campos comerciais
- Princípio #4 respeitado: histórico é imutável
- Precisa adicionar snapshots para novos campos (Slug, etc.)

### 1.4 Repository: ProductRepository

**Arquivo:** `backend/internal/ports/product_repository.go`

```go
type ProductRepository interface {
    // Produto
    CreateProduct(ctx context.Context, p *domain.Product) error
    FindProductByID(ctx context.Context, id uint) (*domain.Product, error)
    ListProducts(ctx context.Context) ([]domain.Product, error)
    ListActiveProducts(ctx context.Context) ([]domain.Product, error)
    UpdateProduct(ctx context.Context, p *domain.Product) error
    DeleteProduct(ctx context.Context, id uint) error

    // Ingrediente
    CreateIngredient(ctx context.Context, i *domain.Ingredient) error
    FindIngredientByID(ctx context.Context, id uint) (*domain.Ingredient, error)
    ListIngredients(ctx context.Context) ([]domain.Ingredient, error)
    UpdateIngredient(ctx context.Context, i *domain.Ingredient) error
    DeleteIngredient(ctx context.Context, id uint) error

    // Ficha técnica
    SetProductIngredients(ctx context.Context, productID uint, items []domain.ProductIngredient) error
    GetProductIngredients(ctx context.Context, productID uint) ([]domain.ProductIngredient, error)

    // Estoque
    DecreaseIngredientStock(ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB, ingredientName string, currentStock float64) error
    IncreaseIngredientStock(ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB) error
}
```

**Observações:**
- Sem métodos para busca avançada (filtros, ordenação)
- Sem método para duplicar produto
- Sem método para arquivar/desarquivar
- Sem método para ações rápidas (toggle Active, Featured, etc.)

### 1.5 Service: ProductService

**Arquivo:** `backend/internal/service/product_service.go`

```go
type ProductService struct {
    repo ports.ProductRepository
}

// Inputs
type CreateProductInput struct {
    Name                   string
    Description            string
    Price                  float64
    IsComposto             bool
    PhotoURL               string
    CategoryID             *uint
    DisplayOrder           int
    PreparationTimeMinutes int
    Featured               bool
    IsNew                  bool
    PromotionPrice         *float64
    PromotionStart         *time.Time
    PromotionEnd           *time.Time
    AvailableFrom          string
    AvailableUntil         string
    SKU                    string
    InternalNotes          string
}

type UpdateProductInput struct {
    // Mesmos campos + Active
}
```

**Observações:**
- Validações básicas com struct tags
- Sem lógica para slug automático
- Sem lógica para compressão de imagens
- Sem lógica para arquivamento

### 1.6 Handler: ProductHandler

**Arquivo:** `backend/internal/handler/product_handler.go`

**Endpoints Atuais:**
- `POST /api/products` - Criar produto
- `GET /api/products` - Listar todos
- `GET /api/products/active` - Listar ativos
- `GET /api/products/{id}` - Buscar por ID
- `PUT /api/products/{id}` - Atualizar
- `DELETE /api/products/{id}` - Deletar (soft)
- `PUT /api/products/{id}/ingredients` - Configurar ficha técnica
- `GET /api/products/{id}/ingredients` - Buscar ficha técnica

**Endpoints Faltando:**
- `POST /api/products/{id}/duplicate` - Duplicar produto
- `PATCH /api/products/{id}/archive` - Arquivar
- `PATCH /api/products/{id}/unarchive` - Desarquivar
- `PATCH /api/products/{id}/toggle-active` - Toggle Active
- `PATCH /api/products/{id}/toggle-featured` - Toggle Featured
- `POST /api/media/upload` - Upload de mídia
- `DELETE /api/media/{id}` - Deletar mídia
- `GET /api/products/search?q=` - Busca avançada
- `GET /api/products?filter=&sort=` - Listagem com filtros e ordenação

---

## 2. ARQUITETURA FRONTEND (SvelteKit 5)

### 2.1 Estrutura

```
frontend/
├── src/
│   ├── lib/
│   │   ├── api/
│   │   │   ├── client.ts           # Cliente HTTP
│   │   │   ├── product.ts          # API de produtos
│   │   │   ├── category.ts         # API de categorias
│   │   │   └── ingredient.ts       # API de ingredientes
│   │   ├── types/
│   │   │   ├── product.ts          # Tipos de produto
│   │   │   ├── category.ts         # Tipos de categoria
│   │   │   └── ingredient.ts       # Tipos de ingrediente
│   │   ├── components/
│   │   │   ├── ui/
│   │   │   │   ├── Button.svelte
│   │   │   │   ├── Input.svelte
│   │   │   │   ├── Card.svelte
│   │   │   │   ├── Modal.svelte
│   │   │   │   ├── Table.svelte
│   │   │   │   ├── TabNavigation.svelte  # NOVO (ÉPICO 1.1)
│   │   │   │   ├── PhotoUpload.svelte    # NOVO (ÉPICO 1.1)
│   │   │   │   └── ...
│   │   │   └── layout/
│   │   │       └── Workspace.svelte
│   │   └── stores/
│   │       └── auth.ts
│   └── routes/
│       ├── (app)/
│       │   ├── products/
│       │   │   ├── +page.svelte           # Listagem
│       │   │   ├── new/+page.svelte       # NOVO (ÉPICO 1.1)
│       │   │   └── [id]/
│       │   │       └── edit/+page.svelte  # NOVO (ÉPICO 1.1)
│       │   ├── categories/+page.svelte
│       │   └── ...
```

### 2.2 Types: Product

**Arquivo:** `frontend/src/lib/types/product.ts`

```typescript
export interface Product {
  ID: number;
  Name: string;
  Description?: string;
  Price: number;
  IsComposto: boolean;
  Active: boolean;
  PhotoURL?: string;
  CategoryID?: number;
  DisplayOrder: number;
  PreparationTimeMinutes: number;
  Featured: boolean;
  IsNew: boolean;
  PromotionPrice?: number;
  PromotionStart?: string;
  PromotionEnd?: string;
  AvailableFrom?: string;
  AvailableUntil?: string;
  SKU?: string;
  InternalNotes?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
  Ingredients?: Ingredient[];
}
```

**Campos Faltando:**
- Slug
- MetaTitle
- MetaDescription
- AltImage
- Canonical
- ExternalID
- MarketplaceID
- SyncStatus
- LastSync
- Archived
- ThumbnailURL (para performance)

### 2.3 API: Product

**Arquivo:** `frontend/src/lib/api/product.ts`

```typescript
export async function getProducts(): Promise<Product[]>
export async function getActiveProducts(): Promise<Product[]>
export async function getProduct(id: number): Promise<Product>
export async function createProduct(payload: ProductCreatePayload): Promise<Product>
export async function updateProduct(id: number, payload: ProductUpdatePayload): Promise<Product>
export async function deleteProduct(id: number): Promise<void>
export async function getProductIngredients(id: number): Promise<Ingredient[]>
export async function updateProductIngredients(id: number, ingredients: ProductIngredientPayload[]): Promise<void>
```

**Funções Faltando:**
- `duplicateProduct(id: number): Promise<Product>`
- `archiveProduct(id: number): Promise<void>`
- `unarchiveProduct(id: number): Promise<void>`
- `toggleProductActive(id: number): Promise<Product>`
- `toggleProductFeatured(id: number): Promise<Product>`
- `uploadMedia(file: File): Promise<MediaUploadResponse>`
- `deleteMedia(mediaId: string): Promise<void>`
- `searchProducts(query: string, filters: ProductFilters): Promise<Product[]>`

### 2.4 Listagem: Products Page

**Arquivo:** `frontend/src/routes/(app)/products/+page.svelte`

**Estado Atual:**
- Tabela simples para listagem
- Busca básica por nome
- Ordenação por nome e preço
- Paginação (12 itens por página)
- Modal de criação/edição (substituído por páginas no ÉPICO 1.1)
- Skeleton loading

**Limitações:**
- Sem cards elegantes
- Sem filtros avançados (Ativos, Arquivados, Promoção, etc.)
- Sem ordenação avançada (Categoria, Tempo, Ordem, Mais vendidos)
- Sem ações rápidas (toggle sem abrir edição)
- Sem preview completo (modal com detalhes)
- Sem lazy loading de imagens
- Sem indicação visual de disponibilidade
- Sem badge de categoria
- Sem preço promocional
- Sem selo de novo/destaque

### 2.5 Cadastro: Product Pages

**Arquivo:** `frontend/src/routes/(app)/products/new/+page.svelte` (NOVO)
**Arquivo:** `frontend/src/routes/(app)/products/[id]/edit/+page.svelte` (NOVO)

**Estado Atual (após ÉPICO 1.1):**
- Estrutura de 3 abas (Informações, Venda, Produção)
- Exposição dos 18 campos
- Helper texts
- Validação visual
- PhotoUpload com preview local
- Cabeçalho executivo

**Limitações:**
- PhotoUpload apenas preview local (sem upload real)
- Sem compressão de imagens
- Sem geração de thumbnail
- Sem slug automático
- Sem campos SEO
- Sem campos de integração iFood

---

## 3. MÓDULO DE MÍDIA

### 3.1 Estado Atual

**Backend:**
- ❌ Sem módulo de mídia
- ❌ Sem serviço de upload
- ❌ Sem compressão de imagens
- ❌ Sem geração de thumbnails
- ❌ Sem estrutura de pastas para uploads

**Frontend:**
- ✅ PhotoUpload component (preview local)
- ❌ Sem upload real
- ❌ Sem compressão no cliente
- ❌ Sem integração com backend

### 3.2 Necessidades

**Backend:**
- Criar `internal/domain/media.go` - Media entity
- Criar `internal/ports/media_repository.go` - MediaRepository interface
- Criar `internal/infra/repository/gorm_media_repository.go` - Implementação
- Criar `internal/service/media_service.go` - Lógica de negócio
- Criar `internal/handler/media_handler.go` - Endpoints HTTP
- Criar pasta `uploads/products/` e `uploads/products/thumbs/`
- Implementar compressão WEBP (1920px)
- Implementar thumbnail (400px)
- Implementar exclusão de arquivos órfãos

**Frontend:**
- Criar `frontend/src/lib/types/media.ts` - Tipos de mídia
- Criar `frontend/src/lib/api/media.ts` - API de mídia
- Atualizar PhotoUpload para upload real
- Implementar lazy loading
- Usar thumbnail na listagem, imagem completa no preview

---

## 4. GAPS IDENTIFICADOS

### 4.1 Backend

| Funcionalidade | Status | Prioridade |
|---------------|--------|------------|
| Media Service | ❌ Não existe | Alta |
| Upload de imagens | ❌ Não existe | Alta |
| Compressão WEBP | ❌ Não existe | Alta |
| Geração de thumbnails | ❌ Não existe | Alta |
| Exclusão de arquivos órfãos | ❌ Não existe | Alta |
| Slug automático | ❌ Não existe | Média |
| Campos SEO | ❌ Não existe | Média |
| Campos iFood | ❌ Não existe | Média |
| Arquivamento | ⚠️ Usa soft delete | Média |
| Duplicação de produto | ❌ Não existe | Média |
| Ações rápidas (toggle) | ❌ Não existe | Média |
| Busca avançada | ⚠️ Apenas básica | Média |
| Filtros avançados | ❌ Não existe | Média |
| Ordenação avançada | ❌ Não existe | Média |
| Snapshots de novos campos | ❌ Não existe | Baixa |

### 4.2 Frontend

| Funcionalidade | Status | Prioridade |
|---------------|--------|------------|
| Upload real de imagens | ❌ Preview local | Alta |
| Compressão no cliente | ❌ Não existe | Alta |
| Cards de produtos elegantes | ❌ Usa tabela | Alta |
| Lazy loading | ❌ Não existe | Alta |
| Filtros avançados | ❌ Não existe | Alta |
| Ordenação avançada | ⚠️ Básica | Alta |
| Ações rápidas | ❌ Não existe | Alta |
| Preview completo (modal) | ❌ Não existe | Alta |
| Indicação visual de disponibilidade | ❌ Não existe | Alta |
| Duplicação de produto | ❌ Não existe | Média |
| Arquivamento visual | ❌ Não existe | Média |
| Slug automático | ❌ Não existe | Média |
| Campos SEO | ❌ Não existe | Média |
| Campos iFood | ❌ Não existe | Média |
| Badge de categoria | ❌ Não existe | Baixa |
| Preço promocional | ❌ Não existe | Baixa |
| Selo novo/destaque | ⚠️ Existe no backend | Baixa |

---

## 5. RISCOS

### 5.1 Arquiteturais

- **Risco:** Adicionar campos ao Product pode quebrar snapshots existentes
- **Mitigação:** Adicionar campos como nullable, migrar banco de dados
- **Risco:** Media Service pode criar dependência com S3/Cloud Storage
- **Mitigação:** Começar com armazenamento local, preparar interface para migração

### 5.2 Performance

- **Risco:** Upload de imagens grandes pode travar a aplicação
- **Mitigação:** Compressão no cliente, limitar tamanho (5MB)
- **Risco:** Lazy loading pode causar layout shift
- **Mitigação:** Usar placeholder com aspect ratio

### 5.3 UX

- **Risco:** Muitos filtros podem confundir usuários
- **Mitigação:** Filtros padrão inteligentes, salvar preferências
- **Risco:** Arquivamento vs Active pode confundir
- **Mitigação:** Terminologia clara, documentação

---

## 6. RECOMENDAÇÕES

### 6.1 Backend

1. **Criar Media Service antes de implementar upload real**
   - Domain: media.go
   - Repository: media_repository.go
   - Service: media_service.go
   - Handler: media_handler.go

2. **Adicionar campos SEO ao Product**
   - Slug (string, unique)
   - MetaTitle (string)
   - MetaDescription (string)
   - AltImage (string)
   - Canonical (string)

3. **Adicionar campos iFood ao Product**
   - ExternalID (string)
   - MarketplaceID (string)
   - SyncStatus (string)
   - LastSync (*time.Time)

4. **Implementar arquivamento**
   - Criar campo Archived (bool)
   - Manter Active para status de disponibilidade
   - Adicionar endpoints de archive/unarchive

5. **Adicionar métodos ao ProductRepository**
   - ListArchivedProducts
   - SearchProducts (com filtros e ordenação)
   - DuplicateProduct

6. **Atualizar OrderItem snapshots**
   - Adicionar snapshots dos novos campos

### 6.2 Frontend

1. **Atualizar PhotoUpload para upload real**
   - Integrar com Media Service
   - Compressão no cliente (opcional)
   - Mostrar progresso

2. **Criar ProductCard component**
   - Card elegante com foto, categoria, preço
   - Badges (novo, destaque, promoção)
   - Indicação de disponibilidade
   - Ações rápidas (menu dropdown)

3. **Implementar filtros avançados**
   - Todos, Ativos, Arquivados, Promoção, Novidades, Destaques, Compostos, Categoria
   - Filtro combinável

4. **Implementar ordenação avançada**
   - Nome, Preço, Categoria, Tempo, Ordem
   - Mais vendidos (preparado com contador)
   - Mais recentes

5. **Criar ProductPreview modal**
   - Foto grande
   - Todos os detalhes
   - Ingredientes
   - Disponibilidade

6. **Implementar lazy loading**
   - Thumbnail na listagem
   - Imagem completa no preview
   - Intersection Observer

### 6.3 Migração

1. **Migração do banco de dados**
   - Adicionar colunas ao products
   - Adicionar tabela media
   - Migrar PhotoURL existente para media

2. **Migração de snapshots**
   - Atualizar OrderItem com novos campos
   - Migrar dados históricos (opcional)

---

## 7. PRÓXIMOS PASSOS

### Ordem Sugerida de Implementação

1. **ETAPA 2:** Media Service (Upload, Delete, Resize, Compress, Thumbnail, Storage)
2. **ETAPA 3-6:** Upload Real, Compressão, Estrutura Física, Exclusão de arquivos órfãos
3. **ETAPA 7-8:** Cards de Produto, Visual da Listagem
4. **ETAPA 9-11:** Busca Instantânea, Filtros, Ordenação
5. **ETAPA 12-14:** Duplicar, Arquivar, Ações Rápidas
6. **ETAPA 15:** Preview Completo
7. **ETAPA 16-17:** Preparação Cardápio Digital, SEO
8. **ETAPA 18:** Disponibilidade Inteligente
9. **ETAPA 19:** Performance (Lazy loading)
10. **ETAPA 20:** Preparação iFood
11. **ETAPA 21:** Responsividade
12. **ETAPA 22:** Quality Gate
13. **ETAPA 23:** Relatório Final

---

## 8. CONCLUSÃO

A arquitetura atual está bem estruturada seguindo Clean Architecture, mas faltam módulos essenciais para um sistema comercial completo:

- **Módulo de Mídia:** Não existe, precisa ser criado do zero
- **Campos SEO:** Não implementados, necessários para Cardápio Digital
- **Campos iFood:** Não implementados, necessários para integração
- **Listagem:** Usa tabela simples, precisa de cards elegantes
- **Filtros/Ordenação:** Básicos, precisam ser avançados
- **Ações Rápidas:** Não existem, necessárias para produtividade
- **Performance:** Sem lazy loading, necessário para muitas imagens

O ÉPICO 1.2 é essencial para transformar o módulo Produto em um sistema comercial de nível profissional, preparado para múltiplos canais de venda (PDV, Cardápio Digital, QR Code, iFood).

---

**Assinatura:**  
Auditoria Técnica - PratoOnline  
**Status:** Auditoria Concluída
