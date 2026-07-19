# PLATAFORMA 2.0 - SPRINT 5: COMPANY SETTINGS FOUNDATION

## Overview

Company Settings Foundation is the fifth sprint of Plataforma PratoOnline 2.0, implementing the first administrative panel for companies. This feature allows each company to configure their own identity and business settings, building upon the infrastructure created in the previous four sprints (Tenant Engine, White Label Foundation, Business Engine, and Tenant Isolation).

## Architecture

### Core Components

1. **CompanySettingsService** (`internal/service/company_settings_service.go`)
   - Dedicated service for company settings management
   - Separated from CompanyService to maintain single responsibility
   - Methods:
     - `GetSettings(ctx, userID)`: Retrieves settings for user's company
     - `UpdateSettings(ctx, userID, input)`: Updates settings for user's company
   - Security: Only allows updating the tenant's own CompanyID
   - Validation: Returns `ErrUserNoCompany` if user has no CompanyID

2. **CompanySettingsHandler** (`internal/handler/company_settings_handler.go`)
   - HTTP handlers for company settings endpoints
   - Methods:
     - `GetSettings`: Handles GET /api/company/settings
     - `UpdateSettings`: Handles PUT /api/company/settings
   - Security: Extracts userID from context, validates tenant ownership
   - Error handling: Returns 403 for users without CompanyID

### Data Flow

```
Frontend → API Client → Handler → Service → Repository → Database
            (companySettings)         (tenant isolation)
```

1. **Request**: Frontend calls API endpoints
2. **Authentication**: AuthMiddleware validates JWT, sets userID
3. **Tenant Resolution**: TenantMiddleware loads CompanyID from user
4. **Business Logic**: Service validates tenant ownership
5. **Data Access**: Repository applies tenant isolation
6. **Response**: Settings data returned or updated

## Endpoints

### GET /api/company/settings

**Description**: Retrieves company settings for the authenticated user's company

**Authentication**: Required (JWT cookie)

**Response** (200 OK):
```json
{
  "name": "Company Name",
  "slug": "company-slug",
  "description": "Company description",
  "logo_url": "https://example.com/logo.png",
  "primary_color": "#3b82f6",
  "secondary_color": "#1e40af",
  "business_type": "restaurant",
  "locale": "pt-BR",
  "currency": "BRL",
  "timezone": "America/Sao_Paulo"
}
```

**Error Response** (403 Forbidden):
```json
{
  "error": "usuário não possui uma empresa associada"
}
```

### PUT /api/company/settings

**Description**: Updates company settings for the authenticated user's company

**Authentication**: Required (JWT cookie)

**Request Body** (all fields optional):
```json
{
  "name": "Updated Company Name",
  "description": "Updated description",
  "logo_url": "https://example.com/new-logo.png",
  "primary_color": "#ff0000",
  "secondary_color": "#00ff00",
  "business_type": "bakery",
  "locale": "en-US",
  "currency": "USD",
  "timezone": "America/New_York"
}
```

**Response** (200 OK):
```json
{
  "message": "configurações atualizadas com sucesso"
}
```

**Error Response** (403 Forbidden):
```json
{
  "error": "usuário não possui uma empresa associada"
}
```

## Security

### Tenant Isolation

- **CompanyID Validation**: Service validates user's CompanyID before allowing access
- **No CompanyID Override**: Frontend cannot send `company_id` in request body
- **Tenant-Only Updates**: Users can only update their own company's settings
- **Core V1 Protection**: Users without CompanyID receive 403 Forbidden

### Security Flow

1. **Authentication**: JWT validated by AuthMiddleware
2. **Tenant Resolution**: TenantMiddleware loads CompanyID from user
3. **Service Validation**: CompanySettingsService checks user has CompanyID
4. **Repository Isolation**: Tenant isolation prevents cross-tenant access
5. **Error Handling**: Returns 403 for unauthorized access attempts

## Frontend Implementation

### Page Structure

**Route**: `/settings/company`

**File**: `frontend/src/routes/(app)/settings/company/+page.svelte`

**Components**:
- Workspace layout with breadcrumb navigation
- Three main sections: General, Branding, Business
- Color preview with live theme updates
- Save button with loading state

### Form Sections

#### 1. Dados Gerais (General)

- **Nome da Empresa**: Required text input
- **Slug**: Read-only identifier (auto-generated)
- **Descrição**: Optional multiline textarea

#### 2. Branding

- **Logo URL**: Optional URL input (no file upload in this sprint)
- **Cor Primária**: Color picker with live preview
- **Cor Secundária**: Color picker with live preview
- **Preview**: Live preview of buttons with selected colors

#### 3. Negócio (Business)

- **Tipo do Negócio**: Select dropdown (restaurant, bakery, cafe, bar, food_truck, catering, other)
- **Idioma**: Select dropdown (pt-BR, en-US, es-ES)
- **Moeda**: Select dropdown (BRL, USD, EUR)
- **Timezone**: Select dropdown (America/Sao_Paulo, America/New_York, Europe/London, Asia/Tokyo)

### Color Preview

The color preview section shows real-time theme updates:

- **Primary Button**: Shows primary color with hover state using secondary color
- **Secondary Button**: Shows secondary color
- **Live Updates**: Preview updates as user changes colors
- **No Persistence**: Changes only persist after clicking "Salvar"

### API Integration

**API Client** (`frontend/src/lib/api/client.ts`):

```typescript
companySettings: {
  getSettings: () => request<SettingsResponse>('/company/settings'),
  updateSettings: (body: UpdateSettingsInput) => 
    request<{ message: string }>('/company/settings', {
      method: 'PUT',
      body: JSON.stringify(body)
    })
}
```

### Theme Engine Integration

After saving settings, the frontend automatically reloads the theme:

```typescript
async function saveSettings() {
  // ... save logic
  await themeStore.loadTheme(); // Reload theme to reflect changes
}
```

This ensures that:
- Primary color changes are immediately applied
- Secondary color changes are immediately applied
- Logo URL changes are immediately applied
- No backend restart required

## Business Engine Integration

The Business Profile automatically reflects settings changes:

**Business Profile Endpoint** (`/api/business/profile`):

After updating company settings, the Business Profile returns updated values:

```json
{
  "CompanyID": 5,
  "CompanyName": "Updated Settings Test Company",
  "BusinessType": "bakery",
  "Locale": "en-US",
  "Currency": "USD",
  "Timezone": "America/New_York",
  "LogoURL": "https://example.com/test-logo.png",
  "PrimaryColor": "#ff0000",
  "SecondaryColor": "#00ff00"
}
```

This ensures consistency across:
- Company Settings panel
- Business Profile API
- Theme Engine
- All tenant-aware features

## Tests Performed

### Test 1: Alter company name

**Setup**: User with CompanyID updates company name

**Result**: ✅ Company name updated successfully

```bash
curl -X PUT /api/company/settings -d '{"name":"Updated Settings Test Company"}'
# Response: {"message":"configurações atualizadas com sucesso"}
```

### Test 2: Alter description

**Setup**: User with CompanyID updates company description

**Result**: ✅ Description updated successfully

```bash
curl -X PUT /api/company/settings -d '{"description":"Updated description"}'
# Response: {"message":"configurações atualizadas com sucesso"}
```

### Test 3: Alter primary color

**Setup**: User with CompanyID updates primary color

**Result**: ✅ Primary color updated successfully

```bash
curl -X PUT /api/company/settings -d '{"primary_color":"#ff0000"}'
# Response: {"message":"configurações atualizadas com sucesso"}
```

### Test 4: Alter secondary color

**Setup**: User with CompanyID updates secondary color

**Result**: ✅ Secondary color updated successfully

```bash
curl -X PUT /api/company/settings -d '{"secondary_color":"#00ff00"}'
# Response: {"message":"configurações atualizadas com sucesso"}
```

### Test 5: Alter logo

**Setup**: User with CompanyID updates logo URL

**Result**: ✅ Logo URL updated successfully

```bash
curl -X PUT /api/company/settings -d '{"logo_url":"https://example.com/test-logo.png"}'
# Response: {"message":"configurações atualizadas com sucesso"}
```

### Test 6: Alter business type

**Setup**: User with CompanyID updates business type

**Result**: ✅ Business type updated successfully

```bash
curl -X PUT /api/company/settings -d '{"business_type":"bakery"}'
# Response: {"message":"configurações atualizadas com sucesso"}
```

### Test 7: Alter locale

**Setup**: User with CompanyID updates locale

**Result**: ✅ Locale updated successfully

```bash
curl -X PUT /api/company/settings -d '{"locale":"en-US"}'
# Response: {"message":"configurações atualizadas com sucesso"}
```

### Test 8: Alter currency

**Setup**: User with CompanyID updates currency

**Result**: ✅ Currency updated successfully

```bash
curl -X PUT /api/company/settings -d '{"currency":"USD"}'
# Response: {"message":"configurações atualizadas com sucesso"}
```

### Test 9: Alter timezone

**Setup**: User with CompanyID updates timezone

**Result**: ✅ Timezone updated successfully

```bash
curl -X PUT /api/company/settings -d '{"timezone":"America/New_York"}'
# Response: {"message":"configurações atualizadas com sucesso"}
```

### Test 10: Theme Engine reflects changes

**Setup**: Update primary and secondary colors, then check theme endpoint

**Result**: ✅ Theme Engine reflects changes immediately

```bash
curl -X GET /api/theme
# Response: {"PrimaryColor":"#ff0000","SecondaryColor":"#00ff00",...}
```

### Test 11: Business Profile reflects changes

**Setup**: Update business settings, then check business profile endpoint

**Result**: ✅ Business Profile reflects changes immediately

```bash
curl -X GET /api/business/profile
# Response: {"BusinessType":"bakery","Locale":"en-US","Currency":"USD",...}
```

### Test 12: Tenant Isolation preserved

**Setup**: User without CompanyID attempts to update settings

**Result**: ✅ Tenant isolation preserved, returns 403

```bash
curl -X PUT /api/company/settings -d '{"name":"Hacked Company"}'
# Response: {"error":"usuário não possui uma empresa associada"}
# HTTP Status: 403
```

### Test 13: Core V1 returns 403

**Setup**: V1 user (no CompanyID) attempts to access settings

**Result**: ✅ Core V1 user receives 403 Forbidden

```bash
curl -X GET /api/company/settings
# Response: {"error":"usuário não possui uma empresa associada"}
# HTTP Status: 403
```

## Integration with Previous Sprints

### Tenant Engine (Sprint 1)

- **Usage**: Company entity and user-company association
- **Integration**: Service validates user's CompanyID from Tenant Engine
- **Benefit**: Leverages existing multi-tenant infrastructure

### White Label Foundation (Sprint 2)

- **Usage**: Theme Engine for color and logo updates
- **Integration**: Settings updates trigger theme reload
- **Benefit**: Immediate visual feedback without backend restart

### Business Engine (Sprint 3)

- **Usage**: Business profile fields (business_type, locale, currency, timezone)
- **Integration**: Settings updates reflected in Business Profile
- **Benefit**: Consistent business identity across all features

### Tenant Isolation (Sprint 4)

- **Usage**: Repository-level tenant filtering
- **Integration**: Company updates respect tenant isolation
- **Benefit**: Security through existing isolation infrastructure

## Risks and Mitigations

### Risk 1: Cross-Tenant Access

**Risk**: Users might access other companies' settings

**Mitigation**:
- Service validates user's CompanyID before access
- Repository applies tenant isolation
- No CompanyID override allowed in request body
- Returns 403 for unauthorized access

### Risk 2: Theme Cache Issues

**Risk**: Theme changes not reflected immediately

**Mitigation**:
- Frontend reloads theme after successful save
- Theme Engine loads fresh data from Company entity
- No backend caching of theme data
- Immediate visual feedback

### Risk 3: Business Profile Inconsistency

**Risk**: Business Profile not reflecting settings changes

**Mitigation**:
- Both endpoints read from same Company entity
- No separate caching layer
- Real-time data consistency
- Single source of truth

### Risk 4: Core V1 Confusion

**Risk**: Core V1 users confused by 403 error

**Mitigation**:
- Clear error message: "usuário não possui uma empresa associada"
- Consistent 403 status code
- Documentation explains requirement
- Future: Add guidance to create company

## Compatibility

### Core V1 Compatibility

**Status**: ✅ Compatible (with expected 403)

**Details**:
- Core V1 users (no CompanyID) receive 403 Forbidden
- Clear error message explains requirement
- No breaking changes to existing Core V1 functionality
- Users must create Company to access settings

### Platform 2.0 Compatibility

**Status**: ✅ Fully Compatible

**Details**:
- Builds upon all previous sprints
- No changes to existing infrastructure
- Leverages Tenant Engine, White Label, Business Engine, Tenant Isolation
- Consistent with multi-tenant architecture

## Migration Path

**For Core V1 Users**:
1. Create a Company entity via `/api/companies`
2. Assign User to Company via `/api/me` endpoint
3. Access Company Settings via `/settings/company`
4. Configure company identity and business settings

**For Platform 2.0 Users**:
1. Access Settings → Company page
2. Configure general, branding, and business settings
3. Save changes
4. Theme and Business Profile update automatically

## Next Steps

### Immediate (Sprint 6)

1. **File Upload Integration**
   - Implement real logo upload via `/api/media/upload`
   - Replace Logo URL input with file upload
   - Add image preview before upload

2. **Settings Expansion**
   - Add more business settings
   - Add notification preferences
   - Add email configuration

3. **RBAC Integration**
   - Restrict settings access to admin roles
   - Add permission checks
   - Implement role-based UI

### Future Enhancements

1. **Advanced Branding**
   - Custom fonts
   - Custom CSS
   - Multiple logo sizes
   - Favicon upload

2. **Business Configuration**
   - Operating hours
   - Delivery zones
   - Payment methods
   - Tax configuration

3. **Settings Validation**
   - Real-time field validation
   - Business type-specific settings
   - Locale-specific formats
   - Currency-specific formatting

## Conclusion

Company Settings Foundation has been successfully implemented with:

- ✅ Dedicated CompanySettingsService for settings management
- ✅ Secure GET/PUT endpoints with tenant validation
- ✅ Frontend Settings/Company page with form sections
- ✅ Color preview with live theme updates
- ✅ Integration with Theme Engine (immediate reflection)
- ✅ Integration with Business Engine (immediate reflection)
- ✅ Tenant isolation preserved (403 for unauthorized)
- ✅ Core V1 compatibility (403 for users without Company)
- ✅ All field updates tested and working
- ✅ No changes to existing infrastructure
- ✅ Builds upon all previous sprints

The implementation provides the first administrative panel for Plataforma PratoOnline 2.0, allowing each tenant to customize their identity and business configuration while maintaining security and consistency across the platform.
