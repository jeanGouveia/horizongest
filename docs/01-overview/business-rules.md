# Business Rules

**HorizonGest Platform - Business Rules Documentation**

---

## Overview

This document defines the business rules that govern the HorizonGest platform. These rules are implemented in the Service Layer and should not be duplicated in the Frontend or Handler layers.

---

## Product Management

### Product Availability

- Products can be active or inactive
- Inactive products cannot be ordered
- Product changes do not affect historical orders

### Product Pricing

- Prices can be changed at any time
- Price changes do not affect historical orders
- Historical orders preserve the price at time of sale

---

## Order Management

### Order Creation

- Orders must have at least one product
- Orders must be associated with a company
- Orders create snapshots of products (not references)

### Order Status

- Orders progress through: Pending → Preparing → Ready → Delivered
- Orders can be cancelled at any time
- Cancelled orders preserve historical data

### Order Snapshots

- Orders contain product snapshots (not references)
- Snapshots preserve product state at time of sale
- Snapshots are immutable

---

## Inventory Management

### Stock Levels

- Stock represents current availability
- Stock decreases only during sales
- Stock increases through purchases and returns

### Stock Movements

- All stock movements must be recorded
- Movements include: In, Out, Adjustment
- Movements are auditable

### Low Stock Alerts

- Low stock alerts are triggered when stock < minimum
- Minimum stock is configurable per ingredient
- Alerts are sent to managers

---

## User Management

### User Roles

- **Admin:** Full access to all features
- **Manager:** Can manage users and settings
- **Employee:** Can manage orders and inventory
- **Viewer:** Read-only access

### User Permissions

- Permissions are role-based
- Permissions are granular per operation
- Permissions are verified in Service Layer

---

## Company Management

### Company Configuration

- Companies can configure their own branding
- Companies can configure their own settings
- Company data is isolated by CompanyID

### Company Plans

- Companies are assigned to plans
- Plans define feature limits
- Plans define user limits

---

## Multi-Tenancy

### Data Isolation

- All tenant data is isolated by CompanyID
- Platform users cannot access tenant data
- Tenant users cannot access platform data

### Tenant Configuration

- Each tenant can have its own branding
- Each tenant can have its own settings
- Each tenant can have its own users

---

## Security

### Authentication

- Users must authenticate to access the platform
- Authentication uses JWT tokens
- Tokens expire after 24 hours

### Authorization

- Users must have appropriate permissions
- Permissions are verified per operation
- RBAC is used for access control

---

**Last Updated:** Fase 2 - Documentation & Knowledge Base
