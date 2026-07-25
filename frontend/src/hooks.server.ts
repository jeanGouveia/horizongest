import type { Handle } from '@sveltejs/kit';

// hooks.server.ts é executado em TODA requisição SSR.
// Aqui validamos o cookie e populamos event.locals.user
// para que os layouts e pages acessem sem refazer fetch.

const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080';

export const handle: Handle = async ({ event, resolve }) => {
  const token = event.cookies.get('auth_token');

  if (token) {
    try {
      // Chama o próprio backend Go para validar o token
      // (o cookie é repassado manualmente pois estamos no servidor Node)
      const res = await fetch(`${backendUrl}/api/me`, {
        headers: { Cookie: `auth_token=${token}` }
      });

      if (res.ok) {
        event.locals.user = await res.json();
      } else {
        event.locals.user = null;
      }
    } catch (error) {
      event.locals.user = null;
    }
  } else {
    event.locals.user = null;
  }

  return resolve(event);
};
