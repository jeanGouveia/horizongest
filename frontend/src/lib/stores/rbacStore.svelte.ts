import { api } from '$lib/api/client';

// Role enum matching backend
export type Role = 'owner' | 'admin' | 'manager' | 'cashier' | 'kitchen' | 'waiter';

// Permission enum matching backend
export type Permission = 
  | 'manage_company'
  | 'manage_products'
  | 'manage_orders'
  | 'manage_users'
  | 'manage_settings'
  | 'view_reports';

interface MeResponse {
  id: number;
  name: string;
  email: string;
  role?: Role | null;
}

interface RBACState {
  role: Role | null;
  permissions: Permission[];
  loaded: boolean;
}

class RBACStore {
  private state = $state<RBACState>({
    role: null,
    permissions: [],
    loaded: false
  });

  get role(): Role | null {
    return this.state.role;
  }

  get permissions(): Permission[] {
    return this.state.permissions;
  }

  get loaded(): boolean {
    return this.state.loaded;
  }

  // Load user role and permissions from backend
  async load(): Promise<void> {
    try {
      const res = await api.auth.me();
      if (res.data) {
        // Role will be included in /api/me response in future sprints
        // For now, we'll set it to null and derive permissions from it
        this.state.role = (res.data as MeResponse).role || null;
        this.state.permissions = this.derivePermissions(this.state.role);
        this.state.loaded = true;
      }
    } catch (error) {
      console.error('Failed to load RBAC data:', error);
      this.state.loaded = true;
    }
  }

  // Check if user has a specific role
  hasRole(role: Role): boolean {
    return this.state.role === role;
  }

  // Check if user has any of the specified roles
  hasAnyRole(roles: Role[]): boolean {
    if (!this.state.role) return false;
    return roles.includes(this.state.role);
  }

  // Check if user has a specific permission
  can(permission: Permission): boolean {
    return this.state.permissions.includes(permission);
  }

  // Derive permissions based on role
  private derivePermissions(role: Role | null): Permission[] {
    if (!role) return [];

    switch (role) {
      case 'owner':
        return [
          'manage_company',
          'manage_products',
          'manage_orders',
          'manage_users',
          'manage_settings',
          'view_reports'
        ];
      case 'admin':
        return [
          'manage_company',
          'manage_products',
          'manage_orders',
          'manage_settings',
          'view_reports'
        ];
      case 'manager':
        return [
          'manage_products',
          'manage_orders',
          'view_reports'
        ];
      case 'cashier':
        return [
          'manage_orders'
        ];
      case 'kitchen':
        return [
          'manage_orders'
        ];
      case 'waiter':
        return [
          'manage_orders'
        ];
      default:
        return [];
    }
  }

  // Reset store (for logout)
  reset(): void {
    this.state = {
      role: null,
      permissions: [],
      loaded: false
    };
  }
}

// Export singleton instance
export const rbacStore = new RBACStore();
