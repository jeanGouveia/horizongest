# SPRINT 4.1 - Test Evidences

**Data:** 2026-07-20  
**Auditor:** Cascade AI

---

## Evidências Coletadas

### 1. Backend Startup
**Comando:** `./server`
**Resultado:** ✅ Sucesso
**Saída:**
```
2026/07/20 17:16:24 ✅ PratoOnline backend iniciado em http://localhost:8080
```
**Status:** OK

---

### 2. Frontend Startup
**Comando:** `npm run dev`
**Resultado:** ✅ Sucesso
**Saída:**
```
VITE v5.4.21  ready in 1177 ms
➜  Local:   http://localhost:3000/
```
**Status:** OK

---

### 3. Health Check
**Comando:** `curl -s http://localhost:8080/api/health`
**Resultado:** ✅ Sucesso
**Resposta:**
```json
{"status":"ok","service":"pratoOnline"}
```
**Status:** OK

---

### 4. Database Tables
**Comando:** `sqlite3 app.db ".tables"`
**Resultado:** ✅ Sucesso
**Tabelas:**
```
categories
companies
gorm_token_blacklists
ingredients
invitations
media
order_items
orders
password_reset_tokens
product_ingredients
products
stock_adjustments_pending
users
```
**Status:** OK

---

### 5. Login API Test (jwtfinal@test.com)
**Comando:** `curl -s -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d '{"email":"jwtfinal@test.com","password":"admin123"}'`
**Resultado:** ❌ Falha
**Resposta:**
```json
{"error":"e-mail ou senha incorretos. Verifique suas credenciais."}
```
**Status:** FALHOU

---

### 6. Login API Test (owner@test.com)
**Comando:** `curl -s -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d '{"email":"owner@test.com","password":"owner123"}'`
**Resultado:** ❌ Falha
**Resposta:**
```json
{"error":"e-mail ou senha incorretos. Verifique suas credenciais."}
```
**Status:** FALHOU

---

### 7. Login API Test (qa@test.com)
**Comando:** `curl -s -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d '{"email":"qa@test.com","password":"password"}'`
**Resultado:** ❌ Falha
**Resposta:**
```json
{"error":"e-mail ou senha incorretos. Verifique suas credenciais."}
```
**Status:** FALHOU

---

### 8. Forgot Password API Test
**Comando:** `curl -s -X POST http://localhost:8080/api/auth/forgot-password -H "Content-Type: application/json" -d '{"email":"qa@test.com"}'`
**Resultado:** ❌ Falha
**Resposta:**
```
404 page not found
```
**Status:** FALHOU

---

### 9. Register API Test
**Comando:** `curl -s -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{"email":"qa@pratoonline.com","password":"qa123456","name":"QA Tester"}'`
**Resultado:** ❌ Falha
**Resposta:**
```
404 page not found
```
**Status:** FALHOU

---

## Observações

- Endpoint `/api/auth/register` não existe (removido no Sprint 3)
- Endpoint `/api/auth/forgot-password` não existe ou não está acessível
- Login via API não funcionou com usuários existentes no banco
- Senhas no banco estão hashadas com bcrypt
- Necessário testar login via frontend

---

## Assinatura

**Auditor:** Cascade AI  
**Data:** 2026-07-20
