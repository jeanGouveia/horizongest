import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	try {
		const response = await fetch('/api/stock-adjustments/pending');
		if (!response.ok) {
			throw new Error('Erro ao carregar ajustes');
		}
		const adjustments = await response.json();

		return {
			adjustments
		};
	} catch (error) {
		console.error('Erro ao carregar ajustes:', error);
		return {
			adjustments: [],
			error: error instanceof Error ? error.message : 'Erro desconhecido'
		};
	}
};
