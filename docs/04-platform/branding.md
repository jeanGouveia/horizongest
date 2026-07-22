# Branding Documentation

**HorizonGest Platform - Branding System**

---

## Overview

HorizonGest supports dynamic branding for both platform and tenant levels.

---

## Platform Branding

### Table

`platform_brand_config`

### Fields

- `platform_name` - Platform name
- `platform_short_name` - Abbreviated name
- `owner_company_name` - Owner company
- `owner_document` - Legal document
- `website` - Website URL
- `support_email` - Support email
- `support_url` - Help center URL
- `logo_path` - Logo path
- `favicon_path` - Favicon path
- `logo_light` - Light theme logo
- `logo_dark` - Dark theme logo
- `icon` - Platform icon
- `login_background` - Login background
- `login_illustration` - Login illustration
- `copyright` - Copyright notice
- `privacy_policy_url` - Privacy policy URL
- `terms_url` - Terms URL
- `instagram_url` - Instagram URL
- `facebook_url` - Facebook URL
- `linkedin_url` - LinkedIn URL
- `youtube_url` - YouTube URL
- `default_language` - Default language
- `default_timezone` - Default timezone
- `maintenance_mode` - Maintenance mode
- `maintenance_message` - Maintenance message
- `primary_color` - Primary color
- `secondary_color` - Secondary color

### Public Endpoint

`GET /api/public/brand`

No authentication required. Returns public-safe branding information.

---

## Tenant Branding

### Tables

- `themes` - Theme configuration
- `business_profiles` - Business profile

### Fields

- Logo
- Colors
- Theme customization

---

## Dynamic Branding

### Backend

100% dynamic - no hardcoded branding in backend code.

### Frontend

95% dynamic - uses brandStore to consume `/api/public/brand`.

### Services

Services use dynamic platform name (EmailService, BackupService, AuthService).

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
