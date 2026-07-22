# Environment Variables Reference

**HorizonGest Platform - Environment Variables**

---

## Required Variables

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

---

## Optional Variables

```bash
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
EMAIL_ENABLED=false
```

---

## Security Notes

- Never commit secrets to version control
- Use different secrets for development and production
- Rotate secrets regularly
- Use strong, random secrets for JWT

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
