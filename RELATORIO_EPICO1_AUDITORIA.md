# RELATÓRIO DE AUDITORIA - ÉPICO 1: PRODUTO COMERCIAL

**Data**: 16/07/2026  
**Objetivo**: Auditoria completa da arquitetura existente antes de implementar o Épico 1

---

## 1. ARQUITETURA BACKEND (Go + GORM)

### 1.1 Estrutura de Camadas (Clean Architecture)

```
backend/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── domain/                     # Entidades de domínio
│   │   ├── product.go              # Product entity
│   │   ├── ingredient.go           # Ingredient entity
│   │   ├── product_ingredient.go   # ProductIngredient (ficha técnica)
│   │   ├── order.go                # Order entity
│   │   ├── order_item.go           # OrderItem entity
│   │   ├── user.go                 # User entity
│   │   └── stock_adjustment_pending.go
│   ├── ports/                      # Interfaces de repositórios
│   │   ├── product_repository.go    # ProductRepository interface
│   │   ├── order_repository.go
│   │   ├── user_repository.go
│   │   └── stock_adjustment_repository.go
│   ├── infra/
│   │   ├── database/
│   │   │   ├── connection.go       # Conexão SQLite
│   │   │   └── migrate.go          # AutoMigrate
│   │   └── repository/
│   │       ├── gorm_product_repository.go    # Implementação GORM
│   │       ├── gorm_order_repository.go
│   │       ├── gorm_user_repository.go
│   │       └── gorm_stock_adjustment_repository.go
│   ├── service/                    # Lógica de negócio
│   │   ├── product_service.go      # ProductService
│   │   ├── order_service.go
│   │   ├── auth_service.go
│   │   └── stock_adjustment_service.go
│   ├── handler/                    # HTTP handlers
│   │   ├── product_handler.go      # ProductHandler
│   │   ├── order_handler.go
│   │   ├── auth_handler.go
│   │   └── stock_adjustment_handler.go
│   └── middleware/                 # Middlewares
│       └── auth_middleware.go
```

**Status**: ✅ Clean Architecture respeitada  
**Avaliação**: Arquitetura sólida, separação clara de responsabilidades

### 1.2 Entidade Product Atual

**Arquivo**: `backend/internal/domain/product.go`

```go
type Product struct {
    ID          uint
    Name        string
    Description string
    Price       float64
    IsComposto  bool
    Active      bool
    DeletedAt   *time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Ingredients []ProductIngredient  // Preenchido sob demanda
}
```

**Campos Existentes**:
- ✅ ID (uint)
- ✅ Name (string)
- ✅ Description (string)
- ✅ Price (float64)
- ✅ IsComposto (bool)
- ✅ Active (bool)
- ✅ DeletedAt (*time.Time) - Soft delete
- ✅ CreatedAt (time.Time)
- ✅ UpdatedAt (time.Time)
- ✅ Ingredients ([]ProductIngredient) - Relacionamento

**Campos FALTANDO para Épico 1**:
- ❌ PhotoURL (string) - Foto principal
- ❌ CategoryID (uint) - Relacionamento com Categoria
- ❌ DisplayOrder (int) - Ordem de exibição
- ❌ PreparationTimeMinutes (int) - Tempo de preparo
- ❌ Featured (bool) - Produto em destaque
- ❌ IsNew (bool) - Produto novo
- ❌ PromotionPrice (float64) - Preço promocional
- ❌ PromotionStart (*time.Time) - Início promoção
- ❌ PromotionEnd (*time.Time) - Fim promoção
- ❌ AvailableFrom (string) - Horário disponibilidade início
- ❌ AvailableUntil (string) - Horário disponibilidade fim
- ❌ SKU (string) - Código SKU (opcional)
- ❌ InternalNotes (string) - Observações internas

### 1.3 GORM Model Atual

**Arquivo**: `backend/internal/infra/repository/gorm_product_repository.go`

```go
type GormProduct struct {
    ID          uint   `gorm:"primaryKey;autoIncrement"`
    Name        string `gorm:"not null"`
    Description string
    Price       float64 `gorm:"not null;default:0"`
    IsComposto  bool    `gorm:"not null;default:false"`
    Active      bool    `gorm:"not null;default:true"`
    DeletedAt   *int64  `gorm:"index"`
    CreatedAt   int64   `gorm:"autoCreateTime"`
    UpdatedAt   int64   `gorm:"autoUpdateTime"`
}
```

**Status**: ✅ GORM model consistente com domain  
**Avaliação**: Precisa adicionar novos campos com tags GORM apropriadas

### 1.4 Repository Interface

**Arquivo**: `backend/internal/ports/product_repository.go`

**Métodos Existentes**:
- ✅ CreateProduct
- ✅ FindProductByID
- ✅ ListProducts
- ✅ ListActiveProducts
- ✅ UpdateProduct
- ✅ DeleteProduct
- ✅ CreateIngredient
- ✅ FindIngredientByID
- ✅ ListIngredients
- ✅ UpdateIngredient
- ✅ DeleteIngredient
- ✅ SetProductIngredients
- ✅ GetProductIngredients
- ✅ DecreaseIngredientStock
- ✅ IncreaseIngredientStock

**Métodos FALTANDO**:
- ❌ UploadProductPhoto
- ❌ DeleteProductPhoto
- ❌ CreateCategory
- ❌ ListCategories
- ❌ UpdateCategory
- ❌ DeleteCategory
- ❌ FindCategoryByID

### 1.5 Service Layer

**Arquivo**: `backend/internal/service/product_service.go`

**Inputs Existentes**:
- ✅ CreateProductInput
- ✅ UpdateProductInput
- ✅ CreateIngredientInput
- ✅ UpdateIngredientInput
- ✅ ProductIngredientInput
- ✅ SetProductIngredientsInput
- ✅ UpdateStockInput

**Inputs FALTANDO**:
- ❌ CreateCategoryInput
- ❌ UpdateCategoryInput
- ❌ UploadPhotoInput
- ❌ UpdateProductCommercialInput (com todos campos comerciais)

### 1.6 Handler Layer

**Arquivo**: `backend/internal/handler/product_handler.go`

**Endpoints Existentes**:
- ✅ POST /api/products
- ✅ GET /api/products
- ✅ GET /api/products/active
- ✅ GET /api/products/{id}
- ✅ PUT /api/products/{id}
- ✅ DELETE /api/products/{id}
- ✅ PUT /api/products/{id}/ingredients
- ✅ GET /api/products/{id}/ingredients
- ✅ POST /api/ingredients
- ✅ GET /api/ingredients
- ✅ GET /api/ingredients/{id}
- ✅ PUT /api/ingredients/{id}
- ✅ DELETE /api/ingredients/{id}
- ✅ PATCH /api/ingredients/{id}/stock

**Endpoints FALTANDO**:
- ❌ POST /api/products/{id}/photo
- ❌ DELETE /api/products/{id}/photo
- ❌ POST /api/categories
- ❌ GET /api/categories
- ❅ GET /api/categories/{id}
- ❌ PUT /api/categories/{id}
- ❌ DELETE /api/categories/{id}

### 1.7 Migrations

**Arquivo**: `backend/internal/infra/database/migrate.go`

```go
models := []interface{}{
    &repository.GormUserModel{},
    &repository.GormProduct{},
    &repository.GormIngredient{},
    &repository.GormProductIngredient{},
    &repository.GormOrder{},
    &repository.GormOrderItem{},
    &repository.GormStockAdjustmentPending{},
}
```

**Status**: ✅ AutoMigrate configurado  
**Avaliação**: Precisa adicionar GormCategory quando criado

---

## 2. ARQUITETURA FRONTEND (SvelteKit + TypeScript)

### 2.1 Estrutura de Camadas

```
frontend/src/
├── routes/
│   ├── (app)/
│   │   ├── products/
│   │   │   ├── +page.svelte          # Lista de produtos
│   │   │   └── [id]/
│   │   │       └── +page.svelte      # Detalhes do produto
│   │   ├── ingredients/
│   │   ├── orders/
│   │   ├── stock-adjustments/
│   │   └── profile/
│   └── (auth)/
├── lib/
│   ├── api/
│   │   ├── client.ts                 # HTTP client
│   │   ├── product.ts                # Product API
│   │   ├── ingredient.ts
│   │   ├── order.ts
│   │   └── auth.ts
│   ├── types/
│   │   ├── product.ts                # Product types
│   │   ├── ingredient.ts
│   │   ├── order.ts
│   │   └── user.ts
│   ├── components/
│   │   ├── ui/                       # UI components
│   │   └── layout/                   # Layout components
│   ├── stores/
│   └── theme/
```

**Status**: ✅ Arquitetura organizada  
**Avaliação**: Boa separação de responsabilidades

### 2.2 Tipos Product Atuais

**Arquivo**: `frontend/src/lib/types/product.ts`

```typescript
export interface Product {
  ID: number;
  Name: string;
  Description?: string;
  Price: number;
  IsComposto: boolean;
  Active: boolean;
  CreatedAt?: string;
  UpdatedAt?: string;
  Ingredients?: Ingredient[];
}

export interface ProductCreatePayload {
  name: string;
  description?: string;
  price: number;
  is_composto: boolean;
  active: boolean;
}
```

**Campos FALTANDO**:
- ❌ PhotoURL
- ❌ CategoryID
- ❌ Category (objeto relacionado)
- ❌ DisplayOrder
- ❌ PreparationTimeMinutes
- ❌ Featured
- ❌ IsNew
- ❌ PromotionPrice
- ❌ PromotionStart
- ❌ PromotionEnd
- ❌ AvailableFrom
- ❌ AvailableUntil
- ❌ SKU
- ❌ InternalNotes

### 2.3 API Client Product

**Arquivo**: `frontend/src/lib/api/product.ts`

**Funções Existentes**:
- ✅ getProducts()
- ✅ getActiveProducts()
- ✅ getProduct(id)
- ✅ createProduct(payload)
- ✅ updateProduct(id, payload)
- ✅ deleteProduct(id)
- ✅ getProductIngredients(id)
- ✅ updateProductIngredients(id, ingredients)

**Funções FALTANDO**:
- ❌ uploadProductPhoto(id, file)
- ❌ deleteProductPhoto(id)
- ❌ getCategories()
- ❌ createCategory(payload)
- ❌ updateCategory(id, payload)
- ❌ deleteCategory(id)

### 2.4 Interface Atual

**Arquivo**: `frontend/src/routes/(app)/products/+page.svelte`

**Funcionalidades Existentes**:
- ✅ Lista de produtos em grid
- ✅ Busca por nome
- ✅ Ordenação (nome, preço)
- ✅ Paginação
- ✅ Modal de criação/edição
- ✅ Indicadores (ativo, composto, estoque baixo)
- ✅ Loading states
- ✅ Empty states

**Funcionalidades FALTANDO**:
- ❌ Upload de foto
- ❌ Preview de foto
- ❌ Drag & drop
- ❌ Seleção de categoria
- ❌ Campos comerciais (destaque, novo, promoção, etc.)
- ❌ Interface com abas
- ❌ Validações avançadas
- ❌ Helper texts

**Arquivo**: `frontend/src/routes/(app)/products/[id]/+page.svelte`

**Funcionalidades Existentes**:
- ✅ Visualização de detalhes
- ✅ Edição de ficha técnica
- ✅ Adição/remoção de ingredientes

**Funcionalidades FALTANDO**:
- ❌ Abas de organização
- ❌ Edição de campos comerciais
- ❌ Upload de foto
- ❌ Edição de categoria
- ❌ Configuração de disponibilidade
- ❌ Configuração de promoção

---

## 3. BANCO DE DADOS

### 3.1 Tabela Products Atual

```sql
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    price REAL NOT NULL DEFAULT 0,
    is_composto INTEGER NOT NULL DEFAULT 0,
    active INTEGER NOT NULL DEFAULT 1,
    deleted_at INTEGER,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE INDEX idx_products_deleted_at ON products(deleted_at);
```

**Colunas FALTANDO**:
- ❌ photo_url TEXT
- ❌ category_id INTEGER
- ❌ display_order INTEGER DEFAULT 0
- ❌ preparation_time_minutes INTEGER
- ❌ featured INTEGER DEFAULT 0
- ❌ is_new INTEGER DEFAULT 0
- ❌ promotion_price REAL
- ❌ promotion_start INTEGER
- ❌ promotion_end INTEGER
- ❌ available_from TEXT
- ❌ available_until TEXT
- ❌ sku TEXT
- ❌ internal_notes TEXT

**Índices FALTANDO**:
- ❌ idx_products_category_id
- ❌ idx_products_display_order
- ❌ idx_products_featured
- ❌ idx_products_sku

### 3.2 Tabela Categories (NÃO EXISTE)

**Status**: ❌ Tabela não existe  
**Necessidade**: Criar tabela categories

```sql
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    active INTEGER NOT NULL DEFAULT 1,
    display_order INTEGER DEFAULT 0,
    deleted_at INTEGER,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);
CREATE INDEX idx_categories_active ON categories(active);
```

### 3.3 Tabelas Relacionadas

**Tabelas Existentes**:
- ✅ products
- ✅ ingredients
- ✅ product_ingredients
- ✅ orders
- ✅ order_items
- ✅ users
- ✅ stock_adjustments_pending

**Compatibilidade**:
- ✅ order_items tem snapshots de produto (ProductName, ProductDescription, ProductIsComposto)
- ⚠️ Precisa adicionar snapshots dos novos campos comerciais para preservar histórico

---

## 4. COMPATIBILIDADE RETROATIVA

### 4.1 Pedidos (Orders)

**Status**: ⚠️ PREOCUPAÇÃO

**Análise**:
- OrderItem atualmente tem snapshots: ProductName, ProductDescription, ProductIsComposto
- Ao adicionar novos campos ao Product, OrderItem PRECISA preservar snapshots dos campos comerciais
- Campos que precisam de snapshot em OrderItem:
  - PhotoURL (para exibir foto no pedido histórico)
  - CategoryID (para saber categoria no momento do pedido)
  - Price (já existe)
  - PromotionPrice (para saber se estava em promoção)
  - Featured (para saber se estava em destaque)
  - IsNew (para saber se estava marcado como novo)

**Impacto**: ALTO - precisa alterar OrderItem

### 4.2 Ficha Técnica (ProductIngredients)

**Status**: ✅ SEM IMPACTO

**Análise**:
- ProductIngredients não será alterado
- Relacionamento Produto-Ingredientes permanece igual

### 4.3 Estoque (Ingredients)

**Status**: ✅ SEM IMPACTO

**Análise**:
- Lógica de estoque não depende de campos comerciais
- DecreaseIngredientStock permanece igual

### 4.4 Ajustes de Estoque

**Status**: ✅ SEM IMPACTO

**Análise**:
- StockAdjustmentPending não depende de campos comerciais

---

## 5. IMPACTO DAS MUDANÇAS

### 5.1 Backend

**Alta Complexidade**:
- ✅ Adicionar campos ao domain.Product
- ✅ Atualizar GormProduct
- ✅ Atualizar mappers (productToDomain)
- ✅ Atualizar inputs (CreateProductInput, UpdateProductInput)
- ✅ Criar domain.Category
- ✅ Criar GormCategory
- ✅ Criar CategoryRepository
- ✅ Criar CategoryService
- ✅ Criar CategoryHandler
- ✅ Implementar upload de fotos
- ✅ Adicionar endpoints de foto
- ✅ Atualizar migrations

**Média Complexidade**:
- ⚠️ Atualizar OrderItem com snapshots dos novos campos
- ⚠️ Atualizar OrderRepository para salvar snapshots
- ⚠️ Atualizar OrderService para incluir snapshots

**Baixa Complexidade**:
- ✅ Atualizar rotas em main.go

### 5.2 Frontend

**Alta Complexidade**:
- ✅ Redesenhar interface com abas
- ✅ Implementar upload de foto com drag & drop
- ✅ Criar componentes de foto upload
- ✅ Adicionar todos campos comerciais ao formulário
- ✅ Implementar validações avançadas
- ✅ Melhorar UX com helpers e placeholders

**Média Complexidade**:
- ⚠️ Atualizar tipos Product e ProductCreatePayload
- ⚠️ Criar tipos Category
- ⚠️ Criar API client para categories
- ⚠️ Criar API client para photo upload
- ⚠️ Atualizar interface de lista para mostrar novos campos

**Baixa Complexidade**:
- ✅ Atualizar product.ts com novos campos

### 5.3 Banco de Dados

**Alta Complexidade**:
- ✅ Executar AutoMigrate para adicionar colunas
- ✅ Criar tabela categories
- ✅ Adicionar FK category_id em products
- ✅ Criar índices necessários
- ✅ Validar constraints

**Média Complexidade**:
- ⚠️ Atualizar tabela order_items com snapshots
- ⚠️ Migration para dados existentes (default values)

---

## 6. RISCOS IDENTIFICADOS

### 6.1 Riscos Altos

1. **Quebra de compatibilidade com pedidos históricos**
   - **Risco**: Pedidos antigos podem não ter snapshots dos novos campos
   - **Mitigação**: Adicionar campos NULLABLE em OrderItem, tratar valores nulos no frontend

2. **Upload de fotos**
   - **Risco**: Implementação de upload pode ser complexa (storage, validação, compressão)
   - **Mitigação**: Começar com upload simples local, evoluir para cloud storage depois

3. **Performance com novos campos**
   - **Risco**: Consultas podem ficar mais lentas com joins de categoria
   - **Mitigação**: Adicionar índices apropriados, usar eager loading

### 6.2 Riscos Médios

1. **Migração de banco de dados**
   - **Risco**: AutoMigrate pode falhar em produção
   - **Mitigação**: Testar extensivamente em ambiente de desenvolvimento

2. **Complexidade da interface**
   - **Risco**: Interface com muitos campos pode ficar confusa
   - **Mitigação**: Usar abas para organizar, seguir Design System

3. **Validações**
   - **Risco**: Campos comerciais podem ter regras complexas
   - **Mitigação**: Começar com validações simples, evoluir gradualmente

### 6.3 Riscos Baixos

1. **Tipos TypeScript**
   - **Risco**: Desync entre backend e frontend types
   - **Mitigação**: Manter tipos sincronizados, documentar mudanças

2. **Testes**
   - **Risco**: Falta de testes para novos campos
   - **Mitigação**: Criar testes unitários e integração

---

## 7. O QUE JÁ EXISTE (PONTO DE PARTIDA)

### 7.1 Backend

✅ **Solid Foundation**:
- Clean Architecture respeitada
- Separação clara de camadas
- GORM configurado
- Soft delete implementado
- Transações funcionando
- Validações com go-playground/validator
- Autenticação JWT funcionando
- Middleware de autenticação

✅ **Product CRUD Completo**:
- Create, Read, Update, Delete
- List com filtros
- Ficha técnica funcionando
- Estoque integrado

✅ **Infraestrutura**:
- SQLite configurado
- AutoMigrate funcionando
- Router chi configurado
- Logger configurado
- Timeout configurado

### 7.2 Frontend

✅ **Solid Foundation**:
- SvelteKit 5 com runes
- TypeScript configurado
- Component library existente
- Design System definido
- API client estruturado
- Layout responsivo

✅ **Product UI**:
- Lista de produtos funcional
- Busca e ordenação
- Paginação
- Modal de criação/edição
- Loading states
- Empty states
- Ficha técnica funcional

✅ **UX**:
- Microinterações
- Feedback visual
- Design profissional

---

## 8. O QUE FALTA (GAP ANALYSIS)

### 8.1 Backend - Domain

❌ **Entity Category**:
- Criar domain.Category
- Criar domain.Category com campos básicos

❌ **Product Fields**:
- PhotoURL
- CategoryID (FK)
- DisplayOrder
- PreparationTimeMinutes
- Featured
- IsNew
- PromotionPrice
- PromotionStart
- PromotionEnd
- AvailableFrom
- AvailableUntil
- SKU
- InternalNotes

### 8.2 Backend - Repository

❌ **CategoryRepository**:
- CreateCategory
- FindCategoryByID
- ListCategories
- UpdateCategory
- DeleteCategory

❌ **ProductRepository**:
- UploadProductPhoto
- DeleteProductPhoto

### 8.3 Backend - Service

❌ **CategoryService**:
- CreateCategory
- GetCategory
- ListCategories
- UpdateCategory
- DeleteCategory

❌ **ProductService**:
- UploadProductPhoto
- DeleteProductPhoto
- UpdateProductCommercial (todos campos)

### 8.4 Backend - Handler

❌ **CategoryHandler**:
- CreateCategory
- GetCategory
- ListCategories
- UpdateCategory
- DeleteCategory

❌ **ProductHandler**:
- UploadProductPhoto
- DeleteProductPhoto

### 8.5 Frontend - Types

❌ **Category Types**:
- Category interface
- CategoryCreatePayload
- CategoryUpdatePayload

❌ **Product Types**:
- Adicionar todos campos comerciais
- ProductCommercialPayload

### 8.6 Frontend - API

❌ **Category API**:
- getCategories()
- createCategory()
- updateCategory()
- deleteCategory()

❌ **Product API**:
- uploadProductPhoto()
- deleteProductPhoto()

### 8.7 Frontend - UI

❌ **Components**:
- PhotoUpload component (drag & drop, preview)
- CategorySelect component
- TabNavigation component

❌ **Pages**:
- Redesenhar products/[id]/+page.svelte com 3 abas
- Criar categories/+page.svelte (CRUD de categorias)

❌ **UX Improvements**:
- Helper texts
- Validações avançadas
- Loading elegante
- Mensagens claras

---

## 9. RECOMENDAÇÕES

### 9.1 Ordem de Implementação Sugerida

1. **Categoria (Foundation)**
   - Criar entity Category
   - Implementar CRUD completo
   - Adicionar FK em Product
   - Testar integração

2. **Campos Comerciais (Core)**
   - Adicionar campos ao Product
   - Atualizar todos layers
   - Atualizar OrderItem snapshots
   - Testar compatibilidade

3. **Foto (Complex)**
   - Implementar upload simples
   - Criar componentes frontend
   - Testar performance
   - Evoluir para drag & drop

4. **Interface (UX)**
   - Redesenhar com abas
   - Adicionar validações
   - Melhorar UX
   - Testar usabilidade

5. **Testes (Quality)**
   - Criar testes unitários
   - Criar testes integração
   - Testar compatibilidade
   - Validar Quality Gate

### 9.2 Boas Práticas a Seguir

✅ **Respeitar Clean Architecture**:
- Não quebrar separação de camadas
- Manter domain puro
- Usar ports para interfaces

✅ **Manter Compatibilidade**:
- Adicionar campos como NULLABLE
- Usar default values apropriados
- Preservar snapshots em OrderItem

✅ **Seguir Design System**:
- Usar componentes existentes
- Manter consistência visual
- Seguir UX guidelines

✅ **Quality First**:
- Testar cada mudança
- Executar Quality Gate
- Não acumular dívida técnica

### 9.3 Atenção Especial

⚠️ **OrderItem Snapshots**:
- CRÍTICO: Adicionar snapshots dos campos comerciais
- Testar pedidos históricos
- Validar exibição em histórico

⚠️ **Upload de Fotos**:
- Começar simples (local)
- Validar tipos e tamanhos
- Considerar compressão
- Evoluir para cloud storage

⚠️ **Migrations**:
- Testar AutoMigrate
- Validar dados existentes
- Backup antes de migrar
- Rollback planejado

---

## 10. CONCLUSÃO

### 10.1 Estado Atual

**Backend**: ✅ **SÓLIDO**
- Clean Architecture respeitada
- Fundação robusta
- Pronto para expansão

**Frontend**: ✅ **SÓLIDO**
- SvelteKit 5 moderno
- Component library funcional
- Design System definido

**Banco**: ✅ **ESTÁVEL**
- SQLite configurado
- AutoMigrate funcionando
- Soft delete implementado

### 10.2 Prontidão para Épico 1

**Avaliação Geral**: ✅ **PRONTO**

O sistema possui uma fundação sólida que permite implementar o Épico 1 sem quebrar a arquitetura existente. Os principais pontos de atenção são:

1. **Compatibilidade retroativa** (OrderItem snapshots)
2. **Upload de fotos** (complexidade técnica)
3. **Interface com abas** (complexidade UX)

### 10.3 Próximos Passos

1. ✅ Auditoria concluída
2. ⏭️ Implementar ETAPA 2 (Foto)
3. ⏭️ Implementar ETAPA 3 (Categoria)
4. ⏭️ Implementar ETAPA 4-11 (Campos comerciais)
5. ⏭️ Implementar ETAPA 12 (Interface)
6. ⏭️ Implementar ETAPA 13 (UX)
7. ⏭️ Executar ETAPA 14 (Migrations)
8. ⏭️ Validar ETAPA 15 (Compatibilidade)
9. ⏭️ Criar ETAPA 16 (Testes)
10. ⏭️ Executar ETAPA 17 (Quality Gate)
11. ⏭️ Gerar ETAPA 18 (Relatório Final)

---

**Assinatura**: Cascade AI Assistant  
**Aprovação**: Pendente revisão do usuário
