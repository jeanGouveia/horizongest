import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

// Guard SSR: qualquer rota dentro de (app)/ sem sessão → /login
export const load: LayoutServerLoad = ({ locals }) => {
  if (!locals.user) {
    redirect(302, '/login');
  }

  return { user: locals.user };
};
