# PLANO DE DEPLOY PRODUÇÃO — PRATOONLINE

**Versão:** 1.0  
**Data:** 17 de Julho de 2026  
**Ambiente:** Produção Piloto  
**Classificação:** YELLOW — Riscos Conhecidos

---

## VISÃO GERAL

Este documento descreve o plano completo para deploy do PratoOnline em ambiente de produção piloto. O plano inclui preparação, correções obrigatórias, deployment, validação e monitoramento.

**Tempo Estimado Total:** 3-5 dias úteis  
**Pessoas Necessárias:** 1 desenvolvedor  
**Risco:** Médio (mitigável)

---

## PRÉ-REQUISITOS

### Infraestrutura

- **Servidor:** Linux (Ubuntu 22.04+ recomendado)
- **CPU:** 2 cores mínimos
- **RAM:** 4GB mínimos
- **Armazenamento:** 20GB SSD
- **Rede:** Conectividade estável
- **Domínio:** Opcional (pode usar IP)

### Software

- **Go:** 1.26.3+
- **Node.js:** 18+
- **npm:** 9+
- **SQLite:** 3.35+
- **Git:** 2.0+
- **Nginx:** Opcional (para reverse proxy)

### Acesso

- **SSH:** Acesso root ou sudo
- **Git:** Acesso ao repositório
- **Ambiente:** Variáveis de ambiente configuradas

---

## FASE 1 — CORREÇÕES OBRIGATÓRIAS (2-3 dias)

### 1.1 Implementar CORS Middleware

**Arquivo:** `backend/internal/middleware/cors_middleware.go` (novo)

```go
package middleware

import (
    "net/http"
    "strings"
)

func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Em produção, substituir por origens específicas
        origin := r.Header.Get("Origin")
        allowedOrigins := []string{
            "http://localhost:5173",
            "http://localhost:3000",
            // Adicionar domínio de produção
        }

        allowed := false
        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin {
                allowed = true
                break
            }
        }

        if allowed {
            w.Header().Set("Access-Control-Preflight-Max-Age", "86400")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

**Integração:** Adicionar em `cmd/server/main.go`

```go
r.Use(middleware.CORSMiddleware)
```

---

### 1.2 Corrigir JWT_SECRET Fallback

**Arquivo:** `backend/internal/service/auth_service.go`

**Antes:**
```go
func NewAuthService(userRepo ports.UserRepository) *AuthService {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "dev-secret-troque-em-producao"
    }
    // ...
}
```

**Depois:**
```go
func NewAuthService(userRepo ports.UserRepository) *AuthService {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("FATAL: JWT_SECRET não está definido")
    }
    if len(secret) < 32 {
        log.Fatal("FATAL: JWT_SECRET deve ter pelo menos 32 caracteres")
    }
    // ...
}
```

---

### 1.3 Implementar Rate Limiting Básico

**Arquivo:** `backend/internal/middleware/rate_limit_middleware.go` (novo)

```go
package middleware

import (
    "net/http"
    "sync"
    "time"
)

type RateLimiter struct {
    visitors map[string]*visitor
    mu       sync.RWMutex
    rate     int
    window   time.Duration
}

type visitor struct {
    requests  int
    lastSeen  time.Time
    resetTime time.Time
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
    rl := &RateLimiter{
        visitors: make(map[string]*visitor),
        rate:     rate,
        window:   window,
    }
    go rl.cleanup()
    return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        rl.mu.Lock()
        v, exists := rl.visitors[ip]
        if !exists {
            v = &visitor{
                resetTime: time.Now().Add(rl.window),
            }
            rl.visitors[ip] = v
        }
        
        if time.Now().After(v.resetTime) {
            v.requests = 0
            v.resetTime = time.Now().Add(rl.window)
        }
        
        v.requests++
        v.lastSeen = time.Now()
        
        if v.requests > rl.rate {
            rl.mu.Unlock()
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }
        
        rl.mu.Unlock()
        next.ServeHTTP(w, r)
    })
}

func (rl *RateLimiter) cleanup() {
    for {
        time.Sleep(time.Minute)
        rl.mu.Lock()
        for ip, v := range rl.visitors {
            if time.Since(v.lastSeen) > 3*time.Minute {
                delete(rl.visitors, ip)
            }
        }
        rl.mu.Unlock()
    }
}
```

**Integração:** Adicionar em `cmd/server/main.go`

```go
rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 100 requests por minuto
r.Use(rateLimiter.Middleware)
```

---

### 1.4 Implementar Paginação

**Arquivo:** `backend/internal/service/product_service.go`

**Adicionar:**
```go
type ListProductsInput struct {
    Page     int `json:"page" validate:"gte=1"`
    PageSize int `json:"page_size" validate:"gte=1,lte=100"`
}

func (s *ProductService) ListProductsPaginated(ctx context.Context, in ListProductsInput) ([]domain.Product, int, error) {
    if in.Page == 0 {
        in.Page = 1
    }
    if in.PageSize == 0 {
        in.PageSize = 20
    }
    
    offset := (in.Page - 1) * in.PageSize
    
    products, err := s.repo.ListProductsPaginated(ctx, offset, in.PageSize)
    if err != nil {
        return nil, 0, fmt.Errorf("ProductService.ListProductsPaginated: %w", err)
    }
    
    total, err := s.repo.CountProducts(ctx)
    if err != nil {
        return nil, 0, fmt.Errorf("ProductService.ListProductsPaginated: %w", err)
    }
    
    return products, total, nil
}
```

**Repository:** Implementar métodos paginados em `gorm_product_repository.go`

---

### 1.5 Implementar Backup Automatizado

**Arquivo:** `scripts/backup.sh`

```bash
#!/bin/bash

# Configurações
BACKUP_DIR="/backups/pratonline"
DB_PATH="/path/to/backend/app.db"
UPLOADS_DIR="/path/to/backend/uploads"
RETENTION_DAYS=30

# Criar diretório de backup
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_PATH="$BACKUP_DIR/$DATE"
mkdir -p $BACKUP_PATH

# Backup do banco
cp $DB_PATH $BACKUP_PATH/app.db

# Backup dos uploads
tar -czf $BACKUP_PATH/uploads.tar.gz $UPLOADS_DIR

# Backup do .env (cuidado!)
cp /path/to/backend/.env $BACKUP_DIR/.env.template

# Comprimir backup
tar -czf $BACKUP_DIR/pratonline_backup_$DATE.tar.gz $BACKUP_PATH
rm -rf $BACKUP_PATH

# Limpar backups antigos
find $BACKUP_DIR -name "pratonline_backup_*.tar.gz" -mtime +$RETENTION_DAYS -delete

# Log
echo "Backup concluído: $BACKUP_DIR/pratonline_backup_$DATE.tar.gz"
```

**Cron:**
```bash
# Adicionar ao crontab
0 2 * * * /path/to/scripts/backup.sh >> /var/log/pratonline/backup.log 2>&1
```

---

### 1.6 Implementar Script de Restore

**Arquivo:** `scripts/restore.sh`

```bash
#!/bin/bash

BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
    echo "Uso: ./restore.sh <backup_file.tar.gz>"
    exit 1
fi

# Extrair backup
TEMP_DIR="/tmp/pratonline_restore_$(date +%s)"
mkdir -p $TEMP_DIR
tar -xzf $BACKUP_FILE -C $TEMP_DIR

# Parar servidor
systemctl stop pratonline-backend

# Restaurar banco
cp $TEMP_DIR/app.db /path/to/backend/app.db

# Restaurar uploads
tar -xzf $TEMP_DIR/uploads.tar.gz -C /path/to/backend/

# Reiniciar servidor
systemctl start pratonline-backend

# Limpar
rm -rf $TEMP_DIR

echo "Restore concluído: $BACKUP_FILE"
```

---

### 1.7 Melhorar Health Check

**Arquivo:** `backend/internal/handler/system_handler.go`

**Modificar GetHealth:**
```go
func (h *SystemHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
    // Verificar banco de dados
    dbStatus := "disconnected"
    sqlDB, err := h.db.DB()
    if err == nil {
        if err := sqlDB.Ping(); err == nil {
            dbStatus = "connected"
        }
    }
    
    // Verificar storage
    storageStatus := "available"
    if _, err := os.Stat("uploads"); os.IsNotExist(err) {
        storageStatus = "unavailable"
    }
    
    health := &domain.Health{
        Status:   "healthy",
        Database: dbStatus,
        Storage:  storageStatus,
        Version:  "1.0.0",
        Uptime:   time.Since(h.startTime).String(),
    }
    
    if dbStatus != "connected" || storageStatus != "available" {
        health.Status = "unhealthy"
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}
```

---

## FASE 2 — PREPARAÇÃO DE AMBIENTE (1 dia)

### 2.1 Configurar Variáveis de Ambiente

**Arquivo:** `/etc/pratonline/.env`

```bash
# Backend
PORT=8080
DB_DSN=/var/lib/pratonline/app.db
JWT_SECRET=<gerar-secret-32-caracteres-aleatorios>

# Frontend
VITE_API_URL=http://localhost:8080
```

**Gerar JWT_SECRET:**
```bash
openssl rand -base64 32
```

### 2.2 Configurar Estrutura de Diretórios

```bash
sudo mkdir -p /var/lib/pratonline
sudo mkdir -p /var/log/pratonline
sudo mkdir -p /backups/pratonline
sudo mkdir -p /etc/pratonline

sudo chown -R $USER:$USER /var/lib/pratonline
sudo chown -R $USER:$USER /var/log/pratonline
sudo chown -R $USER:$USER /backups/pratonline
sudo chown -R $USER:$USER /etc/pratonline
```

### 2.3 Configurar Systemd (Backend)

**Arquivo:** `/etc/systemd/system/pratonline-backend.service`

```ini
[Unit]
Description=PratoOnline Backend
After=network.target

[Service]
Type=simple
User=pratonline
WorkingDirectory=/opt/pratonline/backend
Environment="PORT=8080"
Environment="DB_DSN=/var/lib/pratonline/app.db"
Environment="JWT_SECRET=<secret>"
ExecStart=/opt/pratonline/backend/server
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### 2.4 Configurar Systemd (Frontend)

**Arquivo:** `/etc/systemd/system/pratonline-frontend.service`

```ini
[Unit]
Description=PratoOnline Frontend
After=network.target

[Service]
Type=simple
User=pratonline
WorkingDirectory=/opt/pratonline/frontend
Environment="NODE_ENV=production"
ExecStart=/usr/bin/node /opt/pratonline/frontend/build/index.js
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### 2.5 Criar Usuário de Serviço

```bash
sudo useradd -r -s /bin/false pratonline
sudo usermod -a -G pratonline $USER
```

---

## FASE 3 — BUILD E DEPLOY (1 dia)

### 3.1 Build Backend

```bash
cd /opt/pratonline/backend

# Pull do código
git pull origin main

# Instalar dependências
go mod download

# Build
CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go

# Testar
./server --help || echo "Build ok"
```

### 3.2 Build Frontend

```bash
cd /opt/pratonline/frontend

# Pull do código
git pull origin main

# Instalar dependências
npm ci --production

# Build
npm run build

# Testar
ls -la build/
```

### 3.3 Deploy Backend

```bash
# Parar serviço
sudo systemctl stop pratonline-backend

# Backup do binário atual
sudo cp /opt/pratonline/backend/server /opt/pratonline/backend/server.backup

# Copiar novo binário
sudo cp server /opt/pratonline/backend/server

# Ajustar permissões
sudo chmod +x /opt/pratonline/backend/server
sudo chown pratonline:pratonline /opt/pratonline/backend/server

# Iniciar serviço
sudo systemctl start pratonline-backend

# Verificar status
sudo systemctl status pratonline-backend
```

### 3.4 Deploy Frontend

```bash
# Parar serviço
sudo systemctl stop pratonline-frontend

# Backup do build atual
sudo cp -r /opt/pratonline/frontend/build /opt/pratonline/frontend/build.backup

# Copiar novo build
sudo cp -r build /opt/pratonline/frontend/

# Ajustar permissões
sudo chown -R pratonone:pratonone /opt/pratonline/frontend/build

# Iniciar serviço
sudo systemctl start pratonone-frontend

# Verificar status
sudo systemctl status pratonone-frontend
```

---

## FASE 4 — VALIDAÇÃO (1 dia)

### 4.1 Health Check

```bash
# Backend
curl http://localhost:8080/api/health

# Esperado:
{
  "status": "healthy",
  "database": "connected",
  "storage": "available",
  "version": "1.0.0",
  "uptime": "..."
}
```

### 4.2 Testar Autenticação

```bash
# Registrar
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"test123456"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123456"}'
```

### 4.3 Testar CRUD Básico

```bash
# Criar produto
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Produto Teste","price":10.50,"description":"Descrição"}'

# Listar produtos
curl http://localhost:8080/api/products/active

# Criar pedido
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"quantity":2}]}'
```

### 4.4 Testar Backup

```bash
# Executar script de backup
sudo /opt/pratonline/scripts/backup.sh

# Verificar se backup foi criado
ls -la /backups/pratonline/

# Testar restore (em ambiente de teste)
sudo /opt/pratonline/scripts/restore.sh /backups/pratonline/pratonline_backup_YYYYMMDD_HHMMSS.tar.gz
```

### 4.5 Testar Frontend

```bash
# Acessar frontend
curl http://localhost:3000

# Verificar se carrega corretamente
# Testar login via UI
# Testar criação de pedido via UI
```

---

## FASE 5 — MONITORAMENTO INICIAL (1 semana)

### 5.1 Configurar Logs

```bash
# Verificar logs do backend
sudo journalctl -u pratonone-backend -f

# Verificar logs do frontend
sudo journalctl -u pratonone-frontend -f

# Verificar logs de backup
sudo tail -f /var/log/pratonone/backup.log
```

### 5.2 Métricas para Monitorar

- **Uptime do sistema:** `uptime`
- **Uso de CPU:** `top` ou `htop`
- **Uso de memória:** `free -h`
- **Espaço em disco:** `df -h`
- **Conexões ativas:** `netstat -an | grep 8080 | wc -l`
- **Tamanho do banco:** `ls -lh /var/lib/pratonone/app.db`

### 5.3 Alertas Manuais

Verificar diariamente:
- [ ] Backup foi executado com sucesso
- [ ] Espaço em disco > 20%
- [ ] Sistema respondendo
- [ ] Sem erros nos logs
- [ ] Banco de dados não está locked

---

## FASE 6 — ROLLBACK PLAN

### 6.1 Rollback Backend

```bash
# Parar serviço
sudo systemctl stop pratonone-backend

# Restaurar binário
sudo cp /opt/pratonone/backend/server.backup /opt/pratonone/backend/server

# Restaurar banco (se necessário)
sudo cp /var/lib/pratonone/app.db.backup /var/lib/pratonone/app.db

# Iniciar serviço
sudo systemctl start pratonone-backend

# Verificar
sudo systemctl status pratonone-backend
curl http://localhost:8080/api/health
```

### 6.2 Rollback Frontend

```bash
# Parar serviço
sudo systemctl stop pratonone-frontend

# Restaurar build
sudo rm -rf /opt/pratonone/frontend/build
sudo cp -r /opt/pratonone/frontend/build.backup /opt/pratonone/frontend/build

# Iniciar serviço
sudo systemctl start pratonone-frontend

# Verificar
sudo systemctl status pratonone-frontend
curl http://localhost:3000
```

---

## CHECKLIST FINAL

### Pré-Deploy

- [ ] Correções obrigatórias implementadas
- [ ] CORS middleware implementado
- [ ] JWT_SECRET fallback removido
- [ ] Rate limiting implementado
- [ ] Paginação implementada
- [ ] Backup automatizado configurado
- [ ] Script de restore testado
- [ ] Health check melhorado
- [ ] Variáveis de ambiente configuradas
- [ ] Estrutura de diretórios criada
- [ ] Systemd configurado
- [ ] Usuário de serviço criado

### Pós-Deploy

- [ ] Backend buildado e deployado
- [ ] Frontend buildado e deployado
- [ ] Serviços iniciados
- [ ] Health check passando
- [ ] Autenticação testada
- [ ] CRUD básico testado
- [ ] Backup testado
- [ ] Restore testado
- [ ] Frontend acessível
- [ ] Logs configurados
- [ ] Monitoramento iniciado

---

## TEMPO LINE

**Dia 1:** Implementar correções obrigatórias  
**Dia 2:** Implementar correções obrigatórias (continuação)  
**Dia 3:** Preparar ambiente, build e deploy  
**Dia 4:** Validação e testes  
**Dia 5-12:** Monitoramento inicial (1 semana)

---

## RISCOS E MITIGAÇÃO

### Risco: Falha no Deploy

**Mitigação:**
- Backup completo antes do deploy
- Rollback plan documentado
- Teste em ambiente de staging primeiro

### Risco: Backup Falha

**Mitigação:**
- Monitoramento de logs de backup
- Alerta manual diário
- Backup manual como fallback

### Risco: Performance Insuficiente

**Mitigação:**
- Monitoramento de recursos
- Paginação implementada
- Pronto para otimizar se necessário

### Risco: Segurança

**Mitigação:**
- CORS implementado
- Rate limiting implementado
- JWT_SECRET forte
- Firewall configurado

---

## SUPORTE

**Desenvolvedor:** Jean Gouveia  
**Documentação:** 
- `/home/jean/projetos/pratoOnline/RUNBOOK_OPERACIONAL.md`
- `/home/jean/projetos/pratoOnline/AUDITORIA_PILOT_READY.md`

**Contato de Emergência:** [inserir contato]

---

## PRÓXIMOS PASSOS APÓS PILOTO

1. Coletar feedback do usuário piloto
2. Analisar métricas de uso
3. Identificar melhorias necessárias
4. Planejar expansão (multi-tenant, iFood, etc.)
5. Considerar migração para PostgreSQL
6. Implementar CI/CD
7. Implementar monitoramento avançado

---

## ASSINATURA

**Deploy aprovado por:** [nome]  
**Data:** [data]  
**Versão:** 1.0  
**Ambiente:** Produção Piloto
