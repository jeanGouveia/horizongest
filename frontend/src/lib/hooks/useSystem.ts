import { onMount } from 'svelte';
import { api } from '$lib/api/client';
import type { Health, Version, Capabilities } from '$lib/types/system';

export function useSystem() {
	let health = $state<Health | null>(null);
	let version = $state<Version | null>(null);
	let capabilities = $state<Capabilities | null>(null);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			const [healthRes, versionRes, capabilitiesRes] = await Promise.all([
				api.health(),
				api.version(),
				api.capabilities()
			]);

			if (!healthRes.error) health = healthRes.data;
			if (!versionRes.error) version = versionRes.data;
			if (!capabilitiesRes.error) capabilities = capabilitiesRes.data;
		} catch (e: any) {
			error = e?.message ?? 'Erro ao carregar informações do sistema.';
		} finally {
			loading = false;
		}
	});

	return {
		get health() { return health; },
		get version() { return version; },
		get capabilities() { return capabilities; },
		get loading() { return loading; },
		get error() { return error; }
	};
}
