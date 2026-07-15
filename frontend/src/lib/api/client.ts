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
  const res = await fetch(`${BASE}${path}`, {
    ...options,
    credentials: 'include', // sempre envia o cookie auth_token
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers ?? {})
    }
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
  }
};
