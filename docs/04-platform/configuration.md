# Configuration Documentation

**HorizonGest Platform - Configuration System**

---

## Overview

HorizonGest has three levels of configuration: PlatformBrandConfig, GlobalConfig, and Environment Variables.

---

## PlatformBrandConfig

**Purpose:** Branding/institutional configuration

**Access:** Public endpoint for frontend

**Fields:** See [branding.md](branding.md)

---

## GlobalConfig

**Purpose:** Technical configuration

**Access:** Platform admin only

### Fields

**Localization:**
- `default_timezone` - Platform timezone
- `default_locale` - Default language/locale
- `monetary_format` - Currency format
- `date_format` - Date format
- `time_format` - Time format

**File Uploads:**
- `max_upload_size_mb` - Max file upload size
- `max_image_size_mb` - Max image upload size
- `allowed_image_types` - Allowed image types
- `allowed_file_types` - Allowed file types

**Maintenance:**
- `maintenance_mode` - Maintenance mode
- `maintenance_message` - Maintenance message

**Feature Flags:**
- `enable_finance` - Enable finance module
- `enable_purchasing` - Enable purchasing module
- `enable_inventory` - Enable inventory module
- `enable_crm` - Enable CRM module
- `enable_calendar` - Enable calendar module
- `enable_pos` - Enable POS module
- `enable_ai` - Enable AI features
- `enable_delivery` - Enable delivery module
- `enable_marketplace` - Enable marketplace module

---

## Environment Variables

**Purpose:** Secrets and infrastructure

**Access:** Application startup only

### Required

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=pratodb
JWT_PLATFORM_SECRET=
JWT_TENANT_SECRET=
APP_VERSION=1.0.0
PORT=8080
```

### Optional

```bash
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
EMAIL_ENABLED=false
```

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
