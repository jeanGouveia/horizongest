# API_DOCUMENTATION.md

**HorizonGest Platform API Documentation**  
**Sprint 3.7 - Foundation Alignment**  
**Version:** 1.0.0

---

## Overview

This document describes the REST API for the HorizonGest Platform. All endpoints return JSON and follow RESTful conventions.

**Base URL:** `http://localhost:8080/api`

**Authentication:** JWT Bearer Token (except public endpoints)

**Content-Type:** `application/json`

---

## Public Endpoints

### Health Check

**GET** `/api/health`

Check if the API is running.

**Response:**
```json
{
  "status": "ok",
  "service": "horizongest"
}
```

**Status Codes:**
- `200 OK` - API is running

---

### Public Platform Branding

**GET** `/api/public/brand`

Get public platform branding information. This endpoint does not require authentication and is used by the frontend to display branding dynamically.

**Response:**
```json
{
  "platformName": "HorizonGest",
  "platformShortName": "Horizon",
  "website": "https://horizongest.com",
  "logoPath": "/assets/platform/logo.svg",
  "faviconPath": "/assets/platform/favicon.ico",
  "logoLight": "",
  "logoDark": "",
  "icon": "",
  "loginBackground": "",
  "loginIllustration": "",
  "copyright": "© 2024 HorizonGest Inc. All rights reserved.",
  "primaryColor": "#0f172a",
  "secondaryColor": "#6366f1"
}
```

**Status Codes:**
- `200 OK` - Branding retrieved successfully
- `500 Internal Server Error` - Failed to retrieve branding

---

## Platform Authentication Endpoints

### Platform Login

**POST** `/api/platform/auth/login`

Authenticate as a platform admin and receive a JWT token.

**Request Body:**
```json
{
  "email": "admin@platform.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "admin@platform.com",
    "name": "Platform Admin",
    "role": "admin"
  }
}
```

**Status Codes:**
- `200 OK` - Login successful
- `401 Unauthorized` - Invalid credentials
- `400 Bad Request` - Invalid request format

**Errors:**
```json
{
  "error": "credenciais inválidas"
}
```

---

### Platform Logout

**POST** `/api/platform/auth/logout`

Logout platform admin and invalidate JWT token.

**Headers:**
- `Authorization: Bearer {token}`

**Response:**
```json
{
  "message": "logout bem-sucedido"
}
```

**Status Codes:**
- `200 OK` - Logout successful
- `401 Unauthorized` - Invalid or expired token

---

### Platform Me

**GET** `/api/platform/auth/me`

Get current platform admin information.

**Headers:**
- `Authorization: Bearer {token}`

**Response:**
```json
{
  "id": 1,
  "email": "admin@platform.com",
  "name": "Platform Admin",
  "role": "admin",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**Status Codes:**
- `200 OK` - User retrieved successfully
- `401 Unauthorized` - Invalid or expired token

---

## Platform Branding Endpoints

### Get Platform Brand

**GET** `/api/platform/brand`

Get full platform brand configuration (platform admin only).

**Headers:**
- `Authorization: Bearer {token}`

**Response:**
```json
{
  "id": 1,
  "platformName": "HorizonGest",
  "platformShortName": "Horizon",
  "ownerCompanyName": "HorizonGest Inc.",
  "ownerDocument": "12.345.678/0001-90",
  "website": "https://horizongest.com",
  "supportEmail": "support@horizongest.com",
  "supportUrl": "https://support.horizongest.com",
  "logoPath": "/assets/platform/logo.svg",
  "faviconPath": "/assets/platform/favicon.ico",
  "logoLight": "",
  "logoDark": "",
  "icon": "",
  "loginBackground": "",
  "loginIllustration": "",
  "copyright": "© 2024 HorizonGest Inc. All rights reserved.",
  "privacyPolicyUrl": "https://horizongest.com/privacy",
  "termsUrl": "https://horizongest.com/terms",
  "instagramUrl": "https://instagram.com/horizongest",
  "facebookUrl": "https://facebook.com/horizongest",
  "linkedinUrl": "https://linkedin.com/company/horizongest",
  "youtubeUrl": "https://youtube.com/horizongest",
  "defaultLanguage": "pt-BR",
  "defaultTimezone": "America/Sao_Paulo",
  "maintenanceMode": false,
  "maintenanceMessage": "",
  "primaryColor": "#0f172a",
  "secondaryColor": "#6366f1",
  "updatedAt": "2024-01-01T00:00:00Z",
  "updatedBy": 1
}
```

**Status Codes:**
- `200 OK` - Branding retrieved successfully
- `401 Unauthorized` - Invalid or expired token
- `403 Forbidden` - Not authorized (not platform admin)
- `500 Internal Server Error` - Failed to retrieve branding

---

### Update Platform Brand

**PUT** `/api/platform/brand`

Update platform brand configuration (platform admin only).

**Headers:**
- `Authorization: Bearer {token}`

**Request Body:**
```json
{
  "platformName": "HorizonGest",
  "platformShortName": "Horizon",
  "ownerCompanyName": "HorizonGest Inc.",
  "ownerDocument": "12.345.678/0001-90",
  "website": "https://horizongest.com",
  "supportEmail": "support@horizongest.com",
  "supportUrl": "https://support.horizongest.com",
  "logoPath": "/assets/platform/logo.svg",
  "faviconPath": "/assets/platform/favicon.ico",
  "logoLight": "",
  "logoDark": "",
  "icon": "",
  "loginBackground": "",
  "loginIllustration": "",
  "copyright": "© 2024 HorizonGest Inc. All rights reserved.",
  "privacyPolicyUrl": "https://horizongest.com/privacy",
  "termsUrl": "https://horizongest.com/terms",
  "instagramUrl": "https://instagram.com/horizongest",
  "facebookUrl": "https://facebook.com/horizongest",
  "linkedinUrl": "https://linkedin.com/company/horizongest",
  "youtubeUrl": "https://youtube.com/horizongest",
  "defaultLanguage": "pt-BR",
  "defaultTimezone": "America/Sao_Paulo",
  "maintenanceMode": false,
  "maintenanceMessage": "",
  "primaryColor": "#0f172a",
  "secondaryColor": "#6366f1"
}
```

**Response:**
```json
{
  "message": "configuração de marca atualizada com sucesso"
}
```

**Status Codes:**
- `200 OK` - Branding updated successfully
- `400 Bad Request` - Invalid configuration
- `401 Unauthorized` - Invalid or expired token
- `403 Forbidden` - Not authorized (not platform admin)
- `500 Internal Server Error` - Failed to update branding

**Errors:**
```json
{
  "error": "configuração de marca inválida"
}
```

---

## Global Configuration Endpoints

### Get Global Config

**GET** `/api/platform/global-config`

Get global technical configuration (platform admin only).

**Headers:**
- `Authorization: Bearer {token}`

**Response:**
```json
{
  "defaultTimezone": "America/Sao_Paulo",
  "defaultLocale": "pt-BR",
  "monetaryFormat": "BRL R$ 1.000,00",
  "dateFormat": "DD/MM/YYYY",
  "timeFormat": "HH:mm",
  "maxUploadSizeMb": 10,
  "maxImageSizeMb": 5,
  "allowedImageTypes": "jpg,png,webp,gif",
  "allowedFileTypes": "pdf,doc,docx,xlsx,xls,txt",
  "maintenanceMode": false,
  "maintenanceMessage": "",
  "enableFinance": true,
  "enablePurchasing": true,
  "enableInventory": true,
  "enableCRM": false,
  "enableCalendar": false,
  "enablePOS": false,
  "enableAI": false,
  "enableDelivery": false,
  "enableMarketplace": false,
  "updatedAt": "2024-01-01T00:00:00Z",
  "updatedBy": 1
}
```

**Status Codes:**
- `200 OK` - Configuration retrieved successfully
- `401 Unauthorized` - Invalid or expired token
- `403 Forbidden` - Not authorized (not platform admin)
- `500 Internal Server Error` - Failed to retrieve configuration

---

### Update Global Config

**PUT** `/api/platform/global-config`

Update global technical configuration (platform admin only).

**Headers:**
- `Authorization: Bearer {token}`

**Request Body:**
```json
{
  "defaultTimezone": "America/Sao_Paulo",
  "defaultLocale": "pt-BR",
  "monetaryFormat": "BRL R$ 1.000,00",
  "dateFormat": "DD/MM/YYYY",
  "timeFormat": "HH:mm",
  "maxUploadSizeMb": 10,
  "maxImageSizeMb": 5,
  "allowedImageTypes": "jpg,png,webp,gif",
  "allowedFileTypes": "pdf,doc,docx,xlsx,xls,txt",
  "maintenanceMode": false,
  "maintenanceMessage": "",
  "enableFinance": true,
  "enablePurchasing": true,
  "enableInventory": true,
  "enableCRM": false,
  "enableCalendar": false,
  "enablePOS": false,
  "enableAI": false,
  "enableDelivery": false,
  "enableMarketplace": false
}
```

**Response:**
```json
{
  "defaultTimezone": "America/Sao_Paulo",
  "defaultLocale": "pt-BR",
  "monetaryFormat": "BRL R$ 1.000,00",
  "dateFormat": "DD/MM/YYYY",
  "timeFormat": "HH:mm",
  "maxUploadSizeMb": 10,
  "maxImageSizeMb": 5,
  "allowedImageTypes": "jpg,png,webp,gif",
  "allowedFileTypes": "pdf,doc,docx,xlsx,xls,txt",
  "maintenanceMode": false,
  "maintenanceMessage": "",
  "enableFinance": true,
  "enablePurchasing": true,
  "enableInventory": true,
  "enableCRM": false,
  "enableCalendar": false,
  "enablePOS": false,
  "enableAI": false,
  "enableDelivery": false,
  "enableMarketplace": false,
  "updatedAt": "2024-01-01T00:00:00Z",
  "updatedBy": 1
}
```

**Status Codes:**
- `200 OK` - Configuration updated successfully
- `400 Bad Request` - Invalid configuration
- `401 Unauthorized` - Invalid or expired token
- `403 Forbidden` - Not authorized (not platform admin)
- `500 Internal Server Error` - Failed to update configuration

**Errors:**
```json
{
  "error": "configuração global inválida"
}
```

---

### Get Module Status

**GET** `/api/platform/global-config/module-status`

Check if a specific module is enabled via feature flags.

**Headers:**
- `Authorization: Bearer {token}`

**Query Parameters:**
- `module` (required) - Module key (e.g., "finance", "inventory")

**Example:**
```
GET /api/platform/global-config/module-status?module=finance
```

**Response:**
```json
{
  "enabled": true
}
```

**Status Codes:**
- `200 OK` - Module status retrieved successfully
- `400 Bad Request` - Missing module parameter
- `401 Unauthorized` - Invalid or expired token
- `403 Forbidden` - Not authorized (not platform admin)
- `500 Internal Server Error` - Failed to retrieve module status

---

## Tenant Authentication Endpoints

### Tenant Login

**POST** `/api/auth/login`

Authenticate as a tenant user and receive a JWT token.

**Request Body:**
```json
{
  "email": "user@company.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@company.com",
    "name": "User Name",
    "companyId": 1,
    "role": "admin"
  }
}
```

**Status Codes:**
- `200 OK` - Login successful
- `401 Unauthorized` - Invalid credentials
- `400 Bad Request` - Invalid request format

**Errors:**
```json
{
  "error": "credenciais inválidas"
}
```

---

### Tenant Logout

**POST** `/api/auth/logout`

Logout tenant user and invalidate JWT token.

**Headers:**
- `Authorization: Bearer {token}`

**Response:**
```json
{
  "message": "logout bem-sucedido"
}
```

**Status Codes:**
- `200 OK` - Logout successful
- `401 Unauthorized` - Invalid or expired token

---

### Tenant Me

**GET** `/api/auth/me`

Get current tenant user information.

**Headers:**
- `Authorization: Bearer {token}`

**Response:**
```json
{
  "id": 1,
  "email": "user@company.com",
  "name": "User Name",
  "companyId": 1,
  "role": "admin",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**Status Codes:**
- `200 OK` - User retrieved successfully
- `401 Unauthorized` - Invalid or expired token

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "formato dos dados inválido. Verifique o JSON enviado."
}
```

### 401 Unauthorized
```json
{
  "error": "não autorizado"
}
```

### 403 Forbidden
```json
{
  "error": "acesso negado"
}
```

### 404 Not Found
```json
{
  "error": "recurso não encontrado"
}
```

### 500 Internal Server Error
```json
{
  "error": "erro interno do servidor"
}
```

---

## Authentication

### JWT Token

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer {token}
```

### Token Expiration

- Platform tokens: 24 hours
- Tenant tokens: 24 hours

### Token Refresh

Tokens must be refreshed by logging in again after expiration.

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Platform routes:** 5 requests/minute per IP
- **Tenant routes:** 30 requests/hour per user

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 4
X-RateLimit-Reset: 1640995200
```

When rate limited:
```
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1640995200
```

---

## Security Headers

All responses include security headers:

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
Content-Security-Policy: default-src 'self'
```

---

## Multi-Tenancy

Tenant-specific endpoints require a valid tenant JWT token. The token includes the `companyId` which is used to isolate data between tenants.

**Important:** All tenant data is automatically filtered by `companyId` in the repository layer.

---

## CORS

The API supports CORS for frontend integration. Allowed origins must be configured in environment variables.

---

## Versioning

The current API version is v1. All endpoints are prefixed with `/api`.

Future versions will be prefixed with `/api/v2`, `/api/v3`, etc.

---

## Support

For API support, contact: support@platform.com

---

**Last Updated:** Sprint 3.7 - Foundation Alignment  
**Version:** 1.0.0
