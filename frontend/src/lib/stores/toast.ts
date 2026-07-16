import { writable } from 'svelte/store';

type ToastVariant = 'success' | 'error' | 'warning' | 'info';

interface Toast {
	id: string;
	variant: ToastVariant;
	title: string;
	message?: string;
	duration?: number;
}

function createToastStore() {
	const { subscribe, update } = writable<Toast[]>([]);

	return {
		subscribe,
		add: (toast: Omit<Toast, 'id'>) => {
			const id = crypto.randomUUID();
			update((toasts) => [...toasts, { ...toast, id }]);
		},
		remove: (id: string) => {
			update((toasts) => toasts.filter((t) => t.id !== id));
		},
		clear: () => {
			update(() => []);
		}
	};
}

export const toast = createToastStore();

export function showToast(options: {
	variant: ToastVariant;
	title: string;
	message?: string;
	duration?: number;
}) {
	toast.add(options);
}

export function showSuccess(title: string, message?: string, duration = 5000) {
	showToast({ variant: 'success', title, message, duration });
}

export function showError(title: string, message?: string, duration = 5000) {
	showToast({ variant: 'error', title, message, duration });
}

export function showWarning(title: string, message?: string, duration = 5000) {
	showToast({ variant: 'warning', title, message, duration });
}

export function showInfo(title: string, message?: string, duration = 5000) {
	showToast({ variant: 'info', title, message, duration });
}
