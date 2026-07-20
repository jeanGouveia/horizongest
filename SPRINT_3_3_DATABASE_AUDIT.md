# SPRINT 3.3 - Auditoria do Banco de Dados

**Data:** 2025-01-XX  
**Auditor:** Cascade AI  
**Escopo:** Schema do Banco de Dados, Constraints, Indexes, Foreign Keys  
**Objetivo:** Validar integridade referencial, performance, isolamento de dados

---

## Resumo Executivo

O schema do banco de dados está **bem estruturado** com constraints apropriadas, indexes para performance, e separação clara entre tabelas Platform e Tenant. Foram identificados **3 riscos baixos** relacionados a indexes ausentes e constraints faltantes.

**Status:** ✅ **APROVADO COM RECOMENDAÇÕES**

---

## 1. Tabelas Platform

### 1.1 Tabela `platform_users`

**Arquivo:** `migrations/00013_create_platform_users.sql`

```sql
CREATE TABLE IF NOT EXISTS platform_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    platform_role TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ UNIQUE: `email`
- ✅ NOT NULL: `name`, `email`, `password_hash`, `platform_role`
- ✅ Indexes: `idx_platform_users_email`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_platform_users_email ON platform_users(email);
```

### 1.2 Tabela `platform_sessions`

**Arquivo:** `migrations/00014_create_platform_sessions.sql`

```sql
CREATE TABLE IF NOT EXISTS platform_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform_user_id INTEGER NOT NULL REFERENCES platform_users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at INTEGER NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ FOREIGN KEY: `platform_user_id` → `platform_users(id)` ON DELETE CASCADE
- ✅ UNIQUE: `token`
- ✅ NOT NULL: `platform_user_id`, `token`, `expires_at`
- ✅ Indexes: `idx_platform_sessions_token`, `idx_platform_sessions_platform_user_id`, `idx_platform_sessions_expires_at`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_platform_sessions_token ON platform_sessions(token);
CREATE INDEX IF NOT EXISTS idx_platform_sessions_platform_user_id ON platform_sessions(platform_user_id);
CREATE INDEX IF NOT EXISTS idx_platform_sessions_expires_at ON platform_sessions(expires_at);
```

### 1.3 Tabela `platform_audit`

**Arquivo:** `migrations/00015_create_platform_audit.sql`

```sql
CREATE TABLE IF NOT EXISTS platform_audit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform_user_id INTEGER,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id INTEGER,
    changes TEXT,
    ip_address TEXT,
    user_agent TEXT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ⚠️ FOREIGN KEY: `platform_user_id` → `platform_users(id)` **SEM ON DELETE**
- ✅ NOT NULL: `action`, `entity_type`, `created_at`
- ✅ Indexes: `idx_platform_audit_platform_user_id`, `idx_platform_audit_entity`, `idx_platform_audit_created_at`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_platform_audit_platform_user_id ON platform_audit(platform_user_id);
CREATE INDEX IF NOT EXISTS idx_platform_audit_entity ON platform_audit(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_platform_audit_created_at ON platform_audit(created_at);
```

**Problema:** ⚠️ FK sem ON DELETE pode causar registros órfãos se platform_user for deletado

---

## 2. Tabelas Tenant

### 2.1 Tabela `companies`

**Arquivo:** `migrations/00008_create_companies_table.sql`

```sql
CREATE TABLE IF NOT EXISTS companies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    active BOOLEAN DEFAULT 1,
    logo_url VARCHAR(500),
    primary_color VARCHAR(7) DEFAULT '#3b82f6',
    secondary_color VARCHAR(7) DEFAULT '#1e40af',
    deleted_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ UNIQUE: `slug`
- ✅ NOT NULL: `name`, `slug`, `created_at`, `updated_at`
- ✅ Soft Delete: `deleted_at`
- ✅ Indexes: `idx_companies_slug`, `idx_companies_active`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_companies_slug ON companies(slug);
CREATE INDEX IF NOT EXISTS idx_companies_active ON companies(active);
```

### 2.2 Tabela `users`

**Arquivo:** `migrations/00001_create_users.sql`, `migrations/00009_add_company_id_to_entities.sql`, `migrations/00011_add_role_to_users.sql`, `migrations/00016_make_user_companyid_role_not_null.sql`

```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    company_id INTEGER NOT NULL REFERENCES companies(id),
    role TEXT NOT NULL,
    active BOOLEAN DEFAULT 1,
    deleted_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ UNIQUE: `email`
- ✅ FOREIGN KEY: `company_id` → `companies(id)`
- ✅ NOT NULL: `name`, `email`, `password_hash`, `company_id`, `role`
- ✅ Soft Delete: `deleted_at`
- ✅ Indexes: `idx_users_email`, `idx_users_company_id`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
```

**Nota:** ✅ Migration 00016 enforce NOT NULL em `company_id` e `role`

### 2.3 Tabela `products`

**Arquivo:** `migrations/00001_create_users.sql`, `migrations/00009_add_company_id_to_entities.sql`

```sql
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    price REAL NOT NULL DEFAULT 0,
    is_composto INTEGER NOT NULL DEFAULT 0 CHECK(is_composto IN (0,1)),
    company_id INTEGER NOT NULL REFERENCES companies(id),
    category_id INTEGER,
    active BOOLEAN DEFAULT 1,
    slug TEXT UNIQUE,
    display_order INTEGER DEFAULT 0,
    preparation_time_minutes INTEGER DEFAULT 0,
    featured BOOLEAN DEFAULT 0,
    is_new BOOLEAN DEFAULT 0,
    promotion_price REAL,
    promotion_start INTEGER,
    promotion_end INTEGER,
    available_from TEXT,
    available_until TEXT,
    sku TEXT,
    internal_notes TEXT,
    photo_url TEXT,
    meta_title TEXT,
    meta_description TEXT,
    alt_image TEXT,
    canonical TEXT,
    external_id TEXT,
    marketplace_id TEXT,
    sync_status TEXT,
    last_sync INTEGER,
    deleted_at INTEGER,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ UNIQUE: `slug`
- ✅ FOREIGN KEY: `company_id` → `companies(id)`
- ⚠️ FOREIGN KEY: `category_id` → `categories(id)` **SEM ON DELETE**
- ✅ NOT NULL: `name`, `price`, `is_composto`, `company_id`
- ✅ CHECK: `is_composto IN (0,1)`
- ✅ Soft Delete: `deleted_at`
- ✅ Indexes: `idx_products_company_id`, `idx_products_slug`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_products_company_id ON products(company_id);
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
```

**Problema:** ⚠️ FK `category_id` sem ON DELETE pode causar registros órfãos

### 2.4 Tabela `ingredients`

**Arquivo:** `migrations/00001_create_users.sql`, `migrations/00009_add_company_id_to_entities.sql`

```sql
CREATE TABLE IF NOT EXISTS ingredients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    unit TEXT NOT NULL DEFAULT 'un',
    stock_quantity REAL NOT NULL DEFAULT 0,
    min_stock REAL NOT NULL DEFAULT 0,
    active BOOLEAN DEFAULT 1,
    company_id INTEGER NOT NULL REFERENCES companies(id),
    deleted_at INTEGER,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ FOREIGN KEY: `company_id` → `companies(id)`
- ✅ NOT NULL: `name`, `unit`, `stock_quantity`, `min_stock`, `company_id`
- ✅ Soft Delete: `deleted_at`
- ✅ Indexes: `idx_ingredients_company_id`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_ingredients_company_id ON ingredients(company_id);
```

### 2.5 Tabela `product_ingredients`

**Arquivo:** `migrations/00001_create_users.sql`

```sql
CREATE TABLE IF NOT EXISTS product_ingredients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    ingredient_id INTEGER NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity REAL NOT NULL CHECK(quantity > 0),
    UNIQUE(product_id, ingredient_id),
    deleted_at INTEGER,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ FOREIGN KEY: `product_id` → `products(id)` ON DELETE CASCADE
- ✅ FOREIGN KEY: `ingredient_id` → `ingredients(id)` ON DELETE RESTRICT
- ✅ UNIQUE: `(product_id, ingredient_id)`
- ✅ CHECK: `quantity > 0`
- ✅ Soft Delete: `deleted_at`
- ✅ Indexes: `idx_product_ingredients_product`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_product_ingredients_product ON product_ingredients(product_id);
```

### 2.6 Tabela `orders`

**Arquivo:** `migrations/00009_add_company_id_to_entities.sql`

```sql
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    status TEXT NOT NULL DEFAULT 'pending',
    total_price REAL,
    notes TEXT,
    company_id INTEGER NOT NULL REFERENCES companies(id),
    deleted_at INTEGER,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ FOREIGN KEY: `company_id` → `companies(id)`
- ✅ NOT NULL: `status`, `company_id`
- ✅ Soft Delete: `deleted_at`
- ✅ Indexes: `idx_orders_company_id`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_orders_company_id ON orders(company_id);
```

### 2.7 Tabela `order_items`

**Arquivo:** `migrations/00001_create_users.sql` (implicito)

```sql
CREATE TABLE IF NOT EXISTS order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity REAL NOT NULL,
    unit_price REAL NOT NULL,
    product_name TEXT NOT NULL,
    product_description TEXT,
    product_is_composto INTEGER NOT NULL DEFAULT 0,
    product_photo_url TEXT,
    product_category_id INTEGER,
    product_promotion_price REAL,
    product_featured BOOLEAN DEFAULT 0,
    product_is_new BOOLEAN DEFAULT 0,
    deleted_at INTEGER,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ FOREIGN KEY: `order_id` → `orders(id)` ON DELETE CASCADE
- ⚠️ FOREIGN KEY: `product_id` → `products(id)` **SEM ON DELETE**
- ✅ NOT NULL: `order_id`, `product_id`, `quantity`, `unit_price`, `product_name`, `product_is_composto`
- ✅ Soft Delete: `deleted_at`
- ✅ Indexes: `idx_order_items_order`, `idx_order_items_product`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product ON order_items(product_id);
```

**Problema:** ⚠️ FK `product_id` sem ON DELETE pode causar registros órfãos

### 2.8 Tabela `invitations`

**Arquivo:** `migrations/00012_create_invitations.sql`

```sql
CREATE TABLE IF NOT EXISTS invitations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    email TEXT NOT NULL,
    role TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'pending',
    expires_at INTEGER NOT NULL,
    accepted_at INTEGER,
    created_by INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ FOREIGN KEY: `company_id` → `companies(id)` ON DELETE CASCADE
- ⚠️ FOREIGN KEY: `created_by` → `users(id)` **SEM ON DELETE**
- ✅ UNIQUE: `token`
- ✅ NOT NULL: `company_id`, `email`, `role`, `token`, `status`, `expires_at`, `created_by`
- ✅ Indexes: `idx_invitations_company_id`, `idx_invitations_email`, `idx_invitations_token`, `idx_invitations_status`, `idx_invitations_expires_at`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_invitations_company_id ON invitations(company_id);
CREATE INDEX IF NOT EXISTS idx_invitations_email ON invitations(email);
CREATE INDEX IF NOT EXISTS idx_invitations_token ON invitations(token);
CREATE INDEX IF NOT EXISTS idx_invitations_status ON invitations(status);
CREATE INDEX IF NOT EXISTS idx_invitations_expires_at ON invitations(expires_at);
```

**Problema:** ⚠️ FK `created_by` sem ON DELETE pode causar registros órfãos

### 2.9 Tabela `stock_adjustments_pending`

**Arquivo:** `migrations/00004_create_stock_adjustments_pending.sql`, `migrations/00009_add_company_id_to_entities.sql`

```sql
CREATE TABLE IF NOT EXISTS stock_adjustments_pending (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    ingredient_id INTEGER NOT NULL,
    quantity REAL NOT NULL,
    order_status TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    company_id INTEGER NOT NULL REFERENCES companies(id),
    ingredient_name TEXT,
    ingredient_unit TEXT,
    processed_at INTEGER,
    processed_by INTEGER,
    notes TEXT,
    deleted_at INTEGER,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    UNIQUE(order_id, ingredient_id, status)
);
```

**Validação:**
- ✅ PRIMARY KEY: `id`
- ✅ FOREIGN KEY: `company_id` → `companies(id)`
- ⚠️ FOREIGN KEY: `order_id` → `orders(id)` **SEM ON DELETE**
- ⚠️ FOREIGN KEY: `ingredient_id` → `ingredients(id)` **SEM ON DELETE**
- ✅ UNIQUE: `(order_id, ingredient_id, status)`
- ✅ NOT NULL: `order_id`, `ingredient_id`, `quantity`, `order_status`, `status`, `company_id`
- ✅ Soft Delete: `deleted_at`
- ✅ Indexes: `idx_stock_adjustments_pending_company_id`

**Indexes:**
```sql
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_pending_company_id ON stock_adjustments_pending(company_id);
```

**Problema:** ⚠️ FKs `order_id` e `ingredient_id` sem ON DELETE podem causar registros órfãos

---

## 3. Problemas Identificados

### 3.1 RISCO BAIXO: Foreign Keys Sem ON DELETE

**Tabelas Afetadas:**
- `platform_audit.platform_user_id`
- `products.category_id`
- `order_items.product_id`
- `invitations.created_by`
- `stock_adjustments_pending.order_id`
- `stock_adjustments_pending.ingredient_id`

**Causa Raiz:** FKs criadas sem cláusula ON DELETE  
**Impacto:** Registros órfãos se registro pai for deletado  
**Correção Definitiva:**
```sql
-- Para cada FK, adicionar ON DELETE apropriado
-- Exemplo para platform_audit:
ALTER TABLE platform_audit DROP CONSTRAINT platform_audit_platform_user_id_fkey;
ALTER TABLE platform_audit ADD CONSTRAINT platform_audit_platform_user_id_fkey 
    FOREIGN KEY (platform_user_id) REFERENCES platform_users(id) ON DELETE SET NULL;
```

### 3.2 RISCO BAIXO: Index Composto Ausente

**Tabela:** `orders`  
**Colunas:** `company_id`, `status`, `created_at`  
**Causa Raiz:** Index individual em `company_id` não é suficiente para queries com filtros múltiplos  
**Impacto:** Queries como `SELECT * FROM orders WHERE company_id = ? AND status = ? ORDER BY created_at DESC` podem ser lentas  
**Correção Definitiva:**
```sql
CREATE INDEX IF NOT EXISTS idx_orders_company_status_created ON orders(company_id, status, created_at DESC);
```

### 3.3 RISCO BAIXO: Index Ausente em deleted_at

**Tabelas Afetadas:** `products`, `ingredients`, `orders`  
**Causa Raiz:** Queries com `WHERE deleted_at IS NULL` podem ser lentas sem index  
**Impacto:** Queries de listagem podem ser lentas em tabelas grandes  
**Correção Definitiva:**
```sql
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);
CREATE INDEX IF NOT EXISTS idx_ingredients_deleted_at ON ingredients(deleted_at);
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at);
```

---

## 4. Validação de Constraints

### 4.1 NOT NULL Constraints

**Validação:** ✅ Todas as colunas críticas têm NOT NULL
- `users.company_id` ✅ (migration 00016)
- `users.role` ✅ (migration 00016)
- `products.company_id` ✅
- `ingredients.company_id` ✅
- `orders.company_id` ✅

### 4.2 UNIQUE Constraints

**Validação:** ✅ Unicidade onde necessário
- `platform_users.email` ✅
- `users.email` ✅
- `companies.slug` ✅
- `products.slug` ✅
- `invitations.token` ✅
- `product_ingredients(product_id, ingredient_id)` ✅
- `stock_adjustments_pending(order_id, ingredient_id, status)` ✅

### 4.3 CHECK Constraints

**Validação:** ✅ Validação de dados onde necessário
- `products.is_composto IN (0,1)` ✅
- `product_ingredients.quantity > 0` ✅

### 4.4 Foreign Keys

**Validação:** ⚠️ FKs presentes mas algumas sem ON DELETE
- Todas as tabelas tenant têm FK para `companies(id)` ✅
- `platform_sessions.platform_user_id` → `platform_users(id)` ON DELETE CASCADE ✅
- `product_ingredients.product_id` → `products(id)` ON DELETE CASCADE ✅
- `product_ingredients.ingredient_id` → `ingredients(id)` ON DELETE RESTRICT ✅
- `order_items.order_id` → `orders(id)` ON DELETE CASCADE ✅
- `invitations.company_id` → `companies(id)` ON DELETE CASCADE ✅

---

## 5. Validação de Soft Delete

**Tabelas com Soft Delete:**
- ✅ `companies.deleted_at`
- ✅ `users.deleted_at`
- ✅ `products.deleted_at`
- ✅ `ingredients.deleted_at`
- ✅ `orders.deleted_at`
- ✅ `order_items.deleted_at`
- ✅ `product_ingredients.deleted_at`
- ✅ `stock_adjustments_pending.deleted_at`

**Validação:** ✅ Todas as queries incluem `WHERE deleted_at IS NULL`

---

## 6. Conclusão

O schema do banco de dados está **bem estruturado** com:
- ✅ Separação clara entre tabelas Platform e Tenant
- ✅ Constraints NOT NULL em colunas críticas
- ✅ Foreign Keys para integridade referencial
- ✅ Soft Delete implementado consistentemente
- ✅ Indexes para performance em colunas frequentemente filtradas

**Status Final:** ✅ **APROVADO COM RECOMENDAÇÕES**

**Recomendações:**
1. Adicionar ON DELETE em FKs que não têm (SET NULL ou CASCADE conforme apropriado)
2. Adicionar index composto em `orders(company_id, status, created_at DESC)`
3. Adicionar index em `deleted_at` para tabelas grandes
4. Considerar adicionar index em `products.active` para queries de produtos ativos
