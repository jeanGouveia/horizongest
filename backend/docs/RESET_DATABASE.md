# Database Reset Documentation

## Overview

The HorizonGest database reset mechanism provides a safe way to clear all business data while preserving the complete database schema, including tables, indexes, constraints, triggers, sequences, extensions, and permissions.

## When to Use

Use the database reset when:

- **Development:** Clear test data between development cycles
- **Testing:** Reset to a clean state before running integration tests
- **Data Seeding:** Prepare a fresh database for seeding with new test data
- **Troubleshooting:** Eliminate data-related issues during debugging

**⚠️ DO NOT use in production** - This will permanently delete all business data.

## What is Deleted

The reset script truncates the following tables (all data removed, sequences restarted):

### Domain Tables (Business Data)
- `order_items` - Order line items
- `orders` - Customer orders
- `product_ingredients` - Product-ingredient relationships
- `products` - Product catalog
- `ingredients` - Ingredient inventory
- `categories` - Product categories
- `stock_adjustments_pending` - Pending stock adjustments
- `media` - Media files (images, documents)
- `invitations` - Company user invitations
- `companies` - Tenant companies
- `users` - Company users

### Auth/Session Tables (Temporary Data)
- `gorm_token_blacklists` - JWT token blacklist
- `password_reset_tokens` - Password reset tokens
- `platform_sessions` - Platform user sessions

### Audit Tables (Optional - Can Be Commented Out)
- `impersonation_audit` - User impersonation audit log
- `platform_audit` - Platform audit log

### Event Sourcing Tables (Optional - Can Be Commented Out)
- `outbox_events` - Outbox event log

## What is Preserved

The following tables are **NOT truncated** to maintain system bootstrap data:

### Platform Configuration
- `platform_brand_config` - Platform branding configuration (colors, logos, themes)

### Platform Users
- `platform_users` - Platform administrator users (required for platform bootstrap)

### Schema Structure
- All table definitions (CREATE TABLE)
- All indexes
- All foreign key constraints
- All unique constraints
- All check constraints
- All triggers
- All sequences (reset to start at 1)
- All extensions
- All permissions and roles

## Available Commands

### Using Makefile (Recommended)

```bash
# From backend directory
make db-reset
```

### Using psql Directly

```bash
PGPASSWORD=horizongest_secure_password psql -U horizongest_user -h localhost -d horizongest -f scripts/reset_database.sql
```

### Using Docker Compose

If PostgreSQL is running in Docker:

```bash
docker-compose exec -T postgres psql -U horizongest_user -d horizongest < scripts/reset_database.sql
```

## Examples

### Example 1: Development Reset

```bash
# Clear all test data after development session
cd backend
make db-reset
```

### Example 2: Pre-Test Reset

```bash
# Reset database before running integration tests
cd backend
make db-reset
go test ./...
```

### Example 3: Custom Environment

```bash
# Reset with custom database credentials
PGPASSWORD=your_password psql -U your_user -h your_host -d your_db -f scripts/reset_database.sql
```

## Technical Details

### How It Works

1. **Transaction Wrapper:** All operations run in a single transaction (BEGIN/COMMIT)
2. **Trigger Optimization:** Temporarily disables triggers for faster execution
3. **CASCADE Truncation:** Uses `TRUNCATE ... RESTART IDENTITY CASCADE` to:
   - Delete all data from tables
   - Restart auto-increment sequences to 1
   - Automatically handle foreign key dependencies
4. **Order Independence:** CASCADE ensures correct truncation order regardless of table order in script

### Performance

- **Typical execution time:** < 1 second for empty database
- **Large databases:** May take several seconds depending on data volume
- **No schema changes:** Only data is modified, structure remains intact

## Validation After Reset

After running the reset, verify:

```sql
-- Check all tables are empty
SELECT schemaname, tablename, n_tup_ins, n_tup_upd, n_tup_del 
FROM pg_stat_user_tables 
WHERE schemaname = 'public';

-- Check sequences are reset
SELECT sequence_name, last_value 
FROM information_schema.sequences 
WHERE sequence_schema = 'public';

-- Check constraints are preserved
SELECT conname, contype 
FROM pg_constraint 
WHERE connamespace = 'public'::regnamespace;

-- Check indexes are preserved
SELECT indexname, tablename 
FROM pg_indexes 
WHERE schemaname = 'public';
```

## Troubleshooting

### Error: "relation does not exist"

**Cause:** Database schema not yet created (first-time setup)

**Solution:** Run database setup first:
```bash
make db-setup
```

### Error: "permission denied"

**Cause:** Insufficient database permissions

**Solution:** Ensure user has TRUNCATE privileges on all tables:
```sql
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO horizongest_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO horizongest_user;
```

### Error: "cannot truncate table referenced by foreign key"

**Cause:** CASCADE not working (PostgreSQL version issue)

**Solution:** Ensure PostgreSQL version >= 9.5 (CASCADE support added in 9.5)

## Customization

### Preserve Additional Tables

To preserve additional tables, comment out their TRUNCATE statements in `scripts/reset_database.sql`:

```sql
-- TRUNCATE TABLE orders RESTART IDENTITY CASCADE;  -- Commented to preserve
```

### Add Additional Tables

To truncate additional tables, add them to the script:

```sql
TRUNCATE TABLE your_table_name RESTART IDENTITY CASCADE;
```

## Security Considerations

- **Production Warning:** Never run this in production without explicit approval
- **Backup First:** Always backup database before reset in production-like environments
- **Access Control:** Restrict execution to authorized personnel only
- **Audit Logging:** Consider adding audit logging for reset operations

## Related Documentation

- [PostgreSQL Migration Report](../POSTGRES_MIGRATION_REPORT.md)
- [Database Connection](../internal/infra/database/connection.go)
- [Setup Script](../setup-postgres.sh)
