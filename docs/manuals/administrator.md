# Administrator Manual

**HorizonGest Platform - Administrator Guide**

---

## Overview

This manual is for platform administrators who manage the entire HorizonGest platform. As a platform administrator, you have full control over platform configuration, branding, global settings, and can manage all companies using the platform.

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Platform Branding](#platform-branding)
3. [Global Configuration](#global-configuration)
4. [Feature Flags](#feature-flags)
5. [Company Management](#company-management)
6. [User Management](#user-management)
7. [Plans and Billing](#plans-and-billing)
8. [Security](#security)
9. [Support](#support)

---

## Getting Started

### Access

- **Login URL:** `https://platform.com/signin`
- **Default Email:** `admin@platform.com`
- **Default Password:** (provided during initial setup)

### First Steps

1. **Change your password** immediately after first login
2. **Configure platform branding** (logo, colors, name)
3. **Set up global configuration** (timezone, locale, formats)
4. **Create your first company** (if not already created)
5. **Configure plans** for different tiers

---

## Platform Branding

### Overview

Platform branding controls the visual identity of the entire platform. This includes the platform name, logo, colors, and other institutional information.

### Access

Navigate to: **Settings → Platform Branding**

### Configuration Options

#### Basic Information
- **Platform Name:** The name displayed throughout the platform (e.g., "HorizonGest")
- **Platform Short Name:** Abbreviated name for UI elements (e.g., "Horizon")
- **Owner Company Name:** Your company name (e.g., "HorizonGest Inc.")
- **Owner Document:** Company legal document (e.g., CNPJ, Tax ID)

#### Contact Information
- **Website:** Platform website URL
- **Support Email:** Email for support inquiries
- **Support URL:** Help center URL

#### Visual Assets
- **Logo Path:** Path to platform logo
- **Favicon Path:** Path to favicon
- **Logo Light:** Light theme logo
- **Logo Dark:** Dark theme logo
- **Icon:** Platform icon
- **Login Background:** Background image for login page
- **Login Illustration:** Illustration for login page

#### Legal Information
- **Copyright:** Copyright notice
- **Privacy Policy URL:** Privacy policy page
- **Terms URL:** Terms of service page

#### Social Media
- **Instagram URL:** Instagram profile
- **Facebook URL:** Facebook page
- **LinkedIn URL:** LinkedIn company page
- **YouTube URL:** YouTube channel

#### Localization
- **Default Language:** Default language for the platform
- **Default Timezone:** Default timezone

#### Maintenance
- **Maintenance Mode:** Enable to put platform in maintenance
- **Maintenance Message:** Message shown during maintenance

#### Colors
- **Primary Color:** Main brand color
- **Secondary Color:** Secondary brand color

### Best Practices

- Use high-quality logos (SVG preferred)
- Keep colors consistent with your brand
- Ensure all URLs are valid and accessible
- Test branding changes in staging first

---

## Global Configuration

### Overview

Global configuration controls technical settings that apply to the entire platform.

### Access

Navigate to: **Settings → Global Configuration**

### Configuration Options

#### Localization
- **Default Timezone:** Platform timezone (e.g., "America/Sao_Paulo")
- **Default Locale:** Default language/locale (e.g., "pt-BR")
- **Monetary Format:** Currency display format (e.g., "BRL R$ 1.000,00")
- **Date Format:** Date display format (e.g., "DD/MM/YYYY")
- **Time Format:** Time display format (e.g., "HH:mm")

#### File Uploads
- **Max Upload Size (MB):** Maximum file upload size
- **Max Image Size (MB):** Maximum image upload size
- **Allowed Image Types:** Comma-separated image types (e.g., "jpg,png,webp,gif")
- **Allowed File Types:** Comma-separated file types (e.g., "pdf,doc,docx,xlsx,xls,txt")

#### Maintenance
- **Maintenance Mode:** Enable to put platform in maintenance
- **Maintenance Message:** Message shown users during maintenance

#### Feature Flags
- **Enable Finance:** Enable finance module
- **Enable Purchasing:** Enable purchasing module
- **Enable Inventory:** Enable inventory module
- **Enable CRM:** Enable CRM module
- **Enable Calendar:** Enable calendar module
- **Enable POS:** Enable POS module
- **Enable AI:** Enable AI features
- **Enable Delivery:** Enable delivery module
- **Enable Marketplace:** Enable marketplace module

### Best Practices

- Set timezone to match your primary market
- Configure file size limits based on your infrastructure
- Enable features gradually as they become available
- Test configuration changes in staging first

---

## Feature Flags

### Overview

Feature flags allow you to enable or disable platform modules globally. This is useful for:
- Gradual feature rollout
- A/B testing
- Controlling access to beta features
- Disabling problematic features

### Access

Navigate to: **Settings → Global Configuration → Feature Flags**

### Available Modules

- **Finance:** Financial management, accounting, reports
- **Purchasing:** Supplier management, purchase orders
- **Inventory:** Stock management, inventory tracking
- **CRM:** Customer relationship management
- **Calendar:** Scheduling, appointments
- **POS:** Point of sale system
- **AI:** AI-powered features and recommendations
- **Delivery:** Delivery management
- **Marketplace:** Multi-vendor marketplace

### Best Practices

- Enable features only when they are fully tested
- Disable features immediately if issues arise
- Communicate feature changes to users
- Monitor feature usage after enabling

---

## Company Management

### Overview

Company management allows you to create, configure, and manage companies (tenants) using the platform.

### Access

Navigate to: **Companies → All Companies**

### Creating a Company

1. Click **"New Company"**
2. Fill in company information:
   - Company Name
   - Document (CNPJ, Tax ID)
   - Email
   - Phone
   - Address
3. Select a plan
4. Click **"Create"**

### Company Configuration

#### Basic Information
- Name, document, contact information
- Address details
- Business hours

#### Branding
- Logo
- Colors
- Theme customization

#### Settings
- Timezone
- Locale
- Feature access

#### Users
- Invite users
- Assign roles
- Manage permissions

### Managing Companies

- **View:** Click on a company to view details
- **Edit:** Modify company settings
- **Suspend:** Temporarily suspend company access
- **Delete:** Permanently delete company (with confirmation)

### Best Practices

- Verify company information before approval
- Assign appropriate plans based on company needs
- Monitor company activity
- Provide support for new companies

---

## User Management

### Overview

User management allows you to manage platform administrators and company users.

### Platform Administrators

#### Creating Platform Administrators

1. Navigate to: **Users → Platform Users**
2. Click **"New Platform User"**
3. Fill in user information:
   - Name
   - Email
   - Role (Admin, Manager, Viewer)
4. Click **"Create"**
5. User will receive an email with setup instructions

#### Managing Platform Administrators

- **View:** View user details
- **Edit:** Modify user information
- **Deactivate:** Disable user access
- **Delete:** Remove user (with confirmation)

### Company Users

Company users are managed by company administrators. Platform administrators can:

- View all company users
- Assist with user issues
- Reset passwords (if needed)

### Best Practices

- Use principle of least privilege
- Regularly review user access
- Remove access for inactive users
- Use strong password policies

---

## Plans and Billing

### Overview

Plans define the features and limits available to companies. Billing is based on the selected plan.

### Access

Navigate to: **Plans → All Plans**

### Creating a Plan

1. Click **"New Plan"**
2. Fill in plan information:
   - Name (e.g., "Basic", "Premium", "Enterprise")
   - Description
   - Price
   - Currency
   - Interval (Monthly, Yearly)
   - Max Users
   - Max Products
   - Features (JSON format)
   - Status (Active, Inactive)
3. Click **"Create"**

### Plan Features

Features are defined in JSON format. Example:

```json
{
  "finance": true,
  "purchasing": true,
  "inventory": true,
  "crm": false,
  "calendar": false,
  "pos": false,
  "ai": false,
  "delivery": false,
  "marketplace": false
}
```

### Managing Plans

- **View:** View plan details
- **Edit:** Modify plan settings
- **Activate/Deactivate:** Enable or disable plan
- **Delete:** Remove plan (if not in use)

### Best Practices

- Create plans for different business sizes
- Clearly communicate plan differences
- Regularly review plan usage
- Adjust plans based on feedback

---

## Security

### Overview

Platform security is critical. As a platform administrator, you must ensure the platform remains secure.

### Security Best Practices

#### Password Security
- Enforce strong password policies
- Require password changes periodically
- Use secure password reset flows

#### Access Control
- Use principle of least privilege
- Regularly review user access
- Remove access for inactive users

#### Data Protection
- Regular backups
- Secure data transmission (HTTPS)
- Encrypt sensitive data

#### Monitoring
- Monitor login attempts
- Track suspicious activity
- Review audit logs

#### Updates
- Keep platform updated
- Apply security patches promptly
- Test updates in staging first

---

## Support

### Overview

As a platform administrator, you may need to provide support to companies and users.

### Common Issues

#### Login Issues
- Verify user credentials
- Check if account is active
- Reset password if needed

#### Configuration Issues
- Review platform branding
- Check global configuration
- Verify feature flags

#### Company Issues
- Review company settings
- Check plan limits
- Verify user permissions

### Getting Help

- **Documentation:** Available at `/docs`
- **Technical Support:** support@platform.com
- **Emergency:** emergency@platform.com

---

## Summary

As a platform administrator, you are responsible for:

- ✅ Configuring platform branding
- ✅ Managing global configuration
- ✅ Controlling feature flags
- ✅ Managing companies
- ✅ Managing platform users
- ✅ Configuring plans and billing
- ✅ Ensuring platform security
- ✅ Providing support

For detailed technical information, refer to the [Architecture Documentation](../01-overview/architecture.md).

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
