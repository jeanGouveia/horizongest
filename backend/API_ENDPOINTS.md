# API Endpoints - Backend Hardening (Sprint 11.3)

## Novos Endpoints Criados

### 1. Verificação de Dependências

#### GET /api/products/:id/can-delete
Verifica se um produto pode ser excluído, retornando dependências.

**Response:**
```json
{
  "canDelete": false,
  "reasons": [
    {
      "type": "order",
      "id": 123,
      "name": "Pedido #123",
      "description": "Status: pending, Data: 16/07/2026"
    },
    {
      "type": "recipe",
      "id": 456,
      "name": "Fichas Técnicas",
      "description": "Usado em 2 fichas técnicas de produtos compostos"
    }
  ]
}
```

#### GET /api/ingredients/:id/can-delete
Verifica se um ingrediente pode ser excluído, retornando dependências.

**Response:**
```json
{
  "canDelete": false,
  "reasons": [
    {
      "type": "product",
      "id": 789,
      "name": "Hambúrguer Artesanal",
      "description": "Usado na ficha técnica deste produto composto"
    }
  ]
}
```

#### GET /api/categories/:id/can-delete
Verifica se uma categoria pode ser excluída, retornando dependências.

**Response:**
```json
{
  "canDelete": false,
  "reasons": [
    {
      "type": "product",
      "id": 101,
      "name": "X-Bacon",
      "description": "Produto vinculado a esta categoria"
    }
  ]
}
```

---

### 2. Validação de Estoque

#### POST /api/orders/validate
Valida se há estoque suficiente para os itens do pedido.

**Request:**
```json
{
  "items": [
    {
      "productId": 1,
      "quantity": 2
    }
  ]
}
```

**Response:**
```json
{
  "valid": false,
  "insufficientStock": [
    {
      "ingredientId": 5,
      "ingredientName": "Carne Moída",
      "required": 1.5,
      "available": 0.8,
      "shortage": 0.7,
      "unit": "kg"
    }
  ]
}
```

---

### 3. Dashboard

#### GET /api/dashboard
Retorna dados executivos do dashboard.

**Response:**
```json
{
  "metrics": {
    "todayRevenue": 1250.50,
    "todayOrders": 15,
    "pendingOrders": 3,
    "lowStockCount": 5,
    "activeProducts": 42
  },
  "recentOrders": [
    {
      "id": 123,
      "status": "pending",
      "totalPrice": 85.50,
      "createdAt": "2026-07-16 18:30",
      "itemsCount": 3
    }
  ],
  "lowStock": [
    {
      "id": 10,
      "name": "Queijo Mussarela",
      "stockQuantity": 0.5,
      "minStock": 2.0,
      "unit": "kg"
    }
  ],
  "totalProducts": 50,
  "totalCategories": 8,
  "totalIngredients": 25
}
```

---

### 4. Notifications

#### GET /api/notifications
Retorna contagem de notificações do sistema.

**Response:**
```json
{
  "pendingOrders": 3,
  "lowStockCount": 5,
  "productsWithoutPhoto": 12,
  "expiredPromotions": 2
}
```

---

### 5. System Endpoints

#### GET /api/health
Retorna status de saúde do sistema.

**Response:**
```json
{
  "status": "healthy",
  "database": "connected",
  "storage": "available",
  "version": "1.0.0",
  "uptime": "2h30m15s"
}
```

#### GET /api/version
Retorna informações de versão.

**Response:**
```json
{
  "version": "1.0.0",
  "commit": "dev",
  "build": "20260716-183045",
  "environment": "development"
}
```

#### GET /api/capabilities
Retorna capacidades do sistema.

**Response:**
```json
{
  "upload": true,
  "seo": true,
  "marketplace": false,
  "ifood": false,
  "pix": false,
  "fiscal": false,
  "delivery": false,
  "cardapioDigital": true
}
```

---

### 6. Resposta de Erro Padronizada

Todas as respostas de erro seguem este formato:

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Dados inválidos",
  "details": {
    "field": "email",
    "error": "formato inválido"
  },
  "timestamp": "2026-07-16T21:30:00Z",
  "requestId": "a1b2c3d4e5f6"
}
```

**Headers:**
- `X-Request-ID`: Identificador único da requisição

---

## Arquitetura Preservada

- Clean Architecture mantida
- Repository Pattern preservado
- Domain Layer separado
- Ports/Interfaces definidos
- Nenhuma alteração em arquitetura existente
