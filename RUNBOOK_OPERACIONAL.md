# RUNBOOK OPERACIONAL — PRATOONLINE

**Versão:** 1.0  
**Data:** 17 de Julho de 2026  
**Ambiente:** Produção Piloto

---

## 1. COMO INICIAR O SISTEMA

### Backend

```bash
cd backend

# Carregar variáveis de ambiente
export $(cat .env | xargs)

# Iniciar servidor
go run cmd/server/main.go
```

**Ou usando binário compilado:**

```bash
cd backend
./server
```

**Verificação:**
- Acessar `http://localhost:8080/api/health`
- Deve retornar: `{"status":"ok","service":"pratoOnline"}`

### Frontend

```bash
cd frontend

# Instalar dependências (primeira vez)
npm install

# Iniciar desenvolvimento
npm run dev

# Ou build para produção
npm run build
npm run preview
```

**Verificação:**
- Acessar `http://localhost:5173` (dev) ou `http://localhost:3000` (preview)

---

## 2. COMO PARAR O SISTEMA

### Backend

```bash
# Se estiver rodando com go run
Ctrl + C

# Se estiver rodando como processo em background
pkill -f "server"
# Ou
kill <PID>
```

### Frontend

```bash
# Se estiver rodando com npm run dev
Ctrl + C

# Se estiver rodando com npm run preview
Ctrl + C
```

---

## 3. COMO VERIFICAR SAÚDE

### Health Check Básico

```bash
curl http://localhost:8080/api/health
```

**Resposta esperada:**
```json
{
  "status": "healthy",
  "database": "connected",
  "storage": "available",
  "version": "1.0.0",
  "uptime": "2h30m15s"
}
```

### Health Check Completo

```bash
# Versão
curl http://localhost:8080/api/version

# Capacidades
curl http://localhost:8080/api/capabilities
```

### Diagnóstico de Problemas

**Sistema não responde:**
1. Verificar se processo está rodando: `ps aux | grep server`
2. Verificar logs: `tail -f /var/log/pratonline/backend.log`
3. Verificar porta: `netstat -tlnp | grep 8080`
4. Verificar banco de dados: `ls -la backend/app.db`

**Banco desconectado:**
1. Verificar arquivo do banco: `ls -la backend/app.db`
2. Verificar permissões: `chmod 644 backend/app.db`
3. Verificar espaço em disco: `df -h`

---

## 4. COMO FAZER BACKUP

### Backup do Banco de Dados

```bash
cd backend

# Backup completo
cp app.db app.db.backup.$(date +%Y%m%d_%H%M%S)

# Ou usando sqlite3
sqlite3 app.db ".backup 'app.db.backup.$(date +%Y%m%d_%H%M%S)'"
```

### Backup de Uploads

```bash
cd backend

# Backup do diretório de uploads
tar -czf uploads.backup.$(date +%Y%m%d_%H%M%S).tar.gz uploads/
```

### Backup Completo

```bash
#!/bin/bash
# backup_completo.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/pratonline/$DATE"

mkdir -p $BACKUP_DIR

# Backup banco
cp backend/app.db $BACKUP_DIR/app.db

# Backup uploads
tar -czf $BACKUP_DIR/uploads.tar.gz backend/uploads/

# Backup .env (cuidado com secrets!)
cp backend/.env $BACKUP_DIR/.env

echo "Backup completo: $BACKUP_DIR"
```

---

## 5. COMO RESTAURAR BACKUP

### Restaurar Banco de Dados

```bash
cd backend

# Parar servidor
pkill -f server

# Restaurar backup
cp app.db.backup.YYYYMMDD_HHMMSS app.db

# Reiniciar servidor
./server
```

### Restaurar Uploads

```bash
cd backend

# Parar servidor
pkill -f server

# Restaurar uploads
tar -xzf uploads.backup.YYYYMMDD_HHMMSS.tar.gz

# Reiniciar servidor
./server
```

### Restauração Completa

```bash
#!/bin/bash
# restore_completo.sh BACKUP_DIR

BACKUP_DIR=$1

if [ -z "$BACKUP_DIR" ]; then
  echo "Uso: ./restore_completo.sh /backups/pratonline/YYYYMMDD_HHMMSS"
  exit 1
fi

# Parar servidor
pkill -f server

# Restaurar banco
cp $BACKUP_DIR/app.db backend/app.db

# Restaurar uploads
tar -xzf $BACKUP_DIR/uploads.tar.gz -C backend/

# Restaurar .env (se necessário)
cp $BACKUP_DIR/.env backend/.env

# Reiniciar servidor
cd backend
./server

echo "Restauração completa: $BACKUP_DIR"
```

---

## 6. COMO EXECUTAR MIGRATIONS

### Migrações Automáticas (GORM AutoMigrate)

O sistema usa GORM AutoMigrate por padrão. As migrations são executadas automaticamente ao iniciar o servidor:

```bash
cd backend
go run cmd/server/main.go
```

**Verificação:**
- Verificar logs para: "Migrações executadas com sucesso"
- Verificar tabelas no banco: `sqlite3 app.db ".tables"`

### Migrações Manuais (Goose - Futuro)

Para produção, recomenda-se usar Goose com migrations SQL versionadas:

```bash
# Instalar Goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Executar migrations up
goose -dir migrations sqlite "app.db" up

# Executar migrations down
goose -dir migrations sqlite "app.db" down

# Verificar status
goose -dir migrations sqlite "app.db" status
```

---

## 7. COMO CONSULTAR LOGS

### Logs do Backend

```bash
# Logs em tempo real (se estiver usando systemd)
journalctl -u pratonline-backend -f

# Logs do arquivo (se configurado)
tail -f /var/log/pratonline/backend.log

# Logs das últimas 100 linhas
tail -n 100 /var/log/pratonline/backend.log

# Logs com filtro de erro
grep ERROR /var/log/pratonline/backend.log

# Logs com filtro de pedido
grep "order_id" /var/log/pratonline/backend.log
```

### Logs do Frontend

```bash
# Logs do navegador (desenvolvimento)
# Abrir DevTools > Console

# Logs do servidor Node (produção)
journalctl -u pratonline-frontend -f
```

### Logs Estruturados

O sistema usa logs estruturados com contexto:

```bash
# Filtrar por Request ID
grep "request_id=abc123" /var/log/pratonline/backend.log

# Filtrar por usuário
grep "user_id=42" /var/log/pratonline/backend.log

# Filtrar por endpoint
grep "POST /api/orders" /var/log/pratonline/backend.log
```

---

## 8. COMO ATUALIZAR VERSÃO

### Atualização Backend

```bash
#!/bin/bash
# update_backend.sh

cd /path/to/pratonline/backend

# 1. Fazer backup
./backup_completo.sh

# 2. Parar servidor
pkill -f server

# 3. Pull do código
git pull origin main

# 4. Instalar dependências
go mod download

# 5. Compilar
go build -o server cmd/server/main.go

# 6. Executar migrations (se necessário)
# goose -dir migrations sqlite "app.db" up

# 7. Iniciar servidor
./server

# 8. Verificar saúde
curl http://localhost:8080/api/health
```

### Atualização Frontend

```bash
#!/bin/bash
# update_frontend.sh

cd /path/to/pratonline/frontend

# 1. Fazer backup
cp -r build build.backup.$(date +%Y%m%d_%H%M%S)

# 2. Parar servidor
pkill -f node

# 3. Pull do código
git pull origin main

# 4. Instalar dependências
npm install

# 5. Build
npm run build

# 6. Iniciar servidor
npm run preview
```

---

## 9. COMO FAZER ROLLBACK

### Rollback Backend

```bash
#!/bin/bash
# rollback_backend.sh VERSION

cd /path/to/pratonline/backend

VERSION=$1

# 1. Parar servidor
pkill -f server

# 2. Checkout da versão
git checkout $VERSION

# 3. Compilar
go build -o server cmd/server/main.go

# 4. Iniciar servidor
./server

# 5. Verificar saúde
curl http://localhost:8080/api/health
```

### Rollback Frontend

```bash
#!/bin/bash
# rollback_frontend.sh VERSION

cd /path/to/pratonline/frontend

VERSION=$1

# 1. Parar servidor
pkill -f node

# 2. Checkout da versão
git checkout $VERSION

# 3. Instalar dependências
npm install

# 4. Build
npm run build

# 5. Iniciar servidor
npm run preview
```

### Rollback com Backup de Dados

```bash
#!/bin/bash
# rollback_completo.sh BACKUP_DIR VERSION

BACKUP_DIR=$1
VERSION=$2

# 1. Parar servidores
pkill -f server
pkill -f node

# 2. Restaurar backup de dados
./restore_completo.sh $BACKUP_DIR

# 3. Rollback do código
cd backend
git checkout $VERSION
go build -o server cmd/server/main.go

cd ../frontend
git checkout $VERSION
npm install
npm run build

# 4. Iniciar servidores
cd ../backend
./server

cd ../frontend
npm run preview
```

---

## 10. COMO DIAGNOSTICAR ERROS COMUNS

### Erro: "database is locked"

**Causa:** SQLite com WAL mode e múltiplas conexões

**Solução:**
```bash
# Verificar processos usando o banco
lsof backend/app.db

# Matar processos órfãos
kill -9 <PID>

# Limpar WAL files
rm backend/app.db-wal backend/app.db-shm

# Reiniciar servidor
./server
```

### Erro: "connection refused"

**Causa:** Servidor não está rodando ou porta errada

**Solução:**
```bash
# Verificar se servidor está rodando
ps aux | grep server

# Verificar porta
netstat -tlnp | grep 8080

# Verificar variável PORT
echo $PORT

# Iniciar servidor
./server
```

### Erro: "unauthorized"

**Causa:** Token inválido ou expirado

**Solução:**
```bash
# Fazer login novamente
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Verificar JWT_SECRET
echo $JWT_SECRET

# Gerar novo token se necessário
```

### Erro: "estoque insuficiente"

**Causa:** Ingredientes sem estoque suficiente

**Solução:**
```bash
# Verificar estoque atual
curl http://localhost:8080/api/ingredients

# Atualizar estoque
curl -X PATCH http://localhost:8080/api/ingredients/<id>/stock \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 100}'
```

### Erro: "upload failed"

**Causa:** Permissões ou diretório inexistente

**Solução:**
```bash
# Verificar diretório de uploads
ls -la backend/uploads

# Criar diretório se não existir
mkdir -p backend/uploads

# Verificar permissões
chmod 755 backend/uploads

# Verificar espaço em disco
df -h
```

### Erro: "migration failed"

**Causa:** Schema inconsistente ou migration erro

**Solução:**
```bash
# Verificar tabelas atuais
sqlite3 app.db ".tables"

# Verificar schema de uma tabela
sqlite3 app.db ".schema products"

# Se necessário, restaurar backup
cp app.db.backup.YYYYMMDD_HHMMSS app.db

# Reiniciar servidor
./server
```

---

## MONITORAMENTO

### Métricas Importantes

**Uptime:**
```bash
curl http://localhost:8080/api/health | jq '.uptime'
```

**Versão:**
```bash
curl http://localhost:8080/api/version | jq '.version'
```

**Capacidades:**
```bash
curl http://localhost:8080/api/capabilities
```

### Alertas

Configurar alertas para:
- Servidor não respondendo (health check falhando)
- Erro 500 em endpoints críticos
- Banco de dados locked
- Espaço em disco < 10%
- Uso de CPU > 80%
- Uso de memória > 80%

---

## CONTATO DE SUPORTE

**Desenvolvedor:** Jean Gouveia  
**Email:** [email do desenvolvedor]  
**Documentação:** /home/jean/projetos/pratoOnline/  
**Repositório:** [URL do repositório]

---

## NOTAS

- Este runbook assume ambiente Linux
- Ajustar caminhos conforme seu ambiente
- Manter backups por pelo menos 30 dias
- Testar procedimentos de restore regularmente
- Documentar qualquer alteração neste runbook
