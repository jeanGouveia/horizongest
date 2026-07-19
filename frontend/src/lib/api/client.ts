// Wrapper de fetch que sempre envia cookies e trata erros de forma uniforme.
// Não inclui o token manualmente — o cookie HttpOnly é enviado pelo browser.

const BASE = '/api'; // Proxy Vite / SvelteKit repassa ao Go

interface ApiResponse<T> {
  data: T | null;
  error: string | null;
  status: number;
}

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  try {
    // Não definir Content-Type para FormData (deixa o navegador definir com boundary correto)
    const headers = options.body instanceof FormData
      ? { ...(options.headers ?? {}) }
      : {
          'Content-Type': 'application/json',
          ...(options.headers ?? {})
        };

    const res = await fetch(`${BASE}${path}`, {
      ...options,
      credentials: 'include', // sempre envia o cookie auth_token
      headers
    });

    if (res.status === 204) {
      return { data: null, error: null, status: 204 };
    }

    const json = await res.json().catch(() => ({}));

    if (!res.ok) {
      return {
        data: null,
        error: json?.error ?? `Erro ${res.status}`,
        status: res.status
      };
    }

    return { data: json as T, error: null, status: res.status };
  } catch (e: any) {
    // Tratamento de erros de rede (fetch failed, timeout, network error)
    if (e instanceof TypeError && e.message === 'Failed to fetch') {
      return {
        data: null,
        error: 'Erro de conexão. Verifique sua internet.',
        status: 0
      };
    }
    if (e.name === 'AbortError') {
      return {
        data: null,
        error: 'Tempo esgotado. Tente novamente.',
        status: 0
      };
    }
    return {
      data: null,
      error: e?.message ?? 'Erro desconhecido',
      status: 0
    };
  }
}

export { request };

// --- Auth endpoints ---

export const api = {
  auth: {
    register: (body: { name: string; email: string; password: string }) =>
      request<{ id: number; name: string; email: string }>('/auth/register', {
        method: 'POST',
        body: JSON.stringify(body)
      }),

    login: (body: { email: string; password: string }) =>
      request<{ id: number; name: string; email: string }>('/auth/login', {
        method: 'POST',
        body: JSON.stringify(body)
      }),

    logout: () =>
      request<{ message: string }>('/auth/logout', { method: 'POST' }),

    me: () =>
      request<{ id: number; name: string; email: string }>('/me'),

    updateProfile: (body: { name: string; email: string }) =>
      request<{ id: number; name: string; email: string }>('/me', {
        method: 'PUT',
        body: JSON.stringify(body)
      }),

    changePassword: (body: { current_password: string; new_password: string }) =>
      request<{ message: string }>('/me/change-password', {
        method: 'POST',
        body: JSON.stringify(body)
      })
  },

  // --- System endpoints ---
  dashboard: () =>
    request<any>('/dashboard'),

  notifications: () =>
    request<any>('/notifications'),

  health: () =>
    request<any>('/health'),

  version: () =>
    request<any>('/version'),

  capabilities: () =>
    request<any>('/capabilities'),

  // --- Dependency check endpoints ---
  canDeleteProduct: (id: number) =>
    request<any>(`/products/${id}/can-delete`),

  canDeleteIngredient: (id: number) =>
    request<any>(`/ingredients/${id}/can-delete`),

  canDeleteCategory: (id: number) =>
    request<any>(`/categories/${id}/can-delete`),

  // --- Stock validation endpoint ---
  validateStock: (body: any) =>
    request<any>('/orders/validate', {
      method: 'POST',
      body: JSON.stringify(body)
    }),

  // --- Company Settings endpoints ---
  companySettings: {
    getSettings: () =>
      request<{
        name: string;
        slug: string;
        description: string;
        logo_url: string;
        primary_color: string;
        secondary_color: string;
        business_type: string;
        locale: string;
        currency: string;
        timezone: string;
      }>('/company/settings'),

    updateSettings: (body: {
      name?: string;
      description?: string;
      logo_url?: string;
      primary_color?: string;
      secondary_color?: string;
      business_type?: string;
      locale?: string;
      currency?: string;
      timezone?: string;
    }) =>
      request<{ message: string }>('/company/settings', {
        method: 'PUT',
        body: JSON.stringify(body)
      })
  },

  // --- Theme endpoints (White Label - Platform 2.0) ---
  theme: {
    getTheme: () =>
      request<{
        primary_color: string;
        secondary_color: string;
        logo_url: string;
        font_family: string;
        border_radius: string;
        is_default: boolean;
      }>('/theme'),

    getDefaultTheme: () =>
      request<{
        primary_color: string;
        secondary_color: string;
        logo_url: string;
        font_family: string;
        border_radius: string;
        is_default: boolean;
      }>('/theme/default')
  },

  // --- Company Users endpoints (Sprint 7) ---
  companyUsers: {
    list: () =>
      request<Array<{
        id: number;
        name: string;
        email: string;
        role: string | null;
        active: boolean;
        company_id: number | null;
      }>>('/company/users'),

    add: (body: { email: string }) =>
      request<{
        id: number;
        name: string;
        email: string;
        role: string | null;
        active: boolean;
        company_id: number | null;
      }>('/company/users/add', {
        method: 'POST',
        body: JSON.stringify(body)
      }),

    changeRole: (id: number, body: { role: string }) =>
      request<{ message: string }>(`/company/users/${id}/role`, {
        method: 'PUT',
        body: JSON.stringify(body)
      }),

    remove: (id: number) =>
      request<{ message: string }>(`/company/users/${id}`, {
        method: 'DELETE'
      })
  }
} as const;
